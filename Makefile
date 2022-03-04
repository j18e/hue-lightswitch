# .env should contain:
# TARGET=user@ssh-host
# GOARCH=targets-go-arch
# HUE_BRIDGE=hue-bridge-address
include .env

NAME := hue-lightswitch
SERVICE_TPL := service.tpl
SERVICE_FILE := $(NAME).service

build:
	GOOS=linux GOARCH=$(GOARCH) go build -o ./$(NAME)

deploy:
	ssh $(TARGET) sudo systemctl stop hue-lightswitch
	scp ./$(NAME) ./config.yml ./.token $(TARGET):
	ssh $(TARGET) sudo systemctl start hue-lightswitch

deploy-systemd:
	sed "s/__HUE_HOST__/$(HUE_BRIDGE)/g" $(SERVICE_TPL) > $(SERVICE_FILE)
	scp ./$(NAME).service $(TARGET):
	ssh $(TARGET) sudo mv $(NAME).service /etc/systemd/system
	ssh $(TARGET) sudo chown root:root /etc/systemd/system/$(NAME).service
	ssh $(TARGET) sudo systemctl daemon-reload
	ssh $(TARGET) sudo systemctl start hue-lightswitch
	rm -f $(SERVICE_FILE)

logs:
	ssh $(TARGET) sudo journalctl -u $(NAME) -f

all: build deploy deploy-systemd
