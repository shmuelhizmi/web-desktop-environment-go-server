package types

import react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"

type WindowState struct {
	Width     string   `json:"width"`
	MaxWidth  string   `json:"maxWidth"`
	MinWidth  string   `json:"minWidth"`
	Height    string   `json:"height"`
	MaxHeight string   `json:"maxHeight"`
	MinHeight string   `json:"minHeight"`
	Maximized bool     `json:"maximized"`
	Position  Position `json:"position"`
}

type CreateWindowParameters struct {
	View           react_fullstack_go_server.CreateView
	Name           string
	Title          string
	Icon           Icon
	State          WindowState
	DesktopManager DesktopManager
}

type CreateWindowReturn struct {
	UpdateTitle func(newTitle string)
	Cancel      func()
}
