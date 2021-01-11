TAILWIND_PKGS ?= @tailwindcss/forms tailwindcss-cli

help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
.PHONY: help

setup: ## setup project workspace 
	mkdir -p local
	git clone git@github.com:manifold/qtalk.git
	@echo
	@echo NOTICE: These directories can be replaced with symlinks if already cloned elsewhere.
.PHONY: setup

dev: ## run dev server
	IMAGE_DIR=local go run cmd/zt100/main.go
.PHONY: dev

image: ## build docker container
	mkdir -p local/build
	GOOS=linux GOARCH=amd64 go build -o local/build/zt100 ./cmd/zt100
	docker build -t okta/zt100 .
.PHONY: image

docker: image
	docker run -it --rm -p 8080:8080 okta/zt100
.PHONY: docker

tailwind: ## compile tailwind from config
	mkdir _tailwind
	cd _tailwind && npm init -y > /dev/null && npm install $(TAILWIND_PKGS)
	cp tailwind.config.js _tailwind
	cd _tailwind && ./node_modules/.bin/tailwindcss-cli build -o ../static/vnd/tailwind-2.0.1.css
	rm -rf _tailwind
.PHONY: tailwind