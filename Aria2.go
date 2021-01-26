package main

import (
	"DownloadBot/src/rpc"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var aria2Rpc rpc.Client

// TMMessageChan is set Torrent/Magnet download mode
var TMMessageChan = make(chan string, 3)

// TMAllowDownloads is Allow Torrent/Magnet download List
var TMAllowDownloads = make(map[string]int, 0)
var ariaDisconnectionChan = *(rpc.CreateAriaDisconnectionChan())

func testTMStop() {
	for {
		gid := <-TMMessageChan
		//log.Println(a)
		if _, have := TMAllowDownloads[gid]; !have {
			dlInfo, err := aria2Rpc.TellStatus(gid)
			dropErr(err)
			//log.Printf("%+v\n", dlInfo)
			if dlInfo.BitTorrent.Info.Name != "" {
				aria2Rpc.Pause(gid)
				TMSelectMessageChan <- gid
			}
		} else {
			delete(TMAllowDownloads, gid)
		}

	}
}
func ariaDisconnectionMonitoring() {
	timeout := 1
	for {
		res := <-ariaDisconnectionChan
		log.Println(locText("ariaDisconnect"))
		if res == "websocket: close 1006 (abnormal closure): unexpected EOF" {
			aria2Rpc.Close()
			var err error
			aria2Rpc, err = rpc.New(context.Background(), info.Aria2Server, info.Aria2Key, time.Second*10, &Aria2Notifier{})
			for err != nil {
				log.Printf(locText("reconnectionFailed"), timeout, timeout)
				time.Sleep(time.Second * time.Duration(timeout))
				timeout++
				aria2Rpc, err = rpc.New(context.Background(), info.Aria2Server, info.Aria2Key, time.Second*10, &Aria2Notifier{})
			}
			version, err := aria2Rpc.GetVersion()
			dropErr(err)
			log.Printf(locText("connectSuccess"), version.Version)
			timeout = 1
		}
	}
}

func aria2Load() {
	var err error
	aria2Rpc, err = rpc.New(context.Background(), info.Aria2Server, info.Aria2Key, time.Second*10, &Aria2Notifier{})
	dropErr(err)
	log.Printf(locText("connectTo"), info.Aria2Server)
	version, err := aria2Rpc.GetVersion()
	dropErr(err)
	log.Printf(locText("connectSuccess"), version.Version)
	go testTMStop()

	go ariaDisconnectionMonitoring()
}

func formatTellSomething(info []rpc.StatusInfo, err error) string {
	dropErr(err)
	//log.Printf("%+v\n\n", info)
	res := ""
	var statusFlag = map[string]string{"active": locText("active"), "paused": locText("paused"), "complete": locText("complete"), "removed": locText("removed")}
	for index, Files := range info {
		if Files.BitTorrent.Info.Name != "" {
			m := make(map[string]string)
			//paths, fileName := filepath.Split(files)
			m["GID"] = Files.Gid
			m["Name"] = Files.BitTorrent.Info.Name
			bytes, err := strconv.ParseFloat(Files.TotalLength, 64)
			dropErr(err)
			m["Size"] = byte2Readable(bytes)
			completedLength, err := strconv.ParseFloat(Files.CompletedLength, 64)
			dropErr(err)
			m["CompletedLength"] = byte2Readable(completedLength)
			m["Progress"] = printProgressBar(completedLength * 100.0 / bytes)
			m["Threads"] = "-"
			m["Seeders"] = Files.NumSeeders
			m["Peers"] = Files.Connections
			downloadSpeed, err := strconv.ParseFloat(Files.DownloadSpeed, 64)
			dropErr(err)
			m["Speed"] = byte2Readable(downloadSpeed)
			m["Status"] = statusFlag[Files.Status]
			day, hours, minutes, seconds := resolveTime(int((bytes - completedLength) / downloadSpeed))
			m["remainingTime"] = ""
			//log.Println(hours, minutes, seconds, int((bytes-completedLength)/downloadSpeed))
			if day > 0 {
				m["remainingTime"] += fmt.Sprintf(locText("onlyDays"), day)
			}
			if hours > 0 {
				m["remainingTime"] += fmt.Sprintf(locText("onlyHours"), hours)
			}
			if minutes > 0 {
				m["remainingTime"] += fmt.Sprintf(locText("onlyMinutes"), minutes)
			}
			if seconds > 0 {
				m["remainingTime"] += fmt.Sprintf(locText("onlySeconds"), seconds)
			}
			if m["remainingTime"] == "" {
				m["remainingTime"] = locText("UnableEstimate")
			}

			if Files.Status == "paused" {
				//res += fmt.Sprintf(locText("queryInformationFormat1"), m["GID"], m["Name"], m["Progress"], m["Size"])
				res += fmt.Sprintf(locText("queryInformationFormat1"), m["Name"], m["Progress"], m["CompletedLength"], m["Size"], m["Threads"], m["GID"])
			} else if Files.Status == "complete" || Files.Status == "removed" {
				//res += fmt.Sprintf(locText("queryInformationFormat2"), m["GID"], m["Name"], m["Status"], m["Progress"], m["Size"])
				res += fmt.Sprintf(locText("queryInformationFormat2"), m["Name"], m["Status"], m["Progress"], m["CompletedLength"], m["Size"], m["Threads"], m["GID"])
			} else {
				//res += fmt.Sprintf(locText("queryInformationFormat3"), m["GID"], m["Name"], m["Progress"], m["Size"], m["Speed"])
				res += fmt.Sprintf(locText("queryInformationBTFormat3"), m["Name"], m["Progress"], m["CompletedLength"], m["Size"], m["Speed"], m["remainingTime"], m["Seeders"], m["Peers"], m["GID"])
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
				m["Size"] = byte2Readable(bytes)

				completedLength, err := strconv.ParseFloat(Files.CompletedLength, 64)
				dropErr(err)
				m["CompletedLength"] = byte2Readable(completedLength)
				m["Progress"] = printProgressBar(completedLength * 100.0 / bytes)
				m["Threads"] = fmt.Sprint(len(File.URIs))
				downloadSpeed, err := strconv.ParseFloat(Files.DownloadSpeed, 64)
				dropErr(err)
				m["Speed"] = byte2Readable(downloadSpeed)
				m["Status"] = statusFlag[Files.Status]
				m["remainingTime"] = ""
				day, hours, minutes, seconds := resolveTime(int((bytes - completedLength) / downloadSpeed))
				//log.Println(hours, minutes, seconds, int((bytes-completedLength)/downloadSpeed))
				if day > 0 {
					m["remainingTime"] += fmt.Sprintf(locText("onlyDays"), day)
				}
				if hours > 0 {
					m["remainingTime"] += fmt.Sprintf(locText("onlyHours"), hours)
				}
				if minutes > 0 {
					m["remainingTime"] += fmt.Sprintf(locText("onlyMinutes"), minutes)
				}
				if seconds > 0 {
					m["remainingTime"] += fmt.Sprintf(locText("onlySeconds"), seconds)
				}

				if Files.Status == "paused" {
					//res += fmt.Sprintf(locText("queryInformationFormat1"), m["GID"], m["Name"], m["Progress"], m["Size"])
					res += fmt.Sprintf(locText("queryInformationFormat1"), m["Name"], m["Progress"], m["CompletedLength"], m["Size"], m["Threads"], m["GID"])
				} else if Files.Status == "complete" || Files.Status == "removed" {
					//res += fmt.Sprintf(locText("queryInformationFormat2"), m["GID"], m["Name"], m["Status"], m["Progress"], m["Size"])
					res += fmt.Sprintf(locText("queryInformationFormat2"), m["Name"], m["Status"], m["Progress"], m["CompletedLength"], m["Size"], m["Threads"], m["GID"])
				} else {
					//res += fmt.Sprintf(locText("queryInformationFormat3"), m["GID"], m["Name"], m["Progress"], m["Size"], m["Speed"])
					res += fmt.Sprintf(locText("queryInformationFormat3"), m["Name"], m["Progress"], m["CompletedLength"], m["Size"], m["Speed"], m["remainingTime"], m["Threads"], m["GID"])
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
		res += fmt.Sprintf(locText("systemInfo"), GetCpuPercent(), GetDiskPercent(), GetMemPercent(), byte2Readable(toFloat64(totalSpeed.DownloadSpeed)), byte2Readable(toFloat64(totalSpeed.UploadSpeed)))
	}
	return res
}

func formatGidAndName(info []rpc.StatusInfo, err error) []map[string]string {
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

func tellName(info rpc.StatusInfo, err error) string {
	dropErr(err)
	//log.Printf("%+v\n", info)
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

func download(uri string) bool {
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

func formatTMFiles(gid string) [][]string {
	var fileList [][]string
	rawList, err := aria2Rpc.GetFiles(gid)
	dropErr(err)
	// log.Printf("%+v", rawList)
	for _, file := range rawList {
		fileInfo := make([]string, 0)
		fileInfo = append(fileInfo, file.Path)
		bytes, err := strconv.ParseFloat(file.Length, 64)
		dropErr(err)
		fileInfo = append(fileInfo, byte2Readable(bytes))
		fileInfo = append(fileInfo, file.Selected)
		fileList = append(fileList, fileInfo)
	}
	return fileList
}

func setTMDownloadFilesAndStart(gid string, FilesList [][]string) {
	selectFile := ""
	for i, file := range FilesList {
		if file[2] == "true" {
			selectFile += fmt.Sprint(i+1) + ","
		}
	}
	aria2Rpc.ChangeOption(gid, rpc.Option{
		"select-file": selectFile[:len(selectFile)-1],
	})
	TMAllowDownloads[gid] = 0
	aria2Rpc.Unpause(gid)

}
func selectBigestFile(gid string) int {
	index := 0
	rawList, err := aria2Rpc.GetFiles(gid)
	dropErr(err)
	for i := 0; i < len(rawList); i++ {
		if rawList[i].Length > rawList[index].Length {
			index = i
		}
	}
	return index
}
func selectBigFiles(gid string) []int {
	indexs := make([]int, 0)
	rawList, err := aria2Rpc.GetFiles(gid)
	dropErr(err)
	totalSize, avgSize := 0, 0.0
	for _, file := range rawList {
		totalSize += toInt(file.Length)
	}
	avgSize = float64(totalSize) / float64(len(rawList))
	avgSize -= avgSize * 0.2

	for i := 0; i < len(rawList); i++ {
		dist, err := strconv.ParseFloat(rawList[i].Length, 64)
		dropErr(err)
		if dist > avgSize {
			indexs = append(indexs, i)
		}
	}
	return indexs
}
