package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type addParam struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type addResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data int    `json:"data"`
}

func add(x int, y int) int {
	return x + y
}

func addHanlder(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var p addParam
	json.Unmarshal(b, &p)

	ret := add(p.X, p.Y)

	respBytes, _ := json.Marshal(addResult{
		Code: 0,
		Msg:  "success",
		Data: ret,
	})
	w.Write(respBytes)
}

func main() {
	http.HandleFunc("/add", addHanlder)
	http.ListenAndServe(":8080", nil)
}
