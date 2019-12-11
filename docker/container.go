package docker

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/moby/moby/client"
)

func BuildContainer(cli *client.Client, ctx context.Context, fileType string) (container.ContainerCreateCreatedBody, error) {
	sb.Reset()
	sb.Grow(32)
	sb.WriteString("rcs_")
	sb.WriteString(fileType)
	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image: sb.String(),
	}, nil, nil, "")
	return res, err
}

func RemoveContainer(cli *client.Client, ctx context.Context, resp container.ContainerCreateCreatedBody) {
	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
		log.Println("can't remove container")
	}
}

func DeployPackage(cli *client.Client, ctx context.Context, resp container.ContainerCreateCreatedBody, tarPath string) error {
	f, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	err = cli.CopyToContainer(ctx, resp.ID, "/src", f, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}
	return nil
}

func ExecContainer(cli *client.Client, ctx context.Context, resp container.ContainerCreateCreatedBody) (io.ReadCloser, error) {
	log.Println("debug: ExecContainer()")
	err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}
	log.Println("debug: Waiting conainer...")
	_, err = cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		return nil, err
	}
	log.Println("debug: conainer returns")

	return cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
}
