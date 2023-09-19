**kratos-example 一个用户管理的小例子**

* kratos.bff.api - bff层, api接口, 允许http/rpc调用
* kratos.user.service - 用户服务, 仅允许rpc调用

流程: 浏览器请求->bff.GetUserInfo->user.GetUser

## Uses

| 组件名        | 介绍             | 链接                                          |
|------------|----------------|---------------------------------------------|
| gorm       | ORM库           | https://github.com/go-gorm/gorm             |
| jwt        | JWT库           | https://github.com/golang-jwt/jwt           |
| zap        | 日志库            | https://github.com/uber-go/zap              |
| go-redis   | Redis库         | https://github.com/redis/go-redi            |
| consul     | 服务注册与发现 & 配置中心 | https://github.com/hashicorp/consul         |
| jaeger     | 链路追踪           | https://github.com/jaegertracing/jaeger     |
| prometheus | 服务监控           | https://github.com/prometheus/client_golang |

## 快速启动

```shell
# 拉取仓库
git clone git@github.com:Gekkoou/kratos-example.git

# 初始化
cd kratos-example
make init

# 启动微服务bff与user & 环境服务 (mysql/redis/consul/jaeger/prometheus/grafana)
docker-compose up -d

# 测试
curl 'http://127.0.0.1:8000/v1/user/info?id=1'

# 如开启单个服务, 在项目根目录下执行
kratos run
```

Consul 面板 http://localhost:8500/

Jaeger UI http://localhost:8500/

Grafana 仪表板 http://localhost:3000/