## dgatewayWebBrowser接口文档

### 1. 首页WebSocket推送消息
url: ws://host:port/info/ws

method: 

params: 无

response:
    #总览信息
    {
        "type": 1,
        "highest_block": 6666, #当前区块高度
        "last_block_time": 10， #上一个区块出块速度, 单位s
        "node_num": 30, #节点数量
        "tx_num": 455, #交易数量
    }

    #区块信息
    {
        "type": 2,
        "height": 2312, #区块高度
        "id": "sdsfsdfsrewrewfdewfe", #区块hash
        "miner": "miner0", #矿工
        "created_used": 10, #出块速度
        "tx_cnt": 35, #交易数量
        "time": 2132312312412,  #出块时间戳
    }
    
    #交易信息
        {
        "type": 3,
        "from_tx_hash": "832e68db1af7be152baf0edbe1f148264ecfca7f58b1a16484ca429c5079b6cc",
        "to_tx_hash": "0x777070e7ac962d94df9414751a51bdc49df11ea422567f19c27e42469e8b9c4c",
        "dgw_hash": "7e35157d7221af7a5ff1286f581ab8706088f917855d9968d6a9dbe8b5501d25",
        "from_chain": "bch",
        "to_chain": "eth",
        "time": 0,
        "block": "9939c541f38753a5acc1a8b555708da563a3cbb6933d07935f330ba025a1353e",
        "block_height": 605,
        "amount": 4300000,
        "from_addrs": "",
        "to_addrs": "0xEFfA957170e7845C3F2d489bB0EC68B1E000952d",
        "from_fee": 0,
        "dgw_fee": 0,
        "to_fee": 0
        }

### 2. 首页交易图表
url: /info/txdata

method: GET

params: 

response: 
    {
        "code": 0,
        "msg": "",
        "data": { #返回数据一一对应
            "time": ["1/6", "2/6", "3/6", ...], #时间
            "count": [45, 56, 67, ...],
            "amount": [10100, 12232, 2324231, ...] 
        }
    }

### 3. 获取节点列表(列表详情数据在返回里面取)
url: /node/list

method: GET

params: 无

response:

    {
        "code": 0,
        "msg": "",
        "data": {
            "count": 9999,
            "online": 9898,
            "data":
                [
                {
                "ip": xx.xxxx.xxx, #ip地址
                "host_name": "host0", #主机名
                "is_leader": true, #是否是leader
                "is_online": false, #是否在线
                "fire_cnt": 10, #被替换次数
                "eth_height": 343432, #eth链高度
                "bch_height": 3423424, #bch链高度
                "btc_height": 123232, #btc链高度
            }，
            。。。
        ],
        }
    }

### 4.区块列表
url: /block/list

method: GET

params: 

    {
        "page": 1,
        "page_size": 20,
        "search": 100 #无搜索不传此参
    }

response: 

    {
        "code": 0,
        "msg": "",
        "data": {
            "count": 1234, #总数       
            "data": [
                {
                "height": 23232, 
                "hash": "ffffffffffffff", 
                "pre_id": "ggggggggg", #前一个区块ID
                "tx_cnt": 3434,
                "time": 1234567890,
                "size": 1111,
                "create_used": 20,
                "miner": "miner0", 
                }
        ]
    }
    }

### 5.区块详情(通过height)
url: /block/detail

method: GET

params:

    {
        "height": 1234
    }

response: 

    {
        "code": 0,
        "msg": "",
        "data": {
            "height": 12332,
            "hash": "eeeeeeeeeeee",
            "pre_id": "vvvvvvvvvvv",
            "tx_cnt": 55,
            "time": 1234567890, 
            "size": 2324,
            "created_used": 10,
            "miner": "miner0",
            "trans": [
                {
                    "from_tx_hash": "vvvvvvvvvvvv",
                    "to_tx_hash": "bbbbbbbbbbbbb",
                    "dgw_hash": "nnnnnnnnnnnnnnn",
                    "from": "eth",
                    "to": "btc",
                    "time": 1234567890,
                    "block": "sdferfecrrfregvtgert",
                    "block_height": "vvfdgfrtvfdasaswxas",
                    "amount": 1000000,
                    "from_addrs": "addr1,addr2,addr3",
                    "from_fee": 1000,
                    "dgw_fee": 0,
                    "to_fee": 1222,
                },
                ...
            ]

        }
    }

### 6. 交易列表
url: /tranx/list

method: GET

params: 

    {
        "page": 1,
        "page_size": 10,
        "search": 122, #无搜索传此参数
    }

response: 

    {
        "code": 0,
        "msg": "",
        "data":{
            "count": 12344,
            "data":
             [
            {
                {
                    "from_tx_hash": "vvvvvvvvvvvv",
                    "to_tx_hash": "bbbbbbbbbbbbb",
                    "dgw_hash": "nnnnnnnnnnnnnnn",
                    "from": "eth",
                    "to": "btc",
                    "time": 1234567890,
                    "block": "sdferfecrrfregvtgert",
                    "block_height": "vvfdgfrtvfdasaswxas",
                    "amount": 1000000,
                    "from_addrs": "addr1,addr2,addr3",
                    "to_addrs": "addr4, addr5, addr6"
                    "from_fee": 1000,
                    "dgw_fee": 0,
                    "to_fee": 1222,
                },
                {...},
                ...
            }
        ]
    }
    }


### 7.交易详情（通过dgw_hash查询）
url: /tranx/detail

method: GET

params:

    {
        "dgw_hash": "nnnnnnnnnnnnnnnn", #网关hash
    }

response: 

    {
        "code": 0,
        "msg": "",
        "data": 
                {
                    "from_tx_hash": "vvvvvvvvvvvv",
                    "to_tx_hash": "bbbbbbbbbbbbb",
                    "dgw_hash": "nnnnnnnnnnnnnnn",
                    "from": "eth",
                    "to": "btc",
                    "time": 1234567890,
                    "block": "sdferfecrrfregvtgert",
                    "block_height": "vvfdgfrtvfdasaswxas",
                    "amount": 1000000,
                    "from_addrs": "addr1,addr2,addr3",
                    "to_addrs": "addr4, addr5, addr6"
                    "from_fee": 1000,
                    "dgw_fee": 0,
                    "to_fee": 1222,
                }
    }

### 8.首页搜索（通过dgw_hash, blockHeight）
url: /info/search
method: GET 
params:

    {
        "search": 124, #块高度或者dgw_hash
    }

response: 

    #同交易列表或者区块列表

### 9.获取最新区块详情
url: /block/currentblock

method: GET

params: 无

response:

    {
        "code": 0,
        "msg": "",
        "data": {
            "height": 4979, #当前区块高度
            "id": "49e3f4586958e8d1f8c719ffbd0a2fafb5cba2c78e638559668e96fca932ea0b",
            "pre_id": "8e1206996a065dff43744eb7e9ae0a6f81f175f2664c8353f30f85ee56ceb76d",
            "tx_cnt": 0,
            "txs": [],
            "time": 1531705920,
            "size": 20721,
            "created_used": 19,
            "miner": "server1"
        }
    }


### 10.获取首页总览信息
url /info/overview

method： GET

params : 无

response: 

    {
        "code": 0,
        "msg": "",
        "data": {
            "type": 1,
            "highest_block": 100,
            "last_block_time": 10,
            "node_num": 4,
            "tx_num": 1999,
        }
    }