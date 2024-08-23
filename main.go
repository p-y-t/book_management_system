package main

import (
	"log"
	"net/http"
)

func main() {
	srv := http.NewServeMux()

	// routes
	srv.HandleFunc("GET /{$}", rootHandler)

	srv.HandleFunc("GET /upload", uploadHandlerGET)
	srv.HandleFunc("POST /upload", uploadHandlerPOST)
	srv.HandleFunc("GET /confirmation", confirmationHandlerGET)
	srv.HandleFunc("POST /confirmation", confirmationHandlerPOST)

	log.Println("[+] Start server listening on 8080")

	err := http.ListenAndServe(":8080", srv)
	if err != nil {
		log.Fatalln("[-] Fail to create server")
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

}
func uploadHandlerGET(w http.ResponseWriter, r *http.Request) {

}
func uploadHandlerPOST(w http.ResponseWriter, r *http.Request) {

}
func confirmationHandlerGET(w http.ResponseWriter, r *http.Request) {

}
func confirmationHandlerPOST(w http.ResponseWriter, r *http.Request) {

}
