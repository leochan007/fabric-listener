# fabric-listener

./general_listener -events-address=127.0.0.1:7053 -listen-to-rejections=true -events-from-chaincode=peersafe-yinchengpai

nodejs版本zmq安装(需要先安装libzmq4,见fabric/devenv目录pre-setup.sh)
git clone https://github.com/JustinTulloss/zeromq.node.git
手工修改binding.gyp library和include
npm install -g
增加环境变量NODE_PATH指向.npmrc配置的global/node_modules目录

通讯模型:
PUB/SUB模型，listener为zmq server,　其他程序(nodejs,client_demo(go))做为zmq的client来连接server,并订阅消息

消息定义：
１. 消息均为json格式

２. 3种不同的包头: notfy rejected cEvent
(1). notfy是一般的通知，包括block的任何通知都会在这里 
{
    "Block": {
        "stateHash": "7Lt0UqkLMSxZsqUU6uspT7961w/QXlWSJNGsSk7uUxyAXTeXSfB13Cxb3Gu+XdzQtjc2P8sk9MlwK0zbrwjh3g==",
        "previousBlockHash": "X/HI4H7AtAwxzw9z3kxLoYAOmdOHBOYQDp3iDyIMhaKBCaGTIrHrmSkIfipHZonITMIx2S5uZ8cl/bGx3ooCqg==",
        "nonHashData": {
            "localLedgerCommitTimestamp": {
                "seconds": 1477366255,
                "nanos": 799956629
            },
            "chaincodeEvents": [
                {
                    "chaincodeID": "peersafe-yinchengpai",
                    "txID": "92e62c89-283a-459b-b373-e22b1d30f38e",
                    "eventName": "YCPChainCode",
                    "payload": "eyJFdmVudFR5cGUiOjQsIkVycm9yTXNnIjoiUmVwbGljYXRlSUQifQ=="
                }
            ]
        }
    }
}
(2). rejected一般是错误信息，异常之类
{
    "Rejection": {
        "tx": {
            "type": 2,
            "chaincodeID": "EhRwZWVyc2FmZS15aW5jaGVuZ3BhaQ==",
            "payload": "Ck8IARIWEhRwZWVyc2FmZS15aW5jaGVuZ3BhaRozCgpjcmVhdGVVc2VyCgEwCgNsZW8KCjIwMTYtMTAtMDgKATAKB2FiYy5qcGcKBTEyMzQ1",
            "txid": "92e62c89-283a-459b-b373-e22b1d30f38e",
            "timestamp": {
                "seconds": 1477366254,
                "nanos": 788885110
            }
        },
        "errorMsg": "Transaction or query returned with failure: [In chaincode]ID \u0000 is replicated"
    }
}

(3). cEvent为程序主动抛出的信息
{
    "ChaincodeEvent": {
        "chaincodeID": "peersafe-yinchengpai",
        "txID": "92e62c89-283a-459b-b373-e22b1d30f38e",
        "eventName": "YCPChainCode",
        "payload": "eyJFdmVudFR5cGUiOjQsIkVycm9yTXNnIjoiUmVwbGljYXRlSUQifQ=="
    }
}

３. 主要是处理cEvent,这里消息里面包含了内部定义的消息类型加消息内容 
	在cEvent的payload里面是自定义的消息
	实际的消息格式如下：
	[cEvent]:{"EventType":4,"Content":"ReplicateID","RequestID":"12345"}
	Content为消息数据，json string格式
