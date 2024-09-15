# myepg
### 简介
- 本人学习使用go写的本地epg api 
- 本项目只是学习测试使用，不得商用，所有的法律责任与后果应由使用者自行承担
#### 参考项目：
- Meroser's IPTV: https://github.com/Meroser/IPTV

#### 使用到的资源
- EPG电子节目单: https://epg.erw.cc/

### 使用
根据需要使用的系统编译成指定的执行文件

##### Linux下编译
```
# 编译mac执行文件
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
 
# 编译windows执行文件
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```
##### Windows下编译
```
# windows编译mac
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build main.go
 
# windows编译linux
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go
```
##### Mac下编译
```
# mac编译linux执行文件
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
 
# mac编译windows执行文件
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```
# linux做成服务
/etc/systemd/system/mypeg.service

```
[Unit]
Description=Mypeg (my)
After=syslog.target
After=network.target
After=mysqld.service
After=postgresql.service
After=memcached.service
After=redis.service
[Service]
Type=simple
ExecStart=/opt/mypeg/mypeg
[Install]
WantedBy=multi-user.target
```
设置好之后执行如下命令:
systemctl daemon-reload
systemctl start mypeg.service
systemctl status mypeg.service
systemctl enable mypeg.service

### 问题
- 还得继续学习优化

