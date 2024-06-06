package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func main() {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "example")

	// Pull an image
	image, err := client.Pull(ctx, "docker.io/library/nginx:latest", containerd.WithPullUnpack)
	if err != nil {
		log.Fatal(err)
	}

	containerName := "example-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	snapshotName := "example-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// Create a container
	container, err := client.NewContainer(
		ctx,
		containerName,
		containerd.WithNewSnapshot(snapshotName, image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	// Create a task
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		log.Fatal(err)
	}
	defer task.Delete(ctx)

	// Start the task
	if err := task.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// Define a helper function to execute commands in the container
	execCommand := func(ctx context.Context, task containerd.Task, command []string) error {
		execID := command[0]
		execProcessSpec := specs.Process{
			Args: command,
			Cwd:  "/",
			Env:  []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		}

		execTask, err := task.Exec(ctx, execID, &execProcessSpec, cio.NewCreator(cio.WithStdio))
		if err != nil {
			return err
		}
		defer execTask.Delete(ctx)

		if err := execTask.Start(ctx); err != nil {
			return err
		}

		statusC, err := execTask.Wait(ctx)
		if err != nil {
			return err
		}

		status := <-statusC
		code, _, err := status.Result()
		if err != nil {
			return err
		}

		if code != 0 {
			return fmt.Errorf("command %s exited with status %d", command, code)
		}

		return nil
	}

	// Execute apt-get update
	if err := execCommand(ctx, task, []string{"apt-get", "update"}); err != nil {
		log.Fatal(err)
	}

	// Execute apt-get install -y curl
	if err := execCommand(ctx, task, []string{"apt-get", "install", "-y", "curl"}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully installed curl")
}
