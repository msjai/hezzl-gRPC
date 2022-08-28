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
- Redis
- Kafka
- ClickHouse


