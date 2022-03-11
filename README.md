# Flight Paths Tracker

## Story
 There are over 100,000 flights a day, with millions of people and cargo being transferred around the world. With so many people, and different carrier/agency groups it can be hard to track where a person might be. In order to determine the flight path of a person, we must sort through all of their flight records.

## Goal
To create a microservice API that can help us understand and track how a particular person’s flight path may be queried. The API should accept a request that includes a list of flights, which are defined by a source and destination airport code. These flights may not be listed in order and will need to be sorted to find the total flight paths starting and ending airports.

Requirements:
-----------------

1. go 1.16+ 

How to run:
-----------------

1. Clone the repo

	git clone https://github.com/kumareswaramoorthi/flight-paths-tracker.git

2. Navigate to project directory 

	cd flight-paths-tracker 

3. Build the application by following command

	go build -o flight-paths-tracker main.go

4. Run the application by the following command 

	./flight-paths-tracker 


Alternatively, using docker,


1. Clone the repo

	git clone https://github.com/kumareswaramoorthi/flight-paths-tracker.git

2. Navigate to project directory 

	cd flight-paths-tracker 

3. Build the docker image by following command

	docker build -t flight-paths-tracker:1.0 .

4. Run the application by the following command 

	docker run -p 8080:8080 flight-paths-tracker:1.0


## **Swagger**

Swagger UI can be accessed at http://127.0.0.1:8080/swagger/index.html

## Documentation for API Endpoints


All URIs are relative to *http://127.0.0.1:8080*


## **1.Track Flight Paths**

Method | HTTP request | Description
------------- | ------------- | -------------
**FindSourceAndDestination** | **POST** /track | Finds source and destination from tickets


### Parameters

JSON body containing array of tickets.

### Response 

Array of string containing source and destination


### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### Example request and response

 - **Request**: `curl -H "Content-type: application/json" -d '{"tickets": [["ATL", "EWR"], ["SFO", "ATL"]]}' 127.0.0.1:8080/track`
 - **Response**: `["SFO","EWR"]`


## **2.Health Check**

Method | HTTP request | Description
------------- | ------------- | -------------
**HealthCheck** | **GET** / | Health Check API


### Parameters
This endpoint does not need any parameter.

### Response
HTTP Status 200.




