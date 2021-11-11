# config

配置文件加载、读取是一个服务启动时必不可少的关键步骤，config 是封装的简单、易用的配置管理工具，具有以下特点：
- 支持ini、yaml格式
- 多种指定配置文件位置方法
- 一次加载、内存缓存，读取速度快
- 支持配置项与指定结构体转化
- 支持自定义多配置文件加载、读取

## 安装
```shell script
go get github.com/hy0kl/gconfig
```

## 默认路径：
conf/conf.ini

## 修改配置路径
```go
gconfig.SetConfigPath("/home/dev/conf.ini")
```
> 重新设置路径后，下次再读取配置便会切换到新的配置源。

## 配置示例:
```ini
[goconfig]
name = goconfig
hosts = 127.0.0.1 127.0.0.2 127.0.0.3 

[goconfigStringMap]
name = goconfig
host = 127.0.0.1

[goconfigArrayMap]
name = goconfig1 goconfig2

[goconfigObject]
max=101
port=9099
rate=1.01
hosts=127.0.0.1 127.0.0.2
timeout=5s
```

## 读取示例：

```go
//获取指定section下指定key的值
func GetConf(sec, key string) string
//获取指定section下指定key的值，不存在则返回默认值
func GetConfDefault(sec, key, def string) string
//获取指定section下指定key的slice类型值
func GetConfs(sec, key string) []string
//获取指定section下所有配置，返回map[string]string类型
func GetConfStringMap(sec string) (ret map[string]string)
//获取指定section下所有配置，返回map[string][]string类型
func GetConfArrayMap(sec string) (ret map[string][]string)
//将指定section下配置与传入类型v进行转化
func ConfMapToStruct(sec string, v interface{}) error
```
