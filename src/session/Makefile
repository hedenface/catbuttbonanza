.DEFAULT_GOAL := session
FORCE: ;

copy-resources: FORCE
	cp -r ../../resources .

session: FORCE copy-resources
	go build -o session .

run: session
	./session

docker: session
	docker build -t cbb-session .
