.PHONY:  default  boltdb

default:
	go run main.go runserver

boltdb:
	GSD_SESSION_STORE=boltdb \
	GSD_SESSION_FILE=data/bolt.db \
	go run main.go runserver
