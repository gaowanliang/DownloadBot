package controller

import (
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/internal/config"
	logger "DownloadBot/tool/zap"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func RemoveFiles(deleteFiles []string) {
	//removeContents(1, config.GetDownloadFolder(), fileSelect)
	for _, removePath := range deleteFiles {
		//log.Println(removePath)
		if removePath != config.GetDownloadFolder() && removePath != config.GetDownloadFolder()+"/" {
			err := os.RemoveAll(removePath)
			logger.DropErr(err)
		}
	}
}
func CopyFiles(srcFiles []string, sendAutoUpdateMessage func(text string)) {
	destPath := config.GetMoveFolder()
	downloadFolder := config.GetDownloadFolder()
	if destPath[:len(destPath)-1] != "/" {
		destPath += "/"
	}
	if downloadFolder[:len(downloadFolder)-1] != "/" {
		downloadFolder += "/"
	}
	newMsg := sendAutoUpdateMessage
	for _, srcPath := range srcFiles {
		if srcPath != config.GetDownloadFolder() && srcPath != config.GetDownloadFolder()+"/" {
			newMsg(fmt.Sprintf(i18nLoc.LocText("copyingTo"), srcPath, destPath+path.Base(srcPath)))
			//log.Println(srcPath)
			file1, err := os.Open(srcPath)
			logger.DropErr(err)
			s, err := os.Stat(srcPath)
			if err == nil {
				//log.Println(strings.ReplaceAll(srcPath, downloadFolder, destPath))
				if s.IsDir() {
					_, err := os.Stat(strings.ReplaceAll(srcPath, downloadFolder, destPath))
					if err != nil {
						err = os.MkdirAll(strings.ReplaceAll(srcPath, downloadFolder, destPath), os.ModePerm)
						logger.DropErr(err)
					} else {
						continue
					}
				} else {
					paths, _ := filepath.Split(strings.ReplaceAll(srcPath, downloadFolder, destPath))
					_, err := os.Stat(paths)
					if err != nil {
						err = os.MkdirAll(paths, os.ModePerm)
						logger.DropErr(err)
					}
				}
			}

			file2, err := os.OpenFile(strings.ReplaceAll(srcPath, downloadFolder, destPath), os.O_WRONLY|os.O_CREATE, os.ModePerm)
			logger.DropErr(err)
			defer file1.Close()
			defer file2.Close()
			//拷贝数据
			bs := make([]byte, 1024, 1024)
			n := -1 //读取的数据量
			total := 0
			for {
				n, err = file1.Read(bs)
				if err == io.EOF || n == 0 {
					break
				}
				logger.DropErr(err)
				total += n
				_, err = file2.Write(bs[:n])
			}
		}
	}
	newMsg("close")
}

func CreateDriveInfoFolder(path string) {
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		logger.DropErr(err)
	}

}
func GetAuthInfoJson(path string) ([]string, string, int) {
	text := ""
	files := make([]string, 0)
	rd, err := ioutil.ReadDir(path)
	logger.DropErr(err)
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

// RandStringRunes Generate random string
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
