DOCKER_COMPOSE_FILE = docker-compose.yml
APP_NAME = parallelizing-app
MAIN_FILE = cmd/api/main.go
EMAIL_WORKER_FILE = cmd/workers/send_email/main.go
PRIVATE_KEY_FILE := ec_private_key.pem
PUBLIC_KEY_FILE := ec_public_key.pem

.PHONY: docker-up docker-down run-app docker-clean start docker-rebuild generate-keys migration run-email-worker

docker-up:
	@echo "Subindo os serviços do Docker..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "Serviços do Docker estão rodando."

docker-down:
	@echo "Derrubando os serviços do Docker..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "Serviços do Docker foram parados."

docker-rebuild:
	@echo "Rebuildando os serviços do Docker..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d

docker-clean:
	@echo "Removendo os contêineres e volumes..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans	

run-email-worker:
	@echo "Iniciando worker de envio de e-mails..."
	@go run $(EMAIL_WORKER_FILE)

start:
	@echo "Iniciando aplicação Go..."
	@go run $(MAIN_FILE)

migration:
	@echo "Rodando as migrações..."
	go run database/migrations/main.go	

