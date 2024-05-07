package container

import (
	"context"
	"github.com/containerd/containerd/namespaces"
	"minik8s/pkg/api"
	"testing"
	"time"
)

func TestContainerManager(t *testing.T) {
	container_ := CreateContainer(api.Container{
		Name:            "test-container",
		Image:           "docker.io/library/nginx:latest",
		ImagePullPolicy: api.PullPolicyIfNotPresent,
	}, "test")
	if container_ == nil {
		t.Fatalf("Failed to create container")
		return
	}

	// use channel to wait for the completion of the operation
	done := make(chan bool)

	// start container
	go func() {
		ctx := namespaces.WithNamespace(context.Background(), "test")
		if StartContainer(container_, ctx) {
			done <- true
		} else {
			done <- false
		}
	}()

	// wait for the start operation to complete or timeout
	select {
	case success := <-done:
		if !success {
			t.Fatalf("Failed to start container")
			return
		}
	case <-time.After(5 * time.Second): // 设置超时时间
		t.Fatalf("Timeout: failed to start container")
		return
	}

	// stop container
	go func() {
		ctx := namespaces.WithNamespace(context.Background(), "test")
		if StopContainer(container_, ctx) {
			done <- true
		} else {
			done <- false
		}
	}()

	// wait for the stop operation to complete or timeout
	select {
	case success := <-done:
		if !success {
			t.Fatalf("Failed to stop container")
			return
		}
	case <-time.After(40 * time.Second): // 设置超时时间
		t.Fatalf("Timeout: failed to stop container")
		return
	}

	// delete container
	go func() {
		ctx := namespaces.WithNamespace(context.Background(), "test")
		if RemoveContainer(container_, ctx) {
			done <- true
		} else {
			done <- false
		}
	}()

	// wait for the delete operation to complete or timeout
	select {
	case success := <-done:
		if !success {
			t.Fatalf("Failed to remove container")
			return
		}
	case <-time.After(5 * time.Second): // 设置超时时间
		t.Fatalf("Timeout: failed to remove container")
		return
	}
}
