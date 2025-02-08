# go-movie
Microservices with golang

## Motivation
The purpose of this project is to learn more about microservices with the help of the book `Microservices with Go: Building Scalable and Reliable Microservces with Go`. That book is written by Alexander Shuiskov a lead engineer at Uber and covers topics such as service discovery, synchronous and asynchronous communication, Kubernetes, testing, reliability and availablity, telemetry and monitoring as well as alerting.
The original code can be found [here](https://github.com/PacktPublishing/Microservices-with-Go)

## Overview
The Architecture consists of three services, two databases, a message bus and a service registry. The three services are the movie serivce acting as an API gateway and is responsible to handle incoming requests. The metadata service holds general data about movies and saves it in a Postgres database (currently in-memory). The ratings service manages ratings for movies and saves its data in a postgres database as well (also currently in-memory). The backbone for asynchronous communication is a kafka instance.

![Diagram of the architecture](/diagram.drawio.png)

## Differences to the book authors code
Even though I followed the authors book example there are some major differences between the authors code and mine:
- the author uses the classic N-tier application design while I implemented the domain-centric architecture as described in my [go-social repository](https://github.com/karaMuha/go-social)
- I used testcontainers to test against real databases while the author uses mocks and in-memory solutions
- I integrated Kafka into the systems architecture while the author provides an example that is just in parts integrated to the actual code

## How to run the app

## Feedback
My overall impression is that the target audience for this book are software engineers who have some production experience and want to make their first step into the world of microservices. On the one hand he author discusses and explains some of the key concepts regarding microservice architecutre like observability, alerting, telemetry, async communication etc. but one the other hand fails to implement these concepts in what I would consider a production-grade example.

## Todos
- dockerize the serivces
- implement graceful shutdown in each service and make sure to inform kafka in case of service shutdown to prevent message loss
- implement a solution for distributed tracing
- save logs and telemetric data in a seperate database
- send feedback to author