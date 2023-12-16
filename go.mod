module github.com/lenny-mo/cart

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.5.0 // indirect
	github.com/lenny-mo/order v0.0.0-20231216052347-f82b8aff9430
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/config/source/consul/v2 v2.9.1
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.9.1
	github.com/miekg/dns v1.1.57 // indirect
	github.com/ngaut/log v0.0.0-20221012222132-f3329cba28a5
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.17.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	golang.org/x/tools v0.16.1 // indirect
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)
