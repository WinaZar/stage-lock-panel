# Stage Lock Panel

Мини-сервис, чтобы менять и мониторить статус stage-серверов - занят / не занят.

## Деплой

Перед `docker build` необходимо в директории `interface` создать / поравить файл `.env.local` и указать в нем корретные для прода значения переменных, перечисленных в `.env.local.example`.