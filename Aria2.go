package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"v2/rpc"
)

var aria2Rpc rpc.Client

func aria2Load() {
	var err error
	aria2Rpc, err = rpc.New(context.Background(), info.Aria2Server, info.Aria2Key, time.Second*10, &Aria2Notifier{})
	dropErr(err)
	log.Printf(locText("connectTo"), info.Aria2Server)
	version, err := aria2Rpc.GetVersion()
	dropErr(err)
	log.Printf(locText("connectSuccess"), version.Version)
}

func formatTellSomething(info []rpc.StatusInfo, err error) string {
	dropErr(err)
	res := ""
	log.Printf("%+v\n", info)
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
			m["Progress"] = strconv.FormatFloat(completedLength*100.0/bytes, 'f', 2, 64) + " %"
			// m["Threads"] = fmt.Sprint(len(File.URIs))
			downloadSpeed, err := strconv.ParseFloat(Files.DownloadSpeed, 64)
			dropErr(err)
			m["Speed"] = byte2Readable(downloadSpeed)
			m["Status"] = statusFlag[Files.Status]
			if Files.Status == "paused" {
				res += fmt.Sprintf(locText("queryInformationFormat1"), m["GID"], m["Name"], m["Progress"], m["Size"])
			} else if Files.Status == "complete" || Files.Status == "removed" {
				res += fmt.Sprintf(locText("queryInformationFormat2"), m["GID"], m["Name"], m["Status"], m["Progress"], m["Size"])
			} else {
				res += fmt.Sprintf(locText("queryInformationFormat3"), m["GID"], m["Name"], m["Progress"], m["Size"], m["Speed"])
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
				m["Progress"] = strconv.FormatFloat(completedLength*100.0/bytes, 'f', 2, 64) + " %"
				m["Threads"] = fmt.Sprint(len(File.URIs))
				downloadSpeed, err := strconv.ParseFloat(Files.DownloadSpeed, 64)
				dropErr(err)
				m["Speed"] = byte2Readable(downloadSpeed)
				m["Status"] = statusFlag[Files.Status]
				if Files.Status == "paused" {
					res += fmt.Sprintf(locText("queryInformationFormat1"), m["GID"], m["Name"], m["Progress"], m["Size"])
				} else if Files.Status == "complete" || Files.Status == "removed" {
					res += fmt.Sprintf(locText("queryInformationFormat2"), m["GID"], m["Name"], m["Status"], m["Progress"], m["Size"])
				} else {
					res += fmt.Sprintf(locText("queryInformationFormat3"), m["GID"], m["Name"], m["Progress"], m["Size"], m["Speed"])
				}
			}
		}

		if index != len(info) {
			res += "\n\n"
		}
	}
	return res
}

func formatGidAndName(info []rpc.StatusInfo, err error) []map[string]string {
	dropErr(err)

	m := make([]map[string]string, 0)
	log.Printf("%+v\n", info)
	for _, Files := range info {
		for _, File := range Files.Files {
			ms := make(map[string]string)
			//paths, fileName := filepath.Split(files)
			ms["GID"] = Files.Gid
			countSplit := strings.Split(File.Path, "/")
			ms["Name"] = countSplit[len(countSplit)-1]
			m = append(m, ms)
		}
	}
	return m
}

func tellName(info rpc.StatusInfo, err error) string {
	dropErr(err)
	log.Printf("%+v\n", info)
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
