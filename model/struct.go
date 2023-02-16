package model

// Config 是读入的配置文件的struct
type Config struct {
	Input struct {
		Aria2 struct {
			Aria2Server string `json:"aria2-server"`
			Aria2Key    string `json:"aria2-key"`
		} `json:"aria2"`
	} `json:"input"`
	Output struct {
		Telegram struct {
			BotKey string `json:"bot-key"`
			UserID string `json:"user-id"`
		} `json:"telegram"`
	} `json:"output"`
	MaxIndex       int    `json:"max-index"`
	Sign           string `json:"sign"`
	Language       string `json:"language"`
	DownloadFolder string `json:"downloadFolder"`
	MoveFolder     string `json:"moveFolder"`
	Server         struct {
		IsServer       bool   `json:"isServer"`
		IsMasterServer bool   `json:"isMasterServer"`
		ServerHost     string `json:"serverHost"`
		ServerPort     int    `json:"serverPort"`
	} `json:"server"`
	Log struct {
		LogPath string `json:"logPath"`
		ErrPath string `json:"errPath"`
		Level   string `json:"level"`
	} `json:"log"`
}
