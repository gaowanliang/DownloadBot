# Step by Step Guide for DownloadBot 
Written in [GoLang](https://golang.org) 

Original Author by [gaowanliang](https://github.com/gaowanliang) 

Guide by [reaitten](https://github.com/reaitten)

###  Please note that the project in Beta!   
###  I reccomend you to treat this as a preview only!!

## Useful Links
[Github Repository](https://github.com/gaowanliang/DownloadBot)

[Very useful for getting a VPS and more!](https://free-for.dev/#/)

## Prerequisites
You'll need a Linux system, preferably Debian or Debian based systems such as Ubuntu or Linux Mint

Any other Linux distros should work, although switch the commands for your package manager. 

For Windows users, you may use the [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/install-win10) (you may face problems deploying, I haven't tried yet)

A basic knowledge of Linux commands

A VPS (Not needed if deploying locally) 

## Installing the Go Language

You'll need to install the GoLang in order to build the project. 

Open a terminal or SSH into your VPS.

### Update Linux 
```
sudo apt-get update && sudo apt-get upgrade -y
```
### Install needed dependencies
```
sudo apt-get install nano wget curl ca-certificates git -y
```
### Download GoLang Ver. 1.16.3
#### x64 based System (for newer hardware, common)

Will be saved to user's Downloads folder.
```
wget https://golang.org/dl/go1.16.3.linux-amd64.tar.gz -P "$HOME/Downloads"
```
#### x86-32 based System (older architecture)
```
wget https://golang.org/dl/go1.16.3.linux-386.tar.gz -P "$HOME/Downloads"
```
### Extract from .tar.gz archive /usr/local/

We extract GoLang.tar.gz archive to /usr/local for easy convenience.
```
cd $HOME/Downloads
sudo tar -C /usr/local/ -xzf go1.16.3.linux-amd64.tar.gz
```
### Check if everything has been extracted correctly

If you see files with the following command, everything extracted correctly.
```
ls /usr/local/go
```
### Set the PATH for GoLang
If you don't know what everything means, it's okay since we're just linking GoLang to the system so whever we need to use GoLang, the system will know where the GoLang is e.g: /usr/local/go

You'll be adding a line of text at the bottom of .profile, follow carefully.

Edit .profile using a command-line editor.
```
sudo nano $HOME/.profile
```
Move all the way down at the bottom of the page using DownArrowKey

Then, add the follow line:
```
export PATH=$PATH:/usr/local/go/bin
```
Everything should look like this:
```
# ~/.profile: executed by the command interpreter for login shells.
# This file is not read by bash(1), if ~/.bash_profile or ~/.bash_login
# exists.
# see /usr/share/doc/bash/examples/startup-files for examples.
# the files are located in the bash-doc package.

# the default umask is set in /etc/profile; for setting the umask
# for ssh logins, install and configure the libpam-umask package.
#umask 022

# if running bash
if [ -n "$BASH_VERSION" ]; then
    # include .bashrc if it exists
    if [ -f "$HOME/.bashrc" ]; then
        . "$HOME/.bashrc"
    fi
fi

# set PATH so it includes user's private bin if it exists
if [ -d "$HOME/bin" ] ; then
    PATH="$HOME/bin:$PATH"
fi

# set PATH so it includes user's private bin if it exists
if [ -d "$HOME/.local/bin" ] ; then
    PATH="$HOME/.local/bin:$PATH"
fi

export PATH=$PATH:/usr/local/go/bin
```
Do ``ctrl+o, enter, ctrl+x`` to save, & exit

After, you'll need to apply the changes to the system
```
cd $HOME/
source .profile
```
Now, to make sure GoLang is installed correctly, do
```
go version
```
If you get ```go version go1.16.3 linux/amd64```, you have successfully install GoLang!

## Installing Aria2
We need aria2 for the download process of DownloadBot.

This script that we're going to need to use is in Chinese.

Will be saved in user's Downloads folder
### Download Script
```
wget -N git.io/aria2.sh -P $HOME/Downloads && chmod +x $HOME/Downloads/aria2.sh
```
### Run Script
```
./aria2.sh
```
### To use the script:
Just enter 1 to install and configure aria2 automatically, If you want to use it on azure, pay attention to limit the upload of BT download. The configuration file is in /$HOME/.aria2c/aria2.conf. you can refer to this website to customize the information in it.

### Start aria2
You'll need to start aria2. Use:
```
sudo service aria2 start
```

### Simple translation of Aria2 Script
```
Aria2 one click installation management script enhanced version [v2.7.4] by P3 TERX.COM
0. Upgrade script
———————————————————————
1. Install aria2
2. Update aria2
3. Unload aria2
———————————————————————
4. Start aria2
5. Stop aria2
6. Restart aria2
———————————————————————
7. Modify the configuration
8. View configuration
9. View log
10. Clear log
———————————————————————
11. Update BT tracker manually
12. Update BT tracker automatically
———————————————————————
Aria2 status: installed | started
Auto update BT Tracker: on
Please enter the number [0-12]:
```

## Cloning DownloadBot Repository
Now, you'll need to copy the DownloadBot repo to your computer.
Run the following to clone the repository to your Downloads folder:
```
cd $HOME/Downloads
git clone https://github.com/gaowanliang/DownloadBot.git
```
## Setup Configs

After we clone the repository, we need to configure DownloadBot with ``config.json``.

Use this command, copy everything listed & paste it somewhere that you can refer to.
```
cat $HOME/aria2/aria2.conf
```


### config.json

You'll need to create a config file. 

This command will create and let you edit ``config.json`` file.
```
nano config.json
```
Copy everything below onto ``config.json`` and edit with your own values.

```
{
    "aria2-server": "ws://127.0.0.1:6800/jsonrpc",
    "aria2-key": "xxxxxxxx",
    "bot-key": "123456789:xxxxxxxxx",
    "user-id": "123456789",
    "max-index": 10,
    "sign":"Main Aria2",
    "language":"en",
    "downloadFolder":"/home/user/aria2/Aria2Data",
    "moveFolder":"/home/user/aria2/GoogleDrive"
}
```
Note: aria2-server should be ``ws://127.0.0.1:6800/jsonrpc`` 

You need to only fill out the following:
```
"aria2-key" # from $HOME/aria2/aria2.conf
"bot-key" # from @BotFather (Telegram)
"user-id" # from @userinfobot (Telegram)
"downloadFolder" # create folder
"moveFolder" # create folder
```

You'll also need to create two folders for ``downloadFolder`` & ``moveFolder`` using `mkdir`. e.g:
```
mkdir $HOME/Downloads/Aria2Data # downloadFolder
mkdir $HOME/Downlaods/GoogleDrive # moveFolder
```

#### Corresponding explanations
* aria2-server : Aria2 server address. Websocket connection is used by default. If you want to use websocket to connect to aria2, be sure to set `enable-rpc=true` in `aria2.conf`. If not necessary, please try to **set the local aria2 address**, in order to maximize the use of this program
* aria2-key : The value of `rpc-secret` in `aria2.conf`
* bot-key : ID of telegram Bot, get it by using [@BotFather](https://telegram.me/botfather)
* user-id : The ID of the administrator. You can get this value from [@userinfobot](https://telegram.me/userinfobot). It supports setting multiple users as administrators. Different users are separated by commas `,` . If you want to set the users whose `user-id` are 123465789, 987654321 and 963852741 as administrators, you need to set them as follows:
  ```json
  {
    ···
    "user-id": "123456789,987654321,963852741",
    ···
  }
  ```
* max-index：Maximum display quantity of download information, 10 pieces are recommended (to be improved in the future)
* sign: Identification of this Bot, If multiple servers are required to connect to the same Bot, the specific server can be determined through this item.
* language: Language of Bot output
* downloadFolder: Aria2 download file save address. If you do not use this parameter, enter `""`
* moveFolder： The folder to which you want to move the files for the `downloadFolder`. If you do not use this parameter, enter `""`


## Build DownloadBot
Last part, you'll need to build DownloadBot & run it.

``go build`` will build the program,

``sudo chmod u+x ./DownloadBot`` will add the required permissions to run DownloadBot,

``./DownloadBot`` will run DownloadBot.

```
go build
sudo chmod u+x ./DownloadBot
./DownloadBot
```
Sources:

https://golangdocs.com/install-go-linux

https://github.com/gaowanliang/DownloadBot/issues/12

https://github.com/gaowanliang/DownloadBot#corresponding-explanations
