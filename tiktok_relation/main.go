package main

import (
	"Tiktok_simple/config"
	"Tiktok_simple/dao"
	relation_gorm "Tiktok_simple/kitex_gen/relation_gorm/userservice"
	"Tiktok_simple/middleware/rabbitmq"
	"Tiktok_simple/middleware/redis"

	"github.com/cloudwego/kitex/server"

	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

func main() {

	initDeps()

	//// 返回的 是 registry.Registry 类型，是用于向 Consul 注册和注销的方法
	//checkId := uuid.New().String()
	//r, err := consul.NewConsulRegister(config.ConsulAddress, consul.WithCheck(&consulapi.AgentServiceCheck{
	//	CheckID:                        checkId,
	//	Interval:                       "7s",
	//	HTTP:                           "http://" + "101.34.4.141" + ":" + "8801" + "/health",
	//	Timeout:                        "5s",
	//	DeregisterCriticalServiceAfter: "1m",
	//	Name:                           "relation.server",
	//}))
	//
	//if err != nil {
	//	panic(err)
	//}
	//addr, _ := net.ResolveTCPAddr("tcp", ":8801")
	//
	//var opts []server.Option
	////opts = append(opts, server.WithServiceAddr(addr), server.WithRegistry(r))
	//opts = append(opts, server.WithServiceAddr(addr), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
	//	ServiceName: "relation.server",
	//}), server.WithRegistry(r))

	// 返回的 是 registry.Registry 类型，是用于向 Etcd 注册和注销的方法
	r, err := etcd.NewEtcdRegistry([]string{config.EtcdAddress})
	if err != nil {
		panic(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", ":8801")

	var opts []server.Option
	opts = append(opts, server.WithServiceAddr(addr), server.WithRegistry(r))

	svr := relation_gorm.NewServer(new(UserServiceImpl), opts...)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

// 加载项目依赖
func initDeps() {
	// 初始化数据库
	dao.Init()

	// 初始化redis-DB0的连接，follow选择的DB0.
	redis.InitRedis()
	// 初始化rabbitMQ。
	rabbitmq.InitRabbitMQ()
	// 初始化Follow的相关消息队列，并开启消费。
	rabbitmq.InitFollowRabbitMQ()

}
