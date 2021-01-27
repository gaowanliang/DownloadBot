package main

import (
	"fmt"
	"googledrive"
	"io/ioutil"
	"onedrive"
	"os"
	"strings"
)

func createDriveInfoFolder(path string) {
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		dropErr(err)
	}
}
func getAuthInfoJson(path string) ([]string, string, int) {
	text := ""
	files := make([]string, 0)
	rd, err := ioutil.ReadDir(path)
	dropErr(err)
	index := 1
	for _, fi := range rd {
		if !fi.IsDir() {
			if strings.HasSuffix(strings.ToLower(fi.Name()), ".json") {
				files = append(files, fi.Name())
				text += fmt.Sprintf("%d.%s\n", index, fi.Name())
				index++
			}
		}
	}
	return files, text, index
}

func getNewOneDriveInfo(url string) string {
	return onedrive.ApplyForNewPass(url)
}

func uploadDFToOneDrive(infoPath string) {
	FileControlChan <- "close"
	//log.Println(strings.ReplaceAll(info.DownloadFolder, "\\", "/"))
	go onedrive.Upload(strings.ReplaceAll(infoPath, "\\", "/"), strings.ReplaceAll(info.DownloadFolder, "\\", "/"), 3, func() func(text string) {
		return sendAutoUpdateMessage()
	}, func(text string) string {
		return locText(text)
	})
}

func getGoogleDriveAuthCodeURL() string {
	return googledrive.GetURL()
}
func getNewGoogleDriveInfo(code string) string {
	return googledrive.CreateNewInfo(code)
}

func uploadDFToGoogleDrive(infoPath string) {
	FileControlChan <- "close"
	//log.Println(strings.ReplaceAll(info.DownloadFolder, "\\", "/"))
	go googledrive.Upload(strings.ReplaceAll(infoPath, "\\", "/"), strings.ReplaceAll(info.DownloadFolder, "\\", "/"), func() func(text string) {
		return sendAutoUpdateMessage()
	}, func(text string) string {
		return locText(text)
	})
}