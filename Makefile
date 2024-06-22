# Define variables
APP_NAME = simple-badwords-api
DOCKER_COMPOSE_FILE = docker-compose.yml

# Build the Docker image
build:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

# Run the application using Docker Compose
up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up

# Run the application in detached mode
up-detached:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop and remove the containers
down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Remove the Docker image
clean:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --rmi all

# Restart the application
restart: down up

# Tail the logs of the running containers
logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Show the status of the running containers
status:
	docker-compose -f $(DOCKER_COMPOSE_FILE) ps
