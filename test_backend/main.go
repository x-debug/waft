package main

import (
	"flag"
	"log"
	"net/http"
)

var port string

func main() {
	flag.StringVar(&port, "port", "5378", "port of server")
	flag.Parse()
	http.HandleFunc("/api/test", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Response From :" + port))
	})
	log.Println("Listen On Port " + port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
