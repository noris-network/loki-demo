LOKI_VERSION ?= v0.4.0
BIN_DIR ?= bin
LOG_DIR ?= logs
DATA_DIR ?= data
LOG_GEN_BIN = $(BIN_DIR)/log_gen

install: build download

build: main.go
	go build -o $(LOG_GEN_BIN) main.go

download: download/loki download/promtail

download/loki:
	curl -fSL -o "$(BIN_DIR)/loki.gz" "https://github.com/grafana/loki/releases/download/$(LOKI_VERSION)/loki-linux-amd64.gz"
	gunzip $(BIN_DIR)/loki.gz
	chmod a+x $(BIN_DIR)/loki

download/promtail:
	curl -fSL -o "$(BIN_DIR)/promtail.gz" "https://github.com/grafana/loki/releases/download/$(LOKI_VERSION)/promtail-linux-amd64.gz"
	gunzip $(BIN_DIR)/promtail.gz
	chmod a+x $(BIN_DIR)/promtail

run/log_gen:
	./$(LOG_GEN_BIN)

run/loki:
	./$(BIN_DIR)/loki -config.file loki.yml

run/promtail:
	./$(BIN_DIR)/promtail -config.file promtail.yml

run/docker/up:
	docker-compose up -d

run/docker/down:
	docker-compose down