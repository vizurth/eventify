### 🧾 Описание проекта

**Eventify** — это микросервисное веб-приложение, позволяющее пользователям создавать события, регистрироваться на них, оставлять отзывы и получать уведомления. Проект построен с использованием REST-архитектуры и разбит на четыре микросервиса:

- **Auth Service** — аутентификация и регистрация пользователей.
    
- **Event Service** — создание и просмотр событий.
    
- **User Interaction Service** — отзывы и регистрация на события.
    
- **Notification Service** — отправка и получение уведомлений.
    

---

### ⚙️ Используемые технологии

- **Язык программирования**: Go (Golang)
    
- **Базы данных**: PostgreSQL
    
- **Docker**: контейнеризация микросервисов
    
- **Nginx**: проксирование и маршрутизация
    
- **gRPC/HTTP**: взаимодействие между сервисами
    
- **YAML**: конфигурация
    
- **SQL миграции**: для инициализации базы данных
    

---

### 🗂️ Структура проекта
```
├── cmd
├── config
├── db
│   └── migrations
├── docker
├── internal
│   ├── authservice
│   ├── config
│   ├── eventservice
│   ├── middleware
│   ├── notification-service
│   └── userinteractionservice
├── models
├── nginx
└── pkg
    ├── logger
    └── postgres
```
    

---

### 🚀 Установка и запуск

```
# Клонировать репозиторий
git clone https://github.com/vizurth/eventify.git
cd eventify

# Запуск через docker-compose
docker-compose up --build
```

Все сервисы будут запущены на портах:

- Auth: `8081`
    
- Event: `8082`
    
- User Interaction: `8083`
    
- Notification: `8084`
    

---

### 📡 Примеры запросов к API

#### 🔐 Auth Service (8081)

- **POST /auth/register**
    
    ``` json
    {
        "username": "ivan123",
        "email": "example@gmail.com",
        "password": "12345678",
        "role": "admin"
    }
    ```
    
- **POST /auth/login**
    
    ``` json
    {
        "email": "example@gmail.com",
        "username": "ivan123",
        "password": "12345678"
    }
    ```
    

#### 📅 Event Service (8082)

- **POST /events/**
    
    ``` json
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
    ``` json
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
    
    ``` json
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
    
    ``` json
    { 
        "event_id": 1,
        "user_id": 5,
        "username": "ivan",
        "rating": 5,
        "comment": "Отличное событие!"
    }
    ```
- **GET /reviews/event/1**
    ``` json
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
    
    ``` json
    {
        "rating": 4,
        "comment": "Хорошо, но не идеально"
    }
    ```
    
- **DELETE /reviews/10**
    
- **POST /registration/event/1**
    
    ``` json
    {
        "event_id": 101,
        "user_id": 5,
        "username": "ivan_user"
    }
    ```
    
- **DELETE /registration/event/1?user_id=5**
    
- **GET /registration/event/1**
    ``` json
    [
        {
            "event_id": 101,
            "user_id": 5,
            "username": "ivan_user"
        }
    ]
    ```
    

#### 🔔 Notification Service (8084)

- **POST /notifications/send**
    
    ``` json
    {
        "user_id": 5,
        "message": "Вы зарегистрированы на событие Go Meetup"
    }
    ```
    
- **GET /notifications/user/5**
    ``` json
    [
        {
            "user_id": 5,
            "message": "Вы зарегистрированы на событие Go Meetup"
        }
    ]
    ```
    
- **PUT /notifications/7/read**
- **DELETE /notifications/7**

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

|Код|Значение|
|---|---|
|200|Успешно|
|201|Ресурс создан|
|400|Ошибка в запросе|
|401|Неавторизован|
|404|Не найдено|
|500|Внутренняя ошибка сервера|

---

### 📄 Лицензия

MIT License — свободное использование с указанием авторства.

