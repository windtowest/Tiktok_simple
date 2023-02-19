package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"tiktok_consul/client"
	"tiktok_consul/config"
)

func registerTest() {
	// 实例化一个 Consul 客户端
	// 配置consul中心的ip地址
	consulClient := client.New(config.ConsulIP, config.ConsulPort)
	// 实例失败，停止服务
	if consulClient == nil {
		panic(0)
	}

	// 需要先通过 uuid 为服务创建一个服务实例ID
	instanceId := uuid.New().String()
	logger := log.New(os.Stderr, "", log.LstdFlags)
	// 服务注册
	if !consulClient.Register("Test", instanceId, "/health", 10086, nil, logger) {
		// 注册失败，服务启动失败
		panic(0)
	}
}

func discoverTest() {
	consulClient := client.New(config.ConsulIP, config.ConsulPort)
	if consulClient == nil {
		panic(0)
	}
	// 根据注册时的服务名 进行服务查找
	rsp := consulClient.DiscoverServices("Test")
	/*
		自行解析想要的数据
		http://xxx.xxx.xxx.xxx:8500/v1/health/service/服务名称 查看json结构
	*/
	fmt.Println(rsp)
}

func deregisterTest() {
	consulClient := client.New(config.ConsulIP, config.ConsulPort)
	if consulClient == nil {
		panic(0)
	}
	// instanceId字符串内容 改为DiscoverServices中获取到的服务实例uuid
	consulClient.DeRegister("08b3923d-f6b0-4acb-aacf-327dc2415478")
}

func main() {
	// 三个使用示例
	registerTest()
	//discoverTest()
	//deregisterTest()
}
