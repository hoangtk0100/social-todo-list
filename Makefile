include app.env
export

mysql:
	docker run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):3306 -e MYSQL_ROOT_PASSWORD=$(DB_PASSWORD) -d $(MYSQL_IMAGE)

createdb:
	docker exec -it $(DB_CONTAINER_NAME) mysql -u$(DB_USERNAME) -p$(DB_PASSWORD) -e "CREATE DATABASE \`$(DB_NAME)\` DEFAULT CHARACTER SET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_0900_ai_ci';"

dropdb:
	docker exec -it $(DB_CONTAINER_NAME) mysql -u$(DB_USERNAME) -p$(DB_PASSWORD) -e "DROP DATABASE \`$(DB_NAME)\`;"

db:
	docker exec -it $(DB_CONTAINER_NAME) mysql -u$(DB_USERNAME) -p$(DB_PASSWORD)

server:
	go run .

build:
	go build -o app

outenv:
	./app outenv

outenvfile:
	./app outenv > .env

.PHONY: mysql createdb dropdb db server build outenv outenvfile