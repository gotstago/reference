//hello_handler.go

package handlers

import (
	"net/http"
)

type HelloHandler struct{}

func (e HelloHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	sayParam := r.FormValue("say")

	if sayParam == "Nothing" {
		rw.WriteHeader(404)
	} else {
		rw.Write([]byte(sayParam))
		//rw.Write([]byte("hello!\n"))
	}
}
