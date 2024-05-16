# Техническое задание

## Ежедневник в формате "TO DO"

---

### Общие требования

Система представляет собой HTTP API со следующими требованиями к бизнес-логике:

* регистрация, аутентификация и авторизация пользователей;
* создание групп, приглашения в группу;
* создание, получение, удалиние и изменение записей;

### Абстрактная схема взаимодействия с системой

Ниже представлена абстрактная бизнес-логика взаимодействия пользователя с системой:

1. Пользователь регистрируется/авторизируется в системе.
2. Пользователь создает/не создает группу людей для создания записей.
3. Пользователь совершает работу с заметками и колонками.

### Сводное HTTP API

Система работы с пользователем должна предоставлять следующие HTTP-хендлеры:

* POST /api/users/register — регистрация пользователя;
* POST /api/users/login — аутентификация пользователя;
* GET /api/users/me — получение профиля аутентифицированного пользователя;
* GET /api/users/all — получение списка пользователей;
* PUT /api/users/change-password — изменение пароля пользователя;
* POST /api/users/refresh-token — подписание рефреш токена;

Система работы с записями должна предоставлять следующие HTTP-хендлеры:

* GET api/todos - получение всех записей;
* GET api/todos/{id} - получение записи по идентификатору;
* POST api/todos - создание записи;
* PUT api/todos/{id} - изменение записи по идентификатору;
* DELETE api/todos/{id} - удаление записи по идентификатору;

Система работы с группами должна предоставлять следующие HTTP-хендлеры:

* POST api/groups - создание группы;
* POST api/groups/{id} - получение группы;
* GET api/groups/create_invite - создание ссылки приглашения;
* GET api/groups - получение своих групп;

#### Регистрация пользователя

Хендлер: POST /api/user/register.

Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.

После успешной регистрации должна происходить автоматическая аутентификация пользователя.

Формат запроса:
POST /api/user/register HTTP/1.1
Content-Type: application/json
...

{
"login": "<login>",
"password": "<password>"
}

Возможные коды ответа:

- 201 — пользователь успешно зарегистрирован и аутентифицирован;
- 400 — неверный формат запроса;
- 409 — логин уже занят;
- 500 — внутренняя ошибка сервера.

#### Аутентификация пользователя

Хендлер: POST /api/user/login.

Аутентификация производится по паре логин/пароль.

Формат запроса:
POST /api/user/login HTTP/1.1
Content-Type: application/json
...

{
"login": "<login>",
"password": "<password>"
}

Возможные коды ответа:

- 200 — пользователь успешно аутентифицирован;
- 400 — неверный формат запроса;
- 401 — неверная пара логин/пароль;
- 500 — внутренняя ошибка сервера.

#### Получение профиля аутентифицированного пользователя

Хендлер: GET /api/users/me.

Хендлер доступен только аутентифицированным пользователям.

Формат запроса:
GET /api/users/me HTTP/1.1
Content-Type: application/json
...

{
"id": "<id>"
"login": "<login>"
}

Возможные коды ответа:

- 200 — успешное возвращение данных пользователя;
- 401 — пользователь не аутентифицирован;
- 409 — пользователь не найден;

#### Получение списка пользователей

Хендлер: GET /api/user/orders.

Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым. Формат даты — RFC3339.

Доступные статусы обработки расчётов:

- NEW — заказ загружен в систему, но не попал в обработку;
- PROCESSING — вознаграждение за заказ рассчитывается;
- INVALID — система расчёта вознаграждений отказала в расчёте;
- PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.

Формат запроса:
GET /api/user/orders HTTP/1.1
Content-Length: 0

Возможные коды ответа:

- 200 — успешная обработка запроса.

  Формат ответа:

       200 OK HTTP/1.1
  Content-Type: application/json
  ...

  [
  {
  "number": "9278923470",
  "status": "PROCESSED",
  "accrual": 500,
  "uploaded_at": "2020-12-10T15:15:45+03:00"
  },
  {
  "number": "12345678903",
  "status": "PROCESSING",
  "uploaded_at": "2020-12-10T15:12:01+03:00"
  },
  {
  "number": "346436439",
  "status": "INVALID",
  "uploaded_at": "2020-12-09T16:09:53+03:00"
  }
  ]

- 204 — нет данных для ответа.
- 401 — пользователь не авторизован.
- 500 — внутренняя ошибка сервера.

#### Получение текущего баланса пользователя

Хендлер: GET /api/user/balance.

Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.

Формат запроса:
GET /api/user/balance HTTP/1.1
Content-Length: 0

Возможные коды ответа:

- 200 — успешная обработка запроса.

  Формат ответа:

       200 OK HTTP/1.1
  Content-Type: application/json
  ...

  {
  "current": 500.5,
  "withdrawn": 42
  }

- 401 — пользователь не авторизован.
- 500 — внутренняя ошибка сервера.

#### Запрос на списание средств

Хендлер: POST /api/user/balance/withdraw

Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя в счет оплаты которого списываются баллы.

Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.

Формат запроса:
POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json

{

Даня, [15.05.2024 18:20]
"order": "2377225624",
"sum": 751
}

Здесь order — номер заказа, а sum — сумма баллов к списанию в счёт оплаты.

Возможные коды ответа:

- 200 — успешная обработка запроса;
- 401 — пользователь не авторизован;
- 402 — на счету недостаточно средств;
- 422 — неверный номер заказа;
- 500 — внутренняя ошибка сервера.

#### Получение информации о выводе средств

Хендлер: GET /api/user/balance/withdrawals.

Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым. Формат даты — RFC3339.

Формат запроса:
GET /api/user/withdrawals HTTP/1.1
Content-Length: 0

Возможные коды ответа:

- 200 — успешная обработка запроса.

  Формат ответа:

       200 OK HTTP/1.1
  Content-Type: application/json
  ...

  [
  {
  "order": "2377225624",
  "sum": 500,
  "processed_at": "2020-12-09T16:09:57+03:00"
  }
  ]

- 204 - нет ни одного списания.
- 401 — пользователь не авторизован.
- 500 — внутренняя ошибка сервера.

### Взаимодействие с системой расчёта начислений баллов лояльности

Для взаимодействия с системой доступен один хендлер:

- GET /api/orders/{number} — получение информации о расчёте начислений баллов лояльности.

Формат запроса:
GET /api/orders/{number} HTTP/1.1
Content-Length: 0

Возможные коды ответа:

- 200 — успешная обработка запроса.

  Формат ответа:

       200 OK HTTP/1.1
  Content-Type: application/json
  ...

  {
  "order": "<number>",
  "status": "PROCESSED",
  "accrual": 500
  }

  Поля объекта ответа:

    - order — номер заказа;
    - status — статус расчёта начисления:

        - REGISTERED — заказ зарегистрирован, но не начисление не рассчитано;
        - INVALID — заказ не принят к расчёту, и вознаграждение не будет начислено;
        - PROCESSING — расчёт начисления в процессе;
        - PROCESSED — расчёт начисления окончен;

    - accrual — рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.

- 429 — превышено количество запросов к сервису.

  Формат ответа:

       429 Too Many Requests HTTP/1.1
  Content-Type: text/plain
  Retry-After: 60

  No more than N requests per minute allowed

- 500 — внутренняя ошибка сервера.

Заказ может быть взят в расчёт в любой момент после его совершения. Время выполнения расчёта системой не регламентировано. Статусы INVALID и PROCESSED являются окончательными.

Общее количество запросов информации о начислении не ограничено.

### Конфигурирование сервиса накопительной системы лояльности

Сервис должн поддерживать конфигурирование следующими методами:

- адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a
- адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d
- адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r