package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"fpdental/server/endpoints"
)

type ServerOpts struct {
	Addr string
}

func RunServer(opts *ServerOpts) {

	log.Println("[RunServer]")
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })

	r.Mount("/auth", endpoints.NewAuth().Routes())

	log.Printf("%+#v", opts)
	http.ListenAndServe(opts.Addr, r)
}
