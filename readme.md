# 随机图片 api -- go 实现

- 采用的是 go 来实现的, 采用 gin 框架, 处理图像采用 bild 包,只是使用了调整大小的功能. ~~原本使用 bimg 包, 但是依赖太严重了(因此不能交叉编译)~~

- 图片来源于互联网, 主要来自 c 站

## 使用方法 ( docker 安装)

```yaml
services:
  randomimages:
    container_name: randomimages-server
    image: qqzsc/randomimages:latest
    # image: ghcr.io/zoushucai/randomimages:latest
    restart: always
    ports:
      - "9113:9113"
    volumes:
      # 需要准好图片资源文件夹
      - ./../assets:/app/assets
      - ./logs:/app/logs
```

然后运行 `sudo docker-compose up` 就可以通过 `http://localhost:9113/v1/random` 访问了

## 参数说明

- 假设: `url = http://localhost:9113/v1/random`

- 访问 `url?sub=xxx&width=xxx&height=xxx&index=xxx` 即可获取指定的图片,(图片的类型包括: `jpeg, png, gif`)

  - `sub`: 文件夹名称, 例如 `url?sub=cat` 即可获取 `assets/cat` 文件夹下的图片
  - `width`: 图片宽度, 默认为 0(pc 端), 非 pc 端默认为 300.
  - `height`: 图片高度, 默认为 0

    - 如果 `width` 和 `height` 都为 0, 则默认为原始图片的宽度
    - 如果二者只要一个为 0, 则调整为指定的宽高(保持比例)
    - 如果都大于 0, 则调整到指定的宽高(拉伸)

  - `index`: 索引, 例如 `url?index=1` 即可获取索引为 1 的图片(按顺序), 如果未指定, 则默认为随机获取

  - 上述选项是可以组合使用的

- 新支持 webp, 由于 webp 的格式问题, 只能访问,不能调整大小

```bash

# example -- 直接返回图片
http://localhost:9113/v1/random?sub=dongman&index=0


# 如果要返回 json,
curl -H "Accept: application/json" http://localhost:9113/v1/random

# 返回
{"file":"assets/renwu/00059-3052920575-masterpiece.png","message":"success"}
# 这里的 file 是图片的 url, 直接把 assets 目录下的所有图片进行了挂载, 因此可以通过这个地址直接访问图片
# message 是返回的提示信息


# 还可以通过 swagger 来查看接口文档(这个会使得编译后的体积大 3 倍,算了,还是添加上, 方便了解 api)
# http://localhost:9113/swagger/index.html

```

- 还提供图片上传的功能, 更多 api 细节参考: [http://localhost:9113/swagger/index.html](http://localhost:9113/swagger/index.html)

## 编译安装

当然使用编译安装, 直接采用二进制文件部署,肯定体积更小

1. 首先,准备好 go 环境, 以及必要的依赖.

```bash
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -xzvf go1.23.0.linux-amd64.tar.gz -C /usr/local
echo 'export PATH="$PATH:/usr/local/go/bin"' >> ~/.zshrc
source ~/.zshrc
echo 'export PATH=$(go env GOPATH)/bin:$PATH' >> ~/.zshrc
source ~/.zshrc
rm go1.23.0.linux-amd64.tar.gz

```

2. 准备好图片资源, 把图片放到 `assets` 目录下, 如果要对图片进行分类, 可以在 `assets` 目录下新建文件夹, 把图片放到对应的文件夹下. 例如 `assets/cat` 文件夹下的图片, 就是猫的图片.

3. 准备好 go 环境, 并且安装好 go. 然后就是一些列的初始化工作

```bash
# 1.初始化项目
go mod init mygo

# 2. 下载依赖
go mod tidy

# 3. 编译启动
go run main.go

# 4. 如果不报错, 就可以通过 http://localhost:9113/v1/random 访问图片

# 5. 就可以打包成二进制文件
go build -o main main.go

# 6. 通过二进制文件启动
./main

# 7. 编译的时候添加 -ldflags="-s -w" 可以减少二进制文件大小
go build -ldflags="-s -w" -o main main.go


# 8添加参数, 指定端口和图片路径
./main --port=9113 --imgdir=imgs

```

### 关于 go 的 交叉编译

不想折腾, 推荐第三方程序: [goreleaser](https://github.com/goreleaser/goreleaser)

- 注意: 如果项目底层依赖 c 库,则可能交叉编译的时候会失败

## 关于 go 的热重载

- [air](https://github.com/air-verse/air)

## 关于 api 文档 -- swagger

- [swag](https://github.com/swaggo/swag)
- [gin-swagger](https://github.com/swaggo/gin-swagger)

## 利用 action 进行自动编译

- 参考本仓库的 `.github/workflows` 文件

## 增加 webp 的支持

- 为了不引入 cgo, 因为对图片的裁剪是不成功的, 由于涉及图片拷贝等操作, 镜像 scratch 难以胜任, 因此选 alpine:latest 为基础镜像

## api 接口

- 直接从 swagger 进行查看

- 地址: http://localhost:9113/swagger/index.html
