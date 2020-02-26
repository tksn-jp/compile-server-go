package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func buildImage(ctx context.Context, cli *client.Client, file *os.File, name string) (types.ImageBuildResponse, error) {
	res, err := cli.ImageBuild(ctx, file, types.ImageBuildOptions{
		Tags:        []string{name},
		ForceRemove: true,
	})
	return res, err
}

func build(dfPath string) bool {
	dockerfile := filepath.Join(dfPath, "Dockerfile.tar.gz")
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()
	log.Printf("reading: %s\n", dockerfile)
	file, err := os.Open(dockerfile)
	check(err)
	defer file.Close()

	cwd, _ := os.Getwd()
	lang, _ := filepath.Rel(filepath.Join(cwd, "dockerFiles"), filepath.Dir(dockerfile))
	var sb strings.Builder
	sb.Reset()
	sb.Grow(32)
	sb.WriteString("rcs_")
	sb.WriteString(lang)

	// Dockerfileからイメージ作成
	res, err := buildImage(ctx, cli, file, sb.String())
	check(err)
	defer res.Body.Close()

	// build log
	_, err = ioutil.ReadAll(res.Body)
	check(err)
	log.Printf("Build image \"%s\"\n", sb.String())
	return true
}

func PrepareImage() (sum, success int) {
	cwd, _ := os.Getwd()
	dfPath, _ := filepath.Abs(filepath.Join(cwd, "dockerFiles"))
	files, _ := ioutil.ReadDir(dfPath)
	sum = len(files)
	success = 0
	ch := make(chan bool, sum)
	for _, f := range files {
		fp := filepath.Join(dfPath, f.Name())
		go func() {
			ch <- build(fp)
		}()
	}
	for i := 0; i < sum; i++ {
		if <-ch {
			success++
		}
	}
	return
}
