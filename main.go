package main

import "fpdental/server"

var serverOpts = &server.ServerOpts{
	Addr: ":8000",
}

func main() {
	server.RunServer(serverOpts)
}
