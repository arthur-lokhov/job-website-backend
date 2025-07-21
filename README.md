# Job Website Backend

## Запуск локально

1. Установите Go 1.21+, Docker, PostgreSQL 16+
2. Скопируйте `config/config.yaml` и настройте параметры
3. Запустите миграции:

```
make migrate
```

4. Запустите приложение:

```
make run
```

## Запуск через Docker

```
make docker-up
```

## Структура проекта

- cmd/api/main.go — запуск HTTP API
- internal/handlers — обработчики HTTP
- internal/middleware — middleware (авторизация)
- internal/models — модели данных
- internal/repository — работа с БД
- internal/services — бизнес-логика
- internal/config — конфиги
- migrations — SQL-миграции
- config — конфиги и ключи

## Аутентификация

Аутентификация через https://authservice.cyberzone.dev/docs/

---

_Подробности см. в ТЗ и комментариях к коду._