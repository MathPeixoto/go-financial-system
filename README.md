# go-financial-system
Repository created to practice Go development and DevOps

Project using these technologies

- Gin framework: for creating the REST APIs
- Migrate tool: for migrating databases
- Sqlc: for generating Go code from SQL (Generates models, daos and more)
- Postgres
- Docker


## Installing
### Migrate tool
`sudo apt install python3-migrate`

### Sqlc (Go > 1.17)
`go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`


### To generate Go code from SQL: 

- Write SQL in './db/query' folder
- run: `make sqlc`

The code will be generate within './db/sqlc' folder



