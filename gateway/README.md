# Eventify Gateway Service

Gateway сервис для объединения всех gRPC сервисов Eventify через единый HTTP интерфейс.

## Описание

Gateway сервис предоставляет единую точку входа для всех API запросов к сервисам Eventify. Он проксирует HTTP запросы к соответствующим gRPC сервисам, используя gRPC-Gateway.

## Поддерживаемые сервисы

- **Auth Service** - аутентификация и авторизация
- **Event Service** - управление событиями
- **User-Interact Service** - взаимодействие пользователей (отзывы, регистрации)

## API Endpoints

### Auth Service (публичные)
- `POST /auth/register` - регистрация пользователя
- `POST /auth/login` - вход в систему

### Event Service (защищенные)
- `POST /events` - создание события
- `GET /events` - получение списка событий
- `GET /events/{id}` - получение события по ID

### User-Interact Service (защищенные)
- `POST /user-interact` - создание отзыва
- `GET /user-interact/event/{event_id}` - получение отзывов по событию
- `PUT /user-interact/{review_id}` - обновление отзыва
- `DELETE /user-interact/{review_id}` - удаление отзыва
- `POST /registration/event/{event_id}` - регистрация на событие
- `DELETE /registration/event/{event_id}` - отмена регистрации
- `GET /registration/event/{event_id}` - получение списка участников

### Health Check (публичный)
- `GET /health` - проверка состояния сервиса

## Аутентификация

Gateway использует JWT токены для аутентификации. Все запросы, кроме auth endpoints, требуют валидный токен в заголовке `Authorization`.

### Публичные endpoints (без токена)
- `POST /auth/register` - регистрация пользователя
- `POST /auth/login` - вход в систему
- `GET /health` - проверка состояния сервиса

### Защищенные endpoints (требуют токен)
Все остальные endpoints требуют токен в формате:
```
Authorization: Bearer <your-jwt-token>
```

Токен получается при успешном входе через `/auth/login`.

## Установка и запуск

### Предварительные требования
- Go 1.21+
- Запущенные gRPC сервисы (auth, event, user-interact)

### Установка зависимостей
```bash
cd gateway
make deps
```

### Запуск в режиме разработки
```bash
make dev
```

### Сборка и запуск
```bash
make build
make run
```

## Конфигурация

Конфигурация сервиса находится в файле `configs/gateway.yaml`:

```yaml
server:
  port: 8080

services:
  auth:
    host: "localhost"
    port: 9091
  
  event:
    host: "localhost"
    port: 9092
  
  user-interact:
    host: "localhost"
    port: 9093
```

### Переменные окружения
- `CONFIG_PATH` - путь к файлу конфигурации (по умолчанию: `configs/gateway.yaml`)

## Обработка ошибок

Gateway предоставляет единообразную обработку ошибок для всех сервисов. Ошибки возвращаются в формате JSON:

```json
{
  "error": "описание ошибки",
  "code": 400
}
```

### Коды ошибок аутентификации
- `401` - отсутствует заголовок Authorization
- `401` - неверный формат заголовка Authorization
- `401` - пустой токен

## Логирование

Сервис использует структурированное логирование в формате JSON с помощью logrus.

## Разработка

### Структура проекта
```
gateway/
├── cmd/
│   └── main.go          # Точка входа
├── internal/
│   ├── config/
│   │   └── config.go    # Конфигурация
│   └── service/
│       └── gateway.go   # Основная логика
├── go.mod
├── Makefile
└── README.md
```

### Команды для разработки
```bash
make test    # Запуск тестов
make lint    # Проверка кода
make clean   # Очистка артефактов сборки
```

## Примеры использования

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "role": "user"
  }'
```

### Вход в систему
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### Создание события (с токеном)
```bash
curl -X POST http://localhost:8080/events \
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