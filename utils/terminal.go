package utils

import (
	"encoding/base64"
	"github.com/creack/pty"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"os/exec"
)

type TerminalInstance struct {
	ListenToOutput func(onOutput func(newOutput string)) (stop func())
	Write          func(input string)
	Exit           func()
}

func CreatePtyInstance(process string, arguments string, location string, logger types.Logger) (TerminalInstance, error) {
	c := exec.Command(process, arguments)

	// Start the command with a pty.
	ptyProcess, err := pty.Start(c)
	if err != nil {
		return TerminalInstance{}, err
	}
	instance := TerminalInstance{
		ListenToOutput: func(onOutput func(newOutput string)) (stop func()) {
			keepReading := true
			go func() {
				buffer := make([]byte, 128)
				for keepReading {
					read, readError := ptyProcess.Read(buffer)
					if readError != nil {
						logger.Warn("fail to read from " + process)
						logger.Error(readError.Error())
						return
					}

					out := make([]byte, base64.StdEncoding.EncodedLen(read))
					base64.StdEncoding.Encode(out, buffer[0:read])
					result, _ := base64.StdEncoding.DecodeString(string(out))
					onOutput(string(result))
				}
			}()
			return func() {
				keepReading = false
			}
		},
		Write: func(input string) {
			_, writeError := ptyProcess.Write([]byte(input))
			if writeError != nil {
				logger.Warn("fail to write \"" + input + "\" to " + process)
				logger.Error(writeError.Error())
			}
		},
		Exit: func() {
			failToCloseTerminal := ptyProcess.Close()
			if failToCloseTerminal != nil {
				logger.Error("fail to close terminal instance")
				logger.Error(failToCloseTerminal.Error())
			}
		},
	}
	// TO DO: make better implementation
	instance.Write("cd " + location + "\n")
	instance.Write("cd " + location + "\n")
	return instance, nil
}
