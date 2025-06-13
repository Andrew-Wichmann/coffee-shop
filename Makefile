build: shop

shop:
	go build ./cmd/shop

.PHONY: run-dev
run-dev: shop
	./shop --config=config/config-dev.yaml

.PHONY: run-prod
run-prod: shop
	./shop --config=config/config-prod.yaml
