package main

import (
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/desktop"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/managers"
	"log"
	"net/http"
)

type DesktopApp struct {
}

func main() {

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	desktopManager := managers.CreateDesktopManager()

	desktopManager.SettingsManager.Initialize()

	react_fullstack_go_server.App(server, desktop.CreateDesktop(desktopManager))

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	log.Panic(http.ListenAndServe(":5000", serveMux))
}
