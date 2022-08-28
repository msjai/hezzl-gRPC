# Hezzl
### тестовое задание GO для hezzl - 

1. Описать proto файл с сервисом из 3 методов: добавить пользователя, удалить пользователя, список пользователей
2. Реализовать gRPC сервис на основе proto файла на Go
3. Для хранения данных использовать PostgreSQL
4. на запрос получения списка пользователей данные будут кешироваться в redis на минуту и брать из редиса
5. При добавлении пользователя делать лог в clickHouse
6. Добавление логов в clickHouse делать через очередь Kafka

#### описание решения

на тестовом сервере (ubuntu) 158.160.10.60 развернуты:

- PostgresSQL
  - port: 5432
  - user: "postgres"
  - password: "postgres"
  - db: hezzlusers
- Redis
  - port: 6379
  - password: "wiNNer4000"
  - DB: 0 (default db)
- Kafka
  - port: 9092
  - topic: HezzlLogs
- ClickHouse
  - port: 8123
  - data base: default
  - user: default
  - password: wiNNer4000
  - можно подключиться из браузера: http://158.160.10.60:8123/play
  - tables: logs_queue, logs, logs_consumer

методы gRPC:
- AddUser
- DelUser
- ListUsers

Для тестирования gRPC можно использовать evans:
https://github.com/ktr0731/evans

Или evans.exe лежит в корне проекта. Параметры запуска evans: "evans protogrpc\manipulation.proto -p 8080"

  


