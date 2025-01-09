# GinTalk 论坛项目

关于 GinTalk 论坛项目的介绍详细可以看[这里](./GinTalk/README.md),这里不多赘述

本项目的改进点在于服务发现方面,增加了 etcd-service-exporter 服务,用于将 etcd 服务注册到 prometheus 中

## 项目部署
项目的MySQL,Redis,etcd,kafka,prometheus使用了Docker部署,可以使用[docker-compose文件夹](./docker-compose)下的进行部署,在每个文件夹下执行
```shell
docker-compose up -d
```
即可启动所有服务

也可以自行部署这些服务

项目的具体配置信息在[配置文件](./conf/config.yaml)中,具体配置需求可以修改该文件

**注意:** 如果项目是在MacOS,Windows上运行,并且使用的是Docker Desktop,同时Prometheus是在Docker中运行,则无需修改.否则需要自行修改[配置文件](./conf/config.yaml)中的service_registry的地址,
否则prometheus将会无法监控到服务

## 项目运行
分别在`etcd-service-exporter`,`GinTalk`文件夹下执行
```shell
go run main.go
```