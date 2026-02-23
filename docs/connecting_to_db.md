# Подключение к базе данных

1. В ```src``` проекта создать файл ```.env```
2. Обязательно добавить ```src/.env``` в ```.gitignore``` в корне проекта 
3. В ```.env``` создать:
~~~
DB_HOST=<host>
DB_PORT=5432
DB_USER=<user>
DB_PASSWORD=<password>
DB_NAME=pioneer(или другое, если у вас называется иначе)
~~~
4. Сохранить
5. Если отдает: ```Connected to PostgreSQL successfully!``` - радуемся, всё правильно 