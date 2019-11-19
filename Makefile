LOKI_VERSION ?= v0.4.0
BIN_DIR ?= bin
LOG_DIR ?= logs
DATA_DIR ?= data
LOG_GEN_BIN = $(BIN_DIR)/log_gen

install: build download

build: main.go
	go build -o $(LOG_GEN_BIN) main.go

download: download/loki download/promtail

clean: clean/logs clean/data clean/bin clean/docker

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

clean/logs:
	find $(LOG_DIR) -regex ".*log" -delete

clean/data:
	find $(DATA_DIR)/loki -type d -not -regex "data/loki" -exec rm -rf "{}" \; || /bin/true
	find $(DATA_DIR)/promtail -type f -not -regex ".*gitkeep$$" -exec rm "{}" \;

clean/bin:
	find $(BIN_DIR) -type f -not -regex ".*gitkeep$$" -delete

clean/docker:
	docker volume ls | tail -n+2 | grep loki-demo | awk '{print $$2}' | xargs docker volume rm || /bin/true