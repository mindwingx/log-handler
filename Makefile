run:
	cp .env.shared .env
	@echo "Starting Docker..."
	docker-compose up -d
	@echo "Docker containers started!"

down:
	@echo "Stopping docker containers..."
	docker-compose down
	@echo "Stopped!"

log:
	@echo "checking out the logger service into the docker container..."
	docker exec -it core_service bash