#!/bin/bash

# Остановка и удаление контейнеров, если они существуют
docker compose down

# Сборка новых образов
docker compose --build

# Запуск миграций
docker compose up