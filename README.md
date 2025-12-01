     PZ8 Технологии индустриального программирования.  
     Студент: Выборнов О.А.  
     Группа: ЭФМО-02-25  
    
Цели работы: 
        
        -Понять базовые принципы документной БД MongoDB (документ, коллекция, BSON, _id:ObjectID). 
        -Научиться подключаться к MongoDB из Go с использованием официального драйвера. 
        -Создать коллекцию, индексы и реализовать CRUD для одной сущности (например, notes). 
        -Отработать фильтрацию, пагинацию, обновления (в т.ч. частичные), удаление и обработку ошибок.
        


Как запустить:

1)cклонировать репозиторий:

    git clone -b main --single-branch https://github.com/omnikk/PZ8.git

2)Создать файл .env (если нет):

    cp .env.example .env
3)Запустить контейнеры:
    
    docker compose up -d
4)Проверить доступность api:

    curl http://localhost:8080/health

Примеры запросов и результат их выполнения

Проверка здоровья:
<img width="1090" height="418" alt="image" src="https://github.com/user-attachments/assets/e891891d-c919-4640-b15d-c0620f82bb26" />


Создание заметки
<img width="1089" height="474" alt="image" src="https://github.com/user-attachments/assets/7b27e254-00f0-45b6-beb7-ce6d4f7b3c11" />

Получение списка заметок
<img width="1083" height="437" alt="image" src="https://github.com/user-attachments/assets/1cad9902-9702-4e85-bc3b-2a95e9cf8991" />

Получение заметки по ID
<img width="727" height="410" alt="image" src="https://github.com/user-attachments/assets/ae72d5f3-4cfd-4d71-bf16-1ea931233bb2" />

Обновление заметки
<img width="1083" height="466" alt="image" src="https://github.com/user-attachments/assets/67dcd848-5a89-4fe9-833d-024aa01c9de0" />

Удаление заметки
<img width="956" height="265" alt="image" src="https://github.com/user-attachments/assets/2460a1ed-4aa9-44cb-986d-d2f40df74ea9" />

Поиск заметок по ключевому слову
<img width="1100" height="87" alt="image" src="https://github.com/user-attachments/assets/4c41c0f3-cd6e-4492-93e1-9f2f33a8e4b5" />



Cтруктура проекта:

pz8-mongo/
│   go.mod
│   go.sum
│   docker-compose.yml
├───cmd
│   └───api
│           main.go
│
└───internal
    └───db
    │        mongo.go
    └───cache
            handler.go
            model.go
            repo.go


Описание структуры:


    go.mod - файл модуля Go с зависимостями
    go.sum - контрольные суммы зависимостей проекта
    docker-compose.yml – конфигурация Docker для поднятия MongoDB. 
    main.go – главная точка входа в приложение, создаёт HTTP-сервер и настраивает маршруты /api/v1/notes и /health. 
    mongo.go – внутренний пакет db, реализует подключение к MongoDB и проверку соединения (ping). 
    handler.go – HTTP-обработчики для CRUD-операций по заметкам (notes). 
    model.go – структура данных Note с полями ID, Title, Content, CreatedAt, UpdatedAt. 
    repo.go – репозиторий для работы с MongoDB: методы Create, ByID, List, Update, Delete, работа с индексами и фильтрами.

