package docker

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

var sb strings.Builder

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
	log.Printf("reading: %s\n", dockerfile)
	file, err := os.Open(dockerfile)
	check(err)
	defer file.Close()

	cwd, _ := os.Getwd()
	lang, _ := filepath.Rel(filepath.Join(cwd, "dockerFiles"), filepath.Dir(dockerfile))
	sb.Reset()
	sb.Grow(32)
	sb.WriteString("rcs_")
	sb.WriteString(lang)

	// Dockerfileからイメージ作成
	res, err := buildImage(ctx, cli, file, sb.String())
	check(err)
	defer res.Body.Close()

	// build log
	log.Printf("OSType: %s\n", res.OSType)
	b, err := ioutil.ReadAll(res.Body)
	check(err)
	log.Println(*(*string)(unsafe.Pointer(&b)))
	log.Printf("Build Image. Image's name is %v\n", sb.String())
	return true
}

func PrepareImage() int {
	cwd, _ := os.Getwd()
	dfPath, _ := filepath.Abs(filepath.Join(cwd, "dockerFiles"))
	files, _ := ioutil.ReadDir(dfPath)
	cnt := 0
	for _, f := range files {
		if build(filepath.Join(dfPath, f.Name())) {
			cnt++
		}
	}
	return cnt
}
