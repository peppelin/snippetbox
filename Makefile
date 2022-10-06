
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
	@echo "\n ws-pdf-publish project tasks:\n";
	@sed -n 's/^##/	/p' $<;
	@echo "\n";

##
## # Global tasks ---------------------------------------------------
##

## # Install task ---------------------------------------------------
##

## bd-up	start the mysql from the docker-compose.yaml
db-up:
	@echo "\n> Starting docker-compose";
	@docker-compose up -d
## bd-down stop the mysql from the docker-compose.yaml
db-down:
	@echo "\n> Stopping docker-compose";
	@docker-compose stop