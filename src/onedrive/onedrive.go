package onedrive

import (
	"fmt"
	"log"
	"net/http"
	"onedrive/api/restore/upload"
	"onedrive/fileutil"
	httpLocal "onedrive/graph/net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var bearerToken = ""
var userID = ""

func ApplyForNewPass(url string) string {
	return httpLocal.NewPassCheck(url)
}

func Upload(infoPath string, filePath string, threads int, sendMsg func() func(text string), locText func(text string) string) {
	userID, bearerToken = httpLocal.GetMyIDAndBearer(infoPath)
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
	//Initialize the upload restore service
	restoreSrvc := upload.GetRestoreService(http.DefaultClient)

	//Get the list of files that needs to be restore with the actual backed up path. 获取需要使用实际备份路径还原的文件列表。
	fileInfoToUpload, err := fileutil.GetAllUploadItemsFrmSource(filePath)
	if err != nil {
		log.Fatalf("Failed to Load Files from source :%v", err)
	}

	//Call restore process based on alternate or original location 基于备用或原始位置调用还原过程
	/*if restoreOption == "alt" {
		restoreToAltLoc(restoreSrvc, fileInfoToUpload)
	} else {
		restore(restoreSrvc, fileInfoToUpload, threads)
	}*/
	restore(restoreSrvc, fileInfoToUpload, threads, sendMsg, locText, username)
	err = os.Chdir(oldDir)
	if err != nil {
		log.Panic(err)
	}
}
func changeBlockSize(MB int) {
	fileutil.SetDefaultChunkSize(MB)
}

//Restore to original location
func restore(restoreSrvc *upload.RestoreService, filesToRestore map[string]fileutil.FileInfo, threads int, sendMsg func() func(text string), locText func(text string) string, username string) {
	var wg sync.WaitGroup
	pool := make(chan struct{}, threads)
	for filePath, fileInfo := range filesToRestore {
		wg.Add(1)
		pool <- struct{}{}
		go func(filePath string, fileInfo fileutil.FileInfo) {
			defer wg.Done()
			defer func() {
				<-pool
			}()
			temp := sendMsg()
			temp("`" + filePath + "`" + locText("startUploadOneDrive"))
			_, err := restoreSrvc.SimpleUploadToOriginalLoc(userID, bearerToken, "rename", filePath, fileInfo, temp, locText, username)
			if err != nil {
				log.Panicf("Failed to Restore :%v", err)
			}
			//printResp(resp)
			temp("close")
		}(filePath, fileInfo)
	}
	wg.Wait()
	temp := sendMsg()
	defer temp(locText("uploadOneDriveComplete"))
}

func printResp(resp interface{}) {
	switch resp.(type) {
	case map[string]interface{}:
		fmt.Printf("\n%+v\n", resp)
		break
	case []map[string]interface{}:
		for _, rs := range resp.([]map[string]interface{}) {
			fmt.Printf("\n%+v\n", rs)
		}
	}
}

//Restore to Alternate location 还原到备用位置
func restoreToAltLoc(restoreSrvc *upload.RestoreService, filesToRestore map[string]fileutil.FileInfo, sendMsg func() func(text string), locText func(text string) string) {
	rootFolder := fileutil.GetAlternateRootFolder()
	var wg sync.WaitGroup
	pool := make(chan struct{}, 10)
	for filePath, fileItem := range filesToRestore {
		rootFilePath := fmt.Sprintf("%s/%s", rootFolder, filePath)
		wg.Add(1)
		pool <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() {
				<-pool
			}()
			temp := sendMsg()
			temp(filePath + "开始上传至OneDrive")
			us:=""
			_, err := restoreSrvc.SimpleUploadToAlternateLoc(userID, bearerToken, "rename", rootFilePath, fileItem, temp, locText,us)
			if err != nil {
				log.Panicf("Failed to Restore :%v", err)
			}
		}()
		wg.Wait()
		// fmt.Println(respStr)
	}
}
