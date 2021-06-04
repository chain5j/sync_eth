# sync_eth

sync_eth是一款开源的以太坊同步工具，它具有简单、高效、易用等特点。

## 特性
- 支持ETH系列的区块同步
- 支持自定义clientIdentifier。即rpc中method的前缀更换
- 支持基础余额统计
- 支持erc20余额的统计
- 支持指定区块开始同步
- 支持指定chainId同步
- 支持交易、区块存储elasticsearch，余额信息存入mysql
- 支持指定地址的交易监听

## 项目地址
Github: https://github.com/chain5j/sync_eth

## 运行
- 直接运行：

```shell
go run main.go --config=./conf/config.yaml
```
- docker运行：

```
docker run \
-it \
--name sync_eth \
-v ./conf/config_local.yaml:/data/conf/config_local.yaml \
chain5j/sync_eth:v1.0.0 \
--config=/data/conf/config_local.yaml
```

## LICENSE
Please refer to [LICENSE](LICENSE) file.

Copyright@2020 chain5j