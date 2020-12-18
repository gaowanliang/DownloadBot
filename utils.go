package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	goTree "v2/gotree"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func byte2Readable(bytes float64) string {
	const kb float64 = 1024
	const mb = kb * 1024
	const gb = mb * 1024
	var readable float64
	var unit string
	_bytes := bytes

	if _bytes >= gb {
		// xx GB
		readable = _bytes / gb
		unit = "GB"
	} else if _bytes < gb && _bytes >= mb {
		// xx MB
		readable = _bytes / mb
		unit = "MB"
	} else {
		// xx KB
		readable = _bytes / kb
		unit = "KB"
	}
	return strconv.FormatFloat(readable, 'f', 2, 64) + " " + unit
}

func isDownloadType(uri string) int {
	httpFtp, _ := regexp.MatchString(`^(https?|ftps?)://.*$`, uri)
	magnet, _ := regexp.MatchString(`(?i)magnet:\?xt=urn:[a-z0-9]+:[a-z0-9]{32}`, uri)
	btFile, _ := regexp.MatchString(`\.torrent$`, uri)
	if httpFtp {
		return 1
	} else if magnet {
		return 2
	} else if btFile {
		return 3
	} else {
		return 0
	}
}

var bundle *i18n.Bundle
var loc *i18n.Localizer

func locLan(locLanguage string) {
	_, err := os.Stat(info.DownloadFolder)
	dropErr(err)

	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = os.Stat("i18n")
	if err != nil {
		err := os.Mkdir("i18n", 0666)
		dropErr(err)
	}
	_, err = os.Stat(fmt.Sprintf("i18n/active.%s.json", locLanguage))
	if err != nil {
		resp, err := http.Get(fmt.Sprintf("https://cdn.jsdelivr.net/gh/gaowanliang/DownloadBot/i18n/active.%s.json", locLanguage))
		dropErr(err)
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		dropErr(err)
		ioutil.WriteFile(fmt.Sprintf("i18n/active.%s.json", locLanguage), data, 0644)
	}
	rd, err := ioutil.ReadDir("i18n")
	dropErr(err)
	for _, fi := range rd {
		if !fi.IsDir() && path.Ext(fi.Name()) == ".json" {
			bundle.LoadMessageFile("i18n/" + fi.Name())
		}
	}
	loc = i18n.NewLocalizer(bundle, locLanguage)

}

func locText(MessageIDs ...string) string {
	res := ""
	for _, MessageID := range MessageIDs {
		res += loc.MustLocalize(&i18n.LocalizeConfig{MessageID: MessageID})
	}
	return res
}

func isLocal(uri string) bool {
	return strings.Contains(uri, "127.0.0.1") || strings.Contains(uri, "localhost")
}

func removeContents(boot int, pathname string, fileSelect map[int]bool) int {
	index := boot
	if pathname[:len(pathname)-1] != "/" {
		pathname += "/"
	}
	rd, err := ioutil.ReadDir(pathname)
	dropErr(err)
	res := make([]string, 0)
	res = append(res, pathname)
	fileCount := 0
	for _, fi := range rd {
		if fi.IsDir() {
			index++
			fileCount++
			res = append(res, pathname+fi.Name())
			index += removeContents(index, pathname+fi.Name(), fileSelect)
		} else {
			fileCount++
			res = append(res, pathname+fi.Name())
		}
	}

	for _, removePath := range res {
		index++
		log.Println(index, removePath)
		if fileSelect[index] || index != 1 {
			err = os.RemoveAll(removePath)
			dropErr(err)
		}
	}
	return fileCount
}

func RemoveFiles(deleteFiles []string) {
	//removeContents(1, info.DownloadFolder, fileSelect)
	for _, removePath := range deleteFiles {
		log.Println(removePath)
		if removePath != info.DownloadFolder && removePath != info.DownloadFolder+"/" {
			err := os.RemoveAll(removePath)
			dropErr(err)
		}
	}
}
func GetAllFile(pathname string) ([][]string, int64) {
	if pathname[:len(pathname)-1] != "/" {
		pathname += "/"
	}
	rd, err := ioutil.ReadDir(pathname)
	dropErr(err)
	res := make([][]string, 0)
	//
	totalSize := int64(0)
	for _, fi := range rd {
		if fi.IsDir() {
			ret, subSize := GetAllFile(pathname + fi.Name())
			totalSize += subSize
			res = append(res, ret...)
		} else {
			var fileNameAndSize = []string{pathname + fmt.Sprintln(fi.Name()), fmt.Sprint(fi.Size())}
			res = append(res, fileNameAndSize)
			totalSize += fi.Size()
		}
	}
	var fileNameAndSize = []string{pathname + fmt.Sprintln(pathname), fmt.Sprint(totalSize)}
	res = append(res, fileNameAndSize)
	return res, totalSize
}

func generateFolderTree(pathname string, boot int, fileSelect map[int]bool, selectFileIndex string, parentSelected int8) (goTree.Tree, [][]string, map[int]bool, []string) {
	index := boot
	if pathname[:len(pathname)-1] != "/" {
		pathname += "/"
	}
	rd, err := ioutil.ReadDir(pathname)
	dropErr(err)
	res := make([][]string, 0)
	totalSize := 0
	treeFolder := make([]goTree.Tree, 0)
	treeFiles := make([]string, 0)
	trueFileSelect := make(map[int]bool, 0)
	subList := make(map[int]bool)
	subFilesPath := make([]string, 0)
	deleteFiles := make([]string, 0)
	var artist goTree.Tree
	bootSelect := int8(0)
	if selectFileIndex == fmt.Sprint(boot) || parentSelected != 0 {
		if fileSelect[boot] || parentSelected == -1 {
			bootSelect = -1 // 其下皆不选
			trueFileSelect[boot] = false
			artist = goTree.New(fmt.Sprintf("⬜%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
		} else {
			bootSelect = 1 // 其下皆选
			trueFileSelect[boot] = true
			artist = goTree.New(fmt.Sprintf("✅%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
		}
	}
	for _, fi := range rd {
		//log.Println(fi.Name())
		if fi.IsDir() {
			index++
			ret, filesInfo, subTrueFileSelect, subDeleteFiles := generateFolderTree(pathname+fi.Name(), index, fileSelect, selectFileIndex, bootSelect)
			subFolderSize := 0
			for _, iSize := range filesInfo {
				subFolderSize += toInt(iSize[1])
			}
			for _, iPath := range subDeleteFiles {
				deleteFiles = append(deleteFiles, iPath)
			}
			var folderNameAndSize = []string{pathname + fi.Name(), fmt.Sprint(subFolderSize)}
			res = append(res, folderNameAndSize)
			res = append(res, filesInfo...)
			treeFolder = append(treeFolder, ret)
			totalSize += subFolderSize
			for k, v := range subTrueFileSelect {

				trueFileSelect[k] = v
			}

			subList[index] = subTrueFileSelect[index]
			index += len(filesInfo)
		} else {
			treeFiles = append(treeFiles, fmt.Sprintf("%s * %s", fi.Name(), byte2Readable(toFloat64(fmt.Sprint(fi.Size())))))
			subFilesPath = append(subFilesPath, pathname+fmt.Sprint(fi.Name()))
			var fileNameAndSize = []string{pathname + fmt.Sprint(fi.Name()), fmt.Sprint(fi.Size())}
			res = append(res, fileNameAndSize)
			totalSize += int(fi.Size())
		}
	}

	tempIndex := index
	for _, _ = range treeFiles {
		index++
		if selectFileIndex == "selectAll" {
			trueFileSelect[index] = true
		} else if selectFileIndex == "invert" || selectFileIndex == fmt.Sprint(index) {
			trueFileSelect[index] = !fileSelect[index]
		} else if bootSelect == 1 {
			trueFileSelect[index] = true
		} else if bootSelect == -1 {
			trueFileSelect[index] = false
		} else {
			trueFileSelect[index] = fileSelect[index]
		}
		subList[index] = trueFileSelect[index]
	}
	index = tempIndex

	if selectFileIndex == "selectAll" {
		artist = goTree.New(fmt.Sprintf("✅%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
		trueFileSelect[boot] = true
	} else if selectFileIndex == "invert" {
		artist = goTree.New(fmt.Sprintf("⬜%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
		trueFileSelect[boot] = false
	} else if _, al := trueFileSelect[toInt(selectFileIndex)]; al {
		artist = goTree.New(fmt.Sprintf("⬜%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
		if fileSelect[boot] {
			trueFileSelect[boot] = false
		} else {
			if !trueFileSelect[toInt(selectFileIndex)] {
				trueFileSelect[boot] = false
			} else {
				selectAllOthers := true
				for k, v := range subList {
					if k == toInt(selectFileIndex) {
						continue
					}
					if !v {
						selectAllOthers = false
						break
					}
				}
				if selectAllOthers {
					trueFileSelect[boot] = true
					artist = goTree.New(fmt.Sprintf("✅%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
				} else {
					trueFileSelect[boot] = false
				}
			}
		}
	} else if selectFileIndex == fmt.Sprint(boot) || parentSelected != 0 {
		if fileSelect[boot] || parentSelected == -1 {
			artist = goTree.New(fmt.Sprintf("⬜%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
		} else {
			artist = goTree.New(fmt.Sprintf("✅%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
		}
	} else {
		if fileSelect[boot] {
			artist = goTree.New(fmt.Sprintf("✅%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
			trueFileSelect[boot] = true
		} else {
			artist = goTree.New(fmt.Sprintf("⬜%d:%s * %s", boot, path.Base(pathname), byte2Readable(toFloat64(fmt.Sprint(fmt.Sprint(totalSize))))))
			trueFileSelect[boot] = false
		}
	}

	for _, val := range treeFolder {
		artist.AddTree(val)
	}
	for i, val := range treeFiles {
		index++
		if trueFileSelect[index] {
			artist.Add(fmt.Sprintf("✅%d:%s", index, val))
			deleteFiles = append(deleteFiles, subFilesPath[i])
		} else {
			artist.Add(fmt.Sprintf("⬜%d:%s", index, val))
		}
	}
	if trueFileSelect[boot] {
		deleteFiles = append(deleteFiles, pathname)
	}
	allFalse := true
	for _, v := range trueFileSelect {
		if v {
			allFalse = false
			break
		}
	}
	if boot == 1 && allFalse {
		//log.Println("ss")
		return generateFolderTree(pathname, 1, fileSelect, "0", int8(0))
	} else {
		return artist, res, trueFileSelect, deleteFiles
	}

}
func printFolderTree(pathName string, fileSelect map[int]bool, selectFileIndex string) (string, map[int]bool, []string) {
	tree, _, trueFileSelect, deleteFiles := generateFolderTree(pathName, 1, fileSelect, selectFileIndex, int8(0))
	return tree.Print(), trueFileSelect, deleteFiles
}

func printProgressBar(progress float64) string {
	progressBar := "["
	for i := 0; i < int(progress/7.7); i++ {
		progressBar += "●"
	}
	for i := 0; i < 13-int(progress/7.7); i++ {
		progressBar += "○"
	}
	progressBar += "] " + strconv.FormatFloat(progress, 'f', 2, 64) + " %"
	return progressBar
}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

const (
	// 定义每分钟的秒数
	SecondsPerMinute = 60
	// 定义每小时的秒数
	SecondsPerHour = SecondsPerMinute * 60
	// 定义每天的秒数
	SecondsPerDay = SecondsPerHour * 24
)

func resolveTime(seconds int) (day int, hour int, minute int, second int) {
	day = seconds / SecondsPerDay
	hour = (seconds - day*SecondsPerDay) / SecondsPerHour
	minute = (seconds - day*SecondsPerDay - hour*SecondsPerHour) / SecondsPerMinute
	second = seconds - day*SecondsPerDay - hour*SecondsPerHour - minute*SecondsPerMinute
	return
}

func toInt(text string) int {
	i, err := strconv.Atoi(text)
	dropErr(err)
	return i
}

func toFloat64(text string) float64 {
	res, err := strconv.ParseFloat(text, 64)
	dropErr(err)
	return res
}
func toInt64(text string) int64 {
	i, err := strconv.ParseInt(text, 10, 64)
	dropErr(err)
	return i
}
