//get_file_handler.go

package handlers

import (
	"fmt"
	"net/http"
	"path"
)

type GetFileHandler struct{}

func (handler GetFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Assuming you want to serve a java file at 'templates/FTPSourceProvider.java'
	fmt.Printf("%s \n", "serving file")
	fp := path.Join("handlers", "templates", "FTPSourceProvider.java")
	http.ServeFile(w, r, fp)
}
