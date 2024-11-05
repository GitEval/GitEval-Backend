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
docker tag pull crpi-vgud82zncz8nwfuc.cn-hangzhou.personal.cr.aliyuncs.com/qianchengsijin4869/giteval:llm llm:v1
```

3、执行`docker-compose up -d`运行，执行前请确保拉取（或者构建）`llm`服务的镜像

