package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

func PullImage(cl *client.Client, ctx context.Context, imageName string) error {
	opts := types.ImagePullOptions{}
	if _, err := cl.ImagePull(ctx, imageName, opts); err != nil {
		return err
	} else {
		log.Printf("image pull done: %s", imageName)
		return nil
	}
}
