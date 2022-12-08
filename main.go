package main

import (
	"fpdental/core"
	"fpdental/server"
)

func main() {

	var wo = &core.WorldOpts{
		Paths: &core.Paths{Users: "./users.json", Appointments: "./appointments.json", Credentials: "./credentials.json"},
	}

	var w = core.Init(wo)
	var serverOpts = &server.ServerOpts{
		Addr:  ":8000",
		World: w,
	}
	server.RunServer(serverOpts)
}
