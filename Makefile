.DEFAULT_GOAL := catbuttbonanza
FORCE: ;

catbuttbonanza: FORCE
	cd src/ui && make
	cd src/auth && make
	cd src/session && make
