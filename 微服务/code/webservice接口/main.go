package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"

	"golang.org/x/net/webservice"
)

// 1. 定义 SOAP 请求/响应结构（XML 标签必须和 WSDL 一致）
type AddRequest struct {
	XMLName xml.Name `xml:"http://soap.demo/ Add"`
	A       int      `xml:"a"`
	B       int      `xml:"b"`
}

type AddResponse struct {
	XMLName xml.Name `xml:"http://soap.demo/ AddResponse"`
	Result  int      `xml:"result"`
}

// 2. 真正的业务函数（可同时被 SOAP 和 JSON 复用）
func Add(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}

// 3. 把业务函数包装成 SOAP Handler
type mathSoap struct{}

func (mathSoap) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	r, err := Add(ctx, req.A, req.B)
	if err != nil {
		return nil, err
	}
	return &AddResponse{Result: r}, nil
}

// 4. 把业务函数包装成 JSON Handler（RESTful）
func addJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}
	var in struct{ A, B int }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	out, err := Add(r.Context(), in.A, in.B)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"result": out})
}

// 5. 路由组装
func main() {
	// 5.1 注册 SOAP 服务
	s := webservice.NewServer()
	s.Register(&mathSoap{})
	http.Handle("/ws", s)                     // SOAP 入口
	http.HandleFunc("/ws?wsdl", s.HandleWSDL) // 动态 WSDL

	// 5.2 注册 JSON 入口
	http.HandleFunc("/json", addJSON)

	log.Println("listen :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
