![workflow badge](https://github.com/karaMuha/go-movie/actions/workflows/ci.yaml/badge.svg)
![workflow badge](https://github.com/karaMuha/go-movie/actions/workflows/cd.yaml/badge.svg)

# go-movie
Microservices with golang

## Motivation
The purpose of this project is to learn more about microservices with the help of the book `Microservices with Go: Building Scalable and Reliable Microservces with Go`. That book is written by Alexander Shuiskov a lead engineer at Uber and covers topics such as service discovery, synchronous and asynchronous communication, Kubernetes, testing, reliability and availablity, telemetry and monitoring as well as alerting.
The original code can be found [here](https://github.com/PacktPublishing/Microservices-with-Go)

## Overview
The Architecture consists of three services, two databases, a message bus and a service registry. The three services are the movie serivce acting as an API gateway and is responsible to handle incoming requests. The metadata service holds general data about movies and saves it in a Postgres database. The ratings service manages ratings for movies and saves its data in a postgres database as well. The backbone for asynchronous communication is a kafka instance.
Messages that were failed to publish and consumed messages that were failed to be processed are saved in the database of the respective service. A cronjob loops thorugh the data periodically and tries to process them.

![Diagram of the architecture](/microservice.drawio.png)

## Differences to the books code
Even though I followed the authors example project there are some major differences between the authors code and mine:
- the books example uses the classic N-tier application design while I implemented the domain-centric architecture as described in my [go-social repository](https://github.com/karaMuha/go-social)
- I used testcontainers to test against real databases while the author uses mocks and in-memory solutions
- I integrated Kafka into the systems architecture while the author provides an example that is just in parts integrated to the actual project
- the books example saves each individual rating and calculates the average rating for a movie on each get request. I followed the CQRS principle by implementing a read and a write table. On each rating submission the rating is saved in the write table and the calculated average rating for that movie is updated in the read table. This leads to a performance increase because the average rating does not have to be calculated on each get request.

## How to run the app
#### Requirenments
- Go version >= 1.23.0
- Protoc
- protoc-gen-go
- protoc-gen-go-grpc
- Docker / Docker Desktop

Clone the repo and run `make prepare`. This will generate the folder structure to persists data from the postgres containers and generate the protobuf code needed for gRPC. Then start Docker Desktop and run `make run`.

## Feedback
My overall impression is that the target audience for this book are software engineers who have some production experience and want to make their first step into the world of microservices.
On the one hand he author discusses and explains some of the key concepts regarding microservice architecutre like service discovery, async communication, observability etc. but one the other hand fails to implement these concepts in what I would consider a production-grade example.
However I can still recommend this literature to those with no prior experience with microservices.

## Todos
- implement a solution for observability and tracing
- visualize telemetric data
- send feedback to author
- implement users serivce for registration and login
- implement auth and rate limiting in movie service

## What next?
To get deeper into microservices and learn how to design and implement a complex event-driven architecture, I will study the literature `Event-Driven Architecture in Golang: Building complex systems with asynchronicity and eventual consistency` from the author Michael Stack.