package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mgutz/ansi"
	"k8s.io/kubectl/pkg/util/term"
)

var tty = func() term.TTY {
	t := term.TTY{
		In:  os.Stdin,
		Out: os.Stdout,
	}

	if !t.IsTerminalIn() {
		return t
	}

	t.Raw = true

	return t
}()

type loadingText struct {
	message        string
	stream         io.Writer
	stopChan       chan bool
	startTimestamp int64
	cnt            int
}

func (l *loadingText) start() {
	l.startTimestamp = time.Now().UnixNano()

	if l.stopChan == nil {
		l.stopChan = make(chan bool)
	}

	go func() {
		l.render()
		for {
			select {
			case <-l.stopChan:
				return
			case <-time.After(time.Millisecond * 200):
				l.render()
			}
		}
	}()
}

func (l *loadingText) stop() {
	l.stopChan <- true
	l.stream.Write([]byte("\r"))

	for i := 0; i < len(l.message)+20; i++ {
		l.stream.Write([]byte(" "))
	}

	l.stream.Write([]byte("\r"))
}

func (l *loadingText) render() {
	l.stream.Write([]byte("\r"))

	messagePrefix := []byte("[wait] ")
	prefixLength := len(messagePrefix)
	l.stream.Write([]byte(ansi.Color(string(messagePrefix), "cyan+b")))

	timeElapsed := fmt.Sprintf("%v", (time.Now().UnixNano()-l.startTimestamp)/int64(time.Second))

	message := []byte(l.getLoadingChar() + " " + l.message)

	messageSuffix := " (" + timeElapsed + "s) "
	suffixLength := len(messageSuffix)

	terminalSize := tty.GetSize()

	// if the whole message is too long, then replace last words with ...
	if terminalSize != nil && terminalSize.Width < uint16(prefixLength+len(message)+suffixLength) {
		dots := []byte("...")
		maxMessageLength := int(terminalSize.Width) - (prefixLength + suffixLength + len(dots) + 5)

		if maxMessageLength > 0 {
			message = append(message[:maxMessageLength], dots...)
		}
	}

	message = append(message, messageSuffix...)
	l.stream.Write(message)
}

func (l *loadingText) getLoadingChar() string {
	var loadingChar string

	switch l.cnt {
	case 0:
		loadingChar = "⠋"
	case 1:
		loadingChar = "⠙"
	case 2:
		loadingChar = "⠹"
	case 3:
		loadingChar = "⠸"
	case 4:
		loadingChar = "⠼"
	case 5:
		loadingChar = "⠴"
	case 6:
		loadingChar = "⠦"
	case 7:
		loadingChar = "⠧"
	case 8:
		loadingChar = "⠇"
	case 9:
		loadingChar = "⠏"
	}

	l.cnt += 1

	if l.cnt > 9 {
		l.cnt = 0
	}

	return loadingChar
}
