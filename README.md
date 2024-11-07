# GitEval-Backend

## 如何运行程序？
1、首先，将`conf/config-example.yaml`改为`config.yaml`，然后配置相关信息

2、构建镜像`docker build -t giteval:v1 .`
或者考虑直接从阿里云拉取本服务的镜像（已上传）

```sh
##拉取镜像
docker pull crpi-vgud82zncz8nwfuc.cn-hangzhou.personal.cr.aliyuncs.com/qianchengsijin4869/giteval:app
##修改tag
docker tag crpi-vgud82zncz8nwfuc.cn-hangzhou.personal.cr.aliyuncs.com/qianchengsijin4869/giteval:app giteval:v1
```

如果需要，也可以拉取`llm`服务的镜像（已上传）
```sh
##拉取镜像
docker pull crpi-vgud82zncz8nwfuc.cn-hangzhou.personal.cr.aliyuncs.com/qianchengsijin4869/giteval:llm
##修改tag
docker tag crpi-vgud82zncz8nwfuc.cn-hangzhou.personal.cr.aliyuncs.com/qianchengsijin4869/giteval:llm llm:v1
```

3、执行`docker-compose up -d`运行，执行前请确保拉取（或者构建）`llm`服务的镜像

## 项目结构

```
GitEval-Backend
├── api
│   ├── request
│   │   ├── auth.go
│   │   └── user.go
│   ├── response
│   │   └── response.go
│   └── route
│       └── route.go
├── client
│   ├── gen
│   │   ├── llm.pb.go
│   │   └── llm_grpc.pb.go
│   └── proto
│       ├── llm.proto
│       ├── client.go
│       └── llm.go
├── conf
│   ├── conf.go
│   ├── config.yaml
│   ├── config-example.yaml
│   └── setting.go
├── controller
│   ├── auth.go
│   ├── controller.go
│   └── user.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── middleware
│   ├── jwt.go
│   └── middleware.go
├── model
│   ├── cache
│   │   ├── cache.go
│   │   └── redis.go
│   ├── contactDAO.go
│   ├── data.go
│   ├── domain.go
│   ├── domainDAO.go
│   ├── model.go
│   ├── type.go
│   ├── user.go
│   └── userDAO.go
├── pkg
│   ├── github
│   │   ├── expireMap
│   │   │   └── map.go
│   │   └── github.go
│   └── tool
│       ├── changeType.go
│       └── pkg.go
├── service
│   ├── auth.go
│   ├── service.go
│   └── user.go
├── .dockerignore
├── .gitignore
├── Dockerfile
├── go.mod
├── main.go
├── README.md
├── wire.go
└── wire_gen.go

```

