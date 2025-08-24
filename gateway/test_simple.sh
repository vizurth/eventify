#!/bin/bash

# Простой тест Gateway Service
GATEWAY_URL="http://localhost:9090"

echo "🚀 Тестирование Gateway Service"
echo "================================"

# 1. Health check
echo "1. Health check:"
curl -s "$GATEWAY_URL/health"
echo -e "\n"

# 2. Тест регистрации
echo "2. Регистрация пользователя:"
curl -s -X POST "$GATEWAY_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","role":"user"}'
echo -e "\n"

# 3. Тест входа
echo "3. Вход пользователя:"
LOGIN_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}')
echo "$LOGIN_RESPONSE"
echo -e "\n"

# 4. Извлекаем токен
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo "4. Тест защищенного endpoint (с токеном):"
    curl -s -X GET "$GATEWAY_URL/events/" \
      -H "Authorization: Bearer $TOKEN"
    echo -e "\n"
else
    echo "4. Токен не найден"
fi

# 5. Тест без токена
echo "5. Тест без токена (должен вернуть 401):"
curl -s -X GET "$GATEWAY_URL/events/"
echo -e "\n"

echo "✅ Тест завершен!" 