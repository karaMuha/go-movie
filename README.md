# go-movie
Microservices with golang

## Motivation
The purpose of this project is to learn more about microservices with the help of the book `Microservices with Go: Building Scalable and Reliable Microservces with Go`. That book is written by Alexander Shuiskov a lead engineer at Uber and covers topics such as service discovery, synchronous and asynchronous communication, Kubernetes, testing, reliability and availablity, telemetry and monitoring as well as alerting.
The original code can be found [here](https://github.com/PacktPublishing/Microservices-with-Go)

## Overview
The Architecture consists of three services, two databases, a message bus and a service registry. The three services are the movie serivce acting as an API gateway and is responsible to handle incoming requests. The metadata service holds general data about movies and saves it in a Postgres database (currently in-memory). The ratings service manages ratings for movies and saves its data in a postgres database as well (also currently in-memory). The backbone for asynchronous communication is a kafka instance.

![Diagram of the architecture](/diagram.drawio.png)