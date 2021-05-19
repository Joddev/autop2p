package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func HandleResponse(resp *http.Response, err error) *http.Response {
	if err != nil {
		panic(err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		panic(errors.New(string(data)))
	}

	return resp
}

func DecodeResponse(resp *http.Response, data interface{}) {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		panic(err)
	}
}
