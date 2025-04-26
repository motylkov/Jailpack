# Jailpack

Инструмент для управления jail, вдохновлённый Kubernetes и Docker, но заточенный под FreeBSD и его философию: простота, стабильность, производительность, безопасность.

## Возможности

- **Сборка приложений** в переносимые контейнеры (.cage.tar.gz)
- **Запуск jail** из готовых контейнеров
- **Управление** запущенными jail
- **Автоматическое определение** типа приложения (Go, Python, Node.js)
- **Простая конфигурация** через jailpack.yaml

## Установка

```bash
# Клонирование репозитория
git clone <repository-url>
cd jailpack

# Сборка
go build -o jailpack

# Установка (опционально)
sudo cp jailpack /usr/local/bin/
```

## Использование

### Сборка приложения

```bash
# Создание Cage из директории приложения
jailpack build ./my-app

# Указание имени выходного файла
jailpack build ./my-app -o my-app.cage.tar.gz
```

### Запуск Cage

```bash
# Запуск с параметрами по умолчанию
jailpack run my-app.cage.tar.gz

# Запуск с кастомными параметрами
jailpack run my-app.cage.tar.gz --name my-jail --ip 10.0.0.20
```

### Управление jail

```bash
# Просмотр запущенных jail
jailpack list
```

## Конфигурация

Создайте `jailpack.yaml` в корне вашего проекта для настройки сборки:

```yaml
# jailpack.yaml — конфигурация для сборки Cage
name: my-application
version: 1.0.0
description: "Описание приложения"

# Настройки сборки
build:
  output: my-app.cage.tar.gz
  ignore:
    - .git
    - node_modules
    - *.log

# Настройки запуска
run:
  default_name: my-jail
  default_ip: 10.0.0.10
  ports:
    - 8080:8080
```

## Требования

- FreeBSD 13.0+
- Go 1.22+
- Права администратора для создания jail

## Архитектура

Jailpack создаёт **Cage** — самодостаточный контейнер, содержащий:

- `rootfs/` — минимальная файловая система
- `app/` — ваше приложение
- `app-start.sh` — скрипт запуска
- `config.json` — метаданные контейнера

## Статус разработки

Проект находится в активной разработке. См. [TODO_RU.md](TODO_RU.md) для планов развития.

## Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Создайте Pull Request

## Ссылки

- [FreeBSD Jail Documentation](https://docs.freebsd.org/en/books/handbook/jails/)
