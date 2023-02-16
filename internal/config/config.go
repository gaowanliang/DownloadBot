package config

import (
	"DownloadBot/model"
	"encoding/json"
	"log"
	"os"
)

var info model.Config

func dropErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func InitConfig(configPath string, clientSign string) {
	filePtr, err := os.Open(configPath)
	dropErr(err)
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	dropErr(err)
	if clientSign != "" {
		info.Sign = clientSign
	}
}

func GetLanguage() string {
	return info.Language
}

func GetLogPath() string {
	return info.Log.LogPath
}
func GetErrPath() string {
	return info.Log.ErrPath
}
func GetLogLevel() string {
	return info.Log.Level
}

func GetDownloadFolder() string {
	return info.DownloadFolder
}
func GetMoveFolder() string {
	return info.MoveFolder
}

func GetAria2Server() string {
	return info.Input.Aria2.Aria2Server
}
func GetAria2Key() string {
	return info.Input.Aria2.Aria2Key
}

//OutputToolMethod is the type of tool used to judge the output, and 0 is telegram,-1 is unconfirmed output
func OutputToolMethod() int {
	if info.Output.Telegram.UserID != "" {
		return 0
	}
	return -1
}

//InputToolMethod is the type of tool used to judge the input, and 0 is aria2,-1 is unconfirmed input
func InputToolMethod() int {
	if info.Input.Aria2.Aria2Key != "" {
		return 0
	}
	return -1
}

func GetTelegramBotKey() string {
	return info.Output.Telegram.BotKey
}

func GetTelegramUserID() string {
	return info.Output.Telegram.UserID
}

//GetMaxIndex is the maximum number of shows
func GetMaxIndex() int {
	return info.MaxIndex
}

func GetSign() string {
	return info.Sign
}

func IsServer() bool {
	return info.Server.IsServer
}

func IsMasterServer() bool {
	return info.Server.IsMasterServer
}

func GetServerIP() string {
	return info.Server.ServerHost
}
func GetServerPort() int {
	return info.Server.ServerPort
}
