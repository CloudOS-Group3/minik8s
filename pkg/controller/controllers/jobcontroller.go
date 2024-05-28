package controllers

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/kafka"
	"minik8s/util/httputil"
	"strings"
	"sync"
)

type JobController struct {
	WaitingJob []api.Job
	ready      chan bool
	done       chan bool
	subscriber *kafka.Subscriber
}

type FunctionParam struct {
	uuid   string
	params string
}

func NewJobController() *JobController {
	group := "job-controller"
	Controller := &JobController{
		ready:      make(chan bool),
		done:       make(chan bool),
		subscriber: kafka.NewSubscriber(group),
	}
	URL := config.GetUrlPrefix() + config.JobsURL
	var initialJob []api.Job
	_ = httputil.Get(URL, &initialJob, "data")
	for _, job := range initialJob {
		if job.Status == api.JOB_CREATED {
			initialJob = append(initialJob, job)
		}
	}
	Controller.WaitingJob = initialJob
	return Controller
}

func (s *JobController) Setup(_ sarama.ConsumerGroupSession) error {
	close(s.ready)
	return nil
}

func (s *JobController) Cleanup(_ sarama.ConsumerGroupSession) error {
	s.ready = make(chan bool)
	return nil
}

func (s *JobController) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg.Topic == msg_type.JobTopic {
			sess.MarkMessage(msg, "")
			s.JobHandler(msg.Value)
		}
		if msg.Topic == msg_type.PodTopic {
			sess.MarkMessage(msg, "")
			s.PodHandler(msg.Value)
		}
	}
	return nil
}

func (s *JobController) JobHandler(msg []byte) {
	var message msg_type.JobMsg
	_ = json.Unmarshal(msg, &message)
	if message.Opt == msg_type.Add {
		if message.NewJob.Instance.Status.PodIP != "" {
			s.CallFunction(message.NewJob)
			job := message.NewJob
			job.Status = api.JOB_RUNNING
			URL := config.GetUrlPrefix() + config.JobURL
			URL = strings.Replace(URL, config.NamePlaceholder, job.JobID, -1)
			byteArr, _ := json.Marshal(job)
			httputil.Put(URL, byteArr)
		} else {
			s.WaitingJob = append(s.WaitingJob, message.NewJob)
		}
	}
	if message.Opt == msg_type.Delete {
		for index, job := range s.WaitingJob {
			if job.JobID == message.OldJob.JobID {
				s.WaitingJob = append(s.WaitingJob[:index], s.WaitingJob[index+1:]...)
			}
		}
	}
}

func (s *JobController) PodHandler(msg []byte) {
	var message msg_type.PodMsg
	_ = json.Unmarshal(msg, &message)
	for index, job := range s.WaitingJob {
		if job.Instance.Metadata.NameSpace == message.NewPod.Metadata.NameSpace && job.Instance.Metadata.Name == message.NewPod.Metadata.Name {
			if message.NewPod.Status.PodIP != "" {
				s.CallFunction(job)
				job.Instance = message.NewPod
				job.Status = api.JOB_RUNNING
				URL := config.GetUrlPrefix() + config.JobURL
				URL = strings.Replace(URL, config.NamePlaceholder, job.JobID, -1)
				byteArr, _ := json.Marshal(job)
				httputil.Put(URL, byteArr)
				s.WaitingJob = append(s.WaitingJob[:index], s.WaitingJob[index+1:]...)
			}
		}
	}
}

func (s *JobController) CallFunction(job api.Job) {
	URL := "http://" + job.Instance.Status.PodIP + ":8080/run"
	param := FunctionParam{
		uuid:   job.JobID,
		params: job.Params,
	}
	bytes, _ := json.Marshal(param)
	httputil.Post(URL, bytes)
}

func (s *JobController) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.JobTopic, msg_type.PodTopic}
	s.subscriber.Subscribe(wg, ctx, topics, s)
	<-s.ready
	<-s.done
	cancel()
	wg.Wait()
}
