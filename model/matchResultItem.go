package model

import "github.com/jinzhu/gorm"

type MatchResultItem struct {
	gorm.Model
	GameID             string `gorm:"primary_key"`
	BoardLength        int
	BoardHeight        int
	Player1ID          string
	Player2ID          string
	Player1FirstHand   bool
	MaxThinkingTime    int
	Winner             int
	StartTime          int64	`gorm:"not null"`
	EndTime            int64	`gorm:"not null"`
	Operations         string	`gorm:"size:65533"` // json of []*Operation, longtext
	FoulPlayer         int    // 0: no foul, 1: player1 foul, 2: player2 foul
	ServerError        bool   `gorm:"not null"` // Server failure, game is invalid
	Player1LogFilePath string
	Player2LogFilePath string
	RefereeLogFilePaht string
}
