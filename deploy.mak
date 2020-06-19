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

.PHONY: all
all: deploy-dev
deploy-dev: $(RT_DEV_BIN_CLIENT_RX) \
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

$(BUILD_BINS):

.PHONY: clean
clean:
	@echo "have to delete something"

# EOF
