# Download Service (Test Assignment)

Простой веб-сервис на Go для скачивания файлов по задачам.

## Возможности
- `POST /tasks` — создать задачу с URL-ами файлов для скачивания.
- `GET /tasks/{id}` — получить статус задачи.
- Задачи сохраняются в файл `tasks.json`, чтобы после перезапуска сервис продолжил выполнение.

## Запуск
```bash
go run ./cmd/server
```

Сервис слушает порт `:8080` (по умолчанию).

## Пример использования
Создать задачу:
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"name":"street2","urls":["https://samplelib.com/lib/preview/mp4/sample-30s.mp4"]}'
```

Проверить статус:
```bash
curl http://localhost:8080/tasks/<task_id>
```

Файлы сохраняются в `./downloads/<task_name>/`.
