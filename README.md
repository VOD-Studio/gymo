# Gymo

[![Go](https://github.com/VOD-Studio/gymo/actions/workflows/go.yml/badge.svg)](https://github.com/VOD-Studio/gymo/actions/workflows/go.yml)

Aite backend!

## Build

```
make build
```

## Deploy

在运行之前，确保更新了需要的环境变量。

```bash
cp .env.example .env
```

### Docker compose

```
docker compose up --build
```

### Binary

```
make build
./gymo
```

## Requirements

-   PostgreSQL: 保存所有数据
-   Redis: 聊天以及缓存

## API 文档

https://gymogymo.apifox.cn
