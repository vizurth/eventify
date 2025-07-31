# 🎉 Eventify

**Eventify** — это микросервисное веб-приложение для управления событиями, построенное на Go с использованием современных технологий и архитектурных паттернов.

## 📋 Описание проекта

Eventify позволяет пользователям:
- 🔐 Регистрироваться и аутентифицироваться
- 📅 Создавать и просматривать события
- 👥 Регистрироваться на события и оставлять отзывы
- 🔔 Получать уведомления через Kafka

## 🏗️ Архитектура

Проект построен на микросервисной архитектуре и состоит из четырех основных сервисов:

- **Auth Service** — аутентификация и управление пользователями
- **Event Service** — создание и управление событиями
- **User Interaction Service** — отзывы и регистрация на события
- **Notification Service** — отправка уведомлений через Kafka

## 🛠️ Технологический стек

### Основные технологии
- **Язык программирования**: Go 1.24.0
- **Веб-фреймворк**: Gin
- **База данных**: PostgreSQL с pgx драйвером
- **Миграции**: golang-migrate
- **Конфигурация**: cleanenv + YAML
- **Логирование**: Zap (Uber)

### Дополнительные технологии
- **JWT**: golang-jwt для аутентификации
- **Kafka**: segmentio/kafka-go для асинхронных уведомлений
- **Хеширование**: golang.org/x/crypto для паролей
- **Контейнеризация**: Docker

## 📁 Структура проекта

```
eventify/
├── auth/                          # Сервис аутентификации
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── config/
│       ├── handler/
│       ├── models/
│       ├── repository/
│       └── service/
├── event/                         # Сервис событий
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── config/
│       ├── handler/
│       ├── models/
│       ├── repository/
│       └── service/
├── user-interaction/              # Сервис взаимодействий
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── config/
│       ├── handler/
│       ├── models/
│       ├── repository/
│       └── service/
├── notification/                  # Сервис уведомлений
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       └── handler/
├── common/                        # Общие компоненты
│   ├── jwt/
│   ├── kafka/
│   ├── logger/
│   └── postgres/
├── configs/                       # Конфигурационные файлы
│   └── config.yaml
├── migrations/                    # SQL миграции
│   ├── auth/
│   ├── event/
│   └── user-interact/
├── models/                        # Общие модели
│   ├── models.go
│   └── models.json
├── go.mod
└── go.sum
```

## 🚀 Запуск проекта

### Предварительные требования
- Go 1.24.0+
- PostgreSQL
- Kafka
- Docker (опционально)

### 🚀 Установка и запуск

```
# Клонировать репозиторий
git clone https://github.com/eventify.git
cd eventify

# Запуск через docker-compose
docker-compose up --build
```

Все сервисы будут запущены на портах:

- Auth: `8081`
- Event: `8082`
- User Interaction: `8083`
- Notification: kafka:9092

---

### 📡 Примеры запросов к API

#### 🔐 Auth Service (8081)

- **POST /auth/register**
  ```json
  {
  	"username": "ivan123",
  	"email": "example@gmail.com",
  	"password": "12345678",
  	"role": "admin"
  }
  ```
- **POST /auth/login**
  ```json
  {
  	"email": "example@gmail.com",
  	"username": "ivan123",
  	"password": "12345678"
  }
  ```

#### 📅 Event Service (8082)

- **POST /events/**
  ```json
  {
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
  }
  ```
- **GET /events/**  
   Ответ:
  ```json
  [
  	{
  		"id": 101,
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
  		"status": "active",
  		"created_at": "2025-06-23T12:00:00Z"
  	}
  ]
  ```
- **GET /events/{id}**  
   Ответ:
  ```json
  {
  	"id": 101,
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
  	"status": "active",
  	"created_at": "2025-06-23T12:00:00Z"
  }
  ```

#### 👥 User Interaction Service (8083)

- **POST /reviews/**
  ```json
  {
  	"event_id": 1,
  	"user_id": 5,
  	"username": "ivan",
  	"rating": 5,
  	"comment": "Отличное событие!"
  }
  ```
- **GET /reviews/event/1**
  ```json
  {
  	"id": 501,
  	"event_id": 101,
  	"user_id": 5,
  	"username": "ivan_user",
  	"rating": 4,
  	"comment": "Очень интересное событие!",
  	"created_at": "2025-06-22T15:30:00Z",
  	"updated_at": "2025-06-23T08:00:00Z"
  }
  ```
- **PUT /reviews/10**
  ```json
  {
  	"rating": 4,
  	"comment": "Хорошо, но не идеально"
  }
  ```
- **DELETE /reviews/10**
- **POST /registration/event/1**
  ```json
  {
  	"event_id": 101,
  	"user_id": 5,
  	"username": "ivan_user"
  }
  ```
- **DELETE /registration/event/1?user_id=5**
- **GET /registration/event/1**
  ```json
  [
  	{
  		"event_id": 101,
  		"user_id": 5,
  		"username": "ivan_user"
  	}
  ]
  ```

#### 🔔 Notification Service 
Работает через kafka(в дальнейшем станет доступна работа с email)
---

### 📦 Формат ответов

Все ответы — в формате JSON, ошибки также возвращаются в JSON:

```
{
  "error": "Invalid credentials"
}
```

---

### ✅ Коды состояния

| Код | Значение                  |
| --- | ------------------------- |
| 200 | Успешно                   |
| 201 | Ресурс создан             |
| 400 | Ошибка в запросе          |
| 401 | Неавторизован             |
| 404 | Не найдено                |
| 500 | Внутренняя ошибка сервера |

---

## 📄 Лицензия

Этот проект распространяется под лицензией MIT. См. файл `LICENSE` для получения дополнительной информации.

## 📞 Поддержка

Если у вас есть вопросы или предложения, создайте issue в репозитории проекта.
