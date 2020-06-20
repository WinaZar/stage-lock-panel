# Stage Lock Panel

Мини-сервис, чтобы менять и мониторить статус stage-серверов - занят / не занят.

## Конфигурация

Для деплоя и разработки требуется конфиг файл, который нужно сделать из `config.example.toml` и указать приложению при старте:

```sh
./main.bin --config /path/to/config.toml
```

## Деплой

Приложение рассчитано на деплой c использованием Docker / Swarm / Nomad / Kubernetes.

Перед `docker build` необходимо в директории `interface` создать / поправить файл `.env.local` и указать в нем корректные для прода значения переменных, перечисленных в `.env.local.example`.

Пример для GitLab:

```sh
docker login registry.gitlab.com # если еще не авторизованы
docker build -t registry.gitlab.com/yourname/stage-lock-panel:1.0 .
docker push registry.gitlab.com/yourname/stage-lock-panel
```

## Разработка

Вам потребуется Go > 1.14 и NodeJS > 12 .

Запуск бэкенда на Go:

```sh
go run main.go --config /path/to/config.toml # или просто положить config.toml рядом с main.go
```

Запуск интерфейса (все манипуляции внутри директории `interface`):

1. Создать файл `.env.local` и указать в нем значения переменных, перечисленных в `.env.local.example`
2. Установить пакеты командой `npm install`
2. Выполнить `npm run serve` для запуска dev-сервера

Интерфейс сделан с использованием `vue`, `vue-cli` и `vuerify`.

Чтобы нормально заработал `eslint` для `vscode` нужно добавить в настройки проекта `.vscode/settings.json` следующее:

```json
{
    "eslint.workingDirectories":[
        "./interface"
    ]
}
```

Взято [отсюда](https://github.com/vuejs/eslint-plugin-vue/issues/976#issuecomment-555925022)

## Тесты

Тестов пока нет, но обязательно появятся.