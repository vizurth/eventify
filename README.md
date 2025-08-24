# 🎉 Eventify

**Eventify** — это микросервисное веб-приложение для управления событиями, построенное на Go с использованием современных технологий и архитектурных паттернов.

## 📋 Описание проекта

Eventify позволяет пользователям:
- 🔐 Регистрироваться и аутентифицироваться
- 📅 Создавать и просматривать события
- 👥 Регистрироваться на события и оставлять отзывы
- 🔔 Получать уведомления через Kafka
- 🌐 Использовать единый API через Gateway сервис


### Микросервисы:

- **Auth Service** (Port: 9091) — аутентификация и управление пользователями
- **Event Service** (Port: 9092) — создание и управление событиями
- **User-Interact Service** (Port: 9093) — отзывы и регистрация на события
- **Notification Service** (Port: 9095) — отправка уведомлений через Kafka
- **Gateway Service** (Port: 9097) — единая точка входа для всех API запросов

## 🛠️ Технологический стек

### Основные технологии
- **Язык программирования**: Go 1.24.0
- **База данных**: PostgreSQL с pgx драйвером
- **Конфигурация**: cleanenv + YAML
- **Логирование**: Zap (Uber)

### Межсервисное взаимодействие
- **gRPC** — для внутреннего взаимодействия сервисов
- **gRPC-Gateway** — HTTP прокси для gRPC сервисов
- **Kafka** — асинхронная обработка событий

### Дополнительные технологии
- **JWT**: golang-jwt для аутентификации
- **Kafka**: segmentio/kafka-go для асинхронных уведомлений
- **Хеширование**: golang.org/x/crypto для паролей
- **Контейнеризация**: Docker & Docker Compose

### Инфраструктура
- **Docker** — контейнеризация
- **Docker Compose** — оркестрация сервисов
- **Goose** — миграции базы данных
- **golangci-lint** — статический анализ кода

## 📁 Структура проекта

```
eventify/
├── auth/                          # Сервис аутентификации
│   ├── api/                       # gRPC API определения
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── config/
│       ├── handler/
│       ├── models/
│       ├── repository/
│       └── service/
├── event/                         # Сервис событий
│   ├── api/                       # gRPC API определения
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── config/
│       ├── handler/
│       ├── models/
│       ├── repository/
│       └── service/
├── user-interact/                 # Сервис взаимодействий
│   ├── api/                       # gRPC API определения
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── config/
│       ├── handler/
│       ├── models/
│       ├── repository/
│       └── service/
├── gateway/                       # Gateway сервис
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── middleware/            # HTTP middleware
│   │   └── service/
│   ├── go.mod
│   ├── Makefile
│   ├── README.md
│   └── test_gateway.sh
├── notification/                  # Сервис уведомлений
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── service/
│   │   └── wsserver/              # WebSocket сервер
│   ├── web/                       # WebSocket клиент
│   │   └── html/
│   ├── go.mod
│   └── Makefile
├── common/                        # Общие компоненты
│   ├── jwt/
│   ├── kafka/
│   ├── logger/
│   ├── postgres/
│   └── retry/
│
├── configs/                       # Конфигурационные файлы
│   └── config.yaml
├── migrations/                    # SQL миграции
├── build/                         # Артефакты сборки
│   └── docker/
├── third_party/                   # Внешние зависимости
└── Makefile                       # Основной Makefile
```

## 🚀 Запуск проекта

### Предварительные требования
- Go 1.24.0+
- PostgreSQL 14+
- Kafka 3.0+
- Docker & Docker Compose
- Make

### 🚀 Установка и запуск

#### 1. Клонирование репозитория
```bash
git clone https://github.com/eventify.git
cd eventify
```

#### 2. Установка зависимостей
```bash
# Установка инструментов разработки
make install-deps
make install-golangci-lint

# Установка зависимостей для всех сервисов
cd auth && go mod tidy
cd ../event && go mod tidy
cd ../user-interact && go mod tidy
cd ../gateway && go mod tidy
cd ../common && go mod tidy
```

#### 3. Настройка окружения
```bash
# Копирование конфигурации
cp configs/config.yaml.example configs/config.yaml

# Настройка переменных окружения
cp configs/.env.example configs/.env
```

#### 4. Запуск с Docker Compose
```bash
# Запуск всех сервисов
make up

# Просмотр логов
make logs

# Остановка сервисов
make down
```

#### 5. Локальная разработка
```bash
# Запуск базы данных
docker-compose -f build/docker/docker-compose.yaml up -d postgres kafka

# Применение миграций
make local-migration-up

# Запуск сервисов локально
make run-all
```

## 🔧 Команды Make

### Основные команды
```bash
make up              # Запуск всех сервисов через Docker
make down            # Остановка всех сервисов
make restart         # Перезапуск сервисов
make logs            # Просмотр логов
```

### Разработка
```bash
make lint            # Проверка кода
make install-deps    # Установка зависимостей
make install-golangci-lint  # Установка линтера
```

### Миграции
```bash
make local-migration-status  # Статус миграций
make local-migration-up      # Применение миграций
make local-migration-down    # Откат миграций
```

### Gateway сервис
```bash
make gateway-build    # Сборка gateway
make gateway-run      # Запуск gateway
make gateway-dev      # Разработка gateway
make gateway-clean    # Очистка gateway
```

### Все сервисы
```bash
make build-all        # Сборка всех сервисов
make run-all          # Запуск всех сервисов локально
```

## 📡 API Endpoints

### 🔐 Auth Service (9091) - Публичные endpoints
- **POST /auth/register** - регистрация пользователя
- **POST /auth/login** - вход в систему

### 📅 Event Service (9092) - Защищенные endpoints
- **POST /events** - создание события
- **GET /events** - список событий
- **GET /events/{id}** - получение события по ID

### 👥 User-Interact Service (9093) - Защищенные endpoints
- **POST /user-interact** - создание отзыва
- **GET /user-interact/event/{event_id}** - отзывы события
- **PUT /user-interact/{review_id}** - обновление отзыва
- **DELETE /user-interact/{review_id}** - удаление отзыва
- **POST /registration/event/{event_id}** - регистрация на событие
- **DELETE /registration/event/{event_id}** - отмена регистрации
- **GET /registration/event/{event_id}** - участники события

### 🌐 Gateway Service (9097) - Единая точка входа
Все вышеперечисленные endpoints доступны через Gateway на порту 9097.

### 🔔 Notification Service (9095)
- WebSocket соединения для real-time уведомлений
- Kafka интеграция для асинхронных уведомлений

## 🔐 Аутентификация

Проект использует JWT токены для аутентификации:

### Публичные endpoints (без токена)
- `POST /auth/register`
- `POST /auth/login`

### Защищенные endpoints (требуют токен)
Все остальные endpoints требуют заголовок:
```
Authorization: Bearer <your-jwt-token>
```

## 📡 Примеры запросов к API

### 🔐 Auth Service (через Gateway: 9097)

#### Регистрация пользователя
```bash
curl -X POST http://localhost:9097/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "ivan123",
    "email": "example@gmail.com",
    "password": "12345678",
    "role": "admin"
  }'
```

#### Вход в систему
```bash
curl -X POST http://localhost:9097/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "example@gmail.com",
    "username": "ivan123",
    "password": "12345678"
  }'
```

### 📅 Event Service (через Gateway: 9097)

#### Создание события
```bash
curl -X POST http://localhost:9097/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "title": "Tech Conference 2025",
    "description": "Международная конференция по новым технологиям.",
    "category": "Технологии",
    "location": {
      "address": "ул. Примерная, д. 1",
      "city": "Москва",
      "country": "Россия"
    },
    "start_time": "2025-09-01T10:00:00Z",
    "end_time": "2025-09-01T18:00:00Z",
    "organizer": {
      "id": 10,
      "name": "Иван Иванов",
      "contact": "ivan@example.com"
    },
    "participants": [
      {
        "id": 1,
        "name": "Алексей Смирнов",
        "email": "alex@example.com"
      }
    ],
    "status": "active"
  }'
```

#### Получение списка событий
```bash
curl -X GET http://localhost:9097/events \
  -H "Authorization: Bearer <your-jwt-token>"
```

#### Получение события по ID
```bash
curl -X GET http://localhost:9097/events/101 \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 👥 User-Interact Service (через Gateway: 9097)

#### Создание отзыва
```bash
curl -X POST http://localhost:9097/user-interact \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "event_id": 1,
    "user_id": 5,
    "username": "ivan",
    "rating": 5,
    "comment": "Отличное событие!"
  }'
```

#### Получение отзывов события
```bash
curl -X GET http://localhost:9097/user-interact/event/1 \
  -H "Authorization: Bearer <your-jwt-token>"
```

#### Регистрация на событие
```bash
curl -X POST http://localhost:9097/registration/event/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "event_id": 101,
    "user_id": 5,
    "username": "ivan_user"
  }'
```

### Запуск тестов для всех сервисов
```bash
# Запуск тестов для каждого сервиса
cd auth && go test ./...
cd ../event && go test ./...
cd ../user-interact && go test ./...
cd ../gateway && go test ./...
```

## 📊 База данных

### Основные таблицы
- **users** - пользователи системы
- **events** - события
- **event_participants** - участники событий
- **reviews** - отзывы на события

### Миграции
Проект использует Goose для управления миграциями:
```bash
# Применение миграций
make local-migration-up

# Откат миграций
make local-migration-down

# Статус миграций
make local-migration-status
```

## 🐛 Отладка

### Логи
```bash
# Просмотр логов всех сервисов
make logs

# Логи конкретного сервиса
docker-compose -f build/docker/docker-compose.yaml logs -f auth-service
```

### Подключение к базе данных
```bash
docker exec -it eventify-postgres-1 psql -U eventify-user -d eventify
```

### Проверка Kafka
```bash
docker exec -it eventify-kafka-1 kafka-topics --list --bootstrap-server localhost:9095
```

## 📦 Формат ответов

Все ответы — в формате JSON, ошибки также возвращаются в JSON:

```json
{
  "error": "Invalid credentials",
  "code": 401
}
```

## ✅ Коды состояния

| Код | Значение                  |
| --- | ------------------------- |
| 200 | Успешно                   |
| 201 | Ресурс создан             |
| 400 | Ошибка в запросе          |
| 401 | Неавторизован             |
| 404 | Не найдено                |
| 500 | Внутренняя ошибка сервера |

## 🔄 CI/CD

### Линтинг
```bash
make lint
```

### Сборка
```bash
make build-all
```

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add amazing feature'`)
4. Отправьте в ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📝 Лицензия

Этот проект распространяется под лицензией MIT. См. файл `LICENSE` для получения дополнительной информации.

## 📞 Поддержка

Если у вас есть вопросы или предложения:
1. Проверьте [Issues](https://github.com/your-repo/eventify/issues)
2. Создайте новое Issue с подробным описанием проблемы
3. Обратитесь к команде разработки

## 🔮 Планы развития

- [ ] Добавление WebSocket уведомлений
- [ ] Интеграция с платежными системами
- [ ] Мобильное приложение
- [ ] Аналитика и отчеты
- [ ] Интеграция с социальными сетями
- [ ] Система рекомендаций
- [ ] Многоязычная поддержка
- [ ] Rate limiting для API
- [ ] Кэширование с Redis
- [ ] Мониторинг с Prometheus + Grafana

---

**Eventify** — Создавайте незабываемые события! 🎉
