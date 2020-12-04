package main

import (
	"regexp"
	"strconv"
)

func byte2Readable(bytes float64) string {
	const kb float64 = 1024
	const mb float64 = kb * 1024
	const gb float64 = mb * 1024
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
	httpFtp, _ := regexp.MatchString(`^(https?|ftps?):\/\/.*$`, uri)
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
