package cloudDrive

import (
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/internal/config"
	"googledrive"
	"onedrive"
	"strings"
)

func GetNewOneDriveInfo(url string) string {
	return onedrive.ApplyForNewPass(url)
}

func GetGoogleDriveAuthCodeURL() string {
	return googledrive.GetURL()
}
func GetNewGoogleDriveInfo(code string) string {
	return googledrive.CreateNewInfo(code)
}

func UploadDFToGoogleDrive(infoPath string, uploadComplete func(), sendAutoUpdateMessage func(string)) {
	uploadComplete()
	//log.Println(strings.ReplaceAll(info.DownloadFolder, "\\", "/"))
	go googledrive.Upload(strings.ReplaceAll(infoPath, "\\", "/"), strings.ReplaceAll(config.GetDownloadFolder(), "\\", "/"), func() func(text string) {
		return sendAutoUpdateMessage
	}, func(text string) string {
		return i18nLoc.LocText(text)
	})
}

func UploadDFToOneDrive(infoPath string, uploadComplete func(), sendAutoUpdateMessage func(text string)) {
	uploadComplete()
	//log.Println(strings.ReplaceAll(info.DownloadFolder, "\\", "/"))
	go onedrive.Upload(strings.ReplaceAll(infoPath, "\\", "/"), strings.ReplaceAll(config.GetDownloadFolder(), "\\", "/"), 3, func() func(text string) {
		return sendAutoUpdateMessage
	}, func(text string) string {
		return i18nLoc.LocText(text)
	})
}
