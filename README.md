# Schedule

Schedule — учебный веб-сервис для хранения и просмотра расписания занятий факультета. Проект включает REST API, PostgreSQL, миграции базы данных, Docker Compose и простой веб-интерфейс.

## Возможности

- просмотр расписания группы на конкретный день;
- просмотр недельного расписания группы;
- добавление преподавателей, аудиторий, предметов и групп;
- создание строки расписания;
- удаление строк расписания по заданным параметрам;
- хранение данных в PostgreSQL;
- запуск приложения, базы данных и миграций через Docker Compose;
- логирование работы приложения в директорию `logs`.

## Стек

- Go 1.25
- PostgreSQL 15
- Gorilla Mux
- Zap Logger
- golang-migrate / migrate
- Docker / Docker Compose
- HTML, CSS, JavaScript

## Структура проекта

```text
.
├── cmd/
│   └── app/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── logger/                  # Инициализация логгера
│   ├── models/                  # Модели, DTO и ошибки
│   ├── repository/              # Работа с PostgreSQL
│   │   ├── admin/
│   │   └── schedule/
│   ├── service/                 # Бизнес-логика
│   │   ├── admin/
│   │   └── schedule/
│   └── transport/               # HTTP-обработчики
│       ├── admin/
│       └── schedule/
├── migrations/                  # SQL-миграции
├── web/                         # Веб-интерфейс
│   ├── assets/
│   └── index.html
├── Dockerfile
├── docker-compose.yaml
├── .env.example
└── go.mod
```

## Модель данных

В базе данных используются таблицы:

- `Groups` — учебные группы;
- `Subjects` — предметы;
- `Teachers` — преподаватели;
- `Classrooms` — аудитории;
- `Schedule` — строки расписания.

Таблица `Schedule` связывает группу, предмет, преподавателя и аудиторию. Также в ней хранятся день недели, номер пары, тип недели и подгруппа.

Ограничения:

- `weekday` — от 1 до 7;
- `lesson_number` — от 1 до 10;
- `week_type` — 1 или 2, либо `NULL` для обеих недель;
- `subgroup` — 1 или 2, либо `NULL` для всей группы;
- уникальность строки расписания задаётся по `group_id`, `weekday`, `lesson_number`, `week_type`, `subgroup`.

## Переменные окружения

Создайте `.env` на основе `.env.example`:

```bash
cp .env.example .env
```

Пример `.env`:

```env
APP_PORT=8080

POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=ScheduleDB

POSTGRES_URL=postgres://postgres:postgres@db:5432/ScheduleDB?sslmode=disable
```

При запуске через Docker Compose в `POSTGRES_URL` должен использоваться хост `db`, потому что это имя сервиса PostgreSQL внутри Docker-сети.

Для локального запуска без Docker адрес базы должен быть другим:

```env
POSTGRES_URL=postgres://postgres:postgres@localhost:5432/ScheduleDB?sslmode=disable
```

## Запуск через Docker Compose

Соберите и запустите проект:

```bash
docker compose up --build
```

После запуска будут подняты три сервиса:

- `schedule_db` — PostgreSQL;
- `migrate` — применение SQL-миграций;
- `schedule_app` — Go-приложение.

Приложение будет доступно по адресу:

```text
http://localhost:8080
```

Если порт изменён в `.env`, используйте значение из `APP_PORT`.

## Проверка контейнеров

```bash
docker ps
```

Проверка таблиц в базе данных:

```bash
docker exec -it schedule_db psql -U postgres -d ScheduleDB
```

Внутри `psql`:

```sql
\dt
```

Ожидаемые таблицы:

```text
groups
subjects
teachers
classrooms
schedule
schema_migrations
```

## Остановка проекта

Остановить контейнеры без удаления данных:

```bash
docker compose down
```

Остановить контейнеры и удалить volume с базой данных:

```bash
docker compose down -v
```

Второй вариант удалит все данные PostgreSQL.

## Локальный запуск без Docker

Для локального запуска нужен установленный PostgreSQL и созданная база данных.

Пример создания базы:

```bash
createdb ScheduleDB
```

В `.env` для локального запуска должен быть указан `localhost`:

```env
POSTGRES_URL=postgres://postgres:postgres@localhost:5432/ScheduleDB?sslmode=disable
```

Примените миграции:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/Schedule_DB?sslmode=disable" up
```

Запустите приложение:

```bash
go run cmd/app/main.go
```

## API

### Веб-интерфейс

```http
GET /
```

Открывает главную страницу приложения.

### Получить расписание на день

```http
GET /schedule?group_id=1&week_type=1&weekday=1&subgroup=1
```

Параметры:

- `group_id` — ID группы;
- `week_type` — тип недели: `1` или `2`;
- `weekday` — день недели от `1` до `7`;
- `subgroup` — необязательный параметр, `1` или `2`.

Пример:

```bash
curl "http://localhost:8080/schedule?group_id=1&week_type=1&weekday=1"
```

### Получить расписание на неделю

```http
GET /schedule/week?group_id=1&week_type=1&subgroup=1
```

Параметры:

- `group_id` — ID группы;
- `week_type` — необязательный параметр, `1` или `2`;
- `subgroup` — необязательный параметр, `1` или `2`.

Пример:

```bash
curl "http://localhost:8080/schedule/week?group_id=1&week_type=1"
```

## Справочники

### Получить преподавателей

```http
GET /teachers
```

```bash
curl "http://localhost:8080/teachers"
```

### Добавить преподавателей

```http
POST /teachers
```

```bash
curl -X POST "http://localhost:8080/teachers" \
  -H "Content-Type: application/json" \
  -d '[
    {"fullname":"Иванов Иван Иванович"},
    {"fullname":"Петров Петр Петрович"}
  ]'
```

### Получить предметы

```http
GET /subjects
```

```bash
curl "http://localhost:8080/subjects"
```

### Добавить предметы

```http
POST /subjects
```

```bash
curl -X POST "http://localhost:8080/subjects" \
  -H "Content-Type: application/json" \
  -d '[
    {"name":"Базы данных"},
    {"name":"Программирование"}
  ]'
```

### Получить аудитории

```http
GET /classrooms
```

```bash
curl "http://localhost:8080/classrooms"
```

### Добавить аудитории

```http
POST /classrooms
```

```bash
curl -X POST "http://localhost:8080/classrooms" \
  -H "Content-Type: application/json" \
  -d '[
    {"number":"101"},
    {"number":"202"}
  ]'
```

### Получить группы

```http
GET /groups
```

```bash
curl "http://localhost:8080/groups"
```

### Добавить группы

```http
POST /groups
```

```bash
curl -X POST "http://localhost:8080/groups" \
  -H "Content-Type: application/json" \
  -d '[
    {"name":"ПИ-101"},
    {"name":"ПИ-102"}
  ]'
```

## Работа с расписанием

### Добавить строку расписания

```http
POST /schedule
```

```bash
curl -X POST "http://localhost:8080/schedule" \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "subject_id": 1,
    "teacher_id": 1,
    "classroom_id": 1,
    "weekday": 1,
    "lesson_number": 1,
    "week_type": 1,
    "subgroup": 1
  }'
```

Если занятие проходит каждую неделю или для всей группы, поля `week_type` и `subgroup` можно не указывать:

```bash
curl -X POST "http://localhost:8080/schedule" \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "subject_id": 1,
    "teacher_id": 1,
    "classroom_id": 1,
    "weekday": 1,
    "lesson_number": 1
  }'
```

### Удалить строки расписания

```http
DELETE /schedule
```

```bash
curl -X DELETE "http://localhost:8080/schedule" \
  -H "Content-Type: application/json" \
  -d '{
    "group_name": "ПИ-101",
    "weekday": 1,
    "weektype": 1,
    "subgroup": 1,
    "lesson_number": 1
  }'
```

Важно: в текущей DTO для удаления поле называется `weektype`, без подчёркивания. Поэтому в JSON для удаления нужно использовать именно `weektype`, а не `week_type`.

Можно удалить строки шире, не указывая часть фильтров:

```bash
curl -X DELETE "http://localhost:8080/schedule" \
  -H "Content-Type: application/json" \
  -d '{
    "group_name": "ПИ-101",
    "weekday": 1
  }'
```

## Логи

Логи приложения сохраняются в директорию:

```text
logs/
```

При запуске через Docker Compose директория `./logs` пробрасывается внутрь контейнера как `/app/logs`.

## Важные замечания по текущей версии

В текущей версии проект работает, но есть несколько моментов, которые стоит учитывать:

- `APP_PORT` и `POSTGRES_URL` должны быть заданы в `.env`, потому что приложение читает их из переменных окружения.
- Сейчас нет авторизации для административных операций: добавлять и удалять данные может любой, кто имеет доступ к приложению.
- Нет отдельного endpoint `/health`, поэтому состояние приложения сейчас проверяется через запуск контейнера, логи или доступность главной страницы.

## Возможные доработки

- добавить авторизацию для административных endpoint'ов;
- добавить `/health` для проверки состояния приложения;
- улучшить обработку ошибок при создании расписания, например отдельно возвращать ошибку при конфликте уникальности;
- добавить unit-тесты для service-слоя;
- добавить интеграционные тесты для repository-слоя;
