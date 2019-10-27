package nsqtrans

import (
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
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
	logrus.Info("log handle msg ing...")
	if string(msg.Body) == "EOFEOF" {
		logrus.Info("eof eof ")
		h.TransFinished <- true
	} else {
		h.LogFile.Write(append(msg.Body, '\n'))
	}
	return nil
}
