package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	url := "http://localhost:8080/add"
	param := addParam{
		X: 1,
		Y: 2,
	}

	paramBytes, _ := json.Marshal(param)
	res, _ := http.Post(url, "application/json", bytes.NewReader(paramBytes))
	var resData addResult
	json.NewDecoder(res.Body).Decode(&resData)
	fmt.Println(resData)
}
