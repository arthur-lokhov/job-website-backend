# Тестирование API (cURL)

## Публичные ручки

### Получить список вакансий
```
curl -X GET "http://localhost:8080/vacancies"
```

### Получить детальную информацию о вакансии
```
curl -X GET "http://localhost:8080/vacancies/{id}"
```

### Получить форму отклика на вакансию
```
curl -X GET "http://localhost:8080/vacancies/{id}/form"
```

### Отправить отклик на вакансию
```
curl -X POST "http://localhost:8080/applications" \
  -H "Content-Type: application/json" \
  -d '{
    "vacancyId": "{vacancy_id}",
    "answers": {"fullName": "Иван Иванов", "email": "test@example.com"}
  }'
```

---

## Приватные (админ) ручки

> Для всех запросов используйте заголовок:
> `-H "Authorization: Bearer $TOKEN"`

### Получить все вакансии (админ)
```
curl -X GET "http://localhost:8080/admin/vacancies" -H "Authorization: Bearer $TOKEN"
```

### Создать новую вакансию
```
curl -X POST "http://localhost:8080/admin/vacancies" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Backend Developer",
    "locationId": "{location_id}",
    "departmentId": "{department_id}",
    "levelId": "{level_id}",
    "info": "Описание вакансии...",
    "applicationForm": {"formSchema": {"fullName": {"type": "string"}}},
    "important": false,
    "priority": 1,
    "isActive": true
  }'
```

### Обновить вакансию
```
curl -X PATCH "http://localhost:8080/admin/vacancies/{id}" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Backend Developer Updated"
  }'
```

### Удалить вакансию
```
curl -X DELETE "http://localhost:8080/admin/vacancies/{id}" -H "Authorization: Bearer $TOKEN"
```

### Получить все отклики
```
curl -X GET "http://localhost:8080/admin/applications" -H "Authorization: Bearer $TOKEN"
```

### Обновить статус отклика
```
curl -X PATCH "http://localhost:8080/admin/applications/{id}" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "APPROVED"
  }'
```

### Удалить отклик
```
curl -X DELETE "http://localhost:8080/admin/applications/{id}" -H "Authorization: Bearer $TOKEN"
``` 