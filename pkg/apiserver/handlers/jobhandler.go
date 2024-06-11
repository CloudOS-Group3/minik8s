package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"minik8s/pkg/api"
	msg "minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/gpu"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func GetJobs(context *gin.Context) {
	log.Info("received get jobs request")

	URL := config.EtcdJobPath
	jobs := etcdClient.PrefixGet(URL)

	log.Debug("get all jobs are: %+v", jobs)
	jsonString := stringutil.EtcdResEntryToJSON(jobs)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}

func AddJob(context *gin.Context) {
	log.Info("received add job request")

	var newJob api.Job
	if err := context.ShouldBind(&newJob); err != nil {
		log.Error("decode job failed")
		context.JSON(http.StatusOK, gin.H{
			"status": "wrong",
		})
		return
	}

	jobByteArray, err := json.Marshal(newJob)

	if err != nil {
		log.Error("error marshal new job")
		return
	}

	URL := config.EtcdJobPath + newJob.JobID
	etcdClient.PutEtcdPair(URL, string(jobByteArray))

	context.JSON(http.StatusOK, gin.H{
		"statas": "ok",
	})

	pod := newJob.Instance
	URL = config.EtcdPodPath + pod.Metadata.NameSpace + "/" + pod.Metadata.Name
	str := etcdClient.GetEtcdPair(URL)
	if str == "" {
		pod.Metadata.UUID = uuid.NewString()
		URL = config.GetUrlPrefix() + config.PodsURL
		if pod.Metadata.NameSpace == "" {
			URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		} else {
			URL = strings.Replace(URL, config.NamespacePlaceholder, pod.Metadata.NameSpace, -1)
		}
		byteArr, _ := json.Marshal(pod)
		httputil.Post(URL, byteArr)
	}
	var jobMsg msg.JobMsg
	jobMsg.Opt = msg.Add
	jobMsg.NewJob = newJob
	byteArr, _ := json.Marshal(jobMsg)
	publisher.Publish(msg.JobTopic, string(byteArr))
}

func GetJob(context *gin.Context) {
	log.Info("received get job request")
	name := context.Param(config.NameParam)

	if name == "" {
		log.Error("job id empty")
		return
	}

	URL := config.EtcdJobPath + name
	jobJson := etcdClient.GetEtcdPair(URL)

	var job api.Job
	json.Unmarshal([]byte(jobJson), &job)

	log.Info("job info: %+v", job)

	byteArr, err := json.Marshal(job)

	if err != nil {
		log.Error("error json marshal job: %s", err.Error())
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data": string(byteArr),
	})
}

func DeleteJob(context *gin.Context) {
	log.Info("received delete job request")
	name := context.Param(config.NameParam)

	if name == "" {
		log.Error("job id empty")
		return
	}

	URL := config.EtcdJobPath + name
	etcdClient.DeleteEtcdPair(URL)
}

func UpdateJob(context *gin.Context) {
	log.Info("received update job request")

	var newJob api.Job
	if err := context.ShouldBind(&newJob); err != nil {
		log.Error("decode job failed")
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	log.Info("job info: %+v", newJob)
	jobByteArray, err := json.Marshal(newJob)

	if err != nil {
		log.Error("error marshal newJob to json string")
	}

	URL := config.EtcdJobPath + newJob.JobID
	oldJob := etcdClient.GetEtcdPair(URL)
	etcdClient.PutEtcdPair(URL, string(jobByteArray))

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	var message msg.JobMsg
	if oldJob == "" {
		message = msg.JobMsg{
			Opt:    msg.Add,
			NewJob: newJob,
		}
	} else {
		var job api.Job
		if err = json.Unmarshal([]byte(oldJob), &job); err != nil {
			log.Error("error unmarshal old job")
		}
		message = msg.JobMsg{
			Opt:    msg.Update,
			OldJob: job,
			NewJob: newJob,
		}
	}
	msg_json, _ := json.Marshal(message)
	publisher.Publish(msg.JobTopic, string(msg_json))
}

func JobResultHandler(context *gin.Context) {
	log.Info("received job result request")
	var newResult api.JobResult
	bytes, _ := ioutil.ReadAll(context.Request.Body)
	log.Info("received job result request : %s", string(bytes))
	_ = json.Unmarshal(bytes, &newResult)
	log.Info("job result info: %+v", newResult)
	URL := config.EtcdJobPath + newResult.UUID
	jobJson := etcdClient.GetEtcdPair(URL)
	var job api.Job
	_ = json.Unmarshal([]byte(jobJson), &job)
	log.Info("job info: %+v", job)
	oldJob := job
	job.Status = api.JOB_ENDED
	if newResult.Error == "" {
		job.Result = newResult.Result
		log.Info("job result info: %+v", job.Result)
	} else {
		job.Result = newResult.Error
		log.Info("job result info: %+v", job.Result)
	}
	jobByteArr, _ := json.Marshal(job)
	log.Info("job result: %v", job)
	etcdClient.PutEtcdPair(URL, string(jobByteArr))
	var message msg.JobMsg
	message = msg.JobMsg{
		Opt:    msg.Update,
		OldJob: oldJob,
		NewJob: job,
	}
	msg_json, _ := json.Marshal(message)
	publisher.Publish(msg.JobTopic, string(msg_json))
}

func GpuResultHandler(context *gin.Context) {
	log.Info("received gpu result request")
	var newResult api.JobResult
	bytes, _ := ioutil.ReadAll(context.Request.Body)
	_ = json.Unmarshal(bytes, &newResult)
	log.Info("job result info: %+v", newResult)

	// Here UUID is <job_name>-<uuid>
	URL := config.GPUjobPath + newResult.UUID
	jobJson := etcdClient.GetEtcdPair(URL)
	var job api.GPUJob
	_ = json.Unmarshal([]byte(jobJson), &job)
	log.Info("job info: %+v", job)

	job.Status = api.JOB_ENDED
	if newResult.Error == "" {
		job.Result = newResult.Result
		log.Info("job result info: %+v", job.Result)
	} else {
		job.Result = job.SourcePath
	}
	job.EndTime = time.Now().Format("2006-01-02 15:04:05")

	jobByteArr, _ := json.Marshal(job)
	log.Info("job result: %v", job)
	etcdClient.PutEtcdPair(URL, string(jobByteArr))

	// get pod, pod name is <job_name>-<uuid>
	pod, success := etcdClient.GetPod(gpu.GPUNamespace, newResult.UUID)
	if pod.Spec.NodeName == "node1" {
		_ = exec.Command("scp", "-r", "root@192.168.3.12:"+job.SourcePath, job.SourcePath).Run()
	} else if pod.Spec.NodeName == "node2" {
		_ = exec.Command("scp", "-r", "root@192.168.3.11:"+job.SourcePath, job.SourcePath).Run()
	}
	if success {
		etcdClient.DeletePod(gpu.GPUNamespace, pod.Metadata.Name)
		message := msg.PodMsg{
			Opt:    msg.Delete,
			OldPod: pod,
		}
		msg_json, _ := json.Marshal(message)
		publisher.Publish(msg.PodTopic, string(msg_json))
	}

}
