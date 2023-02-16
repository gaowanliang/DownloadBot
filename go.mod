module DownloadBot

go 1.17

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/gorilla/websocket v1.4.2
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/nicksnyder/go-i18n/v2 v2.1.1
	github.com/shirou/gopsutil v3.20.11+incompatible
	github.com/wxnacy/wgo v1.0.4
	go.uber.org/zap v1.24.0
	golang.org/x/text v0.6.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	googledrive v0.0.0
	onedrive v0.0.0
)

require (
	cloud.google.com/go v0.107.0 // indirect
	cloud.google.com/go/compute v1.15.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/buger/jsonparser v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/oauth2 v0.4.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	google.golang.org/api v0.103.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
)

replace onedrive => ./tool/cloudDrive/onedrive

replace googledrive => ./tool/cloudDrive/googledrive
