
# Env & Vars --------------------------------------------------------

include .env
export $(shell sed 's/=.*//' .env)

go_path = PATH=${PATH}:~/go/bin GOPATH=~/go 
go_env = $(go_path) GO111MODULE=on

docker_images = mysql
# test_packages = $(shell go list ./... | grep -v 'snippetbox$$' | grep -v '/test')

# Tasks -------------------------------------------------------------

## # Help task ------------------------------------------------------
##

## help		Print project tasks help
help: Makefile
	@echo "\n Snippetbox project tasks:\n";
	@sed -n 's/^##/	/p' $<;
	@echo "\n";

##
## # Global tasks ---------------------------------------------------
##

## # Install task ---------------------------------------------------
##

## db-up	start the mysql from the docker-compose.yaml
db-up:
	@echo "\n> Starting docker-compose";
	@docker-compose up -d
## db-down stop the mysql from the docker-compose.yaml
db-down:
	@echo "\n> Stopping docker-compose";
	@docker-compose stop
## db-init deletes any previos snippetbox database, creates a new one and populates it with basic info
db-init:
	@echo "\n> populating database";
	@mysql -h 127.0.0.1 -uroot -p < setup_mysql.sql
## cert-gen generates the tls certificates and stores intop tls folder
cert-gen:
	@mkdir tls; 
	@cd tls && go run /Users/gargallo/.gvm/gos/go1.19.2/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost