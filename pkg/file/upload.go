package file

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Compile(res http.ResponseWriter, req *http.Request) {
	//res.Header().Set("","")
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = res.Write(nil)
		return
	}
	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Printf("error while dumping: %s", err)
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write(nil)
		return
	}
	log.Print(string(reqDump))

	//fileType := req.FormValue("fileType")
	content := req.FormValue("content")
	//formFile, _, err := req.FormFile("file")
	//var buf []byte
	//_, _ = req.Body.Read(buf)
	//formFile, err := zlib.NewReader(bytes.NewBuffer(buf))
	//if err != nil {
	//	log.Printf("error while getting file: %s", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	_, _ = res.Write(nil)
	//	return
	//}
	//defer formFile.Close()

	// create save file
	savePrefix := "uploaded_file"
	currentPath, _ := os.Getwd()
	filename := fmt.Sprintf("uploaded_%d.go", time.Now().UnixNano())
	zipPath := filepath.Join(currentPath, savePrefix, filename)
	saveFile, err := os.Create(zipPath)
	if err != nil {
		log.Printf("errror while create file: %s", err)
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write(nil)
		return
	}
	defer saveFile.Close()

	// write
	//_, err = io.Copy(saveFile, content)
	_, err = saveFile.Write([]byte(content))
	if err != nil {
		log.Printf("error while writing content: %s", err)
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write(nil)
		return
	}

	srcTarget := zipPath
	// unzip
	//srcTarget := filepath.Join(currentPath, "unzip_src", filename + "." + fileType)
	//srcTarget := filepath.Join(currentPath, "unzip_src", filename + ".go")
	//if err = Unzip(zipPath, srcTarget); err != nil {
	//	log.Printf("error while unzipping: %s", err)
	//	res.WriteHeader(http.StatusInternalServerError)
	//	_, _ = res.Write(nil)
	//	return
	//}

	out, err := exec.Command("go", "run", srcTarget).Output()
	msg := fmt.Sprintf("exec result: %s", out)
	log.Printf("%s\n", msg)
	_, _ = res.Write([]byte(msg))

}

