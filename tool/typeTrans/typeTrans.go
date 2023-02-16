package typeTrans

import (
	logger "DownloadBot/tool/zap"
	"strconv"
)

func dropErr(err error) {
	if err != nil {
		logger.Panic("%w", err)
	}
}

func Str2Float64(text string) float64 {
	res, err := strconv.ParseFloat(text, 64)
	dropErr(err)
	return res
}

func Str2Int(text string) int {
	i, err := strconv.Atoi(text)
	dropErr(err)
	return i
}

func Str2Int64(text string) int64 {
	i, err := strconv.ParseInt(text, 10, 64)
	dropErr(err)
	return i
}

func Byte2Readable(bytes float64) string {
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
