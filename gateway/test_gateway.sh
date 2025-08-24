#!/bin/bash

# Тестовый скрипт для Gateway Service
# Убедитесь, что все микросервисы запущены

GATEWAY_URL="http://localhost:9097"

echo "🧪 Тестирование Gateway Service"
echo "================================"

# Тест регистрации
echo "2. Тест регистрации пользователя..."
REGISTER_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "role": "user"
  }')
echo "Ответ: $REGISTER_RESPONSE"
echo -e "\n"

# Тест входа
echo "3. Тест входа пользователя..."
LOGIN_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')
echo "Ответ: $LOGIN_RESPONSE"
echo -e "\n"

# Извлекаем токен из ответа (предполагаем, что ответ содержит поле "token")
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo "4. Тест защищенного endpoint (с токеном)..."
    curl -s -X GET "$GATEWAY_URL/events/" \
      -H "Authorization: Bearer $TOKEN"
    echo -e "\n"
else
    echo "4. Токен не найден, пропускаем тест защищенного endpoint"
fi

# Тест защищенного endpoint без токена
echo "5. Тест защищенного endpoint (без токена)..."
curl -s -X GET "$GATEWAY_URL/events/"
echo -e "\n"

echo "✅ Тестирование завершено!"