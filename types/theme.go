package types

type Theme struct {
	Background        Color  `json:"background"`
	Primary           Color  `json:"primary"`
	Secondary         Color  `json:"secondary"`
	Success           Color  `json:"success"`
	Warning           Color  `json:"warning"`
	Error             Color  `json:"error"`
	WindowBarColor    string `json:"windowBarColor"`
	ShadowColor       string `json:"shadowColor"`
	WindowBorderColor string `json:"windowBorderColor"`
	WindowBorder      bool   `json:"windowBorder"`
}
type Color struct {
	Main             string `json:"main"`
	Transparent      string `json:"transparent"`
	TransparentLight string `json:"transparentLight"`
	TransparentDark  string `json:"transparentDark"`
	Dark             string `json:"dark"`
	Light            string `json:"light"`
	Text             string `json:"text"`
	DarkText         string `json:"darkText"`
}
