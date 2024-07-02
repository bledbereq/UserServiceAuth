# *UserServiceAuth*

## ТЗ на сервис работы с пользователями (user-service)
1. Добавление, удаление, изменение пользователя, выдача информации о ользователе (CRUD, HTTP)
    - Изменять данные о самом себе может каждый, изменять/удалять других может только админ
    - Авторизация и аутентификация
    - Генерация JWT токенов
    - Сохранение информации о пользователях в базу данных (postgress)
    - Выдача публичного ключа для расшифровки JWT токенов другими сервисами (GRPC)

2. Для запуска сервиса через файл настроек задается HTTP порт, GRPC порт, путь до папки с ключами для генерации JWT токенов и др. настройки при необходимости
3. Всё должно быть развернуто в docker-compose

### Для генерация протофайлов вам нужно выполнить следующие команды в консоли:
1. `mkdr gen/go`

Это команда создаст папки gen/go, в которых будут лежать сгенерированные файлы. 

2. `task generate`

Эта команда запустит Taskfile, который содержит в себе необходимый скрипт для генерации протофайлов