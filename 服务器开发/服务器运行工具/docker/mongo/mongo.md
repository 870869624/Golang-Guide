# mongo启动

```bash
docker run -d \
  --name mongo-light \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=root \  # 用户名
  -e MONGO_INITDB_ROOT_PASSWORD=123456 \# 密码
  -v mongo-data:/data/db \              # 数据持久化
  --restart unless-stopped \            # 开机自启
  mongo:6.0-alpine
```
