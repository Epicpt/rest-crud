# REST CRUD Application

REST CRUD Application - это веб-приложение на Go для управления placement и webmasters c кешированием данных in-memory.

## Установка

### Склонируйте репозиторий
```
git clone https://github.com/Epicpt/rest-crud.git
cd rest-crud
go mod tidy
```
### Настройка базы данных
* Создайте таблицу в PostgreSQL.
* Для настройки базы данных измените данные в URL подключения в файле config/config.yaml.
```
database:
  url: "postgres://LOGIN:PASSWORD@localhost:5432/DB_NAME?sslmode=disable"
```
LOGIN - логин в PostgreSQL
PASSWORD - пароль в PostgreSQL
DB_NAME - название таблицы

## Запуск
```
cd cmd
go run main.go
```

Сервер доступен по адресу http://localhost:8080

## Проверка API
Для проверки API используйте Postman. Примеры для каждого вызова CRUD:
### Webmasters
#### Создание webmaster
#### POST /webmasters
* Пример запроса:
```
{
    "name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "status": "active"
}
```
* Пример ответа:
```
{"id":1}
```
#### Получение списка webmasters с пагинацией
#### GET  /webmasters?page=1&limit=10
* Пример ответа:
```
{
    "limit": 10,
    "page": 1,
    "webmasters": [
        {
            "Webmaster": {
                "ID": 1,
                "name": "John",
                "last_name": "Doe",
                "email": "john.doe@example.com",
                "status": "active"
            },
            "Placements": []
        }
    ]
}
```
#### Обновление webmaster
#### PUT /webmasters/id
* Пример запроса:
```
{
    "name": "Janek",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "status": "active"
}
```
* Пример ответа:
```
{"message":"Webmaster обновлён"}
```
#### Удаление webmaster
#### DELETE /webmasters/id
* Пример ответа:
```
{"message":"Webmaster удалён"}
```
### Placements
#### Создание placement
#### POST /placements
* Пример запроса:
```
{
  "user_id": 1,
  "name": "New Placement",
  "description": "Description of the new placement"
}
```
* Пример ответа:
```
{"id":1}
```
После создания, также отображается в GET /webmasters?page=1&limit=10
* Пример
```
{
    "limit": 10,
    "page": 1,
    "webmasters": [
        {
            "Webmaster": {
                "ID": 1,
                "name": "John",
                "last_name": "Doe",
                "email": "john.doe@example.com",
                "status": "active"
            },
            "Placements": [
                {
                    "ID": 1,
                    "user_id": 1,
                    "name": "New Placement",
                    "description": "Description of the new placement"
                }
            ]
        }
    ]
}
```
#### Получение списка placements с пагинацией
#### GET  /placements?page=1&limit=10
* Пример ответа:
```
{
    "limit": 10,
    "page": 1,
    "placements": [
    {
        "ID": 1,
        "user_id": 1,
        "name": "New Placement",
        "description": "Description of the new placement"
    }
  ]
}
```
#### Обновление placement
#### PUT /placements/id
* Пример запроса:
```
{
  "user_id": 1,
  "name": "Updated Placement",
  "description": "Updated description"
}
```
* Пример ответа:
```
{"message":"Placement обновлён"}
```
#### Удаление placement
#### DELETE /placements/id
* Пример ответа:
```
{"message":"Placement удалён"}
```

## Стэк технологий
* Golang
* PostgreSQL
