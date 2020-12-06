# DownloadBot


(目前)🤖 一個控制你的Aria2伺服器的Telegram Bot。

## 實現

#### 下載方式
- [x] Aria2 控制
- [ ] [SimpleTorrent](https://github.com/boypt/simple-torrent) 控制
- [ ] qbittorrent 控制
- [ ] 多下載伺服器同時控制

#### 機器人協定支援
- [x] Telegram Bot
- [ ] 釘釘機器人

#### 功能
- [x] 控制伺服器檔
  - [x] 刪除檔
- [x] 下載檔案
  - [x] 下載 HTTP/FTP 連結
  - [x] 下載 Magnet 連結
  - [x] 下載 BT 文件內的文件
  - [x] 智慧 Torrent/Magnet 下載方式
    - [x] 只下載最大的文件
    - [x] 不下載小於指定大小的檔
  - [ ] 自我調整環境存儲空間的 Torrent/Magnet 下載
    - [ ] 不下載超過存儲空間的檔
    - [ ] 根據存儲空間分塊多次下載 Torrent/Magnet 內的檔
- [ ] 上傳文件
  - [ ] 下載完成後，向 OneDrive 上傳檔
  - [ ] 下載完成後，向 Google Drive 上傳檔
  - [ ] 下載完成後，向 Mega 上傳檔
  - [ ] 下載完成後，向 天翼網盤 上傳文件
- [x] 附加其他功能
  - [x] 多語言支援
    - [x] 簡體中文
    - [x] 英語
    - [x] 繁體中文
    - [ ] 日語
  - [ ] 無人值守的BT站下載
    - [ ] Nyaa
    - [ ] ThePirateBay
  - [ ] 其他功能
    - [ ] 通過演員ID獲取在DMM中使用的所有CID
    - [ ] 查詢 "ikoa "中的影片參數(利用mahuateng)
    - [ ] 通過javlibary演員網址獲得所有演員的編號。
    - [ ] 查詢dmm cid資訊、預覽影片、預覽圖片。
    - [ ] 在sukebei中按關鍵字搜索。
    - [ ] 根據關鍵字在dmm中搜索，最多30項。
    - [ ] 輸入dmm連結，列出所有專案。
    - [ ] 搜索當前dmm熱門和最新電影，限制30條(測試版)

## 目前特點
1. 完全基於觸摸，更容易使用，使用這個機器人基本不需要命令。
2. 即時通知，使用Aria2的Websocket協議進行通信。
3. 更好的設定檔支持。

## 開始

1. 通過[@BotFather](https://telegram.me/botfather)創建您自己的bot並使用。
2. （可選）您所在地區/國家的Telegram被封鎖？一定要有一個 **HTTP** proxy啟動並運行，您可以設置您的系統環境變數`HTTPS_PROXY`為代理位址來進行代理。
3. 下載本程式
4. 在想要執行本程式的根目錄配置`config.json`
5. 運行可執行檔

### 設定檔示例

```json
{
    "aria2-server": "ws://127.0.0.1:5800/jsonrpc",
    "aria2-key": "xxxxxxxx",
    "bot-key": "123456789:xxxxxxxxx",
    "user-id": "123456789",
    "max-index": 10,
    "sign":"Main Aria2",
    "language":"zh-CN",
    "downloadFolder":"C:/aria2/Aria2Data"
}
```
#### 各項對應解釋
* aria2-server：aria2伺服器位址，默認使用websocket連接。如果要使用websocket連接aria2，請務必設置`aria2.conf`內的`enable-rpc=true`。如果不是必須，請儘量設置本地的aria2位址，以便於最大化的使用本程式
* aria2-key：`aria2.conf`中`rpc-secret`的值
* bot-key：Telegram Bot的標識
* user-id：管理員的ID
* max-index：下載資訊最大顯示數量，建議10條（以後會改進）
* sign：此機器人的標識，如果需要多個伺服器連接同一個機器人，通過這一項可以確定具體是哪一台伺服器
* language：機器人輸出的語言
* downloadFolder：Aria2下載檔案保存的位址

#### 目前支援的語言及語言標籤
| 語言     | 標籤  |
|----------|-------|
| 英語     | en    |
| 簡體中文 | zh-CN |
| 繁體中文 | zh-TW |

當您在`config.json`中填寫上面語言的標籤的時候，程式會自動下載語言包

#### 關於user-id
如果您不知道您的 `user-id` ，可以將此項留空，在運行這個機器人後輸入`/myid`，此機器人就會返回您的`user-id`.


