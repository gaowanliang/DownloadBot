package googledrive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3" //v2跟v3不合用要注意
)

// refs https://developers.google.com/drive/v3/web/quickstart/go

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(file string, ctx context.Context, config *oauth2.Config) *http.Client {
	//cacheFile := tokenCacheFile() //取得token的存放位置
	//if err != nil {
	//	log.Fatalf("Unable to get path to cached credential file. %v", err)
	//}
	//log.Println(cacheFile)
	tok, _ := tokenFromFile(file) //呼叫tokenFromFile取得token檔
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
func tokenFromFile(file string) (*oauth2.Token, error) {
	//搜尋路徑，有token檔就開啟並Decode，沒有就回傳nil
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t) //Decode錯誤也會回傳err
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	//fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetAllFile(pathname string, folderIDList []string, srv *drive.Service, startTime int64, username string, sendMsg func(text string), locText func(text string) string) error {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		log.Panic("read dir fail:", err)
		return err
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
			err := GetAllFile(fullDir, tempFolderIDList, srv, startTime, username, sendMsg, locText)
			if err != nil {
				log.Panic("read dir fail:", err)
				return err
			}
			//log.Println("folder", fi.Name(), fullDir)
			sendMsg(fmt.Sprintf(locText("googleDriveUploadTip"), username, fullDir, time.Now().Unix()-startTime))
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
			_, err = srv.Files.Create(&drive.File{Name: fi.Name(), Parents: tempFolderIDList}).Media(f).Do()
			if err != nil {
				log.Panicf("Unable to create file: %v", err)
			}
			sendMsg(fmt.Sprintf(locText("googleDriveUploadTip"), username, fullName, time.Now().Unix()-startTime))

			//log.Printf("file: %+v", driveFile)
		}
	}
	return nil
}

/*func main() {
	GetAllFile("test")
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
	tok, _ := tokenFromFile(infoPath)
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
	filePath = path.Base(filePath)
	temp := sendMsg()
	temp("`" + filePath + "`" + locText("startUploadGoogleDrive"))

	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)

	}
	_ = GetAllFile(filePath, nil, srv, time.Now().Unix(), username, temp, locText)
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

	//GetAllFile("test", nil, srv)
	log.Println("done")
}
