# it-chain ivm 'bank sample'
- deploy and invoke/query test

- ivm Log can be checked by 'docker logs -f [docker_name]' command
```
[root@itchain ~]# docker ps --format 'ID:{{.ID}} CMD:{{.Command}} NAME:{{.Names}}'
ID:6e84078ca76c CMD:"go run /go/src/gi..." NAME:focused_edison
[root@itchain ~]# docker logs -f focused_edison
port : 50002   
```

1. run it-chain
   ```
   [root@itchain engine]# it-chain
   
           ___  _________               ________  ___  ___  ________  ___  ________
           |\  \|\___   ___\            |\   ____\|\  \|\  \|\   __  \|\  \|\   ___  \
           \ \  \|___ \  \_|____________\ \  \___|\ \  \\\  \ \  \|\  \ \  \ \  \\ \  \
            \ \  \   \ \  \|\____________\ \  \    \ \   __  \ \   __  \ \  \ \  \\ \  \
             \ \  \   \ \  \|____________|\ \  \____\ \  \ \  \ \  \ \  \ \  \ \  \\ \  \
              \ \__\   \ \__\              \ \_______\ \__\ \__\ \__\ \__\ \__\ \__\\ \__\
               \|__|    \|__|               \|_______|\|__|\|__|\|__|\|__|\|__|\|__| \|__|
   ```

2. deploy
   - run deploy
   ```
   [root@itchain engine]# it-chain ivm deploy github.com/nesticat/bank-icode
   INFO[2018-10-15T09:52:19+09:00] [Cmd] deploying icode...
   INFO[2018-10-15T09:52:19+09:00] [Cmd] This may take a few minutes
   ```

3. invoke mint
   - run invoke mint
   ```
   [root@itchain engine]# it-chain ivm invoke bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd mint sykwon 10000
   INFO[2018-10-15T10:00:23+09:00] [Cmd] Invoke icode - icodeID: [bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd]
   INFO[2018-10-15T10:00:23+09:00] [Cmd] Transactions are created - ID: [bf1ud9u5apve3hejo6a0]
   [root@itchain engine]#
   [root@itchain engine]# it-chain ivm invoke bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd mint smoon 10000
   INFO[2018-10-15T10:00:30+09:00] [Cmd] Invoke icode - icodeID: [bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd]
   INFO[2018-10-15T10:00:30+09:00] [Cmd] Transactions are created - ID: [bf1udbm5apve3hejo6b0]
   [root@itchain engine]#
   ```
   - ivm log
   ```
   INFO[2018-10-15T01:00:24Z] invoke - mint sykwon/10000
   INFO[2018-10-15T01:00:30Z] invoke - mint smoon/10000
   ```

4. invoke transfer
   - run invoke transfer
   ```
   [root@itchain engine]# it-chain ivm invoke bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd transfer sykwon smoon 1000
   INFO[2018-10-15T10:03:49+09:00] [Cmd] Invoke icode - icodeID: [bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd]
   INFO[2018-10-15T10:03:49+09:00] [Cmd] Transactions are created - ID: [bf1uete5apve3hejo6c0]
   [root@itchain engine]#
   ```
   - ivm log
   ```
   INFO[2018-10-15T01:03:50Z] invoke - transfer sykwon->smoon:1000
   ```

5. query accounts
   - run query accounts
   ```
   [root@itchain engine]# it-chain ivm query bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd accounts
   INFO[2018-10-15T10:05:59+09:00] [Cmd] Querying icode - icodeID: [bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd], func:    [accounts]
   INFO[2018-10-15T10:05:59+09:00] [CMD] Querying result - key: [smoon], value: [11000]
   INFO[2018-10-15T10:05:59+09:00] [CMD] Querying result - key: [sykwon], value: [9000]
   ```
   - ivm log
   ```
   INFO[2018-10-15T01:05:59Z] query - accounts sample_smoon/11000
   INFO[2018-10-15T01:05:59Z] query - accounts sample_sykwon/9000
   ```

6. query balance
   - run query balance
   ```
   [root@itchain engine]# it-chain ivm query bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd balance sykwon
   INFO[2018-10-15T10:06:31+09:00] [Cmd] Querying icode - icodeID: [bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd], func: [balance]
   INFO[2018-10-15T10:06:31+09:00] [CMD] Querying result - key: [sykwon], value: [9000]
   [root@itchain engine]#
   [root@itchain engine]# it-chain ivm query bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd balance smoon
   INFO[2018-10-15T10:06:36+09:00] [Cmd] Querying icode - icodeID: [bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd], func: [balance]
   INFO[2018-10-15T10:06:36+09:00] [CMD] Querying result - key: [smoon], value: [11000]
   [root@itchain engine]#
   ```
   - ivm log
   ```
   INFO[2018-10-15T01:06:31Z] query - balance sykwon:9000
   INFO[2018-10-15T01:06:36Z] query - balance smoon:11000
   ```

7. undeploy
   - run undeploy
   ```
   [root@itchain engine]# it-chain ivm undeploy bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd
   2018/10/15 10:07:03 [bank-icode_cfea3fecf8921495dfa3ae975ad175dbebaaa9cd] icode has undeployed
   ```
