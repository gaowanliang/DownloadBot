package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
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

var bundle *i18n.Bundle
var loc *i18n.Localizer

func locLan(locLanguaged string) {
	_, err := os.Stat(info.DownloadFolder)
	dropErr(err)

	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = os.Stat("i18n")
	if err != nil {
		err := os.Mkdir("i18n", 0666)
		dropErr(err)
	}
	_, err = os.Stat(fmt.Sprintf("i18n/active.%s.json", locLanguaged))
	if err != nil {
		resp, err := http.Get(fmt.Sprintf("https://cdn.jsdelivr.net/gh/gaowanliang/DownloadBot/i18n/active.%s.json", locLanguaged))
		dropErr(err)
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		dropErr(err)
		ioutil.WriteFile(fmt.Sprintf("i18n/active.%s.json", locLanguaged), data, 0644)
	}
	rd, err := ioutil.ReadDir("i18n")
	dropErr(err)
	for _, fi := range rd {
		if !fi.IsDir() && path.Ext(fi.Name()) == ".json" {
			bundle.LoadMessageFile("i18n/" + fi.Name())
		}
	}
	loc = i18n.NewLocalizer(bundle, locLanguaged)

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

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func toInt(text string) int {
	i, err := strconv.Atoi(text)
	dropErr(err)
	return i
}
