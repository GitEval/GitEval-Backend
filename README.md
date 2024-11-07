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
│  Dockerfile		##构建镜像
│  go.mod	## go mod
│  go.sum
│  main.go	##程序的入口
│  README.md
│  wire.go	## wire文件，生成依赖
│  wire_gen.go	## 由wire.go生成
│
├─api	##存放api相关
│  ├─request	##请求相关结构体
│  │      auth.go	
│  │      user.go
│  │
│  ├─response	##响应相关结构体
│  │      response.go	
│  │
│  └─route	##路由
│          route.go
│
├─client	##请求llm
│  │  client.go
│  │  llm.go
│  │
│  ├─gen
│  │      llm.pb.go
│  │      llm_grpc.pb.go
│  │
│  └─proto
│          llm.proto
│
├─conf	##配置相关
│  │  conf.go	##配置相关结构体
│  │  config-example.yaml	##配置信息
│  │  setting.go	##viper
│
├─controller	## controller层
│      auth.go	##鉴权
│      controller.go	##依赖集合
│      user.go	##用户
│
├─docs	##接口文档
│      docs.go
│      swagger.json
│      swagger.yaml
│
├─middleware	##中间件
│      jwt.go	##中间件
│      middleware.go
│
├─model	##结构体
│  │  contactDAO.go
│  │  data.go
│  │  domain.go
│  │  domainDAO.go
│  │  model.go
│  │  type.go
│  │  user.go
│  │  userDAO.go
│  │
│  └─cache	##redis相关
│          cache.go
│          redis.go
│
├─pkg
│  │  pkg.go
│  │
│  ├─github	##请求github相关
│  │  │  github.go
│  │  │
│  │  └─expireMap	##封装sync.Map，有过期功能
│  │          map.go
│  │
│  └─tool	##工具
│          changeType.go
│
└─service	## 服务层
        auth.go
        service.go
        user.go

```

