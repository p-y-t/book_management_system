package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	srv := http.NewServeMux()

	// routes
	srv.HandleFunc("GET /{$}", rootHandler)

	srv.HandleFunc("GET /upload", uploadHandlerGET)
	srv.HandleFunc("POST /upload", uploadHandlerPOST)
	srv.HandleFunc("GET /confirmation/{file_id}", confirmationHandlerGET)
	srv.HandleFunc("POST /confirmation/{file_id}", confirmationHandlerPOST)

	log.Println("[+] Start server listening on 8080")

	err := http.ListenAndServe(":8080", srv)
	if err != nil {
		log.Fatalln("[-] Fail to create server")
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/homePage.html")
}

func uploadHandlerGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/uploadPage.html")
}

func uploadHandlerPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	// Retrieve the file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file type. Please upload a .xlsx or .xls file.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		http.Error(w, "Invalid file type. Please upload a .xlsx or .xls file.", http.StatusBadRequest)
		return
	}

	fileId := time.Now().UnixNano()
	fileName := fmt.Sprintf("%d", fileId)
	dst, err := os.Create("files/" + fileName)
	defer dst.Close()

	// Copy the uploaded file's content to the destination file
	io.Copy(dst, file)

	url := fmt.Sprintf("/confirmation/%d", fileId)
	http.Redirect(w, r, url, http.StatusFound)
}

type PageData struct {
	Tbody template.HTML
}

func confirmationHandlerGET(w http.ResponseWriter, r *http.Request) {
	fileID, _ := strconv.ParseInt(r.PathValue("file_id"), 10, 64)
	filePath := fmt.Sprintf("./files/%d", fileID)
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		http.Error(w, "Invalid file name.", http.StatusBadRequest)
		return
	}
	defer f.Close()

	// Get all the rows in the Sheet1.
	data := ""
	rows, _ := f.GetRows("Sheet1")
	for _, row := range rows {
		data += "<tr>"
		for _, colCell := range row {
			data += fmt.Sprintf("<td>%s</td>", colCell)
		}
		data += "</tr>"
	}

	tbody := PageData{Tbody: template.HTML(data)}
	tmpl, _ := template.ParseFiles("web/confirmationPage.html")

	_ = tmpl.Execute(w, tbody)
}

func confirmationHandlerPOST(w http.ResponseWriter, r *http.Request) {

}
