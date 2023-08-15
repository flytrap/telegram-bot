version=$1
registry=registry.cn-hangzhou.aliyuncs.com/flytrap
echo "build: $registry/telegram-bot:$version"
docker build -t $registry/telegram-bot:$version .

docker push $registry/telegram-bot:$version
