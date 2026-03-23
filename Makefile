include .env
export

service-run:
	@go run main.go

service-deploy:
	docker-compose up -d 

service-undeploy:
	docker compose down 

check-db:
	docker exec -it database_container psql -U user -d postgres -c "SELECT * FROM Tasks;"