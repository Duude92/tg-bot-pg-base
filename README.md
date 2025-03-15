### Basic project for Telegram bot written in Go, using gorm with postgresql.
___
### Usage:
docker run --rm -e ApiKey="\<Telegram bot Api key\>" -e DB_HOST="127.0.0.1" -e DB_PORT=5432 -e DB_USER="admin" -e DB_PASSWORD="adminpass" -e DB_NAME="\<DB Name\>"  -t test-tg-pg-bot
