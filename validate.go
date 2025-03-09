package main

import (
	"encoding/json"
	"net/http"
	"strings"
)



type respVal struct {
	Valid string `json:"cleaned_body"`
}
type reqVal struct {
	Body string `json:"body"`
}

func handlerValidate(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)

	var req reqVal
	var resp respVal
	if err := decoder.Decode(&req); err != nil {
		sendError(w, 500, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		sendError(w, 400, "Chirp is too long")
		return
	}

	resp.Valid = cleanInput(&req)

	sendOK(w, 200, &resp)
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