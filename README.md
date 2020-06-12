# Stage Lock Panel

Мини-сервис, чтобы менять и мониторить статус stage-серверов - занят / не занят.

## Деплой

Перед `docker build` необходимо в директории `interface` создать / поравить файл `.env.local` и указать в нем корретные для прода значения переменных, перечисленных в `.env.local.example`.

## Разработка

Чтобы нормально заработал `eslint` нужно добавить в настройки проекта `.vscode/settings.json` следующее:

```json
{
    "eslint.workingDirectories":[
        "./interface"
    ]
}
```

Взято [отсюда](https://github.com/vuejs/eslint-plugin-vue/issues/976#issuecomment-555925022)