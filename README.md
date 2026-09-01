# HearYou

**HearYou** — realtime-мессенджер с авторизацией, комнатами и обменом сообщениями в реальном времени через WebSocket.

Проект построен с использованием **Go** на backend и **React + Vite** на frontend. Для хранения пользователей, комнат и сообщений используется **PostgreSQL**.

## ✨ Возможности

* 🔐 Регистрация и авторизация пользователей
* 🔑 JWT-аутентификация
* 💬 Обмен сообщениями в реальном времени через WebSocket
* 🏠 Общий чат
* 🚪 Создание отдельных комнат
* 👥 Добавление пользователей в комнату
* 💾 Сохранение сообщений в PostgreSQL
* 📜 Загрузка истории сообщений при подключении
* ❤️ Health-check сервера
* 💻 React-интерфейс мессенджера
* 🐳 PostgreSQL запускается через Docker Compose

## 🛠️ Технологии

### Backend

* **Go 1.26+**
* **net/http**
* **WebSocket**
* **PostgreSQL**
* **JWT**
* **bcrypt**
* **go-playground/validator**
* **pgx**

Основная WebSocket-реализация использует библиотеку `github.com/coder/websocket`.

### Frontend

* **React 19**
* **Vite 8**
* JavaScript
* CSS

Frontend собирается в `web/dist` и затем раздаётся непосредственно Go HTTP-сервером.

### Database

* **PostgreSQL 17**
* SQL migrations

В базе данных используются отдельные таблицы для пользователей, сообщений и комнат.

---

## 🏗️ Архитектура

Проект разделён на несколько основных частей:

```text
Hearyou/
├── cmd/
│   └── main.go
│
├── internel/
│   ├── auth/
│   ├── config/
│   ├── dto/
│   ├── server/
│   ├── storage/
│   ├── validation/
│   └── ws/
│
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   ├── 000002_create_messages_table.up.sql
│   ├── 000003_create_rooms_table.up.sql
│   └── 000004_add_room_to_messages.up.sql
│
├── web/
│   ├── src/
│   ├── index.html
│   └── package.json
│
├── docker-compose.yml
├── go.mod
└── go.sum
```

### Backend

Точка входа приложения находится в `cmd/main.go`.

При запуске приложение:

1. Загружает конфигурацию.
2. Подключается к PostgreSQL.
3. Инициализирует JWT.
4. Создаёт repositories для пользователей, сообщений и комнат.
5. Создаёт `RoomManager`.
6. Запускает HTTP-сервер.

WebSocket-комнаты управляются через `RoomManager`, который создаёт отдельный `Hub` для каждой комнаты.

### Frontend

React-приложение отвечает за:

* авторизацию;
* список комнат;
* создание комнат;
* подключение к WebSocket;
* отправку сообщений;
* отображение истории сообщений;
* состояние подключения.

При подключении к комнате frontend открывает WebSocket-соединение вида:

```text
/ws?token=<JWT>&room=<ROOM_NAME>
```

---

## 🔄 Как работает сообщение

Упрощённая схема:

```text
┌─────────────┐
│    React    │
│   Frontend  │
└──────┬──────┘
       │
       │ WebSocket
       ▼
┌─────────────┐
│  Go Server  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ RoomManager │
│    + Hub    │
└──────┬──────┘
       │
       ├──────────────► connected clients
       │
       ▼
┌─────────────┐
│ PostgreSQL  │
│   Storage   │
└─────────────┘
```

При подключении к комнате сервер:

1. Проверяет JWT.
2. Определяет пользователя.
3. Находит существующую комнату или создаёт её.
4. Подключает пользователя к соответствующему `Hub`.
5. Загружает последние 50 сообщений.
6. Отправляет историю клиенту.
7. Запускает чтение и запись WebSocket-соединения.

---

## 🗄️ Database

### Users

```sql
users
├── id
├── username
├── password_hash
└── created_at
```

`username` является уникальным. Пароль хранится в виде хеша.

### Rooms

```sql
rooms
├── id
├── name
├── created_by
└── created_at
```

Название комнаты является уникальным.

### Messages

```sql
messages
├── id
├── author_id
├── payload
├── created_at
└── room_id
```

Сообщение связано одновременно с пользователем и комнатой.

---

## 🚀 Запуск проекта

### 1. Клонирование

```bash
git clone https://github.com/GanghouDuf/Hearyou.git
cd Hearyou
```

### 2. Запуск PostgreSQL

Для базы данных используется PostgreSQL 17.

```bash
docker compose up -d
```

После запуска PostgreSQL будет доступен на:

```text
localhost:5433
```

Параметры контейнера:

```text
Database: mydb
User:     postgres
Password: postgres
Port:     5433
```

### 3. Настройка переменных окружения

Создайте файл `.env` в корне проекта:

```env
ADDR=:8080

DATABASE_URL=postgres://postgres:postgres@localhost:5433/mydb?sslmode=disable

JWT_SECRET=your-secret-key
```

`JWT_SECRET` является обязательным параметром. Если он отсутствует, сервер не запустится.

> Если PostgreSQL запущен непосредственно на стандартном порту `5432`, измените порт в `DATABASE_URL`.

### 4. Применение миграций

В проекте находятся SQL-миграции для создания:

* пользователей;
* сообщений;
* комнат;
* связи сообщений с комнатами.

Примените их к базе данных перед первым запуском приложения.

### 5. Установка frontend-зависимостей

```bash
cd web
npm install
```

### 6. Сборка frontend

```bash
npm run build
```

После сборки появится:

```text
web/dist
```

Go-сервер использует эту директорию для раздачи frontend-приложения.

### 7. Запуск backend

Вернитесь в корень проекта:

```bash
cd ..
go run ./cmd
```

После запуска приложение будет доступно по адресу:

```text
http://localhost:8080
```

---

## 📡 API

### Health Check

```http
GET /health
```

Проверка состояния сервера.

Ответ:

```text
Ok
```

### Регистрация

```http
POST /register
Content-Type: application/json
```

Request:

```json
{
  "username": "alex",
  "password": "password123"
}
```

### Авторизация

```http
POST /login
Content-Type: application/json
```

Request:

```json
{
  "username": "alex",
  "password": "password123"
}
```

Response:

```json
{
  "token": "<JWT>"
}
```

### WebSocket

```text
/ws?token=<JWT>&room=<ROOM_NAME>
```

Пример:

```text
ws://localhost:8080/ws?token=<JWT>&room=general
```

Для отправки сообщения:

```json
{
  "type": "chat",
  "payload": "Привет!"
}
```

---

## 🔐 Аутентификация

Для доступа к WebSocket клиент передаёт JWT в query-параметре:

```text
/ws?token=<JWT>&room=general
```

Сервер проверяет токен перед установкой WebSocket-соединения. Из JWT извлекаются `user_id` и `username`.

Пароли пользователей не хранятся в открытом виде — в базе хранится `password_hash`.

---

## 🧪 Development

Frontend в режиме разработки можно запустить отдельно:

```bash
cd web
npm run dev
```

Доступные команды:

```bash
npm run dev
npm run build
npm run lint
npm run preview
```

Backend:

```bash
go run ./cmd
```

---

## 📌 Текущий статус

Проект находится в стадии разработки.

