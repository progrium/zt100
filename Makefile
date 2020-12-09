.PHONY: tailwind

TAILWIND_PKGS ?= @tailwindcss/forms tailwindcss-cli

dev:
	echo Dev

tailwind:
	mkdir _tailwind
	cd _tailwind && npm init -y > /dev/null && npm install $(TAILWIND_PKGS)
	cp tailwind.config.js _tailwind
	cd _tailwind && ./node_modules/.bin/tailwindcss-cli build -o ../static/tailwind-2.0.1.css
	rm -rf _tailwind