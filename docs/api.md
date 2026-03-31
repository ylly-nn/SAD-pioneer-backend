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









# Ветка пользователtq тс
### Get /client/orders
Header: Bearer <acsess_token>
Получение заказов для опредеоленного владельца тс 
Пример успешного ответа: nill или
~~~
[
    {
        "order_id": "83817fd0-ffd0-478b-b1ae-b082e8581830",
        "name_company": "ООО \"Ромашка\"",
        "city": "Москва",
        "address": "ул. Тверская, д. 1",
        "service": "Шиномонтаж",
        "start_moment": "2026-03-16T09:30:00+04:00",
        "end_moment": "2026-03-16T11:05:00+04:00",
        "order_details": {
            "Шлифовка": 20,
            "Полировка": 75
        }
    },
    {
        "order_id": "b433a890-5280-4fb1-ad35-e0c94d5909d3",
        "name_company": "ООО \"Ромашка\"",
        "city": "Москва",
        "address": "ул. Тверская, д. 1",
        "service": "Шиномонтаж",
        "start_moment": "2026-03-17T09:00:00+04:00",
        "end_moment": "2026-03-17T10:35:00+04:00",
        "order_details": {
            "Шлифовка": 20,
            "Полировка": 75
        }
    }
]
~~~
### GET /services - защищённый
Header: Authorization: Bearer <токен>
Список всех улуг (общий: шмномотаж, мойка и тд)

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
Ответ ```[]``` - тоже успешный

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
    "city": "<city>"
}
~~~
Успешный ответ - ```200 OK```

---

### Get /branch?city=<city>&service=<service> - защищённый
Получение списка филиалов по городу и id услуги
Header: Authorization: Bearer <токен>


Пример успешного ответа:
~~~
[
    {
        "id_branchserv": "0bb7a20d-4ffc-46cd-b5e5-a549a179ce2a",
        "id_branch": "89d74b8a-8cee-44fa-96ea-6aec1e8ad66b",
        "address": "ул. Тверская, д. 1",
        "org_short_name": "ООО \"Ромашка\""
    },
    {
        "id_branchserv": "e43c3985-9d5c-4658-a4cb-4542d2a38ee3",
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
        "detail": "price",
        "duration_min": 2500
    },
    {
        "detail": "duration_min",
        "duration_min": 60
    }
]
~~~


---
### Get /branch/freetime?branch_id=<branch_id>&date=<yyyy-mm-dd>&duration=<min> - защищённый
Header: Authorization: Bearer <токен>


Получени свободных слотов - время

Пример:
```/branch/freetime?branch_id=89d74b8a-8cee-44fa-96ea-6aec1e8ad66b&date=2026-03-16&duration=300```
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

Получение списка всех заказов

если дата-время не будет попадать в свободные слоты - которе получаются через GET branch/freetime?branch_id=89d74b8a-8cee-44fa-96ea-6aec1e8ad66b&date=2026-03-16&duration=180 - будет выдавать ошибку 
```start moment is not available for the requested duration```

body:
~~~
{
    "service_by_branch": "e456338f-5b1c-49af-84d0-18f248d11b1d",
    "start_moment": "2026-03-16T11:20:00+04:00",
    "order_details": {
        "Полировка": 30,
        "Мойка днища":70
    }
}
~~~
Формат order_details - обязательный (услуга-минуты)

Пример успешного ответа: 
~~~
{
    "id": "244a6a7a-bc6c-4825-ad01-b483d308d77d",
    "users": "ylly_nn@mail.ru",
    "service_by_branch": "e456338f-5b1c-49af-84d0-18f248d11b1d",
    "start_moment": "2026-03-20T05:00:00Z",
    "end_moment": "2026-03-20T05:30:00Z",
    "order_details": {
        "Полировка": 30,
        "Мойка днища":70
    }
}
~~~
---





## /services

### POST /services
Создание новой услуги

body:
~~~
{
    "name": "<название_услуги>"
}
~~~
Пример успешного ответа
~~~
{
    "id": "cdaef1fc-0b99-437e-aba4-3dfb356cfd5c",
    "name": "<название_услуги>"
}
~~~
---

### DELETE /services/{id}
Удаление услуги по id

При успешном удалении статус: ```204 No Content```

---
---
## /company
---
### GET /company/{inn}
Плучение информации о компании по инн
~~~
GET /company/234567890123
~~~
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
### GET /company
Получение списка компаний

Пример успешного ответа:
~~~
[
    {
        "inn": "234567890123",
        "kpp": "234567891",
        "ogrn": "2345678901234",
        "org_name": "АО \"Технопром\"",
        "org_short_name": "Технопром"
    },
    {
        "inn": "345678901234",
        "kpp": "345678912",
        "ogrn": "3456789012345",
        "org_name": "ООО \"Альянс\"",
        "org_short_name": "Альянс"
    },
]
~~~
---
### Get /company/order/{inn}
Получение заказов для определённой компании

Пример успешного ответа: null или 

~~~
[
    {
        "id": "83817fd0-ffd0-478b-b1ae-b082e8581830",
        "users": "ivanov",
        "service_by_branch": "e456338f-5b1c-49af-84d0-18f248d11b1d",
        "inn": "770123456789",
        "name_company": "ООО \"Ромашка\"",
        "city": "Москва",
        "address": "ул. Тверская, д. 1",
        "service": "Шиномонтаж",
        "start_moment": "2026-03-16T09:30:00+04:00",
        "end_moment": "2026-03-16T11:05:00+04:00",
        "order_details": {
            "Шлифовка": 20,
            "Полировка": 75
        }
    },
    {
        "id": "b433a890-5280-4fb1-ad35-e0c94d5909d3",
        "users": "ivanov",
        "service_by_branch": "e456338f-5b1c-49af-84d0-18f248d11b1d",
        "inn": "770123456789",
        "name_company": "ООО \"Ромашка\"",
        "city": "Москва",
        "address": "ул. Тверская, д. 1",
        "service": "Шиномонтаж",
        "start_moment": "2026-03-17T09:00:00+04:00",
        "end_moment": "2026-03-17T10:35:00+04:00",
        "order_details": {
            "Шлифовка": 20,
            "Полировка": 75
        }
    }
]
~~~
---
### POST /company
Создание новой компании

body:
~~~
{
    "inn": "123456789011",
    "kpp": "123456789",
    "ogrn": "1234567890123",
    "org_name": "Общество с ограниченной ответственностью \"Ромашка\"",
    "org_short_name": "ООО \"Ромашка\""
  }
~~~
Пример успешного ответа:
~~~
{
    "inn": "123456789011",
    "kpp": "123456789",
    "ogrn": "1234567890123",
    "org_name": "Общество с ограниченной ответственностью \"Ромашка\"",
    "org_short_name": "ООО \"Ромашка\""
}
~~~
---
### DELETE /company/{inn}
Удаление компании по инн
~~~
DELETE /company/123456789011
~~~
Успешный ответ - ```204 No Content```

---
### Post /company/branch/service
Добавление сервиса - для филиала

body:
~~~
{
    "branch": "7b0fd8a2-0a9e-4004-b973-e36df7cd34e2",
    "service": "38ce95aa-e669-4166-a723-9779f869d894",
    "service_detalis": {
        "Мойка днища": 20,
        "Чистка салона": 30,
        "Полировка": 75
    }
}
~~~
Для service_detalis - формат обязательный(услуга-минуты)

Пример успешного ответа:
~~~
{
    "id": "9da1a2f4-7236-443a-aa4c-0d627dbed89e",
    "branch": "7b0fd8a2-0a9e-4004-b973-e36df7cd34e2",
    "service": "38ce95aa-e669-4166-a723-9779f869d894",
    "service_detalis": {
        "Мойка днища": 20,
        "Чистка салона": 30,
        "Полировка": 75
    }
}
~~~
---

## /client

### POST /client
Создание нового владельца тс - чтобы он создался он должен быть зарегеистрирован (то есть находится в all_users)

body
~~~
{
    "email": "<email>"
}
~~~
Пример успешного ответа
~~~
{
    "id": "<id>",
    "email": "<email>"
}
~~~
---




---

---
---

## /order
### GET /order

Получение списка всех заказов 
Успешный ответ null или 
~~~
[
    {
        "id": "94e416f7-889d-422e-99f1-113ba4c1841b",
        "idusers": "b0debf62-b3d6-4d27-ab66-f20213394be5",
        "users": "email@mail.ru",
        "service_by_branch": "f7148035-fd33-47e2-b380-1eacb4a66128",
        "inn": "345678901234",
        "name_company": "Альянс",
        "city": "Нижний Новгород",
        "address": "ул. Большая Покровская, 5",
        "service": "мойка",
        "start_moment": "2026-03-16T09:30:00+04:00",
        "end_moment": "2026-03-16T11:05:00+04:00",
        "order_details": {
            "Шлифовка": 20,
            "Полировка": 75
        }
    },
    {
        "id": "340e9cd7-c4d5-4063-8444-452ac7c85a59",
        "idusers": "b0debf62-b3d6-4d27-ab66-f20213394be5",
        "users": "email@mail.ru",
        "service_by_branch": "f7148035-fd33-47e2-b380-1eacb4a66128",
        "inn": "345678901234",
        "name_company": "Альянс",
        "city": "Нижний Новгород",
        "address": "ул. Большая Покровская, 5",
        "service": "мойка",
         "start_moment": "2026-03-17T09:00:00+04:00",
        "end_moment": "2026-03-17T10:35:00+04:00",
        "order_details": {
            "Шлифовка": 20,
            "Полировка": 75
        }
    }
]
~~~
---



# Ветка для администраторов + заявок для организаций
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

## /admin
---
### GET /admin/partner-requests/
Получение всех заявок от партнеров

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

