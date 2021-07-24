# ds-go-node
duino coin node

# CMD Line Options
```console
$ MINERNAME=itsame ./ds-go-node -help
Usage of ./ds-go-node:
  -name string
        wallet/miner name. (default "itsame")
  -quiet
        disable logging to console.
  -server string
        addr and port of server. (default "server.duinocoin.com:2817")
  -wait int
        time to wait between task checks. (default 10)
```

# Example Session
```console
$ ./ds-go-node -name itsame
[2021-07-24T10:50:15Z] Connecting to Server: server.duinocoin.com:2817
[2021-07-24T10:50:22Z] Connected to Server Version: 2.5
[2021-07-24T10:50:34Z] Get Job Response: CREATE_JOBS,0045bc6eb038d56732d1c471fa1d0a27c495c73c,6
[2021-07-24T10:50:54Z] Submit Job Response: OK
[2021-07-24T10:51:26Z] Get Job Response: CREATE_JOBS,b126667069ffe96612eb10158dc08b55f132223e,6
[2021-07-24T10:51:43Z] Submit Job Response: OK
```
