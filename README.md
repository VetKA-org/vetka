# Vetka
Сервис для организации работы ветеринарной клиники.

## О проекте
### Архитектура
В проекте организован по шаблону ["чистая архитектура"](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) на основе этих замечательных проектов с открытым исходным кодом:
- [go-clean-template](https://github.com/evrone/go-clean-template)
- [creatly-backend](https://github.com/Creatly/creatly-backend)
- [project-layout](https://github.com/golang-standards/project-layout).

## Настройка и значения по умолчанию
### Опции командной строки
Для вывода списка доступных опций и их значений по умолчанию выполните команду:
```bash
./cmd/app/app --help
```

### Переменные окружения
(!) Переменные окружения имеют приоритет перед опциями командной строки.

```bash
# Адрес и порт, по которым доступно API сервиса:
export RUN_ADDRESS=0.0.0.0:8080

# Полный URL для установления соединения с Postgres:
export DATABASE_URI=postgres://postgres:postgres@127.0.0.1:5432/vetka?sslmode=disable

# Полный URL для установления соединения с Redis:
export REDIS_URI=redis://:redis@127.0.0.1:6379/0

# Секретный ключ для генерации подписи (по умолчанию не задан).
# (!) Необходим для корректной работы сервиса.
export SECRET=

# Уровень логирования:
export LOG_LEVEL=info
```

## Разработка и тестирование
Для получения полного списка доступных команд выполните:
```bash
make help
```

### golangci-lint
В проекте используется `golangci-lint` для локальной разработки. Для установки линтера воспользуйтесь [официальной инструкцией](https://golangci-lint.run/usage/install/).

### pre-commit
В проекте используется `pre-commit` для запуска линтеров перед коммитом. Для установки утилиты воспользуйтесь [официальной инструкцией](https://pre-commit.com/#install), затем выполните команду:
```bash
make install-tools
```

### migrate
Для работы с миграциями БД необходимо установить утилиту [golang-migrate](https://github.com/golang-migrate/migrate):
```bash
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Команды
Для добавления новой миграции выполните:
```bash
migrate create -ext sql -dir ./migrations -seq имя_миграции
```

Для применения миграций выполните команду:
```bash
migrate -database ${DATABASE_URI} -path ./migrations up
```

Для возврата базы данных в первоначальное состояние выполните команду:
```bash
migrate -database ${DATABASE_URI} -path ./migrations down -all
```

## Лицензия
Copyright (c) 2023 Alexander Kurbatov

Лицензировано по [GPLv3](LICENSE).
