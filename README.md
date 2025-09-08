# Backend для сайта с вакансиями

Это бэкенд для сайта с вакансиями, разработанный на Go. Он следует принципам чистой архитектуры, разделяя логику на обработчики (Handlers), сервисы (Services), репозитории (Repositories) и модели (Models). Приложение использует Docker и Docker Compose для контейнеризации, PostgreSQL в качестве основной базы данных и Valkey (форк Redis) для кэширования.

## Сборка и запуск

Для сборки и запуска приложения вам понадобятся установленные Docker и Docker Compose.

Все необходимые параметры конфигурации определены как переменные окружения в `docker-compose.yml` со значениями по умолчанию. Опциональный файл `.env` можно использовать для переопределения этих значений при локальной разработке.

1.  **Запуск приложения:**

    ```bash
    docker-compose up -d
    ```

    Эта команда соберет образ `job-website-backend:latest` (если он еще не собран) и запустит все определенные сервисы (базу данных, кэш и приложение).

2.  **Доступ к API:**

    API задокументировано с помощью Swagger. Вы можете получить доступ к Swagger UI по адресу `http://localhost:8080/swagger/index.html`.

    Ключевые эндпоинты API:
    *   `GET /api/v1/vacancies`: Получить список активных вакансий.
    *   `GET /api/v1/vacancies/{id}`: Получить вакансию по ID.
    *   `GET /api/v1/vacancies/{id}/form`: Получить форму отклика на вакансию.
    *   `POST /api/v1/applications`: Создать новый отклик.
    *   `GET /api/v1/departments`: Получить список отделов.
    *   `GET /api/v1/levels`: Получить список уровней.
    *   `GET /api/v1/locations`: Получить список локаций.

    Эндпоинты администратора (требуют аутентификации):
    *   `GET /api/v1/admin/vacancies`: Получить список всех вакансий.
    *   `POST /api/v1/admin/vacancies`: Создать новую вакансию.
    *   `PATCH /api/v1/admin/vacancies/{id}`: Обновить вакансию.
    *   `DELETE /api/v1/admin/vacancies/{id}`: Удалить вакансию.
    *   `GET /api/v1/admin/applications`: Получить список всех откликов.
    *   `PATCH /api/v1/admin/applications/{id}`: Обновить отклик.
    *   `DELETE /api/v1/admin/applications/{id}`: Удалить отклик.

## Соглашения по разработке

*   **Тестирование:** Для запуска тестов Go используйте команду:
    ```bash
    go test ./...
    ```
*   **Конфигурация:** Конфигурация управляется исключительно через переменные окружения. Приложение использует `viper` для загрузки этих переменных с явной привязкой для надежности. Значения по умолчанию указаны в `docker-compose.yml` и могут быть переопределены локальным файлом `.env`.
*   **Миграции базы данных:** Проект использует Goose для миграций базы данных. Миграции автоматически применяются при запуске приложения.
*   **Поддержка UUID:** Поддержка UUID в PostgreSQL включена через расширение `uuid-ossp`, которое настраивается в `scripts/init_dev.sql` для окружения разработки.

## Аутентификация и авторизация

Приложение интегрируется с внешним сервисом аутентификации (AuthService) для аутентификации и авторизации.

*   **Проверка публичного ключа:** Сервис получает публичный ключ от AuthService для проверки JWT-токенов. Этот ключ периодически обновляется.
*   **Проверка по черному списку:** JWT-токены проверяются по черному списку, который ведется в AuthService. Черный список также периодически обновляется.
*   **Права доступа:** Права пользователя запрашиваются у AuthService и кэшируются на 24 часа для снижения нагрузки на AuthService.
*   **API-токен:** Сервис использует собственный API-токен для взаимодействия с AuthService при выполнении административных задач (например, получение прав доступа).

## Примеры cURL-команд

### Публичные эндпоинты

**Получить все вакансии:**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies"
```

**Получить важные вакансии:**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies?important=true"
```

**Поиск вакансий по названию (например, "Engineer"):**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies?search=Engineer"
```

**Получить вакансию по ID (замените на реальный ID):**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies/YOUR_VACANCY_ID"
```

**Получить форму отклика на вакансию (замените на реальный ID):**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies/YOUR_VACANCY_ID/form"
```

**Создать новый отклик:**
```bash
curl -X POST "http://localhost:8080/api/v1/applications" \
     -H "Content-Type: application/json" \
     -d '{
           "vacancy_id": "YOUR_VACANCY_ID",
           "answers": {
             "fullName": "John Doe",
             "email": "john.doe@example.com",
             "phone": "123-456-7890"
           }
         }'
```

**Получить все отделы:**
```bash
curl -X GET "http://localhost:8080/api/v1/departments"
```

**Получить все уровни:**
```bash
curl -X GET "http://localhost:8080/api/v1/levels"
```

**Получить все локации:**
```bash
curl -X GET "http://localhost:8080/api/v1/locations"
```

### Эндпоинты администратора (требуется заголовок Authorization с JWT-токеном)

**Получить все вакансии администратора:**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/vacancies" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Получить все отклики администратора:**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/applications" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Поиск откликов администратора по ФИО (например, "John Doe"):**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/applications?search=John%20Doe" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```