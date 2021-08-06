# DownloadBot Q&A
This document is a collection of problems and solutions that have arisen during use, providing solutions in a question-and-answer style and from a beginner's perspective.
## Purpose
Since [project author](https://github.com/gaowanliang) does not have enough time to maintain this project and the documentation on its use is not comprehensive, there are many questions in use that are not answered in the documentation. The writer has found solutions to various problems in the process of using the program through his own exploration, so he has written a help file here to help more users who encounter problems. All the problems are actually encountered or issue solved, so the problems you encounter may not be included, I hope more people will come to improve this document.

## Q & A

### Q 1: How to run DownloadBot?

A 1: Execute the command `./DownloadBot` in the working directory

If `-bash:./DownloadBot: Permission denied` is displayed, enter `chmod 777 ./DownloadBot` to grant Permission.

### Q 2:  None of the versions in release are working on my server, what should I do?
A 2: Self-compilation

```bash
# CentOS for example
yum install go # If you are prompted for a package that cannot be found, install epel first: yum install epel-release

go version # If the version number is displayed, the installation is successful and the version needs to be greater than or equal to 1.15.0
cd /root
git clone https://github.com/gaowanliang/DownloadBot && cd DownloadBot
go build
# Compile is complete and you can run DowloadBot
./DownloadBot
```

### Q 3: Prompt `open ./config.json: no such file or directory` when running

A 3: Missing configuration file `config.json`，According to [README](https://github.com/gaowanliang/DownloadBot#example-of-a-profile)

### Q 4: Run error

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

A 4: There are still some problems with multi-user support, so please use a single user first, i.e. the UserID is filled in with an ID first.


### Q 5: How do I download offline?

A 5: Just send a direct link to the file or a BT link directly to your Bot

### Q 6: How do I delete a file from the download directory?

A 6: The current version of the Bot menu does not have a command, you have to manually enter the command: `Delete files from the download folder` and follow the instructions. It seems that in version 0.4 there is a command in the menu, pending a fix.

PS: If you are controlling a remote server to download files through this software, and the downloaded files are on the remote server, naturally this program cannot access the files on the remote server, so it is stipulated that the IP address you set needs to be the local IP address before you are allowed to use the upload/move/delete functions. So only if the IP address is `localhost` or `127.0.0.1` will uploads be allowed.

### Q 7: How do I binding for a OneDrive account?

A 7: Send `Upload all files in the download folder` to Bot, click `1` (OneDrive) in the back menu, click `New` and you will get a URL. open the URL in your browser and log in to authorise your OneDrive account, your browser will then jump to a URL that will not open. Copy it and send it to Bot, the binding is successful.

You can also use your own API (if the web page suggests that the application is not validated), which you need to compile yourself. See [here](https://github.com/gaowanliang/DownloadBot/issues/30#issuecomment-888344140)

### Q 8: How do I unbinding my OneDrive account?

 A 8: Delete `./DownloadBot/info/onedrive/your_onedrive_account.json`

PS: It's meaningless to say binding/unbinding users in this application, you can use multiple accounts at once and they don't interfere with each other

### Q 9: How to run DownloadBot in the background?

A 9: 
1. Add `&` to the end of the command, e.g. `. /DownloadBot &` so that the program will still run when SSH is disconnected. Why not use a script to check the program's running status every certain time?
```bash
#!/bin/bash --posix 

while true  
do   
    procnum=` ps -ef|grep "DownloadBot"|grep -v grep|wc -l`  
   if [ $procnum -eq 0 ]; then  
       cd /home/DownloadBot/ && /home/DownloadBot/DownloadBot &  ## your bot location
   fi  
   sleep 30  ## check status every 30 seconds
done
```
3. Use `screen`, see https://linuxize.com/post/how-to-use-linux-screen/ for details


## Contributors

[DullJZ](https://github.com/DullJZ)

[gaowanliang](https://github.com/gaowanliang)
