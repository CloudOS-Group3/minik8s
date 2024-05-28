package controllers

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/kafka"
	"minik8s/pkg/serverless/function"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"strings"
	"sync"
	"time"
)

const (
	MaxFreeTime   time.Duration = 30 * time.Second
	CheckInterval time.Duration = 10 * time.Second
)

type PodWithStatus struct {
	pod      *api.Pod
	isFree   bool
	freeTime time.Duration
}

type ServerlessController struct {
	freePods   []PodWithStatus //uuid of pod -> isFree
	subscriber *kafka.Subscriber
	ready      chan bool
	done       chan bool
}

func NewServerlessController() *ServerlessController {
	group := "serverless-controller"
	Controller := &ServerlessController{
		ready:      make(chan bool),
		done:       make(chan bool),
		freePods:   []PodWithStatus{},
		subscriber: kafka.NewSubscriber(group),
	}
	return Controller
}

func (this *ServerlessController) Setup(_ sarama.ConsumerGroupSession) error {
	close(this.ready)
	return nil
}

func (this *ServerlessController) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (this *ServerlessController) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg.Topic == msg_type.TriggerTopic {
			sess.MarkMessage(msg, "")
			this.triggerNewJob(msg.Value)
		} else if msg.Topic == msg_type.JobTopic {
			sess.MarkMessage(msg, "")
			this.updateJob(msg.Value)
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
	for index, podWithStatus := range this.freePods {
		if podWithStatus.isFree {
			found = true
			freePod = podWithStatus.pod
			this.freePods[index].isFree = false
		}
	}
	if !found {
		freePod = function.CreatePodFromFunction(&triggerMsg.Function)
		newPodWithStatus := PodWithStatus{pod: freePod, isFree: false, freeTime: 0}
		this.freePods = append(this.freePods, newPodWithStatus)
	}
	if freePod == nil {
		log.Error("freePod shouldn't be nil")
		return
	}

	var job api.Job
	job.JobID = uuid.NewString()
	job.CreateTime = time.Now().String()
	job.Instance = *freePod
	job.Params = triggerMsg.Params
	job.Status = api.JOB_CREATED

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
		for index, podWithStatus := range this.freePods {
			if podWithStatus.pod.Metadata.Name == jobMsg.OldJob.Instance.Metadata.Name {
				switch jobMsg.NewJob.Status {
				case api.JOB_DELETED:
					this.freePods[index].isFree = true
					this.freePods[index].freeTime = 0
				case api.JOB_RUNNING:
					this.freePods[index].isFree = false
					this.freePods[index].freeTime = 0
				}
			}
		}
	case msg_type.Delete:
		for index, podWithStatus := range this.freePods {
			if podWithStatus.pod.Metadata.Name == jobMsg.OldJob.Instance.Metadata.Name {
				this.freePods[index].isFree = true
				this.freePods[index].freeTime = 0
			}
		}
	default:
		log.Warn("unknown operation %v", jobMsg.Opt)
	}
}

func (this *ServerlessController) clearExpirePod() {
	for {
		<-time.After(CheckInterval)
		for index, podWithStatus := range this.freePods {
			if podWithStatus.freeTime >= MaxFreeTime {
				URL := config.GetUrlPrefix() + config.PodURL
				URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
				URL = strings.Replace(URL, config.NamePlaceholder, podWithStatus.pod.Metadata.Name, -1)

				err := httputil.Delete(URL)
				if err != nil {
					log.Error("delete pod err %v", err)
					return
				}
				this.freePods = append(this.freePods[:index], this.freePods[index+1:]...)
			}
			this.freePods[index].freeTime += CheckInterval
		}
	}
}

func (this *ServerlessController) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.TriggerTopic, msg_type.JobTopic}
	this.subscriber.Subscribe(wg, ctx, topics, this)
	<-this.ready
	<-this.done
	cancel()
	wg.Wait()
}
