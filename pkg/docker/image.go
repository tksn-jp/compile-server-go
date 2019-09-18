package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

func PullImage(imageName string) error {
	cl, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	defer cl.Close()
	opts := types.ImagePullOptions{}
	if _, err := cl.ImagePull(context.TODO(), imageName, opts); err != nil {
		return err
	}
	return nil
}
