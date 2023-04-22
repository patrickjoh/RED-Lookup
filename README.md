# Assignment 2 - Renewable Energy Data Lookup and Webhook Notification
This project is a REST web application in Golang that provides the client with the ability to retrieve information about developments related to renewable energy production for and across countries. This is done by using data from two existing services. The service also allows for notification registration using webhooks, and stores thes persistently using Firebase.

# Table of Contents
* [Introduction](#assignment-2---renewable-energy-data-lookup-and-webhook-notification)
* [Deployment](#deployment)
    * [Preparation:](#preparation)
    * [Docker](#docker)
        * [Alternative 1 - Docker-Compose](#alternative-1---docker-compose)
        * [Alternative 2 - Dockerfile](#alternative-2---dockerfile)
    * [Golang](#golang)
* [Usage](#usage)
    * [Endpoints](#endpoints)
        * [Current percentage of renewables](#current-percentage-of-renewables)
        * [Historical percentages of renewables](#historical-percentages-of-renewables)
        * [Notifications](#notifications)
            * [Register a webhook](#register-a-webhook)
            * [Delete a registered webhook](#delete-a-registered-webhook)
            * [View registered webhooks](#view-registered-webhook)
            * [View all registered webhooks](#view-all-registered-webhooks)
            * [Webhook Invocation (Upon Trigger)](#webhook-invocation-upon-trigger)
        * [Status](#status)
* [Dependencies](#dependencies)
* [Contributors](#contributors)


# Deployment

## Preparation:

* Create a Firebase project and enable the Firestore database.
* Create a service account for the project and download the JSON file.
* Copy the JSON file to the root directory of the project and rename it to `firebase.json`.

## Docker
Pre-requisites:
* Docker 20.10.24 or higher

### Alternative 1 - Docker-Compose

**Run the following commands in the root directory of the project:**

*Build and run the project using docker-compose:*
```bash
docker-compose up --build
```
**or**

```bash
docker-compose up --build -d
```
To deploy the container in the background.

### Alternative 2 - Dockerfile

**Run the following commands in the root directory of the project:**

*Build the project:*
```bash
docker build -t REDL .
```

*Create an instance of the image mapped to port 80 on your host:*
```bash
docker run -p 80:8080 REDL
```
**or**
```bash
docker run -p 80:8080 -d REDL
```
To deploy the container in the background.

## Golang

Pre-requisites:
* Golang 1.20 or higher

**Run the following commands in the root directory of the project:**

*Build the project:*
```go
go mod tidy
go build -o ./app ./cmd/main.go
```

*Run the project:*
```go
./app
```

# Usage

## Endpoints

The service has four endpoints:
```http
/energy/v1/renewables/current
/energy/v1/renewables/history
/energy/v1/notifications/
/energy/v1/status/
```

If the web service is running on localhost, port 8080,
the full paths to the resources would look like this:
```http
http://localhost:8080/energy/v1/renewables/current
http://localhost:8080/energy/v1/renewables/history
http://localhost:8080/energy/v1/notifications/
http://localhost:8080/energy/v1/status/
```

## Current percentage of renewables

This endpoint focuses on returning the latest percentage of renewables in the energy mix for a given country, as well as for all countries.

**Requests should be sent in the following format:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

```http
Method: GET
Path: /energy/v1/renewables/current/{country?}{?neighbours=bool?}
```

`{country?}` refers to an optional country 3-letter code.

`{?neighbours=bool?}` refers to an optional parameter indicating whether neighbouring countries' values also should be shown.

*Example requests:*
```http
/energy/v1/renewables/current
/energy/v1/renewables/current/nor
/energy/v1/renewables/current/nor?neighbours=true
```

**- Request:**
```http
Method: GET
Path: /energy/v1/renewables/current
```

**- Response:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

```json
[
    {
        "name": "Algeria",
        "isoCode": "DZA",
        "year": "2021",
        "percentage": 0.26136735
    },
    {
        "name": "Argentina",
        "isoCode": "ARG",
        "year": "2021",
        "percentage": 11.329249
    },
    {
        "name": "Australia",
        "isoCode": "AUS",
        "year": "2021",
        "percentage": 12.933532
    },
    ...
]
```

**- Request:**
```http
Method: GET
Path: /energy/v1/renewables/current/nor
```

**- Response:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

```json
{
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2021",
    "percentage": 71.558365
}
```

**- Request:**
```http
Method: GET
Path: /energy/v1/renewables/current/nor?neighbours=true
```

**- Response:**
```json
[
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "2021",
        "percentage": 71.558365
    },
    {
        "name": "Finland",
        "isoCode": "FIN",
        "year": "2021",
        "percentage": 34.61129
    },
    {
        "name": "Russia",
        "isoCode": "RUS",
        "year": "2021",
        "percentage": 6.6202893
    },
    {
        "name": "Sweden",
        "isoCode": "SWE",
        "year": "2021",
        "percentage": 50.924007
    }
]
```

## Historical percentages of renewables
---

This endpoint returns historical percentages of renewables in the energy mix, including individual levels, as well as mean values for individual or selections of countries.

**Requests should be sent in the following format:**

```http
Method: GET
Path: /energy/v1/renewables/history/{country?}{?begin=year&end=year?}
```

`{country}` refers to an optional country 3-letter code.

`{?begin=year&end=year}` refers to an optional range for the selected country.

*Example requests:*
```http
/energy/v1/renewables/history
/energy/v1/renewables/history/nor
/energy/v1/renewables/history/nor?begin=1970
/energy/v1/renewables/history/nor?begin=1960&end=1970
```

**- Request:**
```http
Method: GET
Path: /energy/v1/renewables/history/nor?begin=1960&end=1970
```

**- Response:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

Show percentage for a single country in the given range
```json
[
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "1965",
        "percentage": 67.87996
    },
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "1966",
        "percentage": 65.3991
    },
    ...
]
```

**- Request:**
```http
Method: GET
Path: /energy/v1/renewables/history/nor?begin=1960&end=1970
```

**- Response:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

Shows mean percentage for all countries
```json
[
    {
        "name": "United Arab Emirates",
        "isoCode": "ARE",
        "percentage": 0.0444305504
    },
    {
        "name": "Argentina",
        "isoCode": "ARG",
        "percentage": 9.131337212280702
    },
    {
        "name": "Australia",
        "isoCode": "AUS",
        "percentage": 5.3000481596491245
    },
    ...
]
```

## Notifications
Users can register webhooks that are triggered by the service based on specified events, specifically if information about given countries (or any country) is invoked, where the minimum frequency can be specified. Users can register multiple webhooks. The service saves these registrations in a Firebase DB backend.

### Registration of Webhook

**- Request:**

```http
Method: POST
Path: /energy/v1/notifications
```
* Content-Type: `application/json`

The body contains:

* The URL to be triggered upon event (the service that should be invoked)
* The country for which the trigger applies (if empty, it applies to any invocation)
* The number of invocations after which a notification is triggered (it should re-occur
* Every number of invocations, i.e., if 5 is specified, it should occur after 5, 10, 15 invocation, and so on, unless the webhook is deleted).

```http
{
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5
}
```

**- Response:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

```
{
  "webhook_id": "OIdksUDwveiwe"
}
```

### Deletion of Webhook

**- Request:**

```http
Method: DELETE
Path: /energy/v1/notifications/{id}
```
* `{id}` is the ID returned during the webhook registration

**- Response:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

```json
{
   "webhook_id": "OIdksUDwveiwe",
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5
}
```


### View registered webhook

**- Request**

```http
Method: GET
Path: /energy/v1/notifications/{id}
```
* `{id}` is the ID returned during the webhook registration

**- Response:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

```json
{
  "webhook_id": "OIdksUDwveiwe",
  "url": "https://localhost:8080/client/",
  "country": "NOR",
  "calls": 5
}
```

### View all registered webhooks

**- Request**

```http
Method: GET
Path: /energy/v1/notifications/
```

**- Response:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

```json
[
   {
      "webhook_id": "OIdksUDwveiwe",
      "url": "https://localhost:8080/client/",
      "country": "NOR",
      "calls": 5
   },
   {
      "webhook_id": "DiSoisivucios",
      "url": "https://localhost:8081/anotherClient/",
      "country": "SWE",
      "calls": 2
   },
   ...
]
```

### Webhook Invocation (Upon Trigger)

When a webhook is triggered, it sends information as follows. Where multiple webhooks are triggered, the information will be sent separately (i.e. one notification per triggered webhook).

```http
Method: POST
Path: <url specified in the corresponding webhook registration>
```
* Content-Type: `application/json`

```json
{
  "webhook_id": "OIdksUDwveiwe",
  "country": "Norway",
  "calls": 10
}
```
note: calls show the number of invocations, not the number specified as part of the webhook registration (i.e. not 5, but the actual invocation upon which the webhook was triggered).

## Status

The status interface indicates the availability of all individual services this service depends on. The reporting occurs based on status codes returned by the dependent services. The status interface further provides information about the number of registered webhooks, and the uptime of the service.

**- Request:**

```http
Method: GET
Path: /energy/v1/status
```

**- Response:**

* Content-Type: `application/json`
* Status: `200 OK` if OK, appropriate error code otherwise

```json
{
   "countries_api": "<http status code for *REST Countries API*>",
   "notification_db": "<http status code for *Notification DB* in Firebase>",
   ...
   "webhooks": <number of registered webhooks>,
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}
```

# Dependencies

## External Services
### REST Countries API
* http://universities.hipolabs.com/
* Documentation/Source under: https://github.com/Hipo/university-domains-list/
### Renewable Energy Dataset
The renewable energy data set is retrieved from: https://ourworldindata.org/energy

## Imported Golang modules
- cloud.google.com/go/firestore v1.9.0
- firebase.google.com/go v3.13.0+incompatible
- google.golang.org/api v0.118.0
- github.com/stretchr/testify/assert

# Contributors

## Developers
This project was developed by:
- Hoa Ben The Nguyen
- Magnus Johannessen
- Patrick Johannessen
- Sara Djordjevic

## Aknowledgements
This project was developed as part of the course PROG2005 Cloud Technologies at NTNU Gjøvik, and was made partially using code from the course.

Lecturer: Christopher Frantz
