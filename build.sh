version=$1
registry=hidden
echo "build: $registry/telegram-bot:$version"
docker build -t $registry/telegram-bot:$version .

docker push $registry/telegram-bot:$version
