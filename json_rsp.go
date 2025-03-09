package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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
	}

	sendOK(w, code, &respBody)
	return
}

func sendOK[T any](w http.ResponseWriter, code int, respVal *T){
	if reflect.TypeOf(*respVal).Kind() != reflect.Struct {
		sendError(w, 400, "Expected struct")
		return
	}

	respBody, err := json.Marshal(respVal)

	if err != nil {
		fmt.Println(err)
		sendError(w, 500, "Something went wrong")
	}

	w.WriteHeader(code)
	w.Write(respBody)
	return
}
