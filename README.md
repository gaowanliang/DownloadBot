[简体中文](docs/README_zh-CN.md) [繁體中文](docs/README_zh-TW.md)

# DownloadBot

[![Go Version](https://img.shields.io/github/go-mod/go-version/gaowanliang/DownloadBot.svg?style=flat-square&label=Go&color=00ADD8)](https://github.com/gaowanliang/DownloadBot/blob/master/go.mod)
[![Release Version](https://img.shields.io/github/v/release/gaowanliang/DownloadBot.svg?style=flat-square&label=Release&color=1784ff)](https://github.com/gaowanliang/DownloadBot/releases/latest)
[![GitHub license](https://img.shields.io/github/license/gaowanliang/DownloadBot.svg?style=flat-square&label=License&color=2ecc71)](https://github.com/gaowanliang/DownloadBot/blob/master/LICENSE)
[![GitHub Star](https://img.shields.io/github/stars/gaowanliang/DownloadBot.svg?style=flat-square&label=Star&color=f39c12)](https://github.com/gaowanliang/DownloadBot/stargazers)
[![GitHub Fork](https://img.shields.io/github/forks/gaowanliang/DownloadBot.svg?style=flat-square&label=Fork&color=8e44ad)](https://github.com/gaowanliang/DownloadBot/network/members)


(Currently) 🤖 A distributed cross-platform Telegram Bot that can control your Aria2 server, control server files and also upload to OneDrive / Google Drive.

## Project significance
> The following is only a vision of what this program will look like when it is completed, the functions described so far are not fully implemented, please refer to the following [Functions realized](#functions-realized) for details of implementation

This project is mainly to use small hard disk server for offline downloading, for large BitTorrent files to be downloaded in sections according to the size of the hard disk, each time downloading a part, then uploading the network disk, delete and then download the other parts, until all the files are downloaded.

At the same time, communication via the bot protocol facilitates use on machines that cannot intranet penetration, and simplifies the usual use of download programs for added convenience.For links, sending a message directly to the Bot will directly identify and download them. It can actually delete files from the download folder, which is not possible with web panels such as AriaNG, and is very convenient as a tool for managing downloads and notifying timely completion of downloads. You can move files, and for users who mount their hard drives via rclone you can copy and paste directly through this program, without having to open an ssh connection to the VPS for `cp` operations, which is also very convenient.


## Functions realized

<text style="color:red;">**Note: This project is still in beta testing and the Release submitted is for testing purposes only. Downloading it now does not guarantee you stable use, nor does it guarantee that the content ticked below has been implemented. The software is only stable when the submitted version is v1.0 (v1.0 will not implement all of the features below, but it will work properly and stably).**</text>

**Only the checked content is currently implemented**

#### Download method
- [x] Aria2 control
  - [x] Persistent monitoring
  - [x] Automatic reconnection after disconnection
- [ ] Multi download server control at the same time
  - [x] Multi server download information notification using GRPC
- [ ] [SimpleTorrent](https://github.com/boypt/simple-torrent) control
- [ ] qbittorrent control


#### The Bot protocol supports
- [x] Telegram Bot
  - [ ] Support multi-user use
  - [ ] Support group use
- [ ] Tencent QQ (Use regular QQ users to interact)
- [ ] DingTalk Bot


#### Function
- [x] Control server files
  - [x] Delete files
  - [x] Move/Copy files
  - [ ] Compressed files
- [x] Download files
  - [x] Download HTTP/FTP link
  - [x] Download Magnet link
  - [x] Download the files in the BitTorrent file
  - [x] Custom BitTorrent/Magnet download
    - [x] Select only the largest file to download
    - [x] Intelligent file selection based on file size, do not select small files in BitTorrent/Magnet.
  - [ ] Download files from OneDrive/SharePoint share links ([Python script currently used as a demo](https://github.com/gaowanliang/OneDriveShareLinkPushAria2))
    - [ ] xxx-my.sharepoint.com Download of share links
      - [ ] Downloading multiple files without password for shared links
      - [ ] Downloading multiple files with password for shared links
      - [ ] Download of files in nested folders
      - [ ] Download any file of your choice
    - [ ] xxx.sharepoint.com Downloads with share links
    - [ ] xxx-my.sharepoint.cn Download of share links (theoretically supported)
  - [ ] Download BitTorrent/Magnet according to the size of storage space
    - [ ] Do not download files that exceed storage space
    - [ ] Download the files in BitTorrent/Magnet several times according to the storage space
  - [ ] Senseless seeding functions
      - [ ] After each BitTorrent/Magnet file download, keep the last downloaded file for seeding until the next download starts.
      - [ ] Can be set to force seeding for a period of time at the end of each download
- [x] Upload a file
  - [x] Upload the file to OneDrive when the download is complete
    - [ ] Resume from break point
    - [ ] Supports 21vianet (CN) version
  - [x] Upload the file to Google Drive when the download is complete
    - [x] Custom upload chunk size
    - [x] Custom number of upload threads
    - [x] Custom timeout time
  - [ ] Upload the file to Mega when the download is complete
  - [ ] Upload the file to 189Cloud when the download is complete
  - [ ] (When communicating via Telegram) Upload the file to Telegram when the download is complete
    - [ ] When the file exceeds 2GB, it is compressed in chunks before uploading
- [x] Additional features
  - [x] Multilingual support
    - [x] Simplified Chinese
    - [x] English
    - [x] Traditional Chinese
    - [ ] Japanese
  - [ ] No human intervention, fully automatic downloads of BitTorrent site
    - [ ] Nyaa
    - [ ] ThePirateBay
  - [ ] Other functions
    - [x] File tree output system
        - [x] File tree output for simple folders
        - [x] Use multi message output for complex folder structures
    - [ ] Get all CIDs used in DMM via actor ID
    - [ ] Query the movie parameters in "ikoa" (using mahuateng).
    - [ ] Get the numbers of all actors via the javlibary actors' website. 
    - [ ] Query the dmm cid information, preview the movie, preview the picture. 
    - [ ] Search by keyword in sukebei. 
    - [ ] Search in dmm based on keywords, up to 30 items. 
    - [ ] Enter the dmm link to list all items. 
    - [ ] Search for current dmm hits and the latest movies, limited to 30 (beta).

## Current features
1. Fully touch based, more easy to use, no command required to use this bot.
2. Real time notification, it's now using Aria2's Websocket protocol to communicate.
3. Better config file support.


## Setup

1. Create your own bot and get its access token by using [@BotFather](https://telegram.me/botfather)
2. (Optional) Telegram blocked in your region/country? be sure to have a HTTP proxy up and running,and You can set your system environment variable `HTTPS_Proxy` is the proxy address.
3. Download this program
4. Configure `config.json` at the root of the program that you want to execute.
5. Run the executable file

## tutorial

**For a more detailed tutorial, please see:** [Step by Step Guide for DownloadBot](docs/DownloadBot_Guide_en.md)

[DownloadBot Q&A](docs/Q&A_en.md)


## Screenshots

<div align="center">
<img src="./img/1.jpg" height="300px" alt=""> <img src="./img/2.jpg" height="300px" alt="" >  
</div>
<br>
<div align="center">
<img src="./img/3.jpg" height="300px" alt=""> <img src="./img/4.jpg" height="300px" alt="" >  </div>


## Example of a profile

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
#### Corresponding explanations
* input : Input method, currently only supports aria2
  
  * aria2-server : Aria2 server address. Websocket connection is used by default. If you want to use websocket to connect to aria2, be sure to set `enable-rpc=true` in `aria2.conf`. If not necessary, please try to **set the local aria2 address**, in order to maximize the use of this program
  * aria2-key : The value of `rpc-secret` in `aria2.conf`
* output : Output method, currently only supports telegram
  * bot-key : ID of telegram Bot, get it by using [@BotFather](https://telegram.me/botfather)
  * user-id : The ID of the administrator. ~~It supports setting multiple users as administrators. Different users are separated by commas `,` . If you want to set the users whose `user-id` are 123465789, 987654321 and 963852741 as administrators, you need to set them as follows:~~
    ```jsonc
    {
      //···
      "user-id": "123456789",
      //···
    }
    ```
* max-index：Maximum display quantity of download information, 10 pieces are recommended (to be improved in the future)
* sign: Identification of this Bot, If multiple servers are required to connect to the same Bot, the specific server can be determined through this item.
* language: Language of Bot output
* downloadFolder: Aria2 download file save address.If you do not use this parameter, enter `""`
* moveFolder： The folder to which you want to move the files for the `downloadFolder`. If you do not use this parameter, enter `""`
* server: Server configuration
  * isServer: Whether to enable server mode, if you want to use this program as a server, set it to `true`(When set to `false`, it means that this machine is a client)
  * isMasterServer: Whether to enable master server mode, if you want to use this program as a master server, set it to `true`(now must be set to `true`)
  * serverHost: if it is a client, this item needs to fill in the server address, if it is the main server, this item is the local address
  * serverPort: If it is a client, this item needs to fill in the server port, if it is the main server, this item is the port provided to the client
* log: Log configuration
  * logPath: Log file path, if you do not use this parameter, enter `""`(now invalid)
  * errPath: Error log file path, if you do not use this parameter, enter `""`(now invalid)
  * level: Log level, `debug`, `info`, `warn`, `error`, `fatal`, `panic` are supported, the default is `info`

#### Currently supported languages and language tags
| Languages           | Tag   |
|---------------------|-------|
| English             | en    |
| Simplified Chinese  | zh-CN |
| Traditional Chinese | zh-TW |

When you fill in the above language tag in `config.json`, the program will automatically download the language pack

#### About user-id
If you don't know your `user-id`, you can leave this field blank and enter `/myid` after running the Bot, and the Bot will return your `user-id`

#### donator
If you want to support this project, you can donate to the following address, thank you very much!

https://ko-fi.com/gaowanliang