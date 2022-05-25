stop-docker-test:
	@echo "Stopping Docker container..."
	@docker stop indexer || true && docker rm indexer || true
.PHONY: stop-docker-test

start-docker-test: stop-docker-test
	@echo "Starting Docker container..."
	@docker run --name indexer -e POSTGRES_USER=plural -e POSTGRES_PASSWORD=plural -e POSTGRES_DB=chain -d -p 5432:5432 postgres
.PHONY: start-docker-test

connect-db: 
	@echo "Connecting to dev db..."
	@psql "postgresql://plural:plural@localhost:5432/chain"
.PHONY: start-docker-test

SCHEMADIR = ./database/schema
deploy-schema: 
	@for f in $(shell ls ${SCHEMADIR}); do psql "postgresql://plural:plural@localhost:5432/chain" -f ${SCHEMADIR}/$${f}; done
