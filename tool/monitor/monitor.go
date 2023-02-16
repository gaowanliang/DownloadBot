package monitor

import (
	"DownloadBot/internal/config"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"os"
	"strings"
	"time"
)

/*
about system information
*/

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetDiskPercent(path string) float64 {
	//parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(path)
	return diskInfo.UsedPercent
}

// IsLocal checks if the uri is local
func IsLocal(uri string) bool {
	_, err := os.Stat(config.GetDownloadFolder())
	if err != nil {
		return false
	}
	return strings.Contains(uri, "127.0.0.1") || strings.Contains(uri, "localhost")
}
