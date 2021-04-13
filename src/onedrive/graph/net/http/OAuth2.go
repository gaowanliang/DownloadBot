package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

type Token struct {
	RefreshToken string `json:"refresh_token"`
	ThreadNum    int    `json:"threadNum"`
	BlockSize    int    `json:"blockSize"`
	TimeOut      int    `json:"timeout"`
}

func NewPassCheck(oauth2URL string) string {
	Bearer := getAccessToken(oauth2URL)

	url := "https://graph.microsoft.com/v1.0/me/"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+Bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var mail string
	if strings.Contains(string(body), "userPrincipalName") {
		mail, err = jsonparser.GetString(body, "userPrincipalName")
	} else {
		mail, err = jsonparser.GetString(body, "mail")
	}

	if err != nil {
		log.Println(string(body))
		log.Panicln(err)
	}
	err = os.Rename("./info/onedrive/amazing.json", "./info/onedrive/"+mail+".json")
	if err != nil {
		log.Panic(err)
	}

	return "./info/onedrive/" + mail + ".json"
}

// GetMyIDAndBearer is get microsoft ID and access token
func GetMyIDAndBearer(infoPath string) (string, string) {
	MyID := ""
	Bearer := ""
	_, err := os.Stat(infoPath)
	Bearer = refreshAccessToken(infoPath)
	url := "https://graph.microsoft.com/v1.0/me/"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+Bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	MyID, err = jsonparser.GetString(body, "id")
	if err != nil {
		log.Println(string(body))
		log.Panicln(err)
	}

	//os.Rename("info.json", mail+".json")
	// log.Println(MyID)

	return MyID, Bearer
}

func getAccessToken(oauth2URL string) string {
	var re *regexp.Regexp

	if len(oauth2URL) > 100 {
		re = regexp.MustCompile(`(?m)code=(.*?)&`)
	} else {
		re = regexp.MustCompile(`(?m)code=(.*?)$`)
	}

	var str = oauth2URL
	/*log.Printf(
		`%s https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=ad5e65fd-856d-4356-aefc-537a9700c137&response_type=code&redirect_uri=http://localhost/onedrive-login&response_mode=query&scope=offline_access%%20User.Read%%20Files.ReadWrite.All`,
		"*请打开下面的网址，登录OneDrive账户后，将跳转后的网址复制后，发送给本Bot*\n注意：本程序不会涉及您的隐私信息，请放心使用，后续会提供更换上传API的方法",
	)*/

	code := re.FindStringSubmatch(str)[1]

	url := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("client_id=ad5e65fd-856d-4356-aefc-537a9700c137&scope=offline_access%%20User.Read%%20Files.ReadWrite.All&code=%s&redirect_uri=http://localhost/onedrive-login&grant_type=authorization_code", code)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "https://login.microsoftonline.com")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	accessToken, err := jsonparser.GetString(body, "access_token")
	if err != nil {
		log.Println(string(body))
		log.Panicln(err)
	}
	//log.Println(accessToken)
	refreshToken, err := jsonparser.GetString(body, "refresh_token")
	if err != nil {
		log.Println(string(body))
		log.Panicln(err)
	}
	//log.Println(refreshToken)

	info := Token{
		RefreshToken: refreshToken,
		ThreadNum:    3,
		BlockSize:    10,
		TimeOut:      60,
	}
	// 创建文件
	filePtr, err := os.Create("./info/onedrive/amazing.json")
	if err != nil {
		log.Panicln(err.Error())
		return ""
	}
	defer filePtr.Close()
	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(info)
	if err != nil {
		log.Panicln(err.Error())
	}
	return accessToken
}

func refreshAccessToken(path string) string {
	filePtr, err := os.Open(path)
	if err != nil {
		log.Panicln(err)
		return ""
	}
	defer filePtr.Close()
	var info Token
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	if err != nil {
		log.Panicln(err.Error())
	}
	url := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("client_id=ad5e65fd-856d-4356-aefc-537a9700c137&scope=offline_access%%20User.Read%%20Files.ReadWrite.All&refresh_token=%s&redirect_uri=http://localhost/onedrive-login&grant_type=refresh_token", info.RefreshToken)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "https://login.microsoftonline.com")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	accessToken, err := jsonparser.GetString(body, "access_token")
	if err != nil {
		log.Panicln(err)
	}
	//log.Println(accessToken)
	refreshToken, err := jsonparser.GetString(body, "refresh_token")
	if err != nil {
		log.Panicln(err)
	}
	// log.Println(refreshToken)

	info = Token{
		RefreshToken: refreshToken,
		ThreadNum:    info.ThreadNum,
		BlockSize:    info.BlockSize,
		TimeOut:      info.TimeOut,
	}
	// 创建文件
	filePtr, err = os.Create(path)
	if err != nil {
		log.Panicln(err.Error())
		return ""
	}
	defer filePtr.Close()
	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(info)
	if err != nil {
		log.Panicln(err.Error())
	}
	return accessToken
}
