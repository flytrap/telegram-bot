# telegram-bot

tg bot

## 生成命令

### 依赖注入

```bash
wire ./internal
```

### 生成 bz 包文件

```bash
protoc -I=pb --go_out=pb/v1 --go_opt=paths=source_relative --go-grpc_out=pb/v1 --go-grpc_opt=paths=source_relative pb/tg.proto

```

### run

```bash
docker run -it --rm -p 7007:7007 -v $PWD/config:/app/config tg-bot:0.0.1
```
