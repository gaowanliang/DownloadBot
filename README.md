[简体中文](README_zh-CN.md) [繁體中文](README_zh-TW.md)

# DownloadBot

(Currently) 🤖  A Telegram Bot that controls your Aria2 server. 

## Functions realized

#### Download method
- [x] Aria2 control
- [ ] [SimpleTorrent](https://github.com/boypt/simple-torrent) control
- [ ] qbittorrent control

#### The robot protocol supports
- [x] Telegram Bot
- [ ] DingTalk robot

#### Function
- [x] Download files
  - [x] Download HTTP/FTP link
  - [x] Download Magnet link
  - [ ] Download the files in the BT file
  - [ ] Custom Torrent/Magnet download
    - [ ] Do not download files smaller than the specified size
  - [ ] Download Torrent/Magnet according to the size of storage space
    - [ ] Do not download files that exceed storage space
    - [ ] Download the files in Torrent/Magnet several times according to the storage space
- [ ] Upload a file
  - [ ] Upload the file to OneDrive when the download is complete
  - [ ] Upload a file to Google Drive when the download is complete
  - [ ] Upload the file to Mega when the download is complete
  - [ ] Upload the file to 189Cloud when the download is complete
- [x] Additional features
  - [x] Multilingual support
    - [x] Simplified Chinese
    - [x] English
    - [x] Traditional Chinese
    - [ ] Japanese
  - [ ] Download of unattended BT station
    - [ ] Nyaa
    - [ ] ThePirateBay
  - [ ] Other functions
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
2. (Optional) Telegram blocked in your region/country? be sure to have a HTTP proxy up and running,and You can set your system environment variable `HTTPS_ Proxy` is the proxy address.
3. Download this program
4. Configure `config.json` at the root of the program that you want to execute.
5. Run the executable file

## 3 ways to pass parameters
You can pass parameters to `DownloadBot in three ways:
* [X] configuration file
* [ ] cli
* [ ] environment variable

Option priorities also follow this order, so cli has the highest priority.

|                             | Aria2 server    | Aria2 key    | Telegram bot key | Telegram user id | Max items in range(default 20) | language    |
|-----------------------------|-----------------|--------------|------------------|------------------|--------------------------------|-------------|
| configuration file option   | aria2-server    | aria2-key    | bot-key          | user-id          | max-index                      | language    |
| cli option                  | --aria2-server  | --aria2-key  | --bot-key        | --user-id        | --max-index                    | --language  |
| environment variable option | ta.aria2-server | ta.aria2-key | ta.bot-key       | ta.user-id       | ta.max-index                   | ta.language |

## Example of a profile

```json
{
  "aria2-server": "ws://192.168.1.154:6800/jsonrpc",
  "aria2-key": "xxx",
  "bot-key": "123456789:xxx",
  "user-id": "123456",
  "max-index": 10,
  "language":"en"
}
```
If you don't know your `user-id`, you can leave this field blank and enter `/myid` after running the robot, and the robot will return your `user-id`


