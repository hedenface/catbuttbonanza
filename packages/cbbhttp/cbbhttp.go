package cbbhttp

import (
	"bytes"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	APICallTimeout = 10
)

func Error(w http.ResponseWriter, code int) {
	http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", http.StatusText(code)), code)
}

func Message(w http.ResponseWriter, message string) {
	fmt.Fprintf(w, "{\"message\":\"%s\"}", message)
}

func ReturnObject(w http.ResponseWriter, object interface{}) {
	json, err := json.Marshal(object)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		// TODO: add some better logging here
		fmt.Printf("ReturnObject: marshalling json failed: %v\n", err)
	}
	w.Write([]byte(json))
}

func GetBody(w http.ResponseWriter, r *http.Request, object any) error {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		// TODO: logging
		fmt.Println("GetBody: SetHandler Unmarshal error %v\n", err)
		return err
	}

	return nil
}

func APICall(host string, port int, method string, endpoint string, sendObject any, rcvObject any) error {

	postBody, err := json.Marshal(sendObject)
	if err != nil {
		return err
	}

	responseBody := bytes.NewBuffer(postBody)

	client := &http.Client{
		Timeout: time.Second * time.Duration(APICallTimeout),
	}

	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%d/%s", host, port, endpoint), responseBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, rcvObject)
	if err != nil {
		return err
	}

	return nil
}
