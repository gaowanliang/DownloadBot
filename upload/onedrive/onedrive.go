package onedrive

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"v2/api/restore/upload"
	"v2/fileutil"
	httpLocal "v2/graph/net/http"
)

var bearerToken = ""
var userID = ""

func onedriveUpload(infoPath string, filePath string, threads int) {
	userID, bearerToken = httpLocal.GetMyIDAndBearer(infoPath)
	restoreOption := "orig"

	//Initialize the upload restore service
	restoreSrvc := upload.GetRestoreService(http.DefaultClient)

	//Get the list of files that needs to be restore with the actual backed up path. 获取需要使用实际备份路径还原的文件列表。
	fileInfoToUpload, err := fileutil.GetAllUploadItemsFrmSource(filePath)
	if err != nil {
		log.Fatalf("Failed to Load Files from source :%v", err)
	}

	//Call restore process based on alternate or original location 基于备用或原始位置调用还原过程
	if restoreOption == "alt" {
		restoreToAltLoc(restoreSrvc, fileInfoToUpload)
	} else {
		restore(restoreSrvc, fileInfoToUpload, threads)
	}
}
func changeBlockSize(MB int) {
	fileutil.SetDefaultChunkSize(MB)
}

//Restore to original location
func restore(restoreSrvc *upload.RestoreService, filesToRestore map[string]fileutil.FileInfo, threads int) {
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
			_, err := restoreSrvc.SimpleUploadToOriginalLoc(userID, bearerToken, "rename", filePath, fileInfo)
			if err != nil {
				log.Panicf("Failed to Restore :%v", err)
			}
			//printResp(resp)
		}(filePath, fileInfo)
	}
	wg.Wait()
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
func restoreToAltLoc(restoreSrvc *upload.RestoreService, filesToRestore map[string]fileutil.FileInfo) {
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
			_, err := restoreSrvc.SimpleUploadToAlternateLoc(userID, bearerToken, "rename", rootFilePath, fileItem)
			if err != nil {
				log.Panicf("Failed to Restore :%v", err)
			}
		}()
		wg.Wait()
		// fmt.Println(respStr)
	}
}
