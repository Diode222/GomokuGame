package nsqtrans

import (
	"GomokuGame/app/conf"
	"GomokuGame/model"
	"GomokuGame/utils/path"
	"errors"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	ospath "path"
	"strconv"
	"sync"
	"time"
)

type NsqTrans struct {
	GameId          string
	ResultConsumer  *nsq.Consumer
	Player1Consumer *nsq.Consumer
	Player2Consumer *nsq.Consumer
	RefereeConsumer *nsq.Consumer
	Player1LogFile  *os.File
	Player2LogFile  *os.File
	RefereeLogFile  *os.File
}

func NewNsqTransChan(gameId int64) (*NsqTrans, error) {
	gameIdStr := strconv.Itoa(int(gameId))
	resultTopic := gameIdStr + "_game_result"
	player1Topic := gameIdStr + "_log_player1"
	player2Topic := gameIdStr + "_log_player2"
	refereeTopic := gameIdStr + "_log_referee"

	createTopic(resultTopic, player1Topic, player2Topic, refereeTopic)

	resultConsumer, err := nsq.NewConsumer(resultTopic, "consume", nsq.NewConfig())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic": resultTopic,
			"err":   err.Error(),
		}).Info("get result nsq consumer failed.")
		return nil, err
	}

	player1Consumer, err := nsq.NewConsumer(player1Topic, "consume", nsq.NewConfig())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic": player1Topic,
			"err":   err.Error(),
		}).Info("get player1 nsq consumer failed.")
		return nil, err
	}

	player2Consumer, err := nsq.NewConsumer(player2Topic, "consume", nsq.NewConfig())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic": player2Topic,
			"err":   err.Error(),
		}).Info("get player2 nsq consumer failed.")
		return nil, err
	}

	refereeConsumer, err := nsq.NewConsumer(refereeTopic, "consume", nsq.NewConfig())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic": refereeTopic,
			"err":   err.Error(),
		})
		return nil, err
	}

	gameLogRootPath := path.GetGameLogDirPath(gameIdStr)
	err = os.Mkdir(gameLogRootPath, 755)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": gameLogRootPath,
			"err":  err.Error(),
		}).Info("mkidr game log path failed.")
		return nil, err
	}

	player1LogPath := ospath.Join(gameLogRootPath, "player1.log")
	player2LogPath := ospath.Join(gameLogRootPath, "player2.log")
	refereeLogPath := ospath.Join(gameLogRootPath, "referee.log")

	player1LogFile, err := os.Create(player1LogPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": player1LogPath,
			"err":  err.Error(),
		}).Info("create player1 log file failed.")
		return nil, err
	}
	player2LogFile, err := os.Create(player2LogPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": player2LogPath,
			"err":  err.Error(),
		}).Info("create player2 log file failed.")
		return nil, err
	}
	refereeLogFile, err := os.Create(refereeLogPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": refereeLogPath,
			"err":  err.Error(),
		}).Info("create referee log file failed.")
		return nil, err
	}

	return &NsqTrans{
		GameId:          gameIdStr,
		ResultConsumer:  resultConsumer,
		Player1Consumer: player1Consumer,
		Player2Consumer: player2Consumer,
		RefereeConsumer: refereeConsumer,
		Player1LogFile:  player1LogFile,
		Player2LogFile:  player2LogFile,
		RefereeLogFile:  refereeLogFile,
	}, nil
}

func (n *NsqTrans) PullGameData() (*model.MatchResultModel, error) {
	var gameReulst *model.MatchResultModel
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		resultHandler := NewResultHandler()
		n.ResultConsumer.AddHandler(resultHandler)
		select {
		case <-resultHandler.TransFinished:
			gameReulst = resultHandler.MatchResult
			return
		case <-time.After(conf.MAX_PULL_DATA_TIME):
			return
		}
	}()
	go func() {
		defer wg.Done()
		player1LogHandler := newLogHandler(n.Player1LogFile)
		n.Player1Consumer.AddHandler(player1LogHandler)
		select {
		case <-player1LogHandler.TransFinished:
			return
		case <-time.After(conf.MAX_PULL_DATA_TIME): // game should finish in 5 minutes
			return
		}
	}()
	go func() {
		defer wg.Done()
		player2LogHandler := newLogHandler(n.Player2LogFile)
		n.Player2Consumer.AddHandler(player2LogHandler)
		select {
		case <-player2LogHandler.TransFinished:
			return
		case <-time.After(conf.MAX_PULL_DATA_TIME):
			return
		}
	}()
	go func() {
		defer wg.Done()
		refereeLogHandler := newLogHandler(n.RefereeLogFile)
		n.RefereeConsumer.AddHandler(refereeLogHandler)
		select {
		case <-refereeLogHandler.TransFinished:
			return
		case <-time.After(conf.MAX_PULL_DATA_TIME):
			return
		}
	}()
	wg.Wait()

	if gameReulst == nil {
		return nil, errors.New("pull match data failed")
	}

	return gameReulst, nil
}

func (n *NsqTrans) Stop() {
	n.ResultConsumer.Stop()
	n.Player1Consumer.Stop()
	n.Player2Consumer.Stop()
	n.RefereeConsumer.Stop()
	n.Player1LogFile.Close()
	n.Player2LogFile.Close()
	n.RefereeLogFile.Close()
}

func DeleteTopic(gameId int64) {
	gameIdStr := strconv.Itoa(int(gameId))
	resultTopic := gameIdStr + "_game_result"
	player1Topic := gameIdStr + "_log_player1"
	player2Topic := gameIdStr + "_log_player2"
	refereeTopic := gameIdStr + "_log_referee"

	deleteApiAddr := conf.NSQD_ADDR + "/topic/delete?topic=" + resultTopic
	_, err := http.PostForm(deleteApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         resultTopic,
			"deleteApiAddr": deleteApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + resultTopic + " failed.")
	}

	deleteApiAddr = conf.NSQD_ADDR + "/topic/delete?topic=" + player1Topic
	_, err = http.PostForm(deleteApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         player1Topic,
			"deleteApiAddr": deleteApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + player1Topic + " failed。")
	}

	deleteApiAddr = conf.NSQD_ADDR + "/topic/delete?topic=" + player2Topic
	_, err = http.PostForm(deleteApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         player2Topic,
			"deleteApiAddr": deleteApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + player2Topic + " failed.")
	}

	deleteApiAddr = conf.NSQD_ADDR + "/topic/delete?topic=" + refereeTopic
	_, err = http.PostForm(deleteApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         refereeTopic,
			"deleteApiAddr": deleteApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + refereeTopic + " failed.")
	}
}

func createTopic(resultTopic, player1Topic, player2Topic, refereeTopic string) {
	createApiAddr := conf.NSQD_ADDR + "/topic/create?topic=" + resultTopic
	_, err := http.PostForm(createApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         resultTopic,
			"createApiAddr": createApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + resultTopic + " failed.")
	}

	createApiAddr = conf.NSQD_ADDR + "/topic/create?topic=" + player1Topic
	_, err = http.PostForm(createApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         player1Topic,
			"createApiAddr": createApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + player1Topic + " failed.")
	}

	createApiAddr = conf.NSQD_ADDR + "/topic/create?topic=" + player2Topic
	_, err = http.PostForm(createApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         player2Topic,
			"createApiAddr": createApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + player2Topic + " failed.")
	}

	createApiAddr = conf.NSQD_ADDR + "/topic/create?topic=" + refereeTopic
	_, err = http.PostForm(createApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         refereeTopic,
			"createApiAddr": createApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + refereeTopic + " failed.")
	}
}
