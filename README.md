# go-proxy

## Запуск

Сгенерировать сертификаты и положить в локальную папку `ssl`  

```sh
sudo make cert
```

Запустить стенд  

```sh
docker-compose -f deployment/docker-compose.yml up -d --build
```
