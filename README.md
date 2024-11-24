# Golang Microservice Project Template

This repository contains the basic layout and example code to quickly develop your microservices.
You should properly install the requisites so that the application could be correctly generated, compiled and deployed.

## Getting Started

Before you start, you should make sure you have already installed GNU Make. 
Unix-like systems can easily install the toolchain via package managers.
For Windows users, you can obtain the pre-compiled GNU toolchain on [their website](https://www.gnu.org/software/make/). 
To verify your installation, type the command `make --version` into your console.

We should first prepare the dependencies that the project requires by running the following commands:

```bash
# Download and update dependencies
make init
# Generate all files
make all
# Automated Initialization (wire)
cd cmd/example && wire
```

After necessary files have been generated, you can code and build with the following commands:

```bash
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go build -o ./bin/ ./...
./bin/server -conf ./configs
```

## Docker

The microservice supports running in Docker containers.
You can build your own image with the builtin Dockerfile and run your container easily by a single command.

```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

