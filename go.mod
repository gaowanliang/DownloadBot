module DownloadBot

go 1.14

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/gorilla/websocket v1.4.2
	github.com/nicksnyder/go-i18n/v2 v2.1.1
	github.com/shirou/gopsutil v3.20.11+incompatible
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	golang.org/x/sys v0.0.0-20201221093633-bc327ba9c2f0 // indirect
	golang.org/x/text v0.3.3
	onedrive v0.0.0
	googledrive v0.0.0
)

replace onedrive => ./src/onedrive
replace googledrive => ./src/googledrive