package i18nLoc

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"
)

var bundle *i18n.Bundle
var loc *i18n.Localizer

// LocLan is download i18n language file to local
func LocLan(locLanguage string) {

	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err := os.Stat("i18n")
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
		ioutil.WriteFile(fmt.Sprintf("i18n/active.%s.json", locLanguage), data, 0666)
		log.Printf("Download i18n language file success, the information will be output later using the language you choose")
	} else {
		url := "https://cdn.jsdelivr.net/gh/gaowanliang/DownloadBot@latest/i18n/"
		j := pageDownload(url)
		var re = regexp.MustCompile(`(?m)i18n/(.*?)"[\s\S]*?<td class="time">(.*?)</td>`)
		var newLanFileTime int64
		for _, val := range re.FindAllStringSubmatch(j, -1) {
			if fmt.Sprintf("active.%s.json", locLanguage) == val[1] {
				t, _ := time.Parse(time.RFC1123, val[2])
				newLanFileTime = t.Unix()
			}

		}
		oldLanFileTime := GetFileModTime(fmt.Sprintf("i18n/active.%s.json", locLanguage))
		if newLanFileTime > oldLanFileTime {
			err := os.RemoveAll(fmt.Sprintf("i18n/active.%s.json", locLanguage))
			dropErr(err)
			resp, err := http.Get(fmt.Sprintf("https://cdn.jsdelivr.net/gh/gaowanliang/DownloadBot/i18n/active.%s.json", locLanguage))
			dropErr(err)
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			dropErr(err)
			ioutil.WriteFile(fmt.Sprintf("i18n/active.%s.json", locLanguage), data, 0644)
			log.Printf("language file is up to date")
		}
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

// LocText is translated json key to target language
func LocText(MessageIDs ...string) string {
	res := ""
	for _, MessageID := range MessageIDs {
		res += loc.MustLocalize(&i18n.LocalizeConfig{MessageID: MessageID})
	}
	return res
}

func dropErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func pageDownload(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return ""
	}
	//函数结束后关闭相关链接
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error", err)
		return ""
	}
	return string(body)
}

func GetFileModTime(path string) int64 {
	f, err := os.Open(path)
	dropErr(err)
	defer f.Close()

	fi, err := f.Stat()
	dropErr(err)

	return fi.ModTime().Unix()
}
