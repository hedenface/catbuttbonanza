.DEFAULT_GOAL := ui
FORCE: ;

html-templates: FORCE
	./html-templates/convert-to-static-functions.sh

ui: FORCE html-templates
	go build -o ui .

run: ui
	./ui
