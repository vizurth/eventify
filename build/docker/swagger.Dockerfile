FROM swaggerapi/swagger-ui:latest

# Копируем все swagger JSON файлы внутрь контейнера
COPY swagger/docs /usr/share/nginx/html/swagger

# Указываем URLS для нескольких файлов
ENV SWAGGER_JSON=""
ENV URLS='[{\"name\":\"Auth\",\"url\":\"/swagger/auth.swagger.json\"}, \
            {\"name\":\"Event\",\"url\":\"/swagger/event.swagger.json\"}, \
            {\"name\":\"UserInteract\",\"url\":\"/swagger/user-interact.swagger.json\"}]'
