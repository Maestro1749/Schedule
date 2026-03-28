migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ScheduleDB?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ScheduleDB?sslmode=disable" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)