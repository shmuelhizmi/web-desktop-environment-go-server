package types

import react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"

type WindowState struct {
	Width     float64  `json:"width"`
	MaxWidth  float64  `json:"maxWidth"`
	MinWidth  float64  `json:"minWidth"`
	Height    float64  `json:"height"`
	MaxHeight float64  `json:"maxHeight"`
	MinHeight float64  `json:"minHeight"`
	Maximized bool     `json:"maximized"`
	Position  Position `json:"position"`
}

type CreateWindowParameters struct {
	Name           string
	Title          string
	Icon           Icon
	State          WindowState
	DesktopManager DesktopManager
	App            react_fullstack_go_server.Component
	AppId          int64
}

type CreateWindowReturn struct {
	UpdateTitle *func(newTitle string)
	Component   react_fullstack_go_server.Component
}
