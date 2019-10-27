package nsqtrans

import (
	"GomokuGame/model"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

type ResultHandler struct {
	MatchResult   *model.MatchResultNsq
	TransFinished chan bool
}

func NewResultHandler() *ResultHandler {
	return &ResultHandler{
		MatchResult:   nil,
		TransFinished: make(chan bool),
	}
}

func (r *ResultHandler) HandleMessage(msg *nsq.Message) error {
	logrus.Info("result handlemsg ing...")
	matchResult := &model.MatchResultNsq{}
	err := json.Unmarshal(msg.Body, matchResult)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info("Unmarshal game result failed.")
		return nil
	}

	r.MatchResult = matchResult
	logrus.Info("result get data.")
	r.TransFinished <- true
	return nil
}
