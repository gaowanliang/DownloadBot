# DownloadBot



(目前)🤖 一个控制你的Aria2服务器的Telegram Bot。


## 实现

#### 下载方式
- [x] Aria2 控制
- [ ] [SimpleTorrent](https://github.com/boypt/simple-torrent) 控制
- [ ] qbittorrent 控制
- [ ] 多下载服务器同时控制

#### 机器人协议支持
- [x] Telegram Bot
- [ ] 钉钉机器人

#### 功能
- [x] 控制服务器文件
  - [x] 删除文件
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

### 配置文件示例

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
#### 各项对应解释
* aria2-server：aria2服务器地址，默认使用websocket连接。如果要使用websocket连接aria2，请务必设置`aria2.conf`内的`enable-rpc=true`。如果不是必须，请尽量设置本地的aria2地址，以便于最大化的使用本程序
* aria2-key：aria2.conf中rpc-secret的值
* bot-key：Telegram Bot的标识
* user-id：管理员的ID
* max-index：下载信息最大显示数量，建议10条（以后会改进）
* sign：此机器人的标识，如果需要多个服务器连接同一个机器人，通过这一项可以确定具体是哪一台服务器
* language：机器人输出的语言
* downloadFolder：Aria2下载文件保存的地址

#### 目前支持的语言及语言标签
| 语言     | 标签  |
|----------|-------|
| 英语     | en    |
| 简体中文 | zh-CN |
| 繁体中文 | zh-TW |

当您在`config.json`中填写上面语言的标签的时候，程序会自动下载语言包

#### 关于user-id
如果您不知道您的 `user-id` ，可以将此项留空，在运行这个机器人后输入`/myid`，此机器人就会返回您的`user-id`.

