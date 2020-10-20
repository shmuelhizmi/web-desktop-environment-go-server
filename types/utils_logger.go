package types

import "github.com/fatih/color"

type Logger struct {
	Info  func(message string)
	Error func(message string)
	Warn  func(message string)
	Mount func(levelName string, levelColor color.Attribute) Logger
}
