# webhook_test

## make build ##
Сборка проекта

## make run ##
Запуск с настройками по-умолчанию из config/config.yml

## запуск с флагами ##
```
Flags:
  -a, --amount int    Requests amount (default 1000)
  -h, --help          help for send
  -r, --rps int       Requests per second limit (default 100)
  -u, --url string    Webhook URL (default "http://yandex.ru")
  -w, --workers int   Workers number (default 1)
```

## make test ##
Запуск автотестов

## make lint ##
Запуск линтера

## dc_build ##
Сбора Docker-образа проекта

## dc_run ##
Запуск собранного Docker-образа
