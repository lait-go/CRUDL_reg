📦 CRUDL Subscription Service
REST API для управления онлайн-подписками пользователей. Проект разработан в рамках тестового задания Junior Golang Developer (Effective Mobile).

🔗 Репозиторий: github.com/lait-go/CRUDL_reg

🚀 Стек технологий:
Go (Golang)

Chi (легковесный HTTP-фреймворк)

PostgreSQL

SQLX (удобная работа с SQL)

Swaggo (Swagger/OpenAPI для документации)

Viper (работа с конфигами)

Docker + Docker Compose

golang-migrate (миграции базы данных)

📑 Возможности API:
CRUDL операции над записями подписок:

Создание

Чтение (по ID)

Обновление (частичное)

Удаление

Получение всех записей

Подсчёт суммарной стоимости подписок за выбранный период с фильтрами:

По user_id

По service_name

Swagger UI для документации и тестирования API.

📁 Формат данных подписки:
```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```
end_date — опционально.

Дата — в формате MM-YYYY.

Цена — целое число рублей.

⚙️ Конфигурация:
Все параметры подключения к БД и приложения указаны в config.yml. Пример:

```yaml
host: localhost
port: 5432
user: lait
password: 123
name: subscription_db
```
📊 Документация:
Swagger UI доступен после запуска по адресу:

```bash
http://localhost:8080/swagger/index.html
```
🐳 Запуск проекта через Docker Compose:
Собери образ и запусти сервис:

```bash
docker-compose up --build
```
Миграции применяются автоматически при запуске.

Приложение доступно по порту 8080.

📋 Основные эндпоинты:
Метод	Путь	Описание
POST	/api/user	Создать подписку
GET	/api/user/{id}	Получить подписку по ID
PUT	/api/user/{id}	Обновить подписку (частично)
DELETE	/api/user/{id}	Удалить подписку
GET	/api/user	Получить все подписки
GET	/api/total-price	Получить общую стоимость

Параметры для /api/total-price:

start_date, end_date (обязательно)

user_id, service_name (опционально)

📦 Особенности:
Логгирование всех операций.

Валидация входящих данных (uuid, даты, цены).

Пустые даты обрабатываются как NULL в базе данных.

Проект масштабируем, легко расширяется.
