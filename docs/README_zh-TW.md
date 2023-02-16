# DownloadBot

[![Go Version](https://img.shields.io/github/go-mod/go-version/gaowanliang/DownloadBot.svg?style=flat-square&label=Go&color=00ADD8)](https://github.com/gaowanliang/DownloadBot/blob/master/go.mod)
[![Release Version](https://img.shields.io/github/v/release/gaowanliang/DownloadBot.svg?style=flat-square&label=Release&color=1784ff)](https://github.com/gaowanliang/DownloadBot/releases/latest)
[![GitHub license](https://img.shields.io/github/license/gaowanliang/DownloadBot.svg?style=flat-square&label=License&color=2ecc71)](https://github.com/gaowanliang/DownloadBot/blob/master/LICENSE)
[![GitHub Star](https://img.shields.io/github/stars/gaowanliang/DownloadBot.svg?style=flat-square&label=Star&color=f39c12)](https://github.com/gaowanliang/DownloadBot/stargazers)
[![GitHub Fork](https://img.shields.io/github/forks/gaowanliang/DownloadBot.svg?style=flat-square&label=Fork&color=8e44ad)](https://github.com/gaowanliang/DownloadBot/network/members)

(目前) 🤖 一個分散式跨平臺的，可以控制你的Aria2伺服器、控制伺服器檔，同時可以上傳到OneDrive/Google Drive的Telegram Bot。

## 意義
> 以下僅為本程式完成後的設想，目前描述的功能並沒有完全實現，實現的詳情請參考下面的功能實現

這個項目主要就是利用吃灰小盤vps進行離線下載，對於大bt檔進行根據硬碟大小分段下載，每次都下載一部分，然後上傳網盤，刪除再下載其他部分，直到下載完所有檔。

同時，通過機器人協議通信，方便在無法進行內網穿透的機器上進行使用，而且簡化了平時使用下載程式的操作，提高了便利性。對於連結，直接向Bot發送消息就可以直接識別並下載，可以真正刪除下載檔案夾裡的檔，是AriaNG等web面板無法做到的，作為管理下載的工具，及時通知下載完成都是非常的方便的。可以移動檔，對於通過rclone掛載硬碟的使用者可以直接通過本程式進行複製粘貼等操作，無需打開ssh連接VPS進行`cp`操作，也非常的方便。

## 實現


<text style="color:red;">**注意：本項目仍處於測試階段，提交的Release僅供測試，現在下載後並不保證您的穩定使用，也不能保證下面所勾選的內容已經被實現。當真正可以正常使用的時候，我會提交 V1.0 版本（V1.0 版本不會實現下面全部功能，但是已經可以正常穩定的使用）**</text>

**只有已勾選的內容為目前實現的內容**

#### 下載方式

- [x] Aria2 控制
  - [x] 持久化監控
  - [x] 斷線重連
- [ ] 多下載伺服器同時控制
  - [x] 使用GRPC實現多伺服器下載資訊通知
- [ ] [SimpleTorrent](https://github.com/boypt/simple-torrent) 控制
- [ ] qbittorrent 控制


#### 機器人協定支援

- [x] Telegram Bot
  - [x] 支援多用戶使用
  - [ ] 支援群組內使用
- [ ] 騰訊QQ（使用普通QQ用戶來進行交互）
- [ ] 釘釘機器人


#### 功能

- [x] 控制伺服器檔
    - [x] 刪除檔
    - [x] 移動文件
    - [ ] 壓縮檔
    - [ ] 解壓文件
- [x] 下載檔案
    - [x] 下載 HTTP/FTP 連結
    - [x] 下載 Magnet 連結
    - [x] 下載 BitTorrent 文件內的文件
    - [x] 智慧 BitTorrent/Magnet 下載方式
        - [x] 只選擇下載最大的文件
        - [x] 根據檔大小智慧選擇檔，不選擇小文件
    - [ ] 下載OneDrive/SharePoint 分享連結中的檔（[Python demo](https://github.com/gaowanliang/OneDriveShareLinkPushAria2)）
      - [ ] xxx-my.sharepoint.com 下載連結的下載
        - [ ] 無下載密碼的多檔推送
        - [ ] 有下載密碼的多檔推送
        - [ ] 嵌套資料夾的文件推送
        - [ ] 任意選擇檔推送
      - [ ] xxx.sharepoint.com 下載連結的下載
      - [ ] xxx-my.sharepoint.cn 下載連結的下載(理論上支持)
    - [ ] 自我調整環境存儲空間的 BitTorrent/Magnet 下載
        - [ ] 不下載超過存儲空間的檔
        - [ ] 根據存儲空間分塊多次下載 BitTorrent/Magnet 內的檔
    - [ ] 根據具體的硬碟容量制定合適的下載列表
    - [ ] 無感覺化的做種功能
      - [ ] 每次下載BitTorrent/Magnet檔後，保留最後一次下載的檔進行做種，直到下一次下載開始。
      - [ ] 可設置每次下載結束後強制做種一段時間
- [x] 上傳文件
    - [x] 下載完成後，向 OneDrive 上傳檔
      - [ ] 中斷點續傳
      - [ ] 支持世紀互聯
    - [x] 下載完成後，向 Google Drive 上傳檔
      - [x] 自訂上傳分塊大小
      - [x] 自訂上傳執行緒數
      - [x] 自訂超時時間
    - [ ] 下載完成後，向 Mega 上傳檔
    - [ ] 下載完成後，向 天翼網盤 上傳文件
    - [ ] (當使用Telegram進行通信時)下載完成後，向 Telegram 上傳檔
      - [ ] 當檔超過2GB時，分塊壓縮後再進行上傳
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
        - [x] 檔樹輸出系統
            - [x] 對於簡單資料夾的文件樹輸出
            - [x] 對於複雜資料夾結構使用多消息輸出
        - [ ] 通過演員ID獲取在DMM中使用的所有CID
        - [ ] 查詢 "ikoa"中的影片參數(利用mahuateng)
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

1. 通過 [@BotFather](https://telegram.me/botfather) 創建您自己的bot並使用。
2. （可選）您所在地區/國家的Telegram被封鎖？一定要有一個 **HTTP** proxy啟動並運行，您可以設置您的系統環境變數`HTTPS_PROXY`為代理位址來進行代理。
3. 下載本程式
4. 在想要執行本程式的根目錄配置`config.json`
5. 運行可執行檔`./DownloadBot`或`./DownloadBot.exe`。

## 教程

[簡體中文問答文檔](docs/Q&A_zh-CN.md)

## 使用截圖

<div align="center">
<img src="./img/1.jpg" height="300px" alt="">  <img src="./img/2.jpg" height="300px" alt="" >  
</div>
<br>

<div align="center">
<img src="./img/3.jpg" height="300px" alt="">  <img src="./img/4.jpg" height="300px" alt="" >  </div>


## 設定檔示例

```json
{
  "input": {
    "aria2": {
      "aria2-server": "ws://127.0.0.1:6800/jsonrpc",
      "aria2-key": "123456"
    }
  },
  "output": {
    "telegram": {
      "bot-key": "",
      "user-id": ""
    }
  },
  "max-index": 10,
  "sign": "Main Aria2",
  "language": "en",
  "downloadFolder": "/root/download",
  "moveFolder": "/root/upload",
  "server": {
    "isServer": true,
    "isMasterServer": true,
    "serverHost": "127.0.0.1",
    "serverPort": 23369
  },
  "log": {
    "logPath": "",
    "errPath": "",
    "level": "info"
  }
}
```

#### 各項對應解釋
* input: 輸入端，目前僅支持aria2
  * aria2-server：aria2伺服器位址，默認使用websocket連接。如果要使用websocket連接aria2，請務必設置`aria2.conf`內的`enable-rpc=true`
      。如果不是必須，請儘量設置本地的aria2位址，以便於最大化的使用本程式
  * aria2-key：`aria2.conf`中`rpc-secret`的值
* output: 輸出端，目前僅支援telegram

* bot-key：Telegram Bot的標識，通過 [@BotFather](https://telegram.me/botfather)進行獲取。
* user-id：管理員的ID~~，支援設置多用戶為管理員，不同的用戶之間使用半形逗號`,`分割。如您要設置`user-id`為123465789、987654321和963852741的使用者為管理員，您需要這樣設置：~~
  ```jsonc
  {
    //···
    "user-id": "123456789",
    //···
  }
  ```
* max-index：下載資訊最大顯示數量，建議10條（以後會改進）
* sign：此機器人的標識，如果需要多個伺服器連接同一個機器人，通過這一項可以確定具體是哪一台伺服器
* language：機器人輸出的語言
* downloadFolder：Aria2下載檔案保存的位址。如果不使用，請輸入`""`
* moveFolder： 要將下載檔案夾的文件移動到的資料夾。如果不使用，請輸入`""`
* server：伺服器配置
  * isServer：是否開啟伺服器（false為用戶端）
  * isMasterServer：是否為主要伺服器
  * serverHost：如果是用戶端，此項需要填寫伺服器位址，如果是主要伺服器，此項為本機地址
  * serverPort：如果是用戶端，此項需要填寫伺服器埠，如果是主要伺服器，此項為提供給用戶端的埠
* log：日誌配置
  * logPath：日誌檔保存位址，如果不使用，請輸入`""`（目前不支援）
  * errPath：錯誤日誌檔保存位址，如果不使用，請輸入`""`（目前不支援）
  * level：日誌等級，可選項為`debug`、`info`、`warn`、`error`、`fatal`，默認為`info`

#### 目前支援的語言及語言標籤

| 語言     | 標籤  |
|----------|-------|
| 英語     | en    |
| 簡體中文 | zh-CN |
| 繁體中文 | zh-TW |

當您在`config.json`中填寫上面語言的標籤的時候，程式會自動下載語言包

#### 關於user-id

如果您不知道您的 `user-id` ，可以將此項留空，在運行這個機器人後輸入`/myid`，此機器人就會返回您的`user-id`.


#### 捐贈

如果您覺得本程式對您有幫助，您可以通過捐贈的方式支持我，謝謝！

https://ko-fi.com/gaowanliang

