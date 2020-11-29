.PHONY:  default  boltdb

default:
	go run main.go runserver

boltdb:
	GSD_SESSION_STORE=boltdb \
	GSD_SESSION_DIR=data/ \
	go run main.go runserver
