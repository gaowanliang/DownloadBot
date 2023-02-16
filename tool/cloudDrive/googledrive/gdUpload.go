package googledrive

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/googleapi"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3" //v2跟v3不合用要注意
)

type Certificate struct {
	Drive        string      `json:"Drive"`
	RefreshToken string      `json:"RefreshToken"`
	ThreadNum    int         `json:"ThreadNum"`
	BlockSize    int         `json:"BlockSize"`
	MainLand     bool        `json:"MainLand"`
	TimeOut      int         `json:"TimeOut"`
	Other        interface{} `json:"Other"`
}

// refs https://developers.google.com/drive/v3/web/quickstart/go
var chunkSize = 10 * 1024 * 1024

func changeChunkSize(block int) {
	chunkSize = block * 1024 * 1024
}

var wg sync.WaitGroup
var threads = 3
var pool = make(chan struct{}, threads)

func changeThread(thread int) {
	threads = thread
	pool = make(chan struct{}, threads)
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(file string, ctx context.Context, config *oauth2.Config) *http.Client {
	//cacheFile := tokenCacheFile() //取得token的存放位置
	//if err != nil {
	//	log.Fatalf("Unable to get path to cached credential file. %v", err)
	//}
	//log.Println(cacheFile)
	tok := &oauth2.Token{} //呼叫tokenFromFile取得token檔
	/*if err != nil {                      //注意這邊的err不一定是沒有token，也有可能是token的Decode錯誤
		tok = getTokenFromWeb(config) //呼叫getTokenFromWeb重新產生一個網址要求複製token並貼上
		saveToken(cacheFile, tok)     //取得新的token就存檔到cacheFile的路徑
	}*/
	return config.Client(ctx, tok) //成功取得token並return
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("在浏览器中进入以下链接，然后输入授权码: \n%v\n", authURL)
	//Go to the following link in your browser then type the authorization code:
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
		//
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func GetURL() string {
	_, config := gdInit()
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() string {
	//此為另一範例，會在當前user資料夾裡建立資料夾
	//usr, err := user.Current()
	//if err != nil {
	//	return "", err
	//}
	//tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	//os.MkdirAll(tokenCacheDir, 0700)
	//return filepath.Join(tokenCacheDir, url.QueryEscape("drive-go-quickstart.json")), err

	//這邊定義token存放的路徑及資料夾，並且建立資料夾
	tokenCacheDir := filepath.Join("./info/", "googleDrive")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, "drive-go-quickstart.json")
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(c *oauth2.Config, file string) (*oauth2.Token, error) {
	//搜尋路徑，有token檔就開啟並Decode，沒有就回傳nil
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	ts := &Certificate{}
	err = json.NewDecoder(f).Decode(ts) //Decode錯誤也會回傳err
	oauth2Token := ts.Other.(map[string]interface{})
	expiry, _ := time.Parse(time.RFC3339, oauth2Token["expiry"].(string))
	t := &oauth2.Token{
		AccessToken:  oauth2Token["access_token"].(string),
		TokenType:    oauth2Token["token_type"].(string),
		RefreshToken: oauth2Token["refresh_token"].(string),
		Expiry:       expiry,
	}
	defer f.Close()

	// log.Printf("%+v\n", t)
	updatedToken, err := c.TokenSource(context.TODO(), t).Token()
	// log.Printf("%+v\n", updatedToken)
	f, err = os.OpenFile(file, os.O_WRONLY|os.O_TRUNC, 0666)
	defer f.Close()
	// log.Printf("%+v\n", t)
	// log.Printf("%+v\n", updatedToken)
	data := Certificate{
		Drive:        "GoogleDrive",
		RefreshToken: ts.RefreshToken,
		ThreadNum:    ts.ThreadNum,
		BlockSize:    ts.BlockSize,
		MainLand:     false,
		TimeOut:      ts.TimeOut,
		Other:        updatedToken,
	}
	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	changeThread(ts.ThreadNum)
	changeChunkSize(ts.BlockSize)
	return updatedToken, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	//fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	data := Certificate{
		Drive:        "GoogleDrive",
		RefreshToken: token.RefreshToken,
		ThreadNum:    3,
		BlockSize:    10,
		MainLand:     false,
		TimeOut:      60,
		Other:        token,
	}
	json.NewEncoder(f).Encode(data)
}

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

func UploadAllFile(pathname string, folderIDList []string, srv *drive.Service, startTime int64, username string, sendMsg func() func(text string), locText func(text string) string) (func(string), error) {
	filePath := path.Base(pathname)
	temp := sendMsg()
	tip := "`" + filePath + "`" + locText("startUploadGoogleDrive")
	temp(tip)

	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		_, fullName := filepath.Split(pathname)
		var tempFolderIDList []string
		if folderIDList == nil {
			tempFolderIDList = folderIDList
		}
		f, err := os.Open(fullName)
		if err != nil {
			log.Fatalf("error opening %q: %v", fullName, err)
		}
		defer f.Close()
		//log.Println("file", fi.Name(), fullName)
		//上傳檔案，create要給定檔案名稱，要傳進資料夾就加上Parents參數給定folderID的array，media傳入我們要上傳的檔案，最後Do

		uploadFile(srv, pathname, fullName, tempFolderIDList, f, sendMsg, locText, username)
		// sendMsg(fmt.Sprintf(locText("googleDriveUploadTip"), username, fullName, time.Now().Unix()-startTime))
	}
	_, foldName := filepath.Split(pathname)
	//log.Println(foldName)
	var createFolder *drive.File
	if folderIDList == nil {
		createFolder, err = srv.Files.Create(&drive.File{Name: foldName, MimeType: "application/vnd.google-apps.folder"}).Do()
	} else {
		createFolder, err = srv.Files.Create(&drive.File{Name: foldName, Parents: folderIDList, MimeType: "application/vnd.google-apps.folder"}).Do()
	}
	if err != nil {
		log.Panicf("Unable to create folder: %v", err)
	}

	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			var tempFolderIDList []string
			if folderIDList == nil {
				tempFolderIDList = folderIDList
			}
			tempFolderIDList = append(tempFolderIDList, createFolder.Id)
			_, err := UploadAllFile(fullDir, tempFolderIDList, srv, startTime, username, sendMsg, locText)
			if err != nil {
				log.Panic("read dir fail:", err)
				return nil, err
			}
			//log.Println("folder", fi.Name(), fullDir)
			// sendMsg(fmt.Sprintf(locText("googleDriveUploadTip"), username, fullDir, time.Now().Unix()-startTime))
		} else {
			fullName := pathname + "/" + fi.Name()
			var tempFolderIDList []string
			if folderIDList == nil {
				tempFolderIDList = folderIDList
			}
			tempFolderIDList = append(tempFolderIDList, createFolder.Id)
			f, err := os.Open(fullName)
			if err != nil {
				log.Fatalf("error opening %q: %v", fullName, err)
			}
			defer f.Close()
			//log.Println("file", fi.Name(), fullName)
			//上傳檔案，create要給定檔案名稱，要傳進資料夾就加上Parents參數給定folderID的array，media傳入我們要上傳的檔案，最後Do
			// _, err = srv.Files.Create(&drive.File{Name: fi.Name(), Parents: tempFolderIDList}).Media(f, googleapi.ChunkSize(chunkSize)).Do()
			uploadFile(srv, fullName, fi.Name(), tempFolderIDList, f, sendMsg, locText, username)
			if err != nil {
				log.Panicf("Unable to create file: %v", err)
			}
			// sendMsg(fmt.Sprintf(locText("googleDriveUploadTip"), username, fullName, time.Now().Unix()-startTime))
			//log.Printf("file: %+v", driveFile)
		}
	}
	wg.Wait()
	return temp, nil
}
func FileSizeFormat(bytes int64, forceBytes bool) string {
	if forceBytes {
		return fmt.Sprintf("%v B", bytes)
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}

	var i int
	value := float64(bytes)

	for value > 1000 {
		value /= 1000
		i++
	}
	return fmt.Sprintf("%.1f %s", value, units[i])
}

func MeasureTransferRate() func(int64) string {
	start := time.Now()

	return func(bytes int64) string {
		seconds := int64(time.Now().Sub(start).Seconds())
		if seconds < 1 {
			return fmt.Sprintf("%s/s", FileSizeFormat(bytes, false))
		}
		bps := bytes / seconds
		return fmt.Sprintf("%s/s", FileSizeFormat(bps, false))
	}
}

func uploadFile(srv *drive.Service, filePath string, filename string, tempFolderIDList []string, f *os.File, sendMsg func() func(text string), locText func(text string) string, username string) {
	wg.Add(1)
	pool <- struct{}{}
	startTime := time.Now().Unix()
	go func() {
		defer wg.Done()
		defer func() {
			<-pool
		}()
		temp := sendMsg()

		tip := "`" + filePath + "`" + locText("startUploadGoogleDrive")
		temp(tip)

		fi, err := os.Stat(filePath)
		if err != nil {
			fi, err = os.Stat(filename)
			if err != nil {
				log.Panicln(err)
			}
		}
		getRate := MeasureTransferRate()
		size := fi.Size()
		// log.Println(byte2Readable(float64(size)))

		showProgress := func(current, total int64) {
			temp(fmt.Sprintf(locText("googleDriveUploadTip1"), username, filePath, byte2Readable(float64(size)), byte2Readable(float64(current)), int(math.Ceil(float64(current)/float64(chunkSize))), int(math.Ceil(float64(size)/float64(chunkSize))), getRate(current), time.Now().Unix()-startTime))
		}

		_, err = srv.Files.Create(&drive.File{Name: filename, Parents: tempFolderIDList}).Media(f, googleapi.ChunkSize(chunkSize)).ProgressUpdater(showProgress).Do()
		if err != nil {
			log.Panicln(err)
		}

		defer temp("close")
	}()

}

/*func main() {
	UploadAllFile("test")
}*/

func gdInit() (ctx context.Context, config *oauth2.Config) {
	ctx = context.Background()
	b := `{
	"installed": {
		"client_id": "413585947449-217qe10vfis0o4od73jjs103p3teih7k.apps.googleusercontent.com",
		"project_id": "quickstart-1609508510551",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_secret": "IJdLG98l4KuTzcNE2KA2qfZ5",
		"redirect_uris": ["urn:ietf:wg:oauth:2.0:oob", "http://localhost"]
	}
}`
	x := (*[2]uintptr)(unsafe.Pointer(&b))
	h := [3]uintptr{x[0], x[1], x[1]}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/drive-go-quickstart.json
	config, err := google.ConfigFromJSON(*(*[]byte)(unsafe.Pointer(&h)), drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return
}
func CreateNewInfo(code string) string {
	ctx, config := gdInit()
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	client := config.Client(ctx, tok)
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	userInfo, err := srv.About.Get().Fields("user").Do()
	if err != nil {
		log.Fatalf("An error occurred: %v\n", err)
	}

	//log.Printf("%+v", userInfo.User.EmailAddress)
	saveToken("./info/googleDrive/"+userInfo.User.EmailAddress+".json", tok) //取得新的token就存檔到cacheFile的路徑
	return "./info/googleDrive/" + userInfo.User.EmailAddress + ".json"
}

func Upload(infoPath string, filePath string, sendMsg func() func(text string), locText func(text string) string) {
	ctx, config := gdInit()
	infoPath, _ = filepath.Abs(infoPath)
	tok, _ := tokenFromFile(config, infoPath)
	client := config.Client(ctx, tok)
	srv, err := drive.New(client)
	username := strings.ReplaceAll(filepath.Base(infoPath), ".json", "")
	// restoreOption := "orig"
	oldDir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	err = os.Chdir(filepath.Dir(filePath))
	if err != nil {
		log.Panic(err)
	}

	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)

	}
	temp, _ := UploadAllFile(filePath, nil, srv, time.Now().Unix(), username, sendMsg, locText)
	temp(locText("uploadGoogleDriveComplete"))
	err = os.Chdir(oldDir)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	//createNewInfo()
	//這邊是用參數的方式開啟，輸入要上傳的檔案名稱
	//if len(os.Args) != 2 {
	//	fmt.Fprintln(os.Stderr, "Usage: drive filename (to upload a file)")
	//	return
	//}
	//filename := os.Args[1]

	//filename := "test/02 数码管.pdf"

	//os.Rename()
	//return
	/*
		f, err := os.Open(filename)
		if err != nil {
			log.Fatalf("error opening %q: %v", filename, err)
		}
		defer f.Close()

		//建立資料夾，給定Name跟MimeType，最後Do會回傳一些資料，型態是dict，資料內容可參考API
		createFolder, err := srv.Files.Create(&drive.File{Name: "testFolder", MimeType: "application/vnd.google-apps.folder"}).Do()
		if err != nil {
			log.Fatalf("Unable to create folder: %v", err)
		}

		//建立array存放資料夾ID，上面建立資料夾Do的回傳內容就包括資料夾ID
		var folderIDList []string
		folderIDList = append(folderIDList, createFolder.Id)
		//上傳檔案，create要給定檔案名稱，要傳進資料夾就加上Parents參數給定folderID的array，media傳入我們要上傳的檔案，最後Do
		driveFile, err := srv.Files.Create(&drive.File{Name: filename, Parents: folderIDList}).Media(f).Do()
		if err != nil {
			log.Fatalf("Unable to create file: %v", err)
		}

		log.Printf("file: %+v", driveFile)*/

	//UploadAllFile("test", nil, srv)
	log.Println("done")
}
