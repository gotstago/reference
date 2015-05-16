//hello_handler.go

package handlers

/*
import (
	"net/http"
)

type GameHandler struct{}

func (e GameHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	sayParam := r.FormValue("say")

	if sayParam == "Nothing" {
		rw.WriteHeader(404)
	} else {
		rw.Write([]byte(sayParam))
		//rw.Write([]byte("hello!\n"))
	}*/
/*
Game Handler will start a game
Game will have state, which is a function that returns a function


}

type StateFn func(*Game) StateFn


type Game struct {
  Name   string
  Input  string
  Events chan Event
  State  StateFn

  Start int
  Pos   int
  Width int
}*/
