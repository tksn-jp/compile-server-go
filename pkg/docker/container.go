package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/moby/moby/client"
)

func BuildContainer(imageName, containerName string) (string, error) {
	cl, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}
	defer cl.Close()
	config := &container.Config{
		Image: imageName,
	}
	hostConfig := &container.HostConfig{}
	netConfig := &network.NetworkingConfig{}
	resp, err := cl.ContainerCreate(context.TODO(), config, hostConfig, netConfig, containerName)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
