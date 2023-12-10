.DEFAULT_GOAL := ui
FORCE: ;

html-templates: FORCE
	./html-templates/convert-to-static-functions.sh

copy-resources: FORCE
	cp -r ../../resources .

ui: FORCE html-templates copy-resources
	go build -o ui .

run: ui
	./ui

docker: ui
	docker build -t cbb-ui .
