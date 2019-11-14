package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Json struct {
	Id string `json:"id"`
	Pw string `json:"pw"`
}

func Login(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var data Json
	errJ := json.Unmarshal(body, &data)
	if errJ != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// data{Id, Pw} を使って認証

}
