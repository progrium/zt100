TAILWIND_PKGS ?= @tailwindcss/forms tailwindcss-cli

help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
.PHONY: help

setup: ## setup project workspace 
	git clone git@github.com:progrium/macdriver.git
	git clone git@github.com:manifold/qtalk.git
	git clone git@github.com:manifold/tractor.git
	@echo
	@echo NOTICE: These directories can be replaced with symlinks if already cloned elsewhere.
.PHONY: setup

dev: ## TBD
	echo Dev
.PHONY: dev



tailwind: ## compile tailwind from config
	mkdir _tailwind
	cd _tailwind && npm init -y > /dev/null && npm install $(TAILWIND_PKGS)
	cp tailwind.config.js _tailwind
	cd _tailwind && ./node_modules/.bin/tailwindcss-cli build -o ../static/tailwind-2.0.1.css
	rm -rf _tailwind
.PHONY: tailwind