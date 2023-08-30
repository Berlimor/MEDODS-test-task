# MEDODS-test-task
Часть сервиса авторизации с двумя эндпоинтами.
```
/auth/token/:guid
/auth/refresh/:refresh
```
- Первый эндпоинт принимает GUID, ищет пользователя с таким GUID в базе данных, после чего возвращает ответ с HTTP статусом 200 OK и два токена access и refresh в виде JSON.
* Второй эндпоинт выполняет refresh операцию токенов. На вход получает refresh токен, ищет по нему пользователя в базе, после чего генерирует новый access токен, новый refresh токен, заменяет данные в базе данных на новые, возвращает HTTP статус 200 OK и возвращает объект этого пользователя с новыми полями. 

# Как запустить проект 
Нужно ввести следующие команды в папке с проектом, предварительно закрыв локально MongoDB (при открытом монго на порте 27017 проект не сможет запуститься):
```
docker-compose up --build
```

# Чтобы протестировать
В файле main.go, в фуннкции main, перед запуском сервера, я создаю нового пользователя в базе. Для тестов работать надо именно с ним, ведь эндпоинта для создания нового пользователя нет.
- GUID: a6bcd248-c9ce-4475-9caf-b3313af3f14c
* access_token: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHVUlEIjoiYTZiY2QyNDgtYzljZS00NDc1LTljYWYtYjMzMTNhZjNmMTRjIiwic3ViIjoiYWNjZXNzX3Rva2VuIiwiZXhwIjoxNjkzMTQ0NDExLCJpYXQiOjE2OTMxNDQ0MTF9.tE_izf1PJcOsqgi5DRI3hL0_V4Sf1HfES68WEJKsRU8MOjW3X52n-HA4RsJMYgxg9sgigmvpJtd-69YY5QyXzg
- refresh-token: YWJjZGVBQkNERTEyMyFAIyQl