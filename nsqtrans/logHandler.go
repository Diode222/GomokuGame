package nsqtrans

import (
	"github.com/nsqio/go-nsq"
	"os"
)

type logHandler struct {
	LogFile       *os.File
	TransFinished chan bool
}

func newLogHandler(logFile *os.File) *logHandler {
	return &logHandler{
		LogFile:       logFile,
		TransFinished: make(chan bool),
	}
}

func (h *logHandler) HandleMessage(msg *nsq.Message) error {
	if string(msg.Body) == "EOFEOF" {
		h.TransFinished <- true
	} else {
		h.LogFile.Write(append(msg.Body, '\n'))
	}
	return nil
}
