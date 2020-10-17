package main

import (
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/shmuelhizmi/react-fullstack-go-server"
	"log"
	"net/http"
)

type DesktopApp struct {
}

func main() {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	react_fullstack_go_server.App(server, func(params *react_fullstack_go_server.ComponentParams) {
		desktop := params.View(0, "Desktop")
		apps := make([]DesktopApp, 0, 0)
		background := "test"
		desktop.Params["apps"] = apps
		desktop.Params["openApps"] = apps
		desktop.Params["background"] = background
		desktop.Start()
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	log.Panic(http.ListenAndServe(":5000", serveMux))
}
