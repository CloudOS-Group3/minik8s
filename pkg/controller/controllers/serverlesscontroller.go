package controllers

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/kafka"
	"minik8s/pkg/serverless/function"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"strings"
	"sync"
	"time"
)

const (
	MaxFreeTime   time.Duration = 60 * time.Second
	CheckInterval time.Duration = 10 * time.Second
)

type PodWithStatus struct {
	pod      *api.Pod
	isFree   bool
	freeTime time.Duration
}

type ServerlessController struct {
	functionFreePods map[string][]PodWithStatus
	subscriber       *kafka.Subscriber
	ready            chan bool
	done             chan bool
}

func NewServerlessController() *ServerlessController {
	group := "serverless-controller"
	Controller := &ServerlessController{
		ready:            make(chan bool),
		done:             make(chan bool),
		functionFreePods: make(map[string][]PodWithStatus),
		subscriber:       kafka.NewSubscriber(group),
	}
	return Controller
}

func (this *ServerlessController) Setup(_ sarama.ConsumerGroupSession) error {
	close(this.ready)
	return nil
}

func (this *ServerlessController) Cleanup(_ sarama.ConsumerGroupSession) error {
	this.ready = make(chan bool)
	return nil
}

func (this *ServerlessController) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info("Message claimed: value %s", string(msg.Value))
		switch msg.Topic {
		case msg_type.TriggerTopic:
			sess.MarkMessage(msg, "")
			this.triggerNewJob(msg.Value)
		case msg_type.JobTopic:
			sess.MarkMessage(msg, "")
			this.updateJob(msg.Value)
		case msg_type.FunctionTopic:
			sess.MarkMessage(msg, "")
			this.DeleteFunction(msg.Value)
		default:
			log.Warn("Unknown msg type in serverless controller")
		}
	}
	return nil
}

func (this *ServerlessController) triggerNewJob(content []byte) {
	var triggerMsg msg_type.TriggerMsg
	err := json.Unmarshal(content, &triggerMsg)
	if err != nil {
		log.Error("json unmarshal err %v", err)
		return
	}

	var freePod *api.Pod
	found := false
	functionName := triggerMsg.Function.Metadata.Name
	for index, functionFreePod := range this.functionFreePods[functionName] {
		if functionFreePod.isFree {
			found = true
			freePod = functionFreePod.pod
			this.functionFreePods[functionName][index].isFree = false
			this.functionFreePods[functionName][index].freeTime = 0
			break
		}
	}
	if !found {
		freePod = function.CreatePodFromFunction(&triggerMsg.Function)
		newPodWithStatus := PodWithStatus{pod: freePod, isFree: false, freeTime: 0}
		this.functionFreePods[functionName] = append(this.functionFreePods[functionName], newPodWithStatus)
	}
	if freePod == nil {
		log.Error("freePod shouldn't be nil")
		return
	}

	var randomString = stringutil.GenerateRandomString(5)
	freePod.Metadata.Name += "-" + randomString
	freePod.Spec.Containers[0].Name += "-" + randomString
	var job api.Job
	job.JobID = triggerMsg.UUID
	job.CreateTime = time.Now().String()
	job.Instance = *freePod
	job.Params = triggerMsg.Params
	job.Status = api.JOB_CREATED
	job.Function = functionName

	URL := config.GetUrlPrefix() + config.JobsURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
	bytes, err := json.Marshal(job)
	if err != nil {
		log.Error("json marshal err %v", err)
		return
	}
	err = httputil.Post(URL, bytes)
	if err != nil {
		log.Error("post err %v", err)
		return
	}
	log.Info("successfully create new job")
}

func (this *ServerlessController) updateJob(content []byte) {
	log.Debug("updating job")
	var jobMsg msg_type.JobMsg
	err := json.Unmarshal(content, &jobMsg)
	if err != nil {
		log.Error("json unmarshal err %v", err)
		return
	}
	switch jobMsg.Opt {
	case msg_type.Add:
		log.Warn("why will there be add msg of job in a place other than serverless controller?")
	case msg_type.Update:
		for functionName, freePods := range this.functionFreePods {
			for index, freePod := range freePods {
				if freePod.pod.Metadata.Name == jobMsg.OldJob.Instance.Metadata.Name {
					switch jobMsg.NewJob.Status {
					case api.JOB_ENDED:
						this.functionFreePods[functionName][index].isFree = true
						this.functionFreePods[functionName][index].freeTime = 0
					case api.JOB_RUNNING:
						this.functionFreePods[functionName][index].isFree = false
						this.functionFreePods[functionName][index].freeTime = 0
					}
				}
			}
		}
	case msg_type.Delete:
		for functionName, freePods := range this.functionFreePods {
			for index, freePod := range freePods {
				if freePod.pod.Metadata.Name == jobMsg.OldJob.Instance.Metadata.Name {
					this.functionFreePods[functionName][index].isFree = true
					this.functionFreePods[functionName][index].freeTime = 0
				}
			}
		}
	default:
		log.Warn("unknown operation %v", jobMsg.Opt)
	}
}

func (this *ServerlessController) DeleteFunction(content []byte) {
	log.Debug("deleting function")
	var functionMsg msg_type.FunctionMsg
	err := json.Unmarshal(content, &functionMsg)
	if err != nil {
		log.Error("json unmarshal err %v", err)
		return
	}

	if functionMsg.Opt != msg_type.Delete {
		log.Error("we don't support operations other than delete")
		return
	}

	functionName := functionMsg.OldFunctionName
	for _, functionFreePod := range this.functionFreePods[functionName] {
		URL := config.GetUrlPrefix() + config.PodURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		URL = strings.Replace(URL, config.NamePlaceholder, functionFreePod.pod.Metadata.Name, -1)

		err := httputil.Delete(URL)
		if err != nil {
			log.Error("delete err %v", err)
			return
		}
	}

	URL := config.GetUrlPrefix() + config.FunctionURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
	URL = strings.Replace(URL, config.NamePlaceholder, functionName, -1)

	err = httputil.Delete(URL)
	if err != nil {
		log.Error("delete err %v", err)
		return
	}

	delete(this.functionFreePods, functionName)
	log.Info("successfully delete function in serverless controller")
}

func (this *ServerlessController) clearExpirePod() {
	for {
		<-time.After(CheckInterval)
		log.Debug("checking expire pod")
		for functionName, _ := range this.functionFreePods {
			for index := 0; index < len(this.functionFreePods[functionName]); index++ {
				log.Info("freepod: %d, array length: %d", index, len(this.functionFreePods[functionName]))
				freePod := this.functionFreePods[functionName][index]
				if freePod.freeTime >= MaxFreeTime && freePod.pod != nil {
					log.Debug("pod %s is expired", freePod.pod.Metadata.Name)
					URL := config.GetUrlPrefix() + config.PodURL
					URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
					URL = strings.Replace(URL, config.NamePlaceholder, freePod.pod.Metadata.Name, -1)

					err := httputil.Delete(URL)
					if err != nil {
						log.Error("delete pod err %v", err)
						return
					}
					if index == len(this.functionFreePods[functionName])-1 {
						log.Info("delete")
						this.functionFreePods[functionName] = this.functionFreePods[functionName][:index]
						break
					} else {
						log.Info("detele")
						this.functionFreePods[functionName] = append(this.functionFreePods[functionName][:index], this.functionFreePods[functionName][index+1:]...)
						index--
						continue
					}
				}
				this.functionFreePods[functionName][index].freeTime += CheckInterval
			}
			if len(this.functionFreePods[functionName]) == 0 {
				delete(this.functionFreePods, functionName)
			}
		}
	}
}

func (this *ServerlessController) Run() {
	go this.clearExpirePod()
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	//topics := []string{msg_type.TriggerTopic, msg_type.JobTopic, msg_type.FunctionTopic}
	topics := []string{msg_type.TriggerTopic, msg_type.JobTopic}
	this.subscriber.Subscribe(wg, ctx, topics, this)
	<-this.ready
	<-this.done
	cancel()
	wg.Wait()
}
