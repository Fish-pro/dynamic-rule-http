package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 用户认证信息
var (
	username = "admin"
	password = "password"
)

// BasicAuth 中间件函数，用于验证请求的基本身份认证信息
func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(parts[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !validateCredentials(pair[0], pair[1]) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 认证通过，继续处理请求
		next.ServeHTTP(w, r)
	})
}

// 验证用户名和密码
func validateCredentials(user, pass string) bool {
	return user == username && pass == password
}

// 处理 "/hello" 路径的请求
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've been authenticated!")
}

// 处理 "/message" 路径的请求，获取查询参数并返回消息
func messageHandler(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 中的查询参数 "msg"
	val := r.URL.Query().Get("key")

	// 返回消息
	fmt.Fprintf(w, "结构化后的返回信息: %s", val)
}

func main() {
	mux := http.NewServeMux()

	// 认证后返回 hello 消息
	mux.Handle("/", BasicAuth(http.HandlerFunc(helloHandler)))

	// 认证后返回 message 消息
	mux.Handle("/message", BasicAuth(http.HandlerFunc(messageHandler)))

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
