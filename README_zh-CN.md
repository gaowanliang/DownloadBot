# DownloadBot



(目前)🤖 一个控制你的Aria2服务器的Telegram Bot。


## 实现

#### 下载方式
- [x] Aria2 控制
- [ ] [SimpleTorrent](https://github.com/boypt/simple-torrent) 控制
- [ ] qbittorrent 控制

#### 机器人协议支持
- [x] Telegram Bot
- [ ] 钉钉机器人

#### 功能
- [x] 下载文件
  - [x] 下载 HTTP/FTP 链接
  - [x] 下载 Magnet 链接
  - [x] 下载 BT 文件内的文件
  - [ ] 自定义的 Torrent/Magnet 下载
    - [ ] 不下载小于指定大小的文件
  - [ ] 自适应环境存储空间的 Torrent/Magnet 下载
    - [ ] 不下载超过存储空间的文件
    - [ ] 根据存储空间分块多次下载 Torrent/Magnet 内的文件
- [ ] 上传文件
  - [ ] 下载完成后，向 OneDrive 上传文件
  - [ ] 下载完成后，向 Google Drive 上传文件
  - [ ] 下载完成后，向 Mega 上传文件
  - [ ] 下载完成后，向 天翼网盘 上传文件
- [x] 附加其他功能
  - [x] 多语言支持
    - [x] 简体中文
    - [x] 英语
    - [x] 繁体中文
    - [ ] 日语
  - [ ] 无人值守的BT站下载
    - [ ] Nyaa
    - [ ] ThePirateBay
  - [ ] 其他功能
    - [ ] 通过演员ID获取在DMM中使用的所有CID
    - [ ] 查询 "ikoa "中的影片参数(利用mahuateng)
    - [ ] 通过javlibary演员网址获得所有演员的编号。
    - [ ] 查询dmm cid信息、预览影片、预览图片。
    - [ ] 在sukebei中按关键词搜索。
    - [ ] 根据关键词在dmm中搜索，最多30项。
    - [ ] 输入dmm链接，列出所有项目。
    - [ ] 搜索当前dmm热门和最新电影，限制30条(测试版)

## 目前特点
1. 完全基于触摸，更容易使用，使用这个机器人基本不需要命令。
2. 实时通知，使用Aria2的Websocket协议进行通信。
3. 更好的配置文件支持。

## 开始

1. 通过[@BotFather](https://telegram.me/botfather)创建您自己的bot并使用。
2. （可选）您所在地区/国家的Telegram被封锁？一定要有一个 **HTTP** proxy启动并运行，您可以设置您的系统环境变量`HTTPS_PROXY`为代理地址来进行代理。
3. 下载本程序
4. 在想要执行本程序的根目录配置`config.json`
5. 运行可执行文件

## 三种方式传递参数
您可以通过三种方式将参数传递给`DownloadBot`：
* [X] 配置文件
* [ ] Cil 命令行
* [ ] 系统环境变量


Option priorities also follow this order, so cli has the highest priority.

|                             	| Aria2 server    	| Aria2 key    	| Telegram bot key 	| Telegram user id 	 |Max items in range(default 20) 	|language|
|-----------------------------	|-----------------	|--------------	|------------------	|------------------	|--------------------------------	|---|
| 配置文件 参数   	| aria2-server    	| aria2-key    	| bot-key          	| user-id          	 |max-index                      	|language|
| Cil 命令行 参数                  	| --aria2-server  	| --aria2-key  	| --bot-key        	| --user-id        |--max-index                    	|--language|
| 系统环境变量参数 	| ta.aria2-server 	| ta.aria2-key 	| ta.bot-key       	| ta.user-id       	|ta.max-index                   	|ta.language|


### 配置文件示例

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
如果您不知道您的 `user-id` ，可以将此项留空，在运行这个机器人后输入`/myid`，此机器人就会返回您的`user-id`.

