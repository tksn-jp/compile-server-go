package docker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/moby/moby/client"
)

func RemoveContainer(containerId string) {
	cl, _ := client.NewEnvClient()
	ctx := context.TODO()
	_ = cl.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
}

func StartContainer(containerId string) {
	cl, _ := client.NewEnvClient()
	ctx := context.TODO()
	_ = cl.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
}

func BuildContainer(imageName, containerName string) (string, error) {
	cl, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}
	ctx := context.Background()

	//opts := types.ImagePullOptions{}
	//_, err = cl.ImagePull(ctx, "docker.io/library/" + imageName, opts)
	//if err != nil {
	//	return "", err
	//} else {
	//	log.Printf("image pull done: %s", imageName)
	//}

	//if err := PullImage(cl, ctx, "docker.io/library/" + imageName); err != nil {
	//	return "", err
	//}

	filtMap := map[string][]string{"name": {containerName}}
	filtBytes, _ := json.Marshal(filtMap)
	filt, _ := filters.FromParam(string(filtBytes))

	listOpts := types.ContainerListOptions{
		All:     true,
		Quiet:   false,
		Filters: filt,
	}
	resp, err := cl.ContainerList(ctx, listOpts)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, val := range resp {
		if imageName == val.Image {
			return val.ID, nil
		}
	}

	config := &container.Config{
		Image: imageName,
	}
	hostConfig := &container.HostConfig{}
	netConfig := &network.NetworkingConfig{}
	createResp, err := cl.ContainerCreate(ctx, config, hostConfig, netConfig, containerName)
	if err != nil {
		return "", err
	}
	return createResp.ID, nil
}
