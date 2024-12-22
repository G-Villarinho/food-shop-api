DOCKER_COMPOSE_FILE = docker-compose.yml
APP_NAME = parallelizing-app
MAIN_FILE = cmd/api/main.go
PRIVATE_KEY_FILE := ec_private_key.pem
PUBLIC_KEY_FILE := ec_public_key.pem

.PHONY: docker-up docker-down run-app docker-clean start docker-rebuild generate-keys

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

run-app:
	@echo "Iniciando a aplicação Go..."
	go run $(MAIN_FILE)

start:
	@echo "Iniciando aplicação Go..."
	make run-app

migration:
	@echo "Rodando as migrações..."
	go run database/migrations/main.go	

generate-keys:
	@if [ ! -f $(PRIVATE_KEY_FILE) ]; then \
		echo "Generating ECDSA private key..."; \
		openssl ecparam -genkey -name prime256v1 -noout -out $(PRIVATE_KEY_FILE); \
		echo "Private key saved in $(PRIVATE_KEY_FILE)"; \
	else \
		echo "Private key already exists: $(PRIVATE_KEY_FILE)"; \
	fi

	@if [ ! -f $(PUBLIC_KEY_FILE) ]; then \
		echo "Extracting public key from the private key..."; \
		openssl ec -in $(PRIVATE_KEY_FILE) -pubout -out $(PUBLIC_KEY_FILE); \
		echo "Public key saved in $(PUBLIC_KEY_FILE)"; \
	else \
		echo "Public key already exists: $(PUBLIC_KEY_FILE)"; \
	fi