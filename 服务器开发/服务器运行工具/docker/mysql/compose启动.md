# docker-compose.yml 模板

1. 宿主机 3306 端口映射到容器 3306 端口；

2. 宿主机目录 ./mysql-data 挂载为 MySQL 数据目录；

3. 宿主机目录 ./mysql-conf 挂载为配置目录（可选，放自定义 *.cnf 文件）；

4. 宿主机目录 ./mysql-log 挂载为日志目录（可选）。

只要装有 Docker & Docker-Compose 即可 docker-compose up -d 一键启动。
文件结构
.
├── docker-compose.yml
├── mysql-data/      # 数据持久化目录（会自动创建）
├── mysql-conf/      # 自定义配置目录（可选）
└── mysql-log/       # 日志目录（可选）

docker-compose.yml

```yaml
version: "3.9"

services:
  mysql:
    image: mysql:8.0           # 想用 5.7 改成 mysql:5.7
    container_name: mysql
    restart: unless-stopped
    ports:
      - "3306:3306"             # 宿主机:容器
    environment:
      MYSQL_ROOT_PASSWORD: root # 按需改复杂密码
      MYSQL_DATABASE: demo      # 初始化时创建的数据库
      TZ: Asia/Shanghai
    volumes:
      - ./mysql-data:/var/lib/mysql     # 数据持久化
      - ./mysql-conf:/etc/mysql/conf.d  # 自定义配置（可选）
      - ./mysql-log:/var/log/mysql      # 日志挂载（可选）
    command: [
      "--character-set-server=utf8mb4",
      "--collation-server=utf8mb4_unicode_ci",
      "--default-authentication-plugin=mysql_native_password"
    ]
```

使用步骤

1. 把以上内容保存为 docker-compose.yml。

2. 确保 3306 端口未被宿主机其他进程占用：
    sudo lsof -i :3306 无输出即可。

3. 启动：
    docker-compose up -d
    首次会拉镜像，稍等片刻。

4. 连接测试：
    mysql -h127.0.0.1 -uroot -proot 能连说明 OK。

常用后续命令

- 关闭： docker-compose down

- 看日志： docker-compose logs -f mysql

- 升级镜像： docker-compose pull && docker-compose up -d

按需调整密码、数据库名、字符集或挂载路径即可。