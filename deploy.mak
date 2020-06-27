#
#
#

SHELL := /bin/bash
COPY  := cp
COPY_DIR := $(COPY) -r
DELETE := rm -f

BUILD_DIR   := build
BUILD_BIN_DIR := $(BUILD_DIR)/bin

DEPLOY_DEV_DIR := $(BUILD_DIR)/develop
SRC_DEV_CONFIG := ./internal/res/configs/develop
DEVELOP_SERVICES := \
								client_rx \
								client_tx \
								server_rx \
								server_tx \
								pktservice

BUILD_BINS :=  \
								$(BUILD_BIN_DIR)/client_rx \
								$(BUILD_BIN_DIR)/client_tx \
								$(BUILD_BIN_DIR)/server_rx \
								$(BUILD_BIN_DIR)/server_tx \
								$(BUILD_BIN_DIR)/pktservice

RT_DEV_BIN_CLIENT_RX := $(DEPLOY_DEV_DIR)/client_rx/anyclient
RT_DEV_BIN_CLIENT_TX := $(DEPLOY_DEV_DIR)/client_tx/anyclient
RT_DEV_BIN_SERVER_RX := $(DEPLOY_DEV_DIR)/server_rx/anyserver
RT_DEV_BIN_SERVER_TX := $(DEPLOY_DEV_DIR)/server_tx/anyserver
RT_DEV_BIN_GSPKTSERVICE  := $(DEPLOY_DEV_DIR)/pktservice/pktservice

RT_DEV_CONF_CLIENT_RX := $(DEPLOY_DEV_DIR)/client_rx/config.yaml
RT_DEV_CONF_CLIENT_TX := $(DEPLOY_DEV_DIR)/client_tx/config.yaml
RT_DEV_CONF_SERVER_RX := $(DEPLOY_DEV_DIR)/server_rx/config.yaml
RT_DEV_CONF_SERVER_TX := $(DEPLOY_DEV_DIR)/server_tx/config.yaml
RT_DEV_CONF_GSPKTSERVICE  := $(DEPLOY_DEV_DIR)/pktservice/config.yaml

.PHONY: all
all: deploy-dev
deploy-dev: deploy-dev-bin deploy-dev-conf
deploy-dev-bin: $(RT_DEV_BIN_CLIENT_RX) \
								$(RT_DEV_BIN_CLIENT_TX) \
								$(RT_DEV_BIN_SERVER_RX) \
								$(RT_DEV_BIN_SERVER_TX) \
								$(RT_DEV_BIN_GSPKTSERVICE)

$(RT_DEV_BIN_CLIENT_RX): $(BUILD_BIN_DIR)/anyclient
	$(COPY) $^ $@

$(RT_DEV_BIN_CLIENT_TX): $(BUILD_BIN_DIR)/anyclient
	$(COPY) $^ $@

$(RT_DEV_BIN_SERVER_RX): $(BUILD_BIN_DIR)/anyserver
	$(COPY) $^ $@

$(RT_DEV_BIN_SERVER_TX): $(BUILD_BIN_DIR)/anyserver
	$(COPY) $^ $@

$(RT_DEV_BIN_GSPKTSERVICE): $(BUILD_BIN_DIR)/pktservice
	$(COPY) $^ $@

deploy-dev-conf: $(RT_DEV_CONF_CLIENT_RX) \
								 $(RT_DEV_CONF_CLIENT_TX) \
								 $(RT_DEV_CONF_SERVER_RX) \
								 $(RT_DEV_CONF_SERVER_TX) \
								 $(RT_DEV_CONF_GSPKTSERVICE)

$(RT_DEV_CONF_CLIENT_RX): $(SRC_DEV_CONFIG)/client_rx/config.yaml
	$(COPY) $^ $@

$(RT_DEV_CONF_CLIENT_TX): $(SRC_DEV_CONFIG)/client_tx/config.yaml
	$(COPY) $^ $@

$(RT_DEV_CONF_SERVER_RX): $(SRC_DEV_CONFIG)/server_rx/config.yaml
	$(COPY) $^ $@

$(RT_DEV_CONF_SERVER_TX): $(SRC_DEV_CONFIG)/server_tx/config.yaml
	$(COPY) $^ $@

$(RT_DEV_CONF_GSPKTSERVICE): $(SRC_DEV_CONFIG)/pktservice/config.yaml
	$(COPY) $^ $@

$(BUILD_BINS):

.PHONY: clean
clean:
	@echo "have to delete something"

# EOF
