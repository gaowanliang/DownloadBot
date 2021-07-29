# DownloadBot 问答文档
本文档为使用过程中所产生的问题和解决方式的集合，以问答的方式、小白的视角提供解决方案。
## 目的
由于[项目作者](https://github.com/gaowanliang)没有足够的时间维护此项目，使用文档也并不全面，有许多使用中的问题并不能在文档中得到答案。笔者在使用过程中也是通过自己的摸索找到了各种问题的解决方法，于是在这里写了一个帮助文档，来帮助更多遇到问题的用户。所有的问题都是实际碰到的或者issue解决的，所以你碰到的问题可能没有收录，希望更多人来完善本文档。
## Q & A

### Q 1：如何运行DownloadBot？

A 1：在工作目录下执行`./DownloadBot`

### Q 2： release 中的各个版本均无法在我的服务器上运行，怎么办？
A 2：自行编译

```bash
# 以CentOS 为例
yum install go # 提示无法找到包的，先安装epel：yum install epel-release
go version # 有显示版本号则安装成功，版本需要大于等于14.0
cd /root
git clone https://github.com/gaowanliang/DownloadBot && cd DownloadBot
go build
# 编译完成，可以运行DowloadBot
./DownloadBot
```

### Q 3：运行提示`open ./config.json: no such file or directory`

A 3：缺少配置文件`config.json`，按照 [README文档](https://github.com/gaowanliang/DownloadBot/blob/main/docs/README_zh-CN.md#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E7%A4%BA%E4%BE%8B) 中创建配置文件即可

### Q 4：运行报错

```bash
2021/07/23 01:51:31 strconv.ParseInt: parsing "": invalid syntax
panic: strconv.ParseInt: parsing "": invalid syntax

goroutine 37 [running]:
log.Panic(0xc00025f3a8, 0x1, 0x1)
D:/program_app/go/src/log/log.go:351 +0xac
main.dropErr(...)
D:/program_data/go/DownloadBot/main.go:14
main.toInt64(0x0, 0x0, 0x4)
D:/program_data/go/DownloadBot/utils.go:514 +0x9d
main.TMSelectMessage(0xc000212c80)
D:/program_data/go/DownloadBot/Telegram.go:62 +0x4e
created by main.tgBot
D:/program_data/go/DownloadBot/Telegram.go:707 +0x85d
```

A 4：目前对于多用户支持还有点问题，请先使用单用户，也就是说，UserID先填写一个ID。参见[issue37](https://github.com/gaowanliang/DownloadBot/issues/37)

### Q 5：如何离线下载？

A 5：直接发送文件直链或BT链接到你的Bot即可

### Q 6：如何删除下载目录中的文件？

A 6：目前版本中Bot菜单没有指令，必须手动输入指令：`Delete files from the download folder`，然后按照提示操作即可。似乎0.4版本的菜单中有指令，等待修复。

PS：如果您是通过本软件控制远程服务器下载文件，下载的文件在远程服务器，本程序自然无法访问远程服务器上的文件，因此规定您设置的IP地址需要是本机IP地址，才允许使用上传/移动/删除等功能。所以只有当IP地址为`localhost`或者是`127.0.0.1`时，才允许上传。

### Q 7：如何绑定OneDrive 网盘账户？

A 7：发送 `Upload all files in the download folder` 给Bot，单击返回菜单中的 `1`（OneDrive），单击添加，你会得到一个URL。在浏览器中打开URL并登录授权你的OneDrive账户，随后你的浏览器会跳转至一个无法打开的URL，复制它并发送给Bot，即绑定成功。

你也可以使用自己的API，需要自己编译。参见[这里](https://github.com/gaowanliang/DownloadBot/issues/30#issuecomment-888344140)

### Q 8：如何解绑OneDrive网盘账户？

 A 8：删除 `./DownloadBot/info/onedrive/你的onedrive账号.json`即可

PS：本程序中不存在绑定/解绑用户的问题，你可以一次使用多个账户，他们之间互不干扰

### Q 9：如何后台运行DownloadBot？

A 9：
1. 在指令后添加`&`，如`./DownloadBot &`，这样可以在SSH断开时依旧运行本程序。
2. 使用`screen`，具体方法参考 https://www.runoob.com/linux/linux-comm-screen.html

## 文档贡献者

[DullJZ](https://github.com/DullJZ)
