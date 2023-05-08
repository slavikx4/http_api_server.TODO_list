package app

import (
	"fmt"
	"github.com/slavikx4/http-api-server/internal/config"
	"github.com/slavikx4/http-api-server/internal/transport"
	"log"
	"net/http"
)

func Run() {
	handler := transport.NewHandler()

	http.HandleFunc("/", handler.HandleIndex)
	http.HandleFunc("/login", handler.HandleLogin)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.Config.ServerPort), nil); err != nil {
		log.Fatalln(err)
	}
}
