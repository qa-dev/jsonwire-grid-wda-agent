# 0.7.0
- Добавлена совместимость с go get/govendor

# 0.6.1
- Исправлена ошибка установки timeout для дробных значений в clients.NewHTTPClient

# 0.6
- Добавлен middleware для проверки HTTP POST

# 0.5
- Изменён content-type с text/json на application/json

# 0.4
Добавлен тип Duration, реализуюший интерфейс Marshaller и Unmarshaller с поддержкой строковых литералов для 
Unmarshal (например, "2h40m32s")

# 0.3.6
- добавлены недостающие типы ошибок при валидации json схемы:
  - unique
  - multiple_of
  - number_gte
  - number_lt
  
  Более подробно об этих типах ошибках написано [здесь](https://github.com/xeipuuv/gojsonschema#working-with-errors)
  
- добавлены новые типы валидации данных json схемы:
  - number
  - boolean

# 0.3.5
- Добавлена поддержка http.Hijacker интерфейса для LoggedResponseWriter
- Добавлена возможность вернуть ошибку из http mock client-a в testutils

# 0.3.4
- Убран вызов panic из /_error

# 0.3.3
- В связи с переходом на kubernetes hostname убран из названий метрик

# 0.3.2
- Исправлен баг с инициализацией fluent логгера (теперь задается MaxRetry)

# 0.3.1
- исправлен вывод ошибки в валидации jsonschema

# 0.3
- validate fabric
- core writer - оборачивает сообщение в core протокол

# 0.2
Добавлена новая функциональность:
- Логгер добавляет имя функции из который вызван и строчку кода
- Вспомогательные функции для работы с request
- Валидация по json schema

# 0.1
Реализована общая функциональность для web-сервисов:
- Отправка логов через fluent.
- Отправка метрик через statsd.
- Разбор конфигурационных файлов.
