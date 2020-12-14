package http

import (
	"bufio"
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

type token struct {
	RefreshToken string `json:"refresh_token"`
}

// GetMyIDAndBearer is get microsoft ID and access token
func GetMyIDAndBearer(infoPath string) (string, string) {
	MyID := ""
	Bearer := ""
	_, err := os.Stat(infoPath)
	if err != nil {
		Bearer = getAccessToken()
	} else {
		Bearer = refreshAccessToken(infoPath)
	}
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
	//fmt.Println(string(body))
	MyID, err = jsonparser.GetString(body, "id")
	if err != nil {
		log.Panicln(err)
	}
	// log.Println(MyID)

	return MyID, Bearer
}

func getAccessToken() string {
	var re = regexp.MustCompile(`(?m)code=(.*?)&`)
	var str = ""
	log.Printf(`%s https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=9d7f43cd-55c8-448c-ba5a-7fafc01f2291&response_type=code&redirect_uri=http://localhost/onedrive-login&response_mode=query&scope=offline_access%%20User.Read%%20Files.ReadWrite.All`, "请打开下面的网址，登录OneDrive账户后，将跳转后的网址粘贴到命令行中")
	inputReader := bufio.NewReader(os.Stdin)
	str, err := inputReader.ReadString('\n')
	if err == nil {
		fmt.Printf("The input was: %s\n", str)
	}

	code := re.FindStringSubmatch(str)[1]
	//fmt.Println(res)
	url := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("client_id=9d7f43cd-55c8-448c-ba5a-7fafc01f2291&scope=offline_access%%20User.Read%%20Files.ReadWrite.All&code=%s&redirect_uri=http://localhost/onedrive-login&grant_type=authorization_code", code)))
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
	//log.Println(refreshToken)

	info := token{
		RefreshToken: refreshToken,
	}
	// 创建文件
	filePtr, err := os.Create("info.json")
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
	var info token
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	if err != nil {
		log.Panicln(err.Error())
	}
	url := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("client_id=9d7f43cd-55c8-448c-ba5a-7fafc01f2291&scope=offline_access%%20User.Read%%20Files.ReadWrite.All&refresh_token=%s&redirect_uri=http://localhost/onedrive-login&grant_type=refresh_token", info.RefreshToken)))
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

	info = token{
		RefreshToken: refreshToken,
	}
	// 创建文件
	filePtr, err = os.Create("info.json")
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
