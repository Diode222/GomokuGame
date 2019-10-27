package format

type GameResultFormat struct {
	GameID           string                 `json:"game_id,omitempty"`
	GameStatus       bool                   `json:"game_status,omitempty"` // not null, should be used as flag to webapp
	BoardLength      int                    `json:"borad_length,omitempty"`
	BoardHeight      int                    `json:"board_height,omitempty"`
	Player1ID        string                 `json:"player1_id,omitempty"`
	Player2ID        string                 `json:"player2_id,omitempty"`
	Player1FirstHand bool                   `json:"player1_first_hand,omitempty"`
	MaxThinkingTime  int                    `json:"max_thinking_time,omitempty"`
	Winner           int                    `json:"winner,omitempty"`
	StartTime        int64                  `json:"start_time,omitempty"`
	EndTime          int64                  `json:"end_time,omitempty"`
	Operations       []*GameOperationFormat `json:"game_operations,omitempty"`
	FoulPlayer       int                    `json:"foul_player,omitempty"` // 0: no foul, 1: player1 foul, 2: player2 foul
	ServerError      bool                   `json:"server_error"`          // Server failure, game is invalid
}

const (
	// foul player
	NO_FOUL      = 0
	PLAYER1_FOUL = 1
	PLAYER2_FOUL = 2

	// operation type
	BLANK = 0
	WHITE = 1
	NONE  = 2
)

type GameOperationFormat struct {
	Player    int `json:"player,omitempty"`
	PositionX int `json:"x,omitempty"`
	PositionY int `json:"y,omitempty"`
	Type      int `json:"piece_type,omitempty"`
}
