package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)


func sendError(w http.ResponseWriter, code int, msg string){
	type errResp struct {
		Error string `json:"error,omitempty"`
	}
	var respVal errResp
	respVal.Error = msg
	respBody, err := json.Marshal(respVal)

	if err != nil {
		fmt.Println(err)
		return
	}

	sendOK(w, code, &respBody)
	return
}

func sendOK[T any](w http.ResponseWriter, code int, respVal *T){
	respBody, err := json.Marshal(respVal)

	if err != nil {
		fmt.Println(err)
		sendError(w, 500, "Something went wrong")
		return
	}

	w.WriteHeader(code)
	w.Write(respBody)
	return
}
