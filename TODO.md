# Jailpack

Инструмент для управления jail, вдохновлённый Kubernetes и Docker, но заточенный под FreeBSD и его философию: простота, стабильность, производительность, безопасность.

### Цели по разработке

- Создать CLI с Cobra: jailpack build, run, deploy, list, logs, shell, push, pull, destroy
- Реализовать jailpack build — упаковка приложения в .jail.tar.gz cage
- Реализовать jailpack run — запуск cage как jail
- Реализовать jailpack list — показ запущенных jail (через jls)
- Реализовать jailpack logs и shell — доступ к логам и оболочке
- Добавить поддержку флагов: --name, --ip, --port, --env, --output

### Сборка и Cage
- Поддержка `jailpack.yaml` для декларативной конфигурации
- Создание **Cage** `.jail.tar.gz` или `.cage.tar.gz` с `rootfs/`, `config.json`, `startup.sh` (или другие варианты)
- Зависимости
- Поддержка `.jailignore` (как `.dockerignore`) если нужно

