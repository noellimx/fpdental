package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"fpdental/core"
	"fpdental/server/endpoints"

	"github.com/go-chi/chi/v5/middleware"
)

type ServerOpts struct {
	Addr  string
	World *core.World
}

func RunServer(opts *ServerOpts) {

	log.Println("[RunServer]")
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/auth", endpoints.NewEndpointServiceAuthentication(opts.World).Routes())
	r.Mount("/appointments", endpoints.NewEndpointServiceUser(opts.World).Routes())
	r.Mount("/admin", endpoints.NewEndpointServiceAdmin(opts.World).Routes())

	log.Printf("%+#v", opts)
	http.ListenAndServe(opts.Addr, r)
}
