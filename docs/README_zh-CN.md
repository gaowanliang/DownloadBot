# DownloadBot

[![Go Version](https://img.shields.io/github/go-mod/go-version/gaowanliang/DownloadBot.svg?style=flat-square&label=Go&color=00ADD8)](https://github.com/gaowanliang/DownloadBot/blob/master/go.mod)
[![Release Version](https://img.shields.io/github/v/release/gaowanliang/DownloadBot.svg?style=flat-square&label=Release&color=1784ff)](https://github.com/gaowanliang/DownloadBot/releases/latest)
[![GitHub license](https://img.shields.io/github/license/gaowanliang/DownloadBot.svg?style=flat-square&label=License&color=2ecc71)](https://github.com/gaowanliang/DownloadBot/blob/master/LICENSE)
[![GitHub Star](https://img.shields.io/github/stars/gaowanliang/DownloadBot.svg?style=flat-square&label=Star&color=f39c12)](https://github.com/gaowanliang/DownloadBot/stargazers)
[![GitHub Fork](https://img.shields.io/github/forks/gaowanliang/DownloadBot.svg?style=flat-square&label=Fork&color=8e44ad)](https://github.com/gaowanliang/DownloadBot/network/members)

(目前) 🤖 一个可以控制你的Aria2服务器、控制服务器文件，同时可以上传到OneDrive的Telegram Bot。

## 意义

这个项目主要就是利用吃灰小盘vps进行离线下载，对于大bt文件进行根据硬盘大小分段下载，每次都下载一部分，然后上传网盘，删除再下载其他部分，直到下载完所有文件。

同时，通过机器人协议通信，方便在无法进行内网穿透的机器上进行使用，而且简化了平时使用下载程序的操作，提高了便利性。对于链接，直接向Bot发送消息就可以直接识别并下载，可以真正删除下载文件夹里的文件，是AriaNG等web面板无法做到的，作为管理下载的工具，及时通知下载完成都是非常的方便的。可以移动文件，对于通过rclone挂载硬盘的用户可以直接通过本程序进行复制粘贴等操作，无需打开ssh连接VPS进行`cp`操作，也非常的方便。

## 实现

<text style="color:red;">**注意：本项目仍处于测试阶段，提交的Release仅供测试，现在下载后并不保证您的稳定使用，也不能保证下面所勾选的内容已经被实现。当真正可以正常使用的时候，我会提交 V1.0 版本（V1.0 版本不会实现下面全部功能，但是已经可以正常稳定的使用）**</text>

#### 下载方式

- [x] Aria2 控制
  - [x] 持久化监控
  - [x] 断线重连
- [ ] 多下载服务器同时控制
  - [ ] 多服务器之间通过有公网IP的服务器进行WebSocket通信
  - [ ] 允许用户建立公共WebSocket中继端，供不方便建立WebSocket通信的用户进行通信
  - [ ] 在heroku单独部署WebSocket中继端进行中继
- [ ] [SimpleTorrent](https://github.com/boypt/simple-torrent) 控制
- [ ] qbittorrent 控制


#### 机器人协议支持

- [x] Telegram Bot
  - [ ] 支持多用户使用
  - [ ] 支持群组内使用
- [ ] 腾讯QQ（使用普通QQ用户来进行交互）
- [ ] 钉钉机器人


#### 功能

- [x] 控制服务器文件
    - [x] 删除文件
    - [x] 移动文件
    - [ ] 压缩文件
- [x] 下载文件
    - [x] 下载 HTTP/FTP 链接
    - [x] 下载 Magnet 链接
    - [x] 下载 BitTorrent 文件内的文件
    - [x] 智能 BitTorrent/Magnet 下载方式
        - [x] 只选择下载最大的文件
        - [x] 根据文件大小智能选择文件，不选择小文件
    - [ ] 自适应环境存储空间的 BitTorrent/Magnet 下载
        - [ ] 不下载超过存储空间的文件
        - [ ] 根据存储空间分块多次下载 BitTorrent/Magnet 内的文件
    - [ ] 无感觉化的做种功能
      - [ ] 每次下载BitTorrent/Magnet文件后，保留最后一次下载的文件进行做种，直到下一次下载开始。
      - [ ] 可设置每次下载结束后强制做种一段时间
- [x] 上传文件
    - [x] 下载完成后，向 OneDrive 上传文件
      - [ ] 断点续传
    - [x] 下载完成后，向 Google Drive 上传文件
    - [ ] 下载完成后，向 Mega 上传文件
    - [ ] 下载完成后，向 天翼网盘 上传文件
    - [ ] (当使用Telegram进行通信时)下载完成后，向 Telegram 上传文件
      - [ ] 当文件超过2GB时，分块压缩后再进行上传
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
        - [x] 文件树输出系统
            - [x] 对于简单文件夹的文件树输出
            - [ ] 对于复杂文件夹结构使用图片代替文字输出
        - [ ] 通过演员ID获取在DMM中使用的所有CID
        - [ ] 查询 "ikoa"中的影片参数(利用mahuateng)
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

1. 通过 [@BotFather](https://telegram.me/botfather) 创建您自己的bot并使用。
2. （可选）您所在地区/国家的Telegram被封锁？一定要有一个 **HTTP** proxy启动并运行，您可以设置您的系统环境变量`HTTPS_PROXY`为代理地址来进行代理。
3. 下载本程序
4. 在想要执行本程序的根目录配置`config.json`
5. 运行可执行文件

## 使用截图

<div align="center">
<img src="./img/1.jpg" height="300px" alt="">  <img src="./img/2.jpg" height="300px" alt="" >  
</div>
<br>

<div align="center">
<img src="./img/3.jpg" height="300px" alt="">  <img src="./img/4.jpg" height="300px" alt="" >  </div>


## 配置文件示例

```json
{
  "aria2-server": "ws://127.0.0.1:5800/jsonrpc",
  "aria2-key": "xxxxxxxx",
  "bot-key": "123456789:xxxxxxxxx",
  "user-id": "123456789",
  "max-index": 10,
  "sign": "Main Aria2",
  "language": "zh-CN",
  "downloadFolder": "C:/aria2/Aria2Data",
  "moveFolder":"C:/aria2/GoogleDrive"
}
```

#### 各项对应解释

* aria2-server：aria2服务器地址，默认使用websocket连接。如果要使用websocket连接aria2，请务必设置`aria2.conf`内的`enable-rpc=true`
  。如果不是必须，请尽量设置本地的aria2地址，以便于最大化的使用本程序
* aria2-key：`aria2.conf`中`rpc-secret`的值
* bot-key：Telegram Bot的标识，通过 [@BotFather](https://telegram.me/botfather)进行获取。
* user-id：管理员的ID，支持设置多用户为管理员，不同的用户之间使用半角逗号`,`分割。如您要设置`user-id`为123465789、987654321和963852741的用户为管理员，您需要这样设置：
  ```jsonc
  {
    //···
    "user-id": "123456789,987654321,963852741",
    //···
  }
  ```
* max-index：下载信息最大显示数量，建议10条（以后会改进）
* sign：此机器人的标识，如果需要多个服务器连接同一个机器人，通过这一项可以确定具体是哪一台服务器
* language：机器人输出的语言
* downloadFolder：Aria2下载文件保存的地址。如果不使用，请输入`""`
* moveFolder： 要将下载文件夹的文件移动到的文件夹。如果不使用，请输入`""`

#### 目前支持的语言及语言标签

| 语言     | 标签  |
|----------|-------|
| 英语     | en    |
| 简体中文 | zh-CN |
| 繁体中文 | zh-TW |

当您在`config.json`中填写上面语言的标签的时候，程序会自动下载语言包

#### 关于user-id

如果您不知道您的 `user-id` ，可以将此项留空，在运行这个机器人后输入`/myid`，此机器人就会返回您的`user-id`.
