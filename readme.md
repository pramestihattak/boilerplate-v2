# Go Boilerplate with gRPC + gRPC Gateway

In this boring weekend, I spent my time creating a boilerplate code to build application with Go, gRPC and gRPC Gateway (to support RESTful JSON API from frontend into gRPC server).

## How to run
1. Create new docker postgres instance `docker run --name boilerplatedb -e POSTGRES_PASSWORD=secret -p 5438:5432 -d postgres`
1. Copy sample.config => config
2. `make migrate-up`
3. `make run-backend`

## How to create new migration
1. `cd migrations`
2. create new migration file with increment prefix under sql folder example: `0003_create_events_table.sql`
3. Write the query, for up and down
4. Run `make migrate-up`
5. If you want to rollback `make migrate-down`

2024
