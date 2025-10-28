WIP 

Social - How to run it 

## Direnv

direnv to be installed depending on OS. For example, on Ubuntu : `sudo apt search direnv`
then everytime an environment variable is modified in .envrc, run:
`direnv allow .` from the project base directory.  
Otherwise a warning such as "direnv: error (..)/social/.envrc is blocked. Run `direnv allow` to approve its content" will appear.
 
## Hot reload

requires Air: 

`go install github.com/air-verse/air@latest`

initiate : `air init` will create .air.toml. Some fields need customisation.
e.g.  
```
tmp_dir = "bin"
bin = "./bin/main"
cmd = "go build -o ./bin/main ./cmd/api"
exclude_dir = ["assets", "bin", "vendor", "testdata", "web", "docs", "scripts"]
pre_cmd = ["make gen-docs"]
```

then run : `air`

Note : the make gen-docs won't work unless the Makefile is properly configured and go swagg installed (see bellow). It should give only a warning.

## PostgreSQL db + redis docker file 

my docker-compose.yml to be customized : 
```
version: '3'
services:
  database:
    image: 'postgres:latest'
    ports: 
      - 5432:5432
    container_name: postgre-db-social
    environment:
      POSTGRES_USER: admin # The PostgreSQL user (useful to connect to the database)
      POSTGRES_PASSWORD: adminpassword # The PostgreSQL password (useful to connect to the database)
      POSTGRES_DB: social # The PostgreSQL default database (automatically created at first launch)
    volumes : 
      - ${PWD}/db-data/:/var/lib/postgresql/data/

  redis:
    image: redis:alpine
    container_name: redis-cache
    ports:
      - "6379:6379"
    restart: unless-stopped
    command: redis-server --save 60 1 --loglevel warning


  redis-commander:
    container_name: redis-commander
    hostname: redis-commander
    image: rediscommander/redis-commander:latest
    environment:
    - REDIS_HOST=redis-cache
    ports:
    - "127.0.0.1:8081:8081"
    depends_on:
    - redis
    restart: unless-stopped
```

 then run :
 `docker compose up -d `
and stop it with:
`docker compose down`

## Migrations & Seeding

migrations rely on a Makefile such as this one:

```
include .envrc
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: seeding
seeding:
@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
@swag init -g main.go -d cmd/api,internal/store && swag fmt

.PHONY: test
test:
@go test -v ./...
```
 
Run migrations via:
`make migrate-up`

Once migrations are done, one can run:
`make seeding`

In order to verify the seeding did happen, one can use either a free solution such as dbeaver-ce or the PostgreSQL interactive terminal : 
̀psql -h localhost -p 5432 -U admin -d social` 
(will require the same POSTGRES_PASSWORD set in the docker-compose)

If the social db is messed up, the table schema-migrations has the dirty column set to true.  Once the root cause is corrected, set the dirty property of version xx back to false :

```
update schema_migrations set dirty =false where version=xx;
```


## Swagger

swagger can be installed with:
`go install github.com/swaggo/swag/cmd/swag@latest `
then run:
`swagg init` (from the project root directory)

Once the handlers are documented properly, documentation can be accessed through:
ip:port/swagger/doc.json
and the endpoints via:
ip:port/v1/swagger/index.html

Note : the `pre_cmd` field in the .air.toml needs to be filled with  `["make gen-docs"]` and the Makefile gen-docs entry  needs to be set accordingly.

## Tips

+ In order to test graceful shutdown via SIGINT, the .air.toml needs to be modified with : 
send_interrupt = true 
and 
kill_delay = "10s"


