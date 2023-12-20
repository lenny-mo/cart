package main

import (
	"fmt"
	"strconv"

	"github.com/lenny-mo/cart/conf"
	"github.com/lenny-mo/cart/domain/models"
	"github.com/lenny-mo/cart/global"
	"github.com/lenny-mo/cart/handler"
	"github.com/lenny-mo/emall-utils/tracer"

	"github.com/lenny-mo/cart/domain/dao"

	"github.com/lenny-mo/cart/domain/services"

	"github.com/lenny-mo/cart/proto/cart"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 配置中心
	// consulCof 是用于获取配置的对象，它连接到配置中心（在Consul上的路径 /micro/config），
	// 以获取应用程序的配置信息。这可以让你的微服务从配置中心获取配置，而不需要硬编码配置信息，
	// 从而实现了灵活性和集中化的配置管理。
	consulCof, err := conf.GetConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		fmt.Println(err)
		fmt.Println("获取配置失败")
		panic(err)
	}

	fmt.Println("获取配置成功")
	// 创建一个注册中心 后续用来注册我们的服务
	// 各个微服务需要注册自己的服务信息（例如服务名称、IP地址、端口等）到注册中心，
	// 并可以从注册中心查询其他微服务的信息。这有助于实现动态的服务发现和负载均衡。
	// 这部分代码使用了Consul作为注册中心，并创建了一个Consul注册器（Registry）
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	serviceName := "go.micro.service.cart"
	// 3 链路追踪
	err = tracer.InitTracer(serviceName, "127.0.0.1:6831")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tracer.Closer.Close()
	opentracing.SetGlobalTracer(tracer.Tracer)

	// 3. 创建服务
	global.GlobalRPCService = micro.NewService(
		micro.Name(serviceName),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8083"), // 服务监听地址
		// 使用consul注册中心
		micro.Registry(consulRegistry),
		// 添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// uber 漏桶 添加限流 每秒处理1000·个请求
		micro.WrapHandler(ratelimit.NewHandlerWrapper(conf.QPS)),
	)

	// 4. 获取mysql配置
	mysqlConf := conf.GetMysqlFromConsul(consulCof, "mysql")

	// 5. 初始化数据库连接
	dsn := mysqlConf.User + ":" + mysqlConf.Password + "@tcp(" + mysqlConf.Host + ":" + strconv.FormatInt(mysqlConf.Port, 10) + ")/" + mysqlConf.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 不需要手动关闭，程序退出时会自动关闭

	// 禁止复表的存在, 如果没有表则创建
	if !db.Migrator().HasTable(&models.Cart{}) {
		db.Migrator().CreateTable(&models.Cart{})
	}

	// 6. service初始化
	global.GlobalRPCService.Init()

	// 7. 创建service 和 handler 并且注册服务
	cartDAO := dao.NewCartDAO(db)
	cartService := services.NewCartService(cartDAO.(*dao.CartDAO))
	err = cart.RegisterCartHandler(global.GlobalRPCService.Server(), &handler.CartHandler{CartService: *cartService.(*services.CartService)})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// 8. 启动service
	if err = global.GlobalRPCService.Run(); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
