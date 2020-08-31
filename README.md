[![Build Status](https://travis-ci.com/Maksimall89/gogland.svg?branch=master)](https://travis-ci.com/Maksimall89/gogland) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Maksimall89_gogland&metric=alert_status)](https://sonarcloud.io/dashboard?id=Maksimall89_gogland)
### Info
Бот telegram для игры в [http://en.cx](http://en.cx). Он умеет играть в схватку, точки, МШ в линейную или заданную последовательность.
### Сборка
Для того, чтобы собрать билд выполните команды:
```go
go get
go build -o gogland.exe
```
### Тестирование
Перед началом тестирования в файле `config.json` заполните поля:
```json
  "TestNickName": "user",
  "TestPassword": "pass",
  "TestURLGame": "http://demo.en.cx/GameDetails.aspx?gid=1"
```
Или же установите переменные среды с этими же названиями. Приоритетным для сборки будут переменные среды.

Для запуска тестов введите:
```go
go test
```
Сейчас отключены тесты: `TestSentCodeJSON` и `TestGetPenaltyJSON`.
### Запуск игры
Для запуска вам необходимо объявить переменные среды или же сконфигурировать файл `config.json` поля:
```json
  "TelegramBotToken": "token",
  "OwnName": "nickOwn"
```
В них надо указать ваш telegram токен под которым будет запускаться бот (получить его необходимо у [@botfather](https://t.me/botfather)), а также ник игрока, которые будет администратором для бота. Или же установите переменные среды с этими же названиями. Приоритетным для сборки будут переменные среды.
### Логирование
Чтобы включить логирование событий в файл, необходимо установить переменную среды `Gogland_logs` в `1`, иначе логирование будет выполняться в console.
