package tiktok_relation

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
