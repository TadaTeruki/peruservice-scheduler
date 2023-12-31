include .env

.PHONY: create-db
create-db:
	mkdir -p postgres/data

.PHONY: delete-db
delete-db:
	rm -r postgres/data 

.PHONY: reset-db
reset-db: delete-db create-db

.PHONY: login-to-db
login-to-db:
	docker exec -it scheduler-db psql -U ${DB_USER} -d ${DB_NAME}
