package main

import (
	"fmt"
	"iptv/api"
	"log"
	"net/http"
)

func main() {

	port := "8899"
	http.HandleFunc("/api/v1", api.Api_Handler)
	fmt.Println("Running epg service at port " + port + " ...")
	err := http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err.Error())
	}
}
