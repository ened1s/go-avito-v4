# go-avito-v4

REST API реализующее методы создания, удаления сегментов, получение информации об активных сегментах пользователя,а также добавление и удаление пользователей в сегмент.
Используемый редактор кода: VS Code
СУБД: PostgreSQL

1.В PostgreSQL создается база данных с параметрами:
- user=postgres 
- password=1
- dbname=avito-tech
PostgreSQL использует адрес localhost:5432

2.В этой БД создается таблица c cегментами "segments":
- SQL запрос находится в файле create_table_segments.sql

3.Далее в этой БД создается таблица "usersegments", в которой будут содержаться пользователи и сегменты, которым он принадлежит:
- SQL запрос находится в файле create_table_usersegments.sql

4.Разворачиваем docker проект с нашим сервисом и запускаем его.

5.Для тестирования запросов/ответов использовался POSTMAN. Сервис располагается на адресе localhost:1234.
Для тестирования методов использовались следующие запросы:

- метод создания сегмента имеет тип запроса POST
    Запрос: localhost:1234/segm?slug=add1
    Ответ: 
    {
        "type": "success",
        "user_id": 0,
        "segments": null,
        "message": "Добавление сегмента add1:success"
    }

- метод удаления сегмента имеет тип запроса DELETE
    Запрос: localhost:1234/segm?slug=add1
    Ответ: 
    {
        "type": "success",
        "user_id": 0,
        "segments": null,
        "message": "Удаление сегмента add1:success"
    }

- метод добавления/удаления сегментов пользователя имеет тип запроса PATCH
    Запрос: localhost:1234/user?add_seg=add1,add2&del_seg=&user_id=1
    Ответ:
    {
        "type": "success",
        "user_id": 1,
        "segments": null,
        "message": "Добавление в пользователя в сегмент add1: success, Добавление в пользователя в сегмент add2: success"
    }

    Запрос: localhost:1234/user?add_seg=&del_seg=add1,add2&user_id=1
    Ответ:
    {
        "type": "success",
        "user_id": 1,
        "segments": null,
        "message": "Удаление пользователя из сегмента add1: success, Удаление пользователя из сегмента add2: success"
    }

- метод получения активных сегментов пользователя имеет тип запроса GET
    Запрос: localhost:1234/user?user_id=1
    {
    "type": "success",
    "user_id": 1,
    "segments": [
        "add1",
        "add2",
        "add3"
    ],
    "message": ""
}




