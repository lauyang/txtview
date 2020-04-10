package main

import (
	"log"
	"net/http"
	"txtview/handlers"
)

func main() {
	// 静态资源服务
	http.Handle("/public/", http.FileServer(http.Dir("./")))

	// 路由
	http.HandleFunc("/new", handlers.NewTxtView)
	http.HandleFunc("/edit", handlers.EditTxtView)
	//http.HandleFunc("/restore", handlers.DelTxtView)
	http.HandleFunc("/delete", handlers.DelTxtView)
	http.HandleFunc("/", handlers.Index)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer error:", err)
	}
}
