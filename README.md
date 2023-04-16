# Assignment 2
This project is a REST web application in Golang that provides the client with the ability to retrieve information about developments related to renewable energy production for and across countries. This is done by using two existing webservice. The service also allows for notification registration using webhooks. The application is dockerized and deployed using an IaaS system.
## Endpoints
The service consists of the following four endpoints:
```
/energy/v1/renewables/current
/energy/v1/renewables/history
/energy/v1/notifications/
/energy/v1/status/
```

### Current percentage of renewables
Path: ```/energy/v1/renewables/current/country?neighbours=bool```

{country?} refers to an optional country 3-letter code. 

{?neighbours=bool?} refers to an optional parameter indicating whether neighbouring countries' values should be shown.

Example requests:
- ```/energy/v1/renewables/current```
- ```/energy/v1/renewables/current/nor```
- ```/energy/v1/renewables/current/nor?neighbours=true```

Body example with country code:
```
{
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2021",
    "percentage": 71.558365
}
```

Body example with country code and neighbour parameter activated:
```
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
### Historical percentages of renewables
Path: ```/energy/v1/renewables/history/{country?}{?begin=year&end=year?}```

This endpoint focuses on returning historical percentages of renewables in the energy mix, including individual levels, as well as mean values for individual or selections of countries.

{country} refers to an optional country 3-letter code.
{?begin=year&end=year} refers to an optional range for country percentages.

Example requests: 
- ```/energy/v1/renewables/history```
- ```/energy/v1/renewables/history/nor```
- ```/energy/v1/renewables/history/nor?begin=1970```
- ```/energy/v1/renewables/history/nor?begin=1960&end=1970```

Body example with country code:
```
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

Body example without country code (returns mean percentages for all countries):
```
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
### Notifications
```/energy/v1/notifications/```

Users can register webhooks that are triggered by the service based on specified events, specifically if information about given countries (or any country) is invoked, where the minimum frequency can be specified. Users can register multiple webhooks. The registrations should survive a service restart.

#### Registration of Webhooks
Body example:
```
{
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5
}
```

The response contains the ID for the registration that can be used to see detail information or to delete the webhook registration. The format of the ID is not prescribed, as long it is unique. Consider best practices for determining IDs.
```
{
  "webhook_id": "OIdksUDwveiwe"
}
```

#### Deletion of Webhooks
Path: ```/energy/v1/notifications/{id}```

{id} is the ID returned during the webhook registration



#### View registered webhook
Path: ```/energy/v1/notifications/{id}```

{id} is the ID for the webhook registration.

The response is similar to the POST request body, but further includes the ID assigned by the server upon adding the webhook.

Body example:
```
{
  "webhook_id": "OIdksUDwveiwe",
  "url": "https://localhost:8080/client/",
  "country": "NOR",
  "calls": 5
}
```

#### View all registered webhooks
Path: ```/energy/v1/notifications/```

The response is a collection of all registered webhooks.

Body example:
```
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

#### Webhook Invocation
Path: ```<url specified in the corresponding webhook registration>```

When a webhook is triggered, it should send information as follows. Where multiple webhooks are triggered, the information should be sent separately (i.e., one notification per triggered webhook).

Body example:
```
{
  "webhook_id": "OIdksUDwveiwe",
  "country": "Norway",
  "calls": 10
}
```

### Status
Path: ```/energy/v1/status/```

The status interface indicates the availability of all individual services this service depends on. The reporting occurs based on status codes returned by the dependent services. The status interface further provides information about the number of registered webhooks, and the uptime of the service.

Body:
```
{
   "countries_api": "<http status code for *REST Countries API*>",
   "notification_db": "<http status code for *Notification DB* in Firebase>",
   ...
   "webhooks": <number of registered webhooks>,
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}

```

## Dependencies
### Third party services
#### REST Countries API

#### Renewable Energy Dataset
The renewable energy data set is retrieved from: https://ourworldindata.org/energy

### Imported modules
- encoding/json
- log
- net/http
- strconv
- strings
- cloud.google.com/go/firestore v1.9.0
- firebase.google.com/go v3.13.0+incompatible
- google.golang.org/api v0.118.0
- github.com/stretchr/testify/assert

## Contributors
### Developers
This project was developed by:
  - Hoa Ben The Nguyen
  - Magnus Johannessen
  - Patrick Johannessen
  - Sara Djordjevic

## Deployment
This service is deployed on an IaaS solution OpenStack using Docker.
