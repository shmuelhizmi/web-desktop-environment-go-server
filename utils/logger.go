package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"strings"
	"time"
)

func CreateRootLogger() types.Logger {
	now := time.Now()
	return CreateLogger(func() {
		_, _ = color.New(color.BgBlack, color.FgWhite).Print("root")
	}, &now)
}

func CreateLogger(printLevel func(), lastMessageTime *time.Time) types.Logger {
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
	return types.Logger{
		Info: func(message string) {
			log(message, color.FgGreen)
		},
		Error: func(message string) {
			log(message, color.FgRed)
		},
		Warn: func(message string) {
			log(message, color.FgYellow)
		},
		Mount: func(levelName string, levelColor color.Attribute) types.Logger {
			return CreateLogger(func() {
				printLevel()
				fmt.Print(":")
				_, _ = color.New(levelColor).Print(levelName)
			}, lastMessageTime)
		},
	}
}
