# it-chain CLI 
 NOTE: this documentation is for it-chain CLI.

## Overview
 The it-chain CLI is a tool to ivm deploy, undeploy, invoke, query, listing..

## Usage
```
USAGE:
   it-chain [global options] command [command options] [arguments...]
```

## GLOBAL OPTIONS
- --config value : explicit set config file 
```
[root@it-chain engine]# it-chain ./conf/config.yaml

        ___  _________               ________  ___  ___  ________  ___  ________
        |\  \|\___   ___\            |\   ____\|\  \|\  \|\   __  \|\  \|\   ___  \
        \ \  \|___ \  \_|____________\ \  \___|\ \  \\\  \ \  \|\  \ \  \ \  \\ \  \
         \ \  \   \ \  \|\____________\ \  \    \ \   __  \ \   __  \ \  \ \  \\ \  \
          \ \  \   \ \  \|____________|\ \  \____\ \  \ \  \ \  \ \  \ \  \ \  \\ \  \
           \ \__\   \ \__\              \ \_______\ \__\ \__\ \__\ \__\ \__\ \__\\ \__\
            \|__|    \|__|               \|_______|\|__|\|__|\|__|\|__|\|__|\|__| \|__|
...
```
- --help, -h : show help
```
[root@it-chain engine]# it-chain -h
NAME:
   it-chain - A new cli application

USAGE:
   it-chain [global options] command [command options] [arguments...]

VERSION:
   0.1.1

AUTHOR:
   it-chain <it-chain@gmail.com>

COMMANDS:
     ivm, i         options for ivm
     connection, c  options for connection
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value  name for config
   --help, -h      show help
   --version, -v   print the version
```
- --version, -v : show version
```
[root@it-chain engine]# it-chain -v
it-chain version 0.1.1
```

## COMMANDS - ivm
- command option
```
[root@it-chain engine]# it-chain i
NAME:
   it-chain ivm - options for ivm

USAGE:
   it-chain ivm command [command options] [arguments...]

COMMANDS:
     deploy    it-chain ivm deploy [icode-git-url] [ssh-path] [password]
     undeploy  it-chain ivm undeploy [icode-id]
     invoke    it-chain ivm invoke [icode-id] [function-name] [...args]
     query     it-chain ivm query [icode-id] [functioniname] [...args]
     list      it-chain ivm list

OPTIONS:
   --help, -h  show help
```
  - deploy : ivm deploy
    - with ssh/password
    ```
    [root@it-chain engine]# it-chain ivm deploy github.com/nesticat/learn-icode ~/.ssh/id_rsa passwd
    INFO[2018-09-27T21:02:08+09:00] [Cmd] deploying icode...
    INFO[2018-09-27T21:02:08+09:00] [Cmd] This may take a few minutes
    ...
    ```
    - without ssh/password
    ```
    [root@it-chain engine]# it-chain ivm deploy github.com/nesticat/learn-icode ~/.ssh/id_rsa
    INFO[2018-09-27T21:00:32+09:00] [Cmd] deploying icode...
    INFO[2018-09-27T21:00:32+09:00] [Cmd] This may take a few minutes
    ...
    ```
  - undeploy : ivm undeploy with icode-id
  ```  
  [root@it-chain engine]# it-chain ivm undeploy bemcj4e5apva4tp7e400
  2018/09/27 21:15:48 [bemcj4e5apva4tp7e400] icode has undeployed
  ```  
  - invoke : ivm invoke with icode-id, function-name and arguments
  ```
  [root@it-chain engine]# it-chain ivm invoke bemcj4e5apva4tp7e400 initA
  INFO[2018-09-27T21:19:55+09:00] [Cmd] Invoke icode - icodeID: [bemcj4e5apva4tp7e400]
  INFO[2018-09-27T21:19:55+09:00] [Cmd] Transactions are created - ID: [bemclqu5apva4g8550kg]
  ```
  - query : ivm query with icode-id, function-name and arguments
  ```
  [root@it-chain engine]# it-chain ivm query bemcj4e5apva4tp7e400 getA
  INFO[2018-09-27T21:21:44+09:00] [Cmd] Querying icode - icodeID: [bemcj4e5apva4tp7e400], func: [getA]
  INFO[2018-09-27T21:21:44+09:00] [CMD] Querying result - key: [A], value: [1]
  ```
  - list : show ivm list
  ```
  [root@it-chain engine]# it-chain ivm list
  Index    ID                      Version         GitUrl
  [0]      [bemcj4e5apva4tp7e400]  [v1.1]          [github.com/nesticat/learn-icode]
  ```
  
## COMMANDS - connection
- command option
```
[root@it-chain engine]# it-chain connection
NAME:
   it-chain connection - options for connection

USAGE:
   it-chain connection command [command options] [arguments...]

COMMANDS:
     dial  it-chain connection dial [node-ip]
     list  it-chain connection list
     join  it-chain connection join [node-ip-of-the-network]

OPTIONS:
   --help, -h  show help
```
  - dial : connection with node-ip
  ```
  [root@it-chain engine]# it-chain connection dial 192.168.56.230:5000
  INFO[2018-09-28T09:32:14+09:00] [Cmd] Creating connection - Address: [192.168.56.230:5000]
  INFO[2018-09-28T09:32:14+09:00] [Cmd] Connection created - gRPC-Address: [192.168.56.230:5000], Id:[B69aLYLeVCeLFTih5fDpuZVkYh4AF78ejZBTcEfkBbz2]
  ```
  - list : show connection list
  ```
  [root@it-chain engine]# it-chain connection list
  Index    ID                                              gRPC-Address    Api-Address
  [0]      [B69aLYLeVCeLFTih5fDpuZVkYh4AF78ejZBTcEfkBbz2]  [192.168.56.230:5000]   [192.168.56.230:4000]
  ```
  - join : connection with node-ip-of-the-network
  ```
  [root@sykwon engine]# it-chain connection join 192.168.56.230:5000
  INFO[2018-09-28T09:56:34+09:00] [Cmd] Joining network - Address: [192.168.56.230:5000]
  INFO[2018-09-28T09:56:34+09:00] [Cmd] Successfully request to join network
  ```
