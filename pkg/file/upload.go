package file

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"
)

func Compile(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(string(reqDump))

	formFile, _, err := req.FormFile("file")
	fileType := req.FormValue("fileType")
	if err != nil {
		log.Print(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	// create save file
	savePrefix := "uploaded_file"
	currentPath, _ := os.Getwd()
	filename := fmt.Sprintf("uploaded_%d.gz", time.Now().UnixNano())
	zipPath := filepath.Join(currentPath, savePrefix, filename)
	saveFile, err := os.Create(zipPath)
	if err != nil {
		log.Print(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer saveFile.Close()

	// write
	_, err = io.Copy(saveFile, formFile)
	if err != nil {
		log.Print(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// unzip
	if Unzip(zipPath, filepath.Join(currentPath, "unzip_src", filename + "." + fileType)) != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// todo: run container & copy file & compile

}

