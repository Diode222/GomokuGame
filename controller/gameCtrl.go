package controller

import (
	"GomokuGame/dao/gameId"
	"GomokuGame/dao/gameResult"
	"GomokuGame/dao/user"
	"GomokuGame/format"
	"GomokuGame/kube"
	"GomokuGame/nsqtrans"
	"GomokuGame/utils/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"net/http"
	"strconv"
)

type GameCtrl struct {
	GameIdDao     gameId.GameIdDaoInterface
	UserDao       user.UserDaoInterface
	GameResultDao gameResult.GameResultDaoInterface
	PodClient     coreV1.PodInterface
}

func NewGameCtrl(gameIdDao gameId.GameIdDaoInterface, userDao user.UserDaoInterface, gameResultDao gameResult.GameResultDaoInterface, podClient coreV1.PodInterface) *GameCtrl {
	return &GameCtrl{
		GameIdDao:     gameIdDao,
		UserDao:       userDao,
		GameResultDao: gameResultDao,
		PodClient:     podClient,
	}
}

func (g *GameCtrl) Start(c *gin.Context) {
	var player1Name string
	var player2Name string
	//var player1WarehouseAddr string
	//var player2WarehouseAddr string
	var player1ImageAddr string
	var player2ImageAddr string

	userToken := c.GetHeader("Authorization")
	player1FirstHand := c.Query("player1_first_hand")
	maxThinkingTime := c.Query("max_thinking_time")
	enemyUserName := c.Query("enemy_user_name")

	// current user info
	curUserInfo, err := g.UserDao.GetUserInfoWithToken(c.Request.Context(), userToken)
	if err != nil || curUserInfo == nil {
		c.String(http.StatusInternalServerError, json.JsonResponse(http.StatusInternalServerError, "Start game failed, get user info with token failed."))
		return
	}

	// FIXME test
	//enemyUserInfo := &model.UserItem{}
	//if enemyUserName == "" {
	//	// get random registered user as enemy
	//	enemyUserInfo, err = g.UserDao.GetRandomEnemyUserInfo(c.Request.Context())
	//	if err != nil {
	//		logrus.WithFields(logrus.Fields{
	//			"err": err.Error(),
	//		}).Info("Get random enemy user info failed.")
	//	}
	//} else {
	//	// enemy user info
	//	enemyUserInfo, err = g.UserDao.GetUserInfoWithUserName(c.Request.Context(), enemyUserName)
	//	if err != nil {
	//		c.String(http.StatusInternalServerError, json.JsonResponse(http.StatusInternalServerError, "Start game failed, get user info with enemyeUserName failed."))
	//		return
	//	}
	//}

	player1Name = curUserInfo.UserName
	// FIXME player2name = enemyUserInfo.UserName, should be store in redis, in case to show status when enemy flash web page
	player2Name = enemyUserName
	// FIXME test
	//player1WarehouseAddr = curUserInfo.WarehouseAddr
	//player2WarehouseAddr = enemyUserInfo.WarehouseAddr
	player1ImageAddr = "registry.cn-hangzhou.aliyuncs.com/gomoku_game/gomoku_game_impl:test"
	player2ImageAddr = "registry.cn-hangzhou.aliyuncs.com/gomoku_game/gomoku_game_impl:test"

	gameID, err := g.GameIdDao.GetNextGameId(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, json.JsonResponse(http.StatusInternalServerError, "Global game id generating failed."))
		return
	}

	// set game status to "gaming" in redis
	g.GameResultDao.SetTempGamingStatusInRedis(c.Request.Context(), strconv.FormatInt(gameID, 10))

	// init nsq transmission consumers and topics
	nsqTransClient, err := nsqtrans.NewNsqTrans(gameID, g.PodClient)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gameID": gameID,
			"err":    err.Error(),
		}).Info("Init new nsq trans consumers and topics failed.")
		c.String(http.StatusInternalServerError, json.JsonResponse(http.StatusInternalServerError, "Init new nsq trans consumers and topics failed."))
		return
	}

	matchPod := kube.CreateMatchPodResourceFile(strconv.FormatInt(gameID, 10), player1FirstHand, maxThinkingTime, player1Name, player2Name, player1ImageAddr, player2ImageAddr)
	_, err = g.PodClient.Create(matchPod)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gameID":           gameID,
			"player1FirstHand": player1FirstHand,
			"maxThinkingTime":  maxThinkingTime,
			"player1Name":      player1Name,
			"player2Name":      player2Name,
			"player1ImageAddr": player1ImageAddr,
			"player2ImageAddr": player2ImageAddr,
			"err":              err.Error(),
		}).Info("Create match pod in k8s failed")
		c.String(http.StatusInternalServerError, json.JsonResponse(http.StatusInternalServerError, "Create match pod in k8s failed"))
		return
	}

	// pull data from match pod with nsq in a goroutine
	go nsqTransClient.PullGameDataAndStore()

	c.String(http.StatusOK, json.JsonResponse(http.StatusOK, "Match started"))
}

func (g *GameCtrl) GetResult(c *gin.Context) {
	gameID := c.Query("game_id")

	gameResultModel, err := g.GameResultDao.GetGameResult(c.Request.Context(), gameID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gameID": gameID,
			"err":    err.Error(),
		}).Info("Get game result failed.")
		c.String(http.StatusInternalServerError, json.JsonResponse(http.StatusInternalServerError, "Get game result failed."))
		return
	}
	gameResultFormat := format.GameResultFormatter(gameResultModel)
	c.String(http.StatusOK, json.JsonResponse(http.StatusOK, gameResultFormat))
}
