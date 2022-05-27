# Stop the docker psql server
stop-docker-test:
	@echo "Stopping Docker container..."
	@docker stop indexer || true && docker rm indexer || true
.PHONY: stop-docker-test

# Restart the docker psql server
start-docker-test: stop-docker-test
	@echo "Starting Docker container..."
	@docker run --name indexer -e POSTGRES_USER=plural -e POSTGRES_PASSWORD=plural -e POSTGRES_DB=chain -d -p 5432:5432 postgres
.PHONY: start-docker-test

# Connect the docker psql server using test creds
connect-db: 
	@echo "Connecting to dev db..."
	@psql "postgresql://plural:plural@localhost:5432/chain"
.PHONY: start-docker-test

# Deploy the schema NOTE: this needs to be run before the parser starts
SCHEMADIR = ./schema
deploy-schema: 
	@for f in $(shell ls ${SCHEMADIR}); do psql "postgresql://plurallabs:GPBNV1Ufg7vq9xkuBOCc@135.181.252.91:5433/parser" -f ${SCHEMADIR}/$${f}; done
.PHONY: deploy-schema

install: go.sum 
	go install -mod=readonly ./cmd/cosmos-snapshot-parser

go.sum: go.mod
	echo "Ensure dependencies have not been modified ..." >&2
	go mod verify
	go mod tidy

