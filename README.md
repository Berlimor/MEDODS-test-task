# MEDODS-test-task
Часть сервиса авторизации с двумя эндпоинтами.
```
/JWT/get-token/:guid
/JWT/refresh-token/:refresh
```
- Первый эндпоинт принимает GUID, ищет человека с таким GUID в базе данных, после чего возвращает ответ с HTTP статусом 200 OK и два токена access и refresh в виде JSON.
* Второй эндпоинт выполняет refresh операцию токенов. На вход получает refresh токен, ищет по нему пользователя в базе, после чего генерирует новый access токен, нвоый refresh токен, заменяет данные в базе данных на новые, возвращает HTTP статус 200 OK и возвращает объект этого пользователя с новыми полями. 