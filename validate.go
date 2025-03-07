package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)



type respVal struct {
	Valid bool `json:"valid,omitempty"`
	Error string `json:"error,omitempty"`
}

func handlerValidate(w http.ResponseWriter, r *http.Request){
	type reqVal struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)

	var req reqVal
	var resp respVal
	if err := decoder.Decode(&req); err != nil {
		sendError(w, 500, "Something went wrong", &resp)
		return
	}

	if len(req.Body) > 140 {
		sendError(w, 400, "Chirp is too long", &resp)
		return
	}

	sendOK(w, &resp)
	return
}


func sendError(w http.ResponseWriter, code int, msg string, respVal *respVal){
	respVal.Error = msg
	respBody, err := json.Marshal(respVal)

	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(code)
	w.Write(respBody)
	return
}

func sendOK(w http.ResponseWriter, respVal *respVal){
	respVal.Valid = true
	respBody, err := json.Marshal(respVal)

	if err != nil {
		fmt.Println(err)
		sendError(w, 500, "Something went wrong", respVal)
	}

	w.WriteHeader(200)
	w.Write(respBody)
	return
}