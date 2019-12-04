package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/moby/moby/client"
	"github.com/tksn-jp/compile-server-go/docker"
)

const UploadFileDir = "uploaded_file"

func createRawCodeFile(fileType string) (*os.File, string, error) {
	entryPoint := map[string]string{
		"golang": "main.go",
		"c":      "main.c",
		"java":   "Main.java",
	}
	savePrefix, _ := filepath.Abs(UploadFileDir)
	uploadedDir := fmt.Sprintf("uploaded_%d", time.Now().UnixNano())
	saveDirPath := filepath.Join(savePrefix, uploadedDir)
	os.MkdirAll(saveDirPath, 0777)
	savePath := filepath.Join(saveDirPath, entryPoint[fileType])
	saveFile, err := os.Create(savePath)
	return saveFile, saveDirPath, err
}

func Server(res http.ResponseWriter, req *http.Request) {
	//res.Header().Set("","")
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = res.Write(nil)
		return
	}

	fileType := req.FormValue("fileType")
	dataType := req.FormValue("dataType")
	var out io.Reader
	if dataType == "rawCode" {
		content := req.FormValue("content")
		log.Print(content)
		a, err := CompileRawCode(fileType, content)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			_, _ = res.Write(nil)
			return
		}
		out = a
	} else if dataType == "archive" {
		file, _, err := req.FormFile("file")
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			_, _ = res.Write(nil)
			return
		}
		defer file.Close()
		a, err := CompileFile(fileType, file)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			_, _ = res.Write(nil)
			return
		}
		out = a
	} else {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write(streamToByte(out))
		return
	}

	msg := streamToString(out)
	log.Println(msg)
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte(msg))
	return
}

func CompileFile(fileType string, file multipart.File) (io.Reader, error) {
	tarDir, err := filepath.Abs("tar")
	if err != nil {
		return nil, err
	}
	tarFilePath := filepath.Join(tarDir, fmt.Sprintf("uploaded_%d.tar.gz", time.Now().UnixNano()))
	tarFile, err := os.Create(tarFilePath)
	if err != nil {
		return nil, err
	}
	defer tarFile.Close()
	if _, err := tarFile.Write(streamToByte(file)); err != nil {
		return nil, err
	}
	return ExecDocker(fileType, tarFilePath)
}

func CompileRawCode(fileType string, content string) (io.Reader, error) {
	// create save file
	saveFile, saveDir, err := createRawCodeFile(fileType)
	if err != nil {
		log.Printf("errror while create file: %s", err)
		return nil, err
	}
	defer saveFile.Close()

	// write
	_, err = saveFile.Write([]byte(content))
	if err != nil {
		log.Printf("error while writing content: %s", err)
		return nil, err
	}

	// tar source
	cwd, _ := os.Getwd()
	cwdAbs, _ := filepath.Abs(cwd)
	tarPath := filepath.Join(cwdAbs, "tar", filepath.Base(saveDir)+".tar")
	log.Println("tarPath: " + tarPath)
	if err := enTar(dirWalk(saveDir), tarPath); err != nil {
		log.Printf("error while archiving uploaded files: %s", err)
		return nil, err
	}

	////////////// docker ///////////////////////
	return ExecDocker(fileType, tarPath)
}

func ExecDocker(fileType string, packagePath string) (io.Reader, error) {
	log.Println("debug: ExecDocker()")
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Printf("error while creating client for docker: %s", err)
		return nil, err
	}

	// build container
	log.Println("debug: Building container")
	resp, err := docker.BuildContainer(cli, ctx, fileType)
	if err != nil {
		log.Printf("error while building container: %s", err)
		return nil, err
	}
	defer docker.RemoveContainer(cli, ctx, resp)

	// insert package into container
	log.Println("debug: Deploying package")
	err = docker.DeployPackage(cli, ctx, resp, packagePath)
	if err != nil {
		log.Printf("error while deploying package: %s", err)
		return nil, err
	}

	// exec
	return docker.ExecContainer(cli, ctx, resp)
}

func enTar(paths []string, tarPath string) error {
	w, err := os.OpenFile(tarPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil
	}
	defer w.Close()
	tw := tar.NewWriter(w)
	defer tw.Close()
	for _, fp := range paths {
		body, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Println("debug: readfile err")
			return err
		}
		if body != nil {
			hdr := &tar.Header{
				Name: filepath.Base(fp),
				Size: int64(len(body)),
				Mode: 0644,
			}
			if err := tw.WriteHeader(hdr); err != nil {
				log.Println("debug: writeheader err")
				return err
			}
			if _, err := tw.Write(body); err != nil {
				log.Println("debug: writebody err")
				return err
			}
		}
	}
	log.Println("debug: ok tar")
	return nil
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.Bytes()
}

func streamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}

func dirWalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirWalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
