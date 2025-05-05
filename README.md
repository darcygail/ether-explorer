# ether-explorer
a demo of view ether accounts's assets
## 功能
1. 通过 worker 模块获取区块，提取出账户信息，获取账户余额
2. 通过 api 模块获取账户的资产信息
3. 资产信息包括 账户余额和nft资产


## 项目结构
```bash
.
├── LICENSE
├── README.md
├── cmd
│   ├── api.go
│   ├── root.go
│   └── worker.go
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   └── router.go
│   ├── eth
│   │   └── fetcher.go
│   ├── parser
│   │   └── ntf_parser.go
│   ├── store
│   │   ├── asset.go
│   │   └── mongo.go
│   └── worker
│       └── worker.go
├── main.go
└── schema
    ├── account_asset.go
    └── worker_config.go
```




## TODO
1. 通过多协程提升 worker 速度，解析区块，并进行存储
2. 需要加入锁机制，防止同时更新
3. 增加 fetch 进度来检查已同步区块， 方便进行增量更新和后续代码拓展
4. 需要记录 fetcher 状态，用于同步更新
5. 增加缓存来提升查询性能