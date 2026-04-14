# Документация по api - реализованные эндпоинты
# Логика  входа
## /auth
---
### POST /auth/register
Регистрация нового пользователя (отправка кода подтверждения на email)

body:
~~~
{
    "email": "email@mail.ru",
    "password": "Password123!"
}
~~~
Пример успешного ответа
~~~
{
    "email": "email@mail.ru",
    "message": "Verification code sent to email"
}
~~~
---
### POST /auth/verify
Подтверждение кода регистрации и создание пользователя

body:
~~~
{
    "email": "email@mail.ru",
    "code": "a1b2c3"
}
~~~
Пример успешного ответа
~~~
{
    "message": "User successfully registered"
}
~~~
---
### POST /auth/login
Вход в систему, получение токенов доступа

body:
~~~
{
    "email": "email@mail.ru",
    "password": "Password123!"
}
~~~
Пример успешного ответа
~~~
{
     "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
}
~~~
---
### POST /auth/forgot-password
Запрос на отправление кода для восстановления пароля

body:
~~~
{
    "email": "email@mail.ru"
}
~~~
Пример успешного ответа
~~~
{
    "email": "email@mail.ru",
    "message":"Reset code sent to email"
}
~~~
---
### POST /auth/verify-reset-code
Подтверждение кода для восстановления пароля

body:
~~~
{
    "email": "email@mail.ru",
    "code": "000000"
}
~~~
Пример успешного ответа
~~~
{
    "message":"Code verified successfully"
}
~~~
---
### POST /auth/set-password
Установка нового пароля

body:
~~~
{
    "email": "email@mail.ru",
    "new_password": "new_password"
}
~~~
Пример успешного ответа
~~~
{
    "message":"Password reset successfully"
}
~~~
---
### POST /auth/logout
Выход из системы

body:
~~~
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
~~~
Пример успешного ответа
~~~
{
    "message": "Successfully logged out"
}
~~~
---
### POST /auth/refresh
Обновление пары токенов по refresh токену

body:
~~~
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
~~~
Пример успешного ответа
~~~
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
}
~~~
---
# Ветка пользователя тс
### Get /client/orders
Header: Bearer <acсess_token>
Получение заказов авторизованного клиента

Query параметры:
- timezone (опционально) — часовой пояс, например Europe/Moscow

Успешный ответ (200): массив заказов или null
~~~
[
    {
        "order_id": "f5affaa5-cf2c-4e4d-a261-cb93ccb72855",
        "name_company": "ООО \"Ромашка\"",
        "city": "Москва",
        "address": "г. Москва, ул. Тверская, д. 1",
        "service": "автомойка",
        "start_moment": "2026-04-01T06:00:00+00:00",
        "end_moment": "2026-04-01T06:30:00+00:00",
        "status": "create",
        "order_details": [
            {
                "detail": "мойка двигателя",
                "duration_min": 30,
                "price": 30
            }
        ],
        "sum": 30
    }
]

~~~
---
### GET /services - защищённый
Header: Authorization: Bearer <токен>
Список всех услуг (общий: шмномотаж, мойка и тд)

Пример успешного ответа:
~~~
[
    {
        "id": "27ddc0e1-5db0-4c3e-9406-b0d37fa6b4b6",
        "name": "Автомойка"
    },
    {
        "id": "878b0701-600b-409b-b8fe-3453887fc039",
        "name": "Шиномонтаж"
    }
]
~~~


---
### GET /client/city - защищённый
Header: Authorization: Bearer <токен>
Получение грода клиента по email - для автозаполнения

Пример успешного ответа:
~~~
{
    "city": "Южно-Сахалинск"
}
~~~
---
### PUT /client/city - защищённый
Header: Authorization: Bearer <токен>
Обновление города клиента - запоминается каждый раз при выборе города

body:
~~~
{
    "city": "Санкт-Петербург"
}
~~~
Успешный ответ - ```200 OK```

---

### Get /branch?city=<city>&service=<service> - защищённый
Получение списка филиалов по городу и id услуги
Header: Authorization: Bearer <токен>

Query параметры:
- city (обязательный) — город
- service (обязательный) — ID услуги

Успешный ответ (200):

~~~
[
    {
        "id_branchserv": "0bb7a20d-4ffc-46cd-b5e5-a549a179ce2a",
        "id_branch": "89d74b8a-8cee-44fa-96ea-6aec1e8ad66b",
        "address": "ул. Тверская, д. 1",
        "org_short_name": "ООО \"Ромашка\""
    }
]
~~~


---
### Get branch/service/details/{id_branchserv} - защищённый
Получение деталей услуги определённого филиала 

id_branchserv - получали на предыдущем шаге

Header: Authorization: Bearer <токен>


Пример успешного ответа
~~~
[
    {
        "detail": "Мойка салона",
        "duration_min": 40,
        "price": 560.12
    }
]
~~~


---
### Get /branch/freetime
Header: Authorization: Bearer <токен>

Получени свободных слотов - время

Query параметры:
- branch_id (обязательный) — UUID филиала
- date (обязательный) — дата начала недели yyyy-mm-dd
- duration (обязательный) — длительность услуги в минутах
- timezone (опционально) — часовой пояс

Успешный ответ (200):
Ответ:
~~~
[
    {
        "date": "2026-03-16T00:00:00Z",
        "intervals": null
    },
    {
        "date": "2026-03-17T00:00:00Z",
        "intervals": [
            "2026-03-17T10:35:00+04:00",
            "2026-03-17T10:50:00+04:00",
            "2026-03-17T11:05:00+04:00",
            "2026-03-17T11:20:00+04:00",
            "2026-03-17T11:35:00+04:00",
            "2026-03-17T11:50:00+04:00"
        ]
    },
    {
        "date": "2026-03-18T00:00:00Z",
        "intervals": null
    },
    {
        "date": "2026-03-19T00:00:00Z",
        "intervals": null
    },
    {
        "date": "2026-03-20T00:00:00Z",
        "intervals": [
            "2026-03-20T09:00:00+04:00",
            "2026-03-20T09:15:00+04:00",
            "2026-03-20T09:30:00+04:00",
            "2026-03-20T09:45:00+04:00",
            "2026-03-20T10:00:00+04:00",
            "2026-03-20T10:15:00+04:00"
        ]
    },
    {
        "date": "2026-03-21T00:00:00Z",
        "intervals": [
            "2026-03-21T09:00:00+04:00",
            "2026-03-21T09:15:00+04:00",
            "2026-03-21T09:30:00+04:00",
            "2026-03-21T09:45:00+04:00",
            "2026-03-21T10:00:00+04:00",
            "2026-03-21T10:15:00+04:00",
            "2026-03-21T10:30:00+04:00",
            "2026-03-21T10:45:00+04:00",
            "2026-03-21T11:00:00+04:00",
            "2026-03-21T11:15:00+04:00",
            "2026-03-21T11:30:00+04:00",
            "2026-03-21T11:45:00+04:00",
            "2026-03-21T12:00:00+04:00"
        ]
    },
    {
        "date": "2026-03-22T00:00:00Z",
        "intervals": null
    }
]
~~~
---
### Post /order - защищённый
Header: Authorization: Bearer <токен>

Создание заказа

если дата-время не будет попадать в свободные слоты - которые получаются через GET branch/freetime?branch_id=89d74b8a-8cee-44fa-96ea-6aec1e8ad66b&date=2026-03-16&duration=180 - будет выдавать ошибку 
```start moment is not available for the requested duration```

body:
~~~
{
    "service_by_branch": "e456338f-5b1c-49af-84d0-18f248d11b1d",
    "start_moment": "2026-03-16T11:20:00+04:00",
    "order_details": [
        {"detail": "Полировка"},
        {"detail": "Мойка днища"}
    ]
}
~~~
Формат order_details - обязательный (услуга-минуты)

Пример успешного ответа(201): 
~~~
{
    "id": "244a6a7a-bc6c-4825-ad01-b483d308d77d",
    "users": "client@mail.ru",
    "service_by_branch": "e456338f-5b1c-49af-84d0-18f248d11b1d",
    "start_moment": "2026-03-20T05:00:00Z",
    "end_moment": "2026-03-20T05:30:00Z",
    "order_details": [
        {"detail": "Полировка"},
        {"detail": "Мойка днища"}
    ]
}
~~~
---

---
# Ветка для компаний
## /company
---
### GET /company
Получение данных компании текущего авторизованного партнёра

Успешный ответ (200):

Пример успешного ответа: 
~~~
{
    "inn": "234567890123",
    "kpp": "234567891",
    "ogrn": "2345678901234",
    "org_name": "АО \"Технопром\"",
    "org_short_name": "Технопром"
}
~~~

---

### GET /company/branches
Список филиалов компании

Успешный ответ (200): массив филиалов или null
~~~
[
    {
        "branch_id": "9eebb3b9-5b35-4007-9d4f-2f4141786b45",
        "city": "Москва",
        "address": "ул. Тверская, 1",
        "open_time": "10:00:00+00:00",
        "close_time": "18:00:00+00:00",
        "inn_company": "123456789012"
    }
]
~~~
---

### GET /company/branches/{branch_id}
Детальная информация о филиале с услугами

Успешный ответ (200):
~~~
{
    "city": "Санкт-Петербург",
    "address": "Невский пр., 10",
    "open_time": "10:00:00+00:00",
    "close_time": "18:00:00+00:00",
    "services": [
        {
            "branch_serv_id": "6fdd2352-ffc4-4140-b54c-67657f841c1c",
            "service_id": "03db1f58-2bbd-481c-8d93-b2828871b376",
            "service_name": "мойка"
        }
    ]
}
~~~
---

### POST /company/branch
Добавить новый филиал компании

Body:

~~~
{
    "city": "Москва",
    "address": "Улица тверская дом 1",
    "open_time": "09:00:00+03:00",
    "close_time": "17:00:00+03:00"
}
~~~
Успешный ответ (201):
~~~
{
    "message": "Branch added to company successfully",
    "city": "Москва",
    "address": "Улица тверская дом 1",
    "open_time": "10:00:00+00:00",
    "close_time": "18:00:00+00:00"
}
~~~
---
### POST /company/branch/service
Добавить существующую услугу в филиал

Body:

~~~
{
    "branch_id": "9eebb3b9-5b35-4007-9d4f-2f4141786b45",
    "service_id": "03db1f58-2bbd-481c-8d93-b2828871b376"
}
~~~
Успешный ответ (201):

~~~
{
    "message": "Service added to branch successfully",
    "branch_id": "9eebb3b9-5b35-4007-9d4f-2f4141786b45",
    "service_id": "03db1f58-2bbd-481c-8d93-b2828871b376"
}
~~~
---
### DELETE /company/branch/service/detail/{branchServID}
Удалить деталь услуги филиала

Параметры:
- branchServID (path) — UUID записи услуги филиала
- detail (query) — название детали

Успешный ответ (200): обновлённый список деталей
---


### POST /company/branch/service/detail
Добавить деталь (конкретную работу) к услуге филиала

Body:
~~~
{
    "branchserv_id": "6fdd2352-ffc4-4140-b54c-67657f841c1c",
    "detail": "Мойка салона",
    "duration": 40,
    "price": 700.5
}
~~~
Успешный ответ (201):

~~~
[
    {
        "detail": "Мойка салона",
        "duration_min": 40,
        "price": 700.5
    }
]
~~~
---

### GET /company/branch/service/{branchServID}
Получить детали услуги филиала

Успешный ответ (200):
~~~
{
    "detail": "Мойка салона",
    "duration_min": 40,
    "price": 700.5
}
~~~
---
### POST /company/users
Добавить существующего пользователя в компанию как партнёра

Body:
~~~
{
    "email": "newuser@example.com"
}
~~~
Успешный ответ (201):
~~~
{
    "message": "User added to company successfully",
    "email": "newuser@example.com"
}
~~~
---
### GET /company/orders
Получить заказы компании (сгруппированные по филиалам)

Успешный ответ (200):

~~~
[
    {
        "branch_id": "...",
        "city": "Москва",
        "address": "ул. Тверская, 1",
        "orders": [
            {
                "id": "05544774-f958-42a9-8f9b-7b56f5a43b52",
                "users": "aigizshai@gmail.com",
                "service_by_branch": "aa082c35-bc51-40b4-a2c1-e84b25e69b95",
                "name_service": "автомойка",
                "start_moment": "2026-04-04T11:15:00Z",
                "end_moment": "2026-04-04T11:55:00Z",
                "status": "approve",
                "order_details": [
                    {
                        "detail": "Мойка кузова",
                        "duration_min": 40,
                        "price": 700.5
                    }
                ],
                "sum": 700.5
            }
        ]
    }
]
---
~~~
### PUT /company/order/status
Подтвердить или отклонить заказ (доступно только для партнёров)

Header: Authorization: Bearer <токен>

Query параметры:
- orderID (обязательный) — UUID заказа
- status (обязательный) — approve или reject

Успешный ответ (200): обновлённый объект заказа
~~~
{
    "id": "05544774-f958-42a9-8f9b-7b56f5a43b52",
    "users": "client@example.com",
    "service_by_branch": "aa082c35-bc51-40b4-a2c1-e84b25e69b95",
    "name_service": "автомойка",
    "start_moment": "2026-04-04T11:15:00Z",
    "end_moment": "2026-04-04T11:55:00Z",
    "status": "approve",
    "order_details": [...],
    "sum": 700.5
}
~~~
---

---
# Заявки для организаций
## /partner
---
### POST /partner/request
Создание заявки

Header: Authorization: Bearer <токен>

body:
~~~
{
    "inn": "<inn>", 
    "kpp": "<kpp>", 
    "ogrn": "<ogrn>", 
    "org_name": "<org_name>", 
    "org_short_name": "<org_short_name>", 
    "name": "<name>", 
    "surname": "<surname>", 
    "patronymic": "<patronymic>", 
    "email": "<email>", 
    "phone": "<phone>", 
    "info": "<info>"
}
~~~
Пример успешного ответа
~~~
{
    "message": "Partner request created successfully"
}
~~~
---
### GET /partner/request
Получение информации по заявке

Header: Authorization: Bearer <токен>

~~~
GET /partner/request
~~~

Пример успешного ответа
~~~
{
   "status":"<status>",
   "user_email":<user_email>,
   "inn":"<inn>",
   "kpp":"<kpp>",
   "ogrn":"<ogrn>",
   "org_name":"<org_name>",
   "org_short_name":"<org_short_name>",
   "name":"<name>",
   "surname":"<surname>",
   "patronymic":"<patronymic>",
   "email":"<email>",
   "phone":"<phone>",
   "info":"<info>",
   "created_at":"<created_at>,
   "last_used":"last_used"
}
~~~
---

# Ветка для администраторов
## /admin
---
### GET /admin/partner-requests/
Получение всех заявок от партнеров

Header: Authorization: Bearer <токен>

~~~
Пример успешного ответа
~~~
{
    [
        {
        "id":"<uuid>",
        "status":"<status>",
        "user_email":"<user_email>",
        "inn":"<inn>",
        "kpp":"<kpp>",
        "ogrn":"<ogrn>",
        "org_name":"<org_name>",
        "org_short_name":"<org_short_name>",
        "name":"<name>",
        "surname":"<surname>",
        "patronymic":"<patronymic>",
        "email":"<email>",
        "phone":"<phone>",
        "info":"<info>",
        "created_at":"<created_at>,
        "last_used":"last_used"
        },
    ]
}
~~~
---
### GET /admin/partner-requests/new
Получение всех новых заявок от партнеров

Header: Authorization: Bearer <токен>

body:
~~~
{
}
~~~
Пример успешного ответа
~~~
{
    [
        {
        "id":"<uuid>",
        "status":"new",
        "user_email":"<user_email>",
        "inn":"<inn>",
        "kpp":"<kpp>",
        "ogrn":"<ogrn>",
        "org_name":"<org_name>",
        "org_short_name":"<org_short_name>",
        "name":"<name>",
        "surname":"<surname>",
        "patronymic":"<patronymic>",
        "email":"<email>",
        "phone":"<phone>",
        "info":"<info>",
        "created_at":"<created_at>,
        "last_used":"last_used"
        },
    ]
}
~~~
---
### GET /admin/partner-requests/pending
Получение заявок от партнеров, находящихся в работе

Header: Authorization: Bearer <токен>

body:
~~~
{
}
~~~
Пример успешного ответа
~~~
{
    [
        {
        "id":"<uuid>",
        "status":"pending",
        "user_email":"<user_email>",
        "inn":"<inn>",
        "kpp":"<kpp>",
        "ogrn":"<ogrn>",
        "org_name":"<org_name>",
        "org_short_name":"<org_short_name>",
        "name":"<name>",
        "surname":"<surname>",
        "patronymic":"<patronymic>",
        "email":"<email>",
        "phone":"<phone>",
        "info":"<info>",
        "created_at":"<created_at>,
        "last_used":"last_used"
        },
    ]
}
~~~
---
---
### GET /admin/partner-requests/{id}
Получение информации об одной заявке по id

Header: Authorization: Bearer <токен>

body:
~~~
{
}
~~~
Пример успешного ответа
~~~
{
    [
        {
        "id":"<uuid>",
        "status":"<status>",
        "user_email":"<user_email>",
        "inn":"<inn>",
        "kpp":"<kpp>",
        "ogrn":"<ogrn>",
        "org_name":"<org_name>",
        "org_short_name":"<org_short_name>",
        "name":"<name>",
        "surname":"<surname>",
        "patronymic":"<patronymic>",
        "email":"<email>",
        "phone":"<phone>",
        "info":"<info>",
        "created_at":"<created_at>,
        "last_used":"last_used"
        },
    ]
}
~~~
---
### GET /admin/partner-requests/approved
Получение всех принятых заявок от партнеров

Header: Authorization: Bearer <токен>

body:
~~~
{
}
~~~
Пример успешного ответа
~~~
{
    [
        {
        "id":"<uuid>",
        "status":"approved",
        "user_email":"<user_email>",
        "inn":"<inn>",
        "kpp":"<kpp>",
        "ogrn":"<ogrn>",
        "org_name":"<org_name>",
        "org_short_name":"<org_short_name>",
        "name":"<name>",
        "surname":"<surname>",
        "patronymic":"<patronymic>",
        "email":"<email>",
        "phone":"<phone>",
        "info":"<info>",
        "created_at":"<created_at>,
        "last_used":"last_used"
        },
    ]
}
~~~
---
### GET /admin/partner-requests/rejected
Получение всех отклоненных заявок от партнеров

Header: Authorization: Bearer <токен>

body:
~~~
{
}
~~~
Пример успешного ответа
~~~
{
    [
        {
        "id":"<uuid>",
        "status":"rejected",
        "user_email":"<user_email>",
        "inn":"<inn>",
        "kpp":"<kpp>",
        "ogrn":"<ogrn>",
        "org_name":"<org_name>",
        "org_short_name":"<org_short_name>",
        "name":"<name>",
        "surname":"<surname>",
        "patronymic":"<patronymic>",
        "email":"<email>",
        "phone":"<phone>",
        "info":"<info>",
        "created_at":"<created_at>,
        "last_used":"last_used"
        },
    ]
}
~~~
---
### POST /admin/partner-requests/take
Смена статуса заявки с "новая" на "в работе"

Header: Authorization: Bearer <токен>

body:
~~~
{
    "id": "<uuid>", 
}
~~~

Пример успешного ответа
~~~
{
   "id":"<uuid>",
   "message":"Request taken to work",
   "status":"pending"
}
~~~
---
### POST /admin/partner-requests/approve
Смена статуса заявки с "в работе" на "принята". Создание компании и первого пользователя в ней

Header: Authorization: Bearer <токен>

body:
~~~
{
    "id": "<uuid>", 
}
~~~

Пример успешного ответа
~~~
{
   "message":"Partner request approved"
}
~~~
---
### POST /admin/partner-requests/reject
Смена статуса заявки с "в работе" на "отклонена"

Header: Authorization: Bearer <токен>

body:
~~~
{
    "id": "<uuid>", 
}
~~~

Пример успешного ответа
~~~
{
   "id":"<uuid>",
   "message":"Request rejected",
   "status":"rejected"
}
~~~
---
---
### POST /admin/create-admin
Добавление нового админа

Header: Authorization: Bearer <токен>

body:
~~~
{
    "email":"<email>",
    "name":"<name>",
    "surname":"<surname>"
}
~~~

Пример успешного ответа
~~~
{
   "email":"<email>",
   "message":"Admin created successfully"
}
~~~
---

