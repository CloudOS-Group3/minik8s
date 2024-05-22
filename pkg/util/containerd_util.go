package util

import "github.com/containerd/containerd"

func CreateClient() (*containerd.Client, error) {
	return containerd.New("/run/containerd/containerd.sock")
}
