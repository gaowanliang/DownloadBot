package aria2

import (
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/internal/config"
	rpc2 "DownloadBot/tool/input/aria2/rpc"
	"DownloadBot/tool/monitor"
	"DownloadBot/tool/typeTrans"
	logger "DownloadBot/tool/zap"
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var aria2Rpc rpc2.Client

// TMMessageChan is set Torrent/Magnet download mode
var TMMessageChan = make(chan string, 3)

// TMAllowDownloads is Allow Torrent/Magnet download List
var TMAllowDownloads = make(map[string]int, 0)
var ariaDisconnectionChan = *(rpc2.CreateAriaDisconnectionChan())

var SuddenlyMsgChan = make(chan string, 3)

type Aria2 struct {
	Notifier rpc2.Notifier
}

func dropErr(err error) {
	if err != nil {
		logger.Panic("%w", err)
	}
}

func testTMStop(TMStop func(gid string)) {
	for {
		gid := <-TMMessageChan
		//log.Println(a)
		if _, have := TMAllowDownloads[gid]; !have {
			dlInfo, err := aria2Rpc.TellStatus(gid)
			dropErr(err)
			//log.Printf("%+v\n", dlInfo)
			if dlInfo.BitTorrent.Info.Name != "" {
				aria2Rpc.Pause(gid)
				TMStop(gid)
			}
		} else {
			delete(TMAllowDownloads, gid)
		}

	}
}
func disconnectionMonitoring(notifier rpc2.Notifier) {
	timeout := 1
	for {
		res := <-ariaDisconnectionChan
		log.Println(i18nLoc.LocText("ariaDisconnect"))
		if res == "websocket: close 1006 (abnormal closure): unexpected EOF" {
			aria2Rpc.Close()
			var err error
			aria2Rpc, err = rpc2.New(context.Background(), config.GetAria2Server(), config.GetAria2Key(), time.Second*10, notifier)
			for err != nil {
				log.Printf(i18nLoc.LocText("reconnectionFailed"), timeout, timeout)
				time.Sleep(time.Second * time.Duration(timeout))
				timeout++
				aria2Rpc, err = rpc2.New(context.Background(), config.GetAria2Server(), config.GetAria2Key(), time.Second*10, notifier)
			}
			version, err := aria2Rpc.GetVersion()
			dropErr(err)
			logger.Info(fmt.Sprintf(i18nLoc.LocText("connectSuccess"), version.Version))
			timeout = 1
		}
	}
}

func (a Aria2) Load(notifier rpc2.Notifier, TMStop func(gid string), needWait bool) {
	var err error
	var wg sync.WaitGroup
	aria2Rpc, err = rpc2.New(context.Background(), config.GetAria2Server(), config.GetAria2Key(), time.Second*10, notifier)
	if err != nil {
		logger.Panic(i18nLoc.LocText("Aria2 RPC connection failed"))
	}
	// dropErr(err)
	logger.Info(fmt.Sprintf(i18nLoc.LocText("connectTo"), config.GetAria2Server()))
	version, err := aria2Rpc.GetVersion()
	if err != nil {
		logger.Panic(i18nLoc.LocText("Aria2 Version query failed"))
	}
	// dropErr(err)
	logger.Info(fmt.Sprintf(i18nLoc.LocText("connectSuccess"), version.Version))
	if needWait {

		wg.Add(1)
	}

	go testTMStop(TMStop) // Torrent/Magnet download mode
	// wg.Add(1)
	go disconnectionMonitoring(notifier)
	if needWait {
		wg.Wait()
	}

}
func (a Aria2) Close() {
	defer aria2Rpc.Close()
}

func formatTellSomething(info []rpc2.StatusInfo, err error) string {
	dropErr(err)
	//log.Printf("%+v\n\n", info)
	res := ""
	var statusFlag = map[string]string{"active": i18nLoc.LocText("active"), "paused": i18nLoc.LocText("paused"), "complete": i18nLoc.LocText("complete"), "removed": i18nLoc.LocText("removed")}
	for index, Files := range info {
		if Files.BitTorrent.Info.Name != "" {
			m := make(map[string]string)
			//paths, fileName := filepath.Split(files)
			m["GID"] = Files.Gid
			m["Name"] = Files.BitTorrent.Info.Name
			bytes, err := strconv.ParseFloat(Files.TotalLength, 64)
			dropErr(err)
			m["Size"] = typeTrans.Byte2Readable(bytes)
			completedLength, err := strconv.ParseFloat(Files.CompletedLength, 64)
			dropErr(err)
			m["CompletedLength"] = typeTrans.Byte2Readable(completedLength)
			m["Progress"] = printProgressBar(completedLength * 100.0 / bytes)
			m["Threads"] = "-"
			m["Seeders"] = Files.NumSeeders
			m["Peers"] = Files.Connections
			downloadSpeed, err := strconv.ParseFloat(Files.DownloadSpeed, 64)
			dropErr(err)
			m["Speed"] = typeTrans.Byte2Readable(downloadSpeed)
			m["Status"] = statusFlag[Files.Status]
			day, hours, minutes, seconds := resolveTime(int((bytes - completedLength) / downloadSpeed))
			m["remainingTime"] = ""
			//log.Println(hours, minutes, seconds, int((bytes-completedLength)/downloadSpeed))
			if day > 0 {
				m["remainingTime"] += fmt.Sprintf(i18nLoc.LocText("onlyDays"), day)
			}
			if hours > 0 {
				m["remainingTime"] += fmt.Sprintf(i18nLoc.LocText("onlyHours"), hours)
			}
			if minutes > 0 {
				m["remainingTime"] += fmt.Sprintf(i18nLoc.LocText("onlyMinutes"), minutes)
			}
			if seconds > 0 {
				m["remainingTime"] += fmt.Sprintf(i18nLoc.LocText("onlySeconds"), seconds)
			}
			if m["remainingTime"] == "" {
				m["remainingTime"] = i18nLoc.LocText("UnableEstimate")
			}

			if Files.Status == "paused" {
				//res += fmt.Sprintf(locText("queryInformationFormat1"), m["GID"], m["Name"], m["Progress"], m["Size"])
				res += fmt.Sprintf(i18nLoc.LocText("queryInformationFormat1"), m["Name"], m["Progress"], m["CompletedLength"], m["Size"], m["Threads"], m["GID"])
			} else if Files.Status == "complete" || Files.Status == "removed" {
				//res += fmt.Sprintf(locText("queryInformationFormat2"), m["GID"], m["Name"], m["Status"], m["Progress"], m["Size"])
				res += fmt.Sprintf(i18nLoc.LocText("queryInformationFormat2"), m["Name"], m["Status"], m["Progress"], m["CompletedLength"], m["Size"], m["Threads"], m["GID"])
			} else {
				//res += fmt.Sprintf(locText("queryInformationFormat3"), m["GID"], m["Name"], m["Progress"], m["Size"], m["Speed"])
				res += fmt.Sprintf(i18nLoc.LocText("queryInformationBTFormat3"), m["Name"], m["Progress"], m["CompletedLength"], m["Size"], m["Speed"], m["remainingTime"], m["Seeders"], m["Peers"], m["GID"])
			}
		} else {
			for _, File := range Files.Files {
				m := make(map[string]string)
				//paths, fileName := filepath.Split(files)
				m["GID"] = Files.Gid
				countSplit := strings.Split(File.Path, "/")
				m["Name"] = countSplit[len(countSplit)-1]
				bytes, err := strconv.ParseFloat(Files.TotalLength, 64)
				dropErr(err)
				m["Size"] = typeTrans.Byte2Readable(bytes)

				completedLength, err := strconv.ParseFloat(Files.CompletedLength, 64)
				dropErr(err)
				m["CompletedLength"] = typeTrans.Byte2Readable(completedLength)
				m["Progress"] = printProgressBar(completedLength * 100.0 / bytes)
				m["Threads"] = fmt.Sprint(len(File.URIs))
				downloadSpeed, err := strconv.ParseFloat(Files.DownloadSpeed, 64)
				dropErr(err)
				m["Speed"] = typeTrans.Byte2Readable(downloadSpeed)
				m["Status"] = statusFlag[Files.Status]
				m["remainingTime"] = ""
				day, hours, minutes, seconds := resolveTime(int((bytes - completedLength) / downloadSpeed))
				//log.Println(hours, minutes, seconds, int((bytes-completedLength)/downloadSpeed))
				if day > 0 {
					m["remainingTime"] += fmt.Sprintf(i18nLoc.LocText("onlyDays"), day)
				}
				if hours > 0 {
					m["remainingTime"] += fmt.Sprintf(i18nLoc.LocText("onlyHours"), hours)
				}
				if minutes > 0 {
					m["remainingTime"] += fmt.Sprintf(i18nLoc.LocText("onlyMinutes"), minutes)
				}
				if seconds > 0 {
					m["remainingTime"] += fmt.Sprintf(i18nLoc.LocText("onlySeconds"), seconds)
				}

				if Files.Status == "paused" {
					//res += fmt.Sprintf(locText("queryInformationFormat1"), m["GID"], m["Name"], m["Progress"], m["Size"])
					res += fmt.Sprintf(i18nLoc.LocText("queryInformationFormat1"), m["Name"], m["Progress"], m["CompletedLength"], m["Size"], m["Threads"], m["GID"])
				} else if Files.Status == "complete" || Files.Status == "removed" {
					//res += fmt.Sprintf(locText("queryInformationFormat2"), m["GID"], m["Name"], m["Status"], m["Progress"], m["Size"])
					res += fmt.Sprintf(i18nLoc.LocText("queryInformationFormat2"), m["Name"], m["Status"], m["Progress"], m["CompletedLength"], m["Size"], m["Threads"], m["GID"])
				} else {
					//res += fmt.Sprintf(locText("queryInformationFormat3"), m["GID"], m["Name"], m["Progress"], m["Size"], m["Speed"])
					res += fmt.Sprintf(i18nLoc.LocText("queryInformationFormat3"), m["Name"], m["Progress"], m["CompletedLength"], m["Size"], m["Speed"], m["remainingTime"], m["Threads"], m["GID"])
				}
			}
		}

		if index != len(info) {
			res += "\n\n"
		}
	}
	if res != "" {
		totalSpeed, err := aria2Rpc.GetGlobalStat()
		dropErr(err)
		res += fmt.Sprintf(i18nLoc.LocText("systemInfo"), monitor.GetCpuPercent(), monitor.GetDiskPercent(config.GetDownloadFolder()), monitor.GetMemPercent(), typeTrans.Byte2Readable(typeTrans.Str2Float64(totalSpeed.DownloadSpeed)), typeTrans.Byte2Readable(typeTrans.Str2Float64(totalSpeed.UploadSpeed)))
	}
	return res
}

func (a Aria2) FormatTellActive() string {
	return formatTellSomething(aria2Rpc.TellActive())
}
func (a Aria2) FormatTellWaiting() string {
	return formatTellSomething(aria2Rpc.TellWaiting(0, config.GetMaxIndex()))
}
func (a Aria2) FormatTellStopped() string {
	return formatTellSomething(aria2Rpc.TellStopped(0, config.GetMaxIndex()))
}

//
// FormatGidAndName
//  @Description: Provide formatted GID and name according to method
//  @param method int "0: active, 1: waiting"
//  @return []map[string]string
//
func (a Aria2) FormatGidAndName(method int) []map[string]string {
	var info []rpc2.StatusInfo
	var err error
	switch method {
	case 0:
		info, err = aria2Rpc.TellActive()
		break
	case 1:
		info, err = aria2Rpc.TellWaiting(0, config.GetMaxIndex())
		break
	default:
		logger.Panic("method error")
	}
	dropErr(err)

	m := make([]map[string]string, 0)
	//log.Printf("%+v\n", info)
	for _, Files := range info {
		// log.Println(Files.BitTorrent.Info.Name, Files.BitTorrent.Info.Name != "")
		if Files.BitTorrent.Info.Name != "" {
			ms := make(map[string]string)
			ms["GID"] = Files.Gid
			ms["Name"] = Files.BitTorrent.Info.Name
			m = append(m, ms)
		} else {
			for _, File := range Files.Files {
				ms := make(map[string]string)
				//paths, fileName := filepath.Split(files)
				ms["GID"] = Files.Gid
				countSplit := strings.Split(File.Path, "/")
				ms["Name"] = countSplit[len(countSplit)-1]
				m = append(m, ms)
			}
		}
	}
	return m
}

func (a Aria2) TellName(gid string) string {
	info, err := aria2Rpc.TellStatus(gid)
	dropErr(err)
	logger.Debug("GID info:%v", info)
	Name := ""
	if info.BitTorrent.Info.Name != "" {
		Name = info.BitTorrent.Info.Name
	} else {
		for _, File := range info.Files {
			if File.Path != "" {
				countSplit := strings.Split(File.Path, "/")
				Name = countSplit[len(countSplit)-1]
			} else {
				Name = info.Gid
			}
		}
	}
	return Name
}

func tellName(gid string) string {

	info, err := aria2Rpc.TellStatus(gid)
	dropErr(err)
	logger.Debug("GID info:%v", info)
	Name := ""
	if info.BitTorrent.Info.Name != "" {
		Name = info.BitTorrent.Info.Name
	} else {
		for _, File := range info.Files {
			if File.Path != "" {
				countSplit := strings.Split(File.Path, "/")
				Name = countSplit[len(countSplit)-1]
			} else {
				Name = info.Gid
			}
		}
	}
	return Name
}

func (a Aria2) Download(uri string) bool {
	uriType := isDownloadType(uri)
	if uriType == 0 {
		return false
	}
	uriData := make([]string, 0)
	uriData = append(uriData, uri)
	switch uriType {
	case 1:
		aria2Rpc.AddURI(uriData)
	case 2:
		aria2Rpc.AddURI(uriData)
	case 3:
		aria2Rpc.AddTorrent(uri)
	}
	return true
}

// FormatTMFiles is a function that can format the file information of torrent/magnet file, the return value is [][2]string,0 is file name,1 is file size
func (a Aria2) FormatTMFiles(gid string) [][]string {
	var fileList [][]string
	rawList, err := aria2Rpc.GetFiles(gid)

	dropErr(err)
	// log.Printf("%+v", rawList)
	for _, file := range rawList {
		fileInfo := make([]string, 0)
		fileInfo = append(fileInfo, file.Path)
		bytes, err := strconv.ParseFloat(file.Length, 64)
		dropErr(err)
		fileInfo = append(fileInfo, typeTrans.Byte2Readable(bytes))
		fileInfo = append(fileInfo, file.Selected)
		fileList = append(fileList, fileInfo)
	}
	return fileList
}

func (a Aria2) SetTMDownloadFilesAndStart(gid string, FilesList [][2]int) {
	selectFile := ""
	for _, file := range FilesList {
		if file[0] == 1 && file[1] >= 0 {
			selectFile += fmt.Sprint(file[1]) + ","
		}
	}
	aria2Rpc.ChangeOption(gid, rpc2.Option{
		"select-file": selectFile[:len(selectFile)-1], // remove the last comma
	})
	TMAllowDownloads[gid] = 0
	aria2Rpc.Unpause(gid)

}
func (a Aria2) SelectBiggestFile(gid string) int {
	index := 0
	rawList, err := aria2Rpc.GetFiles(gid)
	dropErr(err)
	for i := 0; i < len(rawList); i++ {
		if typeTrans.Str2Int(rawList[i].Length) > typeTrans.Str2Int(rawList[index].Length) {
			index = i
		}
	}
	return index + 1
}
func (a Aria2) SelectBigFiles(gid string) []int {
	index := make([]int, 0)
	rawList, err := aria2Rpc.GetFiles(gid)
	dropErr(err)
	totalSize, avgSize := 0, 0.0
	for _, file := range rawList {
		totalSize += typeTrans.Str2Int(file.Length)
	}
	avgSize = float64(totalSize) / float64(len(rawList))
	avgSize -= avgSize * 0.2

	for i := 0; i < len(rawList); i++ {
		dist, err := strconv.ParseFloat(rawList[i].Length, 64)
		dropErr(err)
		if dist > avgSize {
			index = append(index, i+1)
		}
	}
	return index
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

func printProgressBar(progress float64) string {
	progressBar := "["
	if progress != 100 {
		for i := 0; i < int(progress/7.7); i++ {
			progressBar += "●"
		}
		for i := 0; i < 13-int(progress/7.7); i++ {
			progressBar += "○"
		}
	} else {
		progressBar += "●●●●●●●●●●●●●"
	}
	progressBar += "] " + strconv.FormatFloat(progress, 'f', 2, 64) + " %"
	return progressBar
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

func (a Aria2) Pause(gid string) {
	aria2Rpc.Pause(gid)
}
func (a Aria2) Unpause(gid string) {
	aria2Rpc.Unpause(gid)
}

func (a Aria2) ForceRemove(gid string) {
	aria2Rpc.ForceRemove(gid)
}
func (a Aria2) PauseAll() {
	aria2Rpc.PauseAll()
}
func (a Aria2) UnpauseAll() {
	aria2Rpc.UnpauseAll()
}

func (a Aria2) GetVersion() string {
	version, err := aria2Rpc.GetVersion()
	dropErr(err)
	return version.Version
}
