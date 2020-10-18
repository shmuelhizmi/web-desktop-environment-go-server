package utils

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"time"
)

type Logger struct {
	Info  func(message string)
	Error func(message string)
	Warn  func(message string)
	Mount func(levelName string, levelColor color.Attribute) Logger
}

func CreateRootLogger() Logger {
	now := time.Now()
	return CreateLogger(func() {
		_, _ = color.New(color.FgBlack).Print("root")
	}, &now)
}

func CreateLogger(printLevel func(), lastMessageTime *time.Time) Logger {
	formatMessage := func(message string) string {

		if strings.Contains(message, "\n") {
			fullMessage := `
			 -------------------message-----------------
			`
			lines := strings.Split(message, "\n")
			for _, line := range lines {
				fullMessage += "|" + line + "\n\t\t\t"
			}
			fullMessage += ` -------------------------------------------`
			return fullMessage
		}
		return message
	}
	log := func(message string, messageColor color.Attribute) {
		newTime := time.Now()
		timeSinceLastMessage := newTime.Sub(*lastMessageTime).String()
		lastMessageTime = &newTime
		fmt.Print("[ ")
		printLevel()
		fmt.Print(" ]: ")
		_, _ = color.New(messageColor).Print(formatMessage(message))
		_, _ = color.New(color.FgBlue).Println(" " + timeSinceLastMessage + " ")
	}
	return Logger{
		Info: func(message string) {
			log(message, color.FgGreen)
		},
		Error: func(message string) {
			log(message, color.FgRed)
		},
		Warn: func(message string) {
			log(message, color.FgYellow)
		},
		Mount: func(levelName string, levelColor color.Attribute) Logger {
			return CreateLogger(func() {
				printLevel()
				fmt.Print(":")
				_, _ = color.New(levelColor).Print(levelName)
			}, lastMessageTime)
		},
	}
}
