# go-financial-system

go-financial-system is a Go-based project that provides REST and gRPC services for fictional banking transactions. The project is built on the [Gin](https://github.com/gin-gonic/gin) framework and utilizes [Postgres](https://www.postgresql.org/) and [Redis](https://redis.io/) for data storage.


## Installation

To install the project dependencies, use the package manager Makefile.


## Usage
The project can be started locally using the make run command or using the docker-compose.yml file.

`make run`

`docker-compose up`


## Deployment
The project is deployed on AWS using EKS and utilizes CI/CD for running unit tests and deploying updates.

## Contribution
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
