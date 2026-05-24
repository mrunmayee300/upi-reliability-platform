.PHONY: help tree docs validate-contracts infra-up infra-down topics

help:
	@echo "UPI Transaction Intelligence Platform"
	@echo ""
	@echo "  make tree              - Print monorepo structure"
	@echo "  make docs              - Open architecture index (prints path)"
	@echo "  make infra-up          - Start Phase 2 infra stack"
	@echo "  make topics            - Create Kafka topics"
	@echo "  make infra-down        - Stop infra stack"
	@echo "  make validate-contracts - Validate OpenAPI/AsyncAPI (Phase 3)"

tree:
	@powershell -Command "Get-ChildItem -Recurse -Directory | Where-Object { $$_.FullName -notmatch 'node_modules|\.git|\.terraform' } | ForEach-Object { $$_.FullName.Replace((Get-Location).Path + '\\', '') }"

docs:
	@echo "Architecture: docs/architecture/README.md"
	@echo "Kafka:       docs/kafka/topics.md"
	@echo "Contracts:   shared/contracts/"

validate-contracts:
	@echo "Contract validation wired in Phase 3 (spectral/redocly)"

infra-up:
	docker compose up -d

topics:
	bash ./scripts/kafka/create-topics.sh

infra-down:
	docker compose down
