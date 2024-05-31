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
	"minik8s/pkg/serverless/function/function_util"
	workflow_util "minik8s/pkg/serverless/workflow"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"strings"
	"sync"
	"time"
)

type WorkflowStatus struct {
	workflow   api.Workflow
	currNode   api.Graph
	waitJob    string
	retRes     []api.Template
	resultUUID string
}

type WorkflowController struct {
	publisher  kafka.Publisher
	jobList    map[string]*WorkflowStatus // // map job uuid to workflow status
	subscriber *kafka.Subscriber
	ready      chan bool
	done       chan bool
}

func NewWorkflowController() *WorkflowController {
	group := "workflow-controller"
	KafkaURL := config.Remotehost + ":9092"
	Controller := &WorkflowController{
		ready:      make(chan bool),
		done:       make(chan bool),
		jobList:    make(map[string]*WorkflowStatus),
		subscriber: kafka.NewSubscriber(group),
		publisher:  *kafka.NewPublisher([]string{KafkaURL}),
	}
	return Controller
}

func (this *WorkflowController) Setup(_ sarama.ConsumerGroupSession) error {
	close(this.ready)
	return nil
}

func (this *WorkflowController) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (this *WorkflowController) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info("Message claimed: value %s", string(msg.Value))
		switch msg.Topic {
		case msg_type.TriggerWorkflowTopic:
			sess.MarkMessage(msg, "")
			this.triggerNewWorkflow(msg.Value)
		case msg_type.JobTopic:
			sess.MarkMessage(msg, "")
			var jobMsg msg_type.JobMsg
			err := json.Unmarshal(msg.Value, &jobMsg)
			if err != nil {
				log.Error("json unmarshal err %v", err)
				continue
			}
			if jobMsg.Opt == msg_type.Update && jobMsg.NewJob.Status == api.JOB_ENDED {
				this.execNextNode(jobMsg.NewJob)
			}
		default:
			log.Warn("Unknown msg type in serverless controller")
		}
	}
	return nil
}

func (this *WorkflowController) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.TriggerTopic, msg_type.JobTopic}
	this.subscriber.Subscribe(wg, ctx, topics, this)
	<-this.ready
	<-this.done
	cancel()
	wg.Wait()
}

func (this *WorkflowController) triggerNewWorkflow(value []byte) {
	var msg msg_type.WorkflowTriggerMsg
	_ = json.Unmarshal(value, &msg)
	log.Info("Trigger workflow: %v", msg)
	workflow := msg.Workflow

	// Get first function
	function, err := workflow_util.GetFunction(workflow.Graph.Function.NameSpace, workflow.Graph.Function.Name)
	if err != nil {
		log.Error("Can't find function: %s. %s", workflow.Graph.Function.Name, err.Error())
		return
	}

	trigger_uuid := uuid.NewString()
	// create trigger msg to exec function
	triggerMsg := &msg_type.TriggerMsg{
		Function: *function,
		UUID:     trigger_uuid,
		Params:   msg.Params,
	}
	workflowStatus := &WorkflowStatus{
		workflow:   workflow,
		currNode:   workflow.Graph,
		waitJob:    trigger_uuid,
		retRes:     function.Result,
		resultUUID: msg.UUID,
	}
	this.jobList[trigger_uuid] = workflowStatus
	jsonString, _ := json.Marshal(triggerMsg)
	this.publisher.Publish(msg_type.TriggerTopic, string(jsonString))
}

func (this *WorkflowController) execNextNode(job api.Job) {
	if this.jobList[job.JobID] == nil {
		log.Error("Can't find workflow uuid by job uuid: %s", job.JobID)
		return
	}
	workflowStatus := this.jobList[job.JobID]
	// remove job from jobList
	delete(this.jobList, job.JobID)

	// deal with ret result
	var result []interface{}
	if err := json.Unmarshal([]byte(job.Result), &result); err != nil {
		log.Error("Error unmarshaling Result: %s", err.Error())
		return
	}
	result_str := function_util.ConvertToStringList(result)
	resWithName, _ := function_util.CheckParams(workflowStatus.retRes, result_str)
	// get next node
	succssor := workflow_util.CheckRule(workflowStatus.currNode.Rule, resWithName)
	if succssor == nil {
		// get result
		res := &api.WorkflowResult{}
		URL := config.GetUrlPrefix() + config.TriggerResultURL
		URL = strings.Replace(URL, config.UUIDPlaceholder, workflowStatus.waitJob, -1)
		err := httputil.Get(URL, res, "data")
		if err != nil {
			log.Error("Error get result: %s", err.Error())
			return
		}

		res.EndTime = time.Now().Format("2006-01-02 15:04:05")
		res.Result = result_str

		// store result
		byteArr, _ := json.Marshal(res)
		err = httputil.Put(URL, byteArr)
		if err != nil {
			log.Error("Error put result: %s", err.Error())
			return
		}

		// workflow end
		delete(this.jobList, workflowStatus.waitJob)
		return
	}
}
