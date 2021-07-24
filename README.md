# ds-go-node
duino coin node

# Compile
```console
CGO_ENABLED=0 go build
```

# Unit/Bench Test
```console
go test
go test -bench=.
```

# CMD Line Options
```console
$ MINERNAME=dsminer ./ds-go-node -help
Usage of ./ds-go-node:
  -batch int
        how many jobs to create. (default 10)
  -debug
        console log send/receive messages.
  -name string
        wallet/miner name. (default "dsminer")
  -quiet
        disable logging to console.
  -server string
        addr and port of server. (default "server.duinocoin.com:2817")
  -wait int
        time to wait between task checks. (default 10)
```

# Example Session
```console
$ ./ds-go-node -name dsminer -batch 3
[2021-07-24T18:35:52Z] Connecting to Server: server.duinocoin.com:2817
[2021-07-24T18:36:07Z] Connected to Server Version: 2.5
[2021-07-24T18:36:23Z] Get Job Response: CREATE_JOBS,18f626d76567d218943c760fe23f1f4d513248cd,6
[2021-07-24T18:36:40Z] Submit Job Response: OK
[2021-07-24T18:36:49Z] Get Job Response: CREATE_JOBS,90595d6dc0e1eab0dc3b9753a75e4dff29bf52a7,6
[2021-07-24T18:37:03Z] Submit Job Response: OK
[2021-07-24T18:37:20Z] Get Job Response: NO_TASK,,0
[2021-07-24T18:37:20Z] no_task sleep for 10s
[2021-07-24T18:37:40Z] Get Job Response: CREATE_JOBS,3d220adde819d7b69a8e46adf8dfa61fe9dbcf06,6
[2021-07-24T18:37:56Z] Submit Job Response: OK
[2021-07-24T18:38:25Z] Get Job Response: CREATE_JOBS,9c00c09edb352b2ac7dc6a0845d23b62aeb1d48c,6
[2021-07-24T18:38:46Z] Submit Job Response: OK
```
