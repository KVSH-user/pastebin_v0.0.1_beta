Реализовал сырой проект "Аля Pastebin" на скорую руку, просто потестировать.

Миграцию можно накатить вручную, используя ```goose -dir db/migrations postgres "postgresql://postgres:qwerty@127.0.0.1:5436/postgres?sslmode=disable" up```

Стандартный ```POST``` запрос для создания записки:
  ```address:port/add```
```json
    {
    "text" : "ТЕКСТ ЗАПИСКИ",
    "only_one" : false     // это поле необязательно. Стандартным значением является false. Если же прописать true, то записка будет "удалена" после первого прочтения.
    }
```

Запрос для прочтения записки ``GET``
  ```address:port/alias   //Алиас генерируется автоматически и выдается при создании записки, в ответ на запрос.```

Запрос для удаления записки ```DELETE```
  ```address:port/alias_for_del   //Алиас для удаления генерируется автоматически и выдается при создании записки, в ответ на запрос.```

В данной версии программы реализовано "мнимое" удаление содержимого записки, в условиях существования одной базы данных.
Если пользователь решает удалить записку, то в БД добавляется запись с датой удаления, записка приобретает тег ```[DEL]``` и становится недоступна по первоначальному ```alias```.
