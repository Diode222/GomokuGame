package nsqtrans

import (
	"GomokuGame/app/conf"
	"GomokuGame/db"
	"GomokuGame/model"
	"GomokuGame/utils/path"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"net/http"
	"os"
	ospath "path"
	"strconv"
	"sync"
	"time"
)

type NsqTrans struct {
	GameId             string
	ResultConsumer     *nsq.Consumer
	Player1Consumer    *nsq.Consumer
	Player2Consumer    *nsq.Consumer
	RefereeConsumer    *nsq.Consumer
	Player1LogFilePath string
	Player2LogFilePath string
	RefereeLogFilePath string
	Player1LogFile     *os.File
	Player2LogFile     *os.File
	RefereeLogFile     *os.File
	DBInstance         *db.DB
	PodClient          corev1.PodInterface
}

func NewNsqTrans(gameId int64, podClient corev1.PodInterface) (*NsqTrans, error) {
	gameIdStr := strconv.Itoa(int(gameId))
	resultTopic := gameIdStr + "_game_result"
	player1Topic := gameIdStr + "_log_player1"
	player2Topic := gameIdStr + "_log_player2"
	refereeTopic := gameIdStr + "_log_referee"

	err := createTopic(resultTopic, player1Topic, player2Topic, refereeTopic)
	if err != nil {
		return nil, err
	}

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
		GameId:             gameIdStr,
		ResultConsumer:     resultConsumer,
		Player1Consumer:    player1Consumer,
		Player2Consumer:    player2Consumer,
		RefereeConsumer:    refereeConsumer,
		Player1LogFilePath: player1LogPath,
		Player2LogFilePath: player2LogPath,
		RefereeLogFilePath: refereeLogPath,
		Player1LogFile:     player1LogFile,
		Player2LogFile:     player2LogFile,
		RefereeLogFile:     refereeLogFile,
		DBInstance:         db.GetDB(),
		PodClient:          podClient,
	}, nil
}

func (n *NsqTrans) PullGameDataAndStore() {
	var gameReulst *model.MatchResultNsq
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		resultHandler := NewResultHandler()
		n.ResultConsumer.AddHandler(resultHandler)
		if err := n.ResultConsumer.ConnectToNSQLookupd(conf.NSQLOOKUPD_ADDR); err != nil {
			logrus.WithFields(logrus.Fields{
				"topic":           n.GameId + "_game_result",
				"NSQLOOKUPD_ADDR": conf.NSQLOOKUPD_ADDR,
				"err":             err.Error(),
			}).Info("Result nsq consumer connected to nsqLookupd failed.")
			return
		}
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
		if err := n.Player1Consumer.ConnectToNSQLookupd(conf.NSQLOOKUPD_ADDR); err != nil {
			logrus.WithFields(logrus.Fields{
				"topic":           n.GameId + "_log_player1",
				"NSQLOOKUPD_ADDR": conf.NSQLOOKUPD_ADDR,
				"err":             err.Error(),
			}).Info("Player1 log nsq consumer connected to nsqLookupd failed.")
			return
		}
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
		if err := n.Player2Consumer.ConnectToNSQLookupd(conf.NSQLOOKUPD_ADDR); err != nil {
			logrus.WithFields(logrus.Fields{
				"topic":           n.GameId + "_log_player2",
				"NSQLOOKUPD_ADDR": conf.NSQLOOKUPD_ADDR,
				"err":             err.Error(),
			}).Info("Player2 log nsq consumer connected to nsqLookupd failed.")
			return
		}
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
		if err := n.RefereeConsumer.ConnectToNSQLookupd(conf.NSQLOOKUPD_ADDR); err != nil {
			logrus.WithFields(logrus.Fields{
				"topic":           n.GameId + "_log_referee",
				"NSQLOOKUPD_ADDR": conf.NSQLOOKUPD_ADDR,
				"err":             err.Error(),
			}).Info("Referee log nsq consumer connected to nsqLookupd failed.")
			return
		}
		select {
		case <-refereeLogHandler.TransFinished:
			return
		case <-time.After(conf.MAX_PULL_DATA_TIME):
			return
		}
	}()
	wg.Wait()

	if gameReulst == nil {
		n.storeServerErrorInPullDataFromNsq()
		logrus.WithFields(logrus.Fields{
			"gameId": n.GameId,
		}).Info("Game result == nil when pull data from nsq.")
		n.Stop()
		return
	}

	operationsJson, err := json.Marshal(&gameReulst.Operations)
	if err != nil {
		n.storeServerErrorInPullDataFromNsq()
		logrus.WithFields(logrus.Fields{
			"gameId": n.GameId,
			"err":    err.Error(),
		}).Info("Marshal operations failed when pull data from nsq.")
		n.Stop()
		return
	}
	logrus.Println("length of operationsJson: " + strconv.Itoa(len(operationsJson)))
	gameReulstItem := &model.MatchResultItem{
		GameID:             gameReulst.GameID,
		GameStatus:         false,
		BoardLength:        gameReulst.BoardLength,
		BoardHeight:        gameReulst.BoardHeight,
		Player1ID:          gameReulst.Player1ID,
		Player2ID:          gameReulst.Player2ID,
		Player1FirstHand:   gameReulst.Player1FirstHand,
		MaxThinkingTime:    gameReulst.MaxThinkingTime,
		Winner:             gameReulst.Winner,
		StartTime:          gameReulst.StartTime,
		EndTime:            gameReulst.EndTime,
		Operations:         string(operationsJson),
		FoulPlayer:         gameReulst.FoulPlayer,
		ServerError:        gameReulst.ServerError,
		Player1LogFilePath: n.Player1LogFilePath,
		Player2LogFilePath: n.Player2LogFilePath,
		RefereeLogFilePath: n.RefereeLogFilePath,
	}

	gameResultItemJson, err := json.Marshal(gameReulstItem)
	if err != nil {
		n.storeServerErrorInPullDataFromNsq()
		logrus.WithFields(logrus.Fields{
			"gameId": n.GameId,
			"err":    err.Error(),
		}).Info("Game result mysql item masrshal failed.")
	}

	n.DBInstance.Lock.Lock()
	n.DBInstance.Redis.Set("game_result_"+n.GameId, gameResultItemJson, conf.GAME_RESULT_REDIS_STORE_TIME)
	n.DBInstance.Mysql.Table(conf.MATCH_RESULT_TABLE_NAME).Create(gameReulstItem)
	n.DBInstance.Lock.Unlock()

	n.Stop()
}

func (n *NsqTrans) storeServerErrorInPullDataFromNsq() {
	serverErrorGameResult := &model.MatchResultItem{
		GameID:      n.GameId,
		GameStatus:  false,
		StartTime:   time.Now().Unix(),
		EndTime:     time.Now().Unix(),
		ServerError: true,
	}
	serverErrorGameResultBinary, _ := json.Marshal(serverErrorGameResult)
	n.DBInstance.Lock.Lock()
	n.DBInstance.Redis.Set("game_result_"+n.GameId, serverErrorGameResultBinary, conf.GAME_RESULT_REDIS_STORE_TIME)
	n.DBInstance.Mysql.Table(conf.MATCH_RESULT_TABLE_NAME).Create(serverErrorGameResult)
	n.DBInstance.Lock.Unlock()
}

func (n *NsqTrans) Stop() {
	deleteTopic(n.GameId)
	n.ResultConsumer.Stop()
	n.Player1Consumer.Stop()
	n.Player2Consumer.Stop()
	n.RefereeConsumer.Stop()
	n.Player1LogFile.Close()
	n.Player2LogFile.Close()
	n.RefereeLogFile.Close()

	// close match pod
	deletePolicyPod := metav1.DeletePropagationForeground
	if err := n.PodClient.Delete("match-game-"+n.GameId, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicyPod,
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"gameId": n.GameId,
			"err":    err.Error(),
		}).Info("Delete match pod failed.")
	}
}

// TODO use waitgroup
func deleteTopic(gameId string) {
	resultTopic := gameId + "_game_result"
	player1Topic := gameId + "_log_player1"
	player2Topic := gameId + "_log_player2"
	refereeTopic := gameId + "_log_referee"

	deleteApiAddr := "http://" + conf.NSQD_PULL_ADDR + "/topic/delete?topic=" + resultTopic
	_, err := http.PostForm(deleteApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         resultTopic,
			"deleteApiAddr": deleteApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + resultTopic + " failed.")
	}

	deleteApiAddr = "http://" + conf.NSQD_PULL_ADDR + "/topic/delete?topic=" + player1Topic
	_, err = http.PostForm(deleteApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         player1Topic,
			"deleteApiAddr": deleteApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + player1Topic + " failed。")
	}

	deleteApiAddr = "http://" + conf.NSQD_PULL_ADDR + "/topic/delete?topic=" + player2Topic
	_, err = http.PostForm(deleteApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         player2Topic,
			"deleteApiAddr": deleteApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + player2Topic + " failed.")
	}

	deleteApiAddr = "http://" + conf.NSQD_PULL_ADDR + "/topic/delete?topic=" + refereeTopic
	_, err = http.PostForm(deleteApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         refereeTopic,
			"deleteApiAddr": deleteApiAddr,
			"err":           err.Error(),
		}).Info("delete topic " + refereeTopic + " failed.")
	}
}

// TODO Use errgroup
func createTopic(resultTopic, player1Topic, player2Topic, refereeTopic string) error {
	createApiAddr := "http://" + conf.NSQD_PULL_ADDR + "/topic/create?topic=" + resultTopic
	_, err := http.PostForm(createApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         resultTopic,
			"createApiAddr": createApiAddr,
			"err":           err.Error(),
		}).Info("create topic " + resultTopic + " failed.")
		return err
	}

	createApiAddr = "http://" + conf.NSQD_PULL_ADDR + "/topic/create?topic=" + player1Topic
	_, err = http.PostForm(createApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         player1Topic,
			"createApiAddr": createApiAddr,
			"err":           err.Error(),
		}).Info("create topic " + player1Topic + " failed.")
		return err
	}

	createApiAddr = "http://" + conf.NSQD_PULL_ADDR + "/topic/create?topic=" + player2Topic
	_, err = http.PostForm(createApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         player2Topic,
			"createApiAddr": createApiAddr,
			"err":           err.Error(),
		}).Info("create topic " + player2Topic + " failed.")
		return err
	}

	createApiAddr = "http://" + conf.NSQD_PULL_ADDR + "/topic/create?topic=" + refereeTopic
	_, err = http.PostForm(createApiAddr, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":         refereeTopic,
			"createApiAddr": createApiAddr,
			"err":           err.Error(),
		}).Info("create topic " + refereeTopic + " failed.")
		return err
	}

	return nil
}
