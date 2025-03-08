package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)



type respVal struct {
	Valid string `json:"cleaned_body,omitempty"`
	Error string `json:"error,omitempty"`
}
type reqVal struct {
	Body string `json:"body"`
}

func handlerValidate(w http.ResponseWriter, r *http.Request){
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

	sendOK(w, &resp, &req)
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

func sendOK(w http.ResponseWriter, respVal *respVal, reqVal *reqVal){
	respVal.Valid = cleanInput(reqVal)
	respBody, err := json.Marshal(respVal)

	if err != nil {
		fmt.Println(err)
		sendError(w, 500, "Something went wrong", respVal)
	}

	w.WriteHeader(200)
	w.Write(respBody)
	return
}

func cleanInput(req *reqVal) string {
	profainWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	wordList := strings.Split(req.Body, " ")

	for i, word := range wordList{
		if _, ok := profainWords[strings.ToLower(word)]; ok {
			wordList[i] = "****"
		}
	}
	return strings.Join(wordList, " ")
}