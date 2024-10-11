# build app image (remove previous image version)
build:
	- docker compose stop apod
	- docker rmi -f apod-apod
	docker compose build

# start services
up:
	docker compose up -d

# start app service
up-app:
	docker compose up -d apod

down:
	docker compose down

# stop all services
stop:
	docker compose stop

# stop app service
stop-app:
	docker compose stop apod

# get services logs
logs:
	docker compose logs -f

# get app service logs
logs-app:
	docker compose logs -f apod

# build swagger api documentation ('swag' tool should be available in PATH)
# link: https://github.com/swaggo/swag
api:
	swag init -g ./cmd/main.go --output=docs --parseInternal=true

# build swagger api documentation ('swag' tool should be available in PATH)
api-fmt:
	swag fmt -dir ./cmd
