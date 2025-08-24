# Gateway Service

API Gateway для микросервисной архитектуры Eventify. Объединяет все микросервисы под одним HTTP портом с аутентификацией.

## Функциональность

- **Единая точка входа** - все запросы проходят через порт 9090
- **JWT аутентификация** - проверка токенов для защищенных endpoints
- **HTTP клиенты** - прямые запросы к микросервисам без proxy
- **Health Check** - endpoint для проверки состояния сервиса

## Структура роутинга

### Публичные endpoints (без аутентификации)
- `POST /auth/register` - регистрация пользователя
- `POST /auth/login` - вход в систему

### Защищенные endpoints (требуют JWT токен)
- `GET /events/` - список событий
- `POST /events/` - создание события
- `GET /events/{id}` - получение события по ID
- `POST /user-interact/` - создание отзыва
- `GET /user-interact/event/{event_id}` - отзывы по событию
- `PUT /user-interact/{review_id}` - обновление отзыва
- `DELETE /user-interact/{review_id}` - удаление отзыва
- `POST /registration/event/{event_id}` - регистрация на событие
- `DELETE /registration/event/{event_id}` - отмена регистрации
- `GET /registration/event/{event_id}` - список участников

### Системные endpoints
- `GET /health` - проверка состояния сервиса

## Аутентификация

Для доступа к защищенным endpoints необходимо передавать JWT токен в заголовке:

```
Authorization: Bearer <your-jwt-token>
```

Токен получается при успешном входе через `/auth/login`.

## Запуск

```bash
go run ./gateway/cmd/main.go
```

Сервис запустится на порту 9090.

## Конфигурация

Настройки в `configs/config.yaml`:

```yaml
gateway:
  port: 9090
  shutdown_timeout: 30s

auth:
  secret_key: "your-secret-key"
  url: "http://auth:9091"

event:
  url: "http://event:9092"

user-interact:
  url: "http://user-interaction:9093"

notification:
  url: "http://notification:9095"
```

## Архитектура

Gateway использует прямые HTTP клиенты для связи с микросервисами:

```
Client → Gateway (9090) → Auth Service (9091)
                      → Event Service (9092)  
                      → User-Interaction Service (9093)
                      → Notification Service (9095)
```

## Примеры использования

### Регистрация
```bash
curl -X POST http://localhost:9090/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "role": "user"
  }'
```

### Вход
```bash
curl -X POST http://localhost:9090/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### Создание события (с токеном)
```bash
curl -X POST http://localhost:9090/events/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "title": "Test Event",
    "description": "Test Description",
    "category": "test",
    "start_time": "2024-01-01T10:00:00Z",
    "end_time": "2024-01-01T12:00:00Z"
  }'
```

## Тестирование

Запустите тестовый скрипт:
```bash
chmod +x gateway/test_simple.sh
./gateway/test_simple.sh
``` 