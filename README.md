## Запуск тестов
Запуск тестов из корня проекта
```shell
go test -v ./internal/computer_club
```
Также можно добавить свои тестовые файлы **inputID.txt** и **expected_outputID.txt** 
(ID -- любой уникальный в рамках директории идентификатор теста, например, число) в директорию **testdata**

## Запуск приложения
Запуск приложения из корня проекта
```shell
go run ./cmd/app/main.go input.txt
```
Вместо **input.txt** может быть любой другой файл с входящими запросами

## Запуск приложения в docker-контейнере
Сборка образа **computer-club-image**
```shell
docker build -t computer-club-image .
```
Запуск контейнера c названием **computer-club-container** для файла **input.txt**
```shell
docker run -d --name computer-club-container computer-club-image input.txt
```
