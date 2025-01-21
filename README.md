# Сервис отправки метрик
 
## Основные команды
- `make start` - запуск проекта
- `make stop` - остановка проекта
- `make down` - остановка и удаление всех артефактов проекта
- `make sh` - провалиться в терминал контейнера
---

- go get github.com/golang-migrate/migrate/v4
- go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
- migrate -database postgres://user:pass@localhost:5432/postgres\?sslmode=disable -path migrations up


## Flow работы приложения
1. Фабрика `db` выбирает драйвер бд, согласно env, потом с этим драйвером использует пакет sqlx для создания провайдера бд.
2. Готовый провайдер передается в адаптер бд.