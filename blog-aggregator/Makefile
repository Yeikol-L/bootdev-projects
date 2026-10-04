ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

run:
	go run . $(ARGS)
migrate:
	goose postgres -dir ./sql/schema "postgres://yeikollorenzo:@localhost:5432/gator?sslmode=disable" $(ARGS)

%:
	@:
