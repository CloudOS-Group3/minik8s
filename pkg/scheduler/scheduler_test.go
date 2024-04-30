package scheduler

import (
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/kafka"
	"sync"
	"testing"
)

func TestScheduler(t *testing.T) {
	s := NewScheduler()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Run()
	}()
	<-s.ready
	addr := []string{"127.0.0.1:9092"}
	publisher := kafka.NewPublisher(addr)
	newPod := &api.Pod{
		Metadata: api.ObjectMeta{
			Name:      "test-pod",
			NameSpace: "default",
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				{
					Name:            "test-container1",
					Image:           "docker.io/library/nginx:latest",
					ImagePullPolicy: api.PullPolicyIfNotPresent,
				},
				{
					Name:            "test-container2",
					Image:           "docker.io/library/nginx:latest",
					ImagePullPolicy: api.PullPolicyIfNotPresent,
				},
			},
		},
	}
	podJson, err := json.Marshal(newPod)
	if err != nil {
		t.Errorf("marshal pod error: %s", err.Error())
	}
	err = publisher.Publish("pod", string(podJson))
	if err != nil {
		t.Errorf("kafka send message error: %s", err.Error())
	}
	for {
		if s.count >= 1 {
			close(s.done)
			break
		}
	}
	wg.Wait()
}
