.DEFAULT_GOAL := catbuttbonanza
FORCE: ;

html-templates: FORCE
	./html-templates/convert-to-static-functions.sh

catbuttbonanza: FORCE html-templates
	go build -o catbuttbonanza .

run: catbuttbonanza
	./catbuttbonanza
