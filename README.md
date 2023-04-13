# Assignment 2

## Endpoints
The service consists of the following four endpoints:
```
/energy/v1/renewables/current
/energy/v1/renewables/history
/energy/v1/notifications/
/energy/v1/status/
```

### Current percentage of renewables
```
/energy/v1/renewables/current
```

{country?} refers to an optional country 3-letter code. 

{?neighbours=bool?} refers to an optional parameter indicating whether neighbouring countries' values should be shown.

Example request: ```/energy/v1/renewables/current/nor, /energy/v1/renewables/current/nor?neighbours=true```

Body (Exemplary message based on schema) - with country code:
```
{
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2021",
    "percentage": 71.558365
}
```

Body (Exemplary message based on schema) - with country code and neighbour parameter activated:
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
```/energy/v1/renewables/history```

This endpoint focuses on returning historical percentages of renewables in the energy mix, including individual levels, as well as mean values for individual or selections of countries.

{country?} refers to an optional country 3-letter code.
Example request: ```/energy/v1/renewables/history/nor```

Body (Exemplary message based on schema) - with country code:
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

Body (Exemplary message based on schema) - without country code (returns mean percentages for all countries):
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

Users can register webhooks that are triggered by the service based on specified events, specifically if information about given countries (or any country) is invoked, where the minimum frequency can be specified. Users can register multiple webhooks. The registrations should survive a service restart (i.e., be persistent using a Firebase DB as backend).

#### Registration of Webhooks
Body (Exemplary message based on schema):
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
- Request

Method: DELETE
Path: /energy/v1/notifications/{id}

{id} is the ID returned during the webhook registration


- Response
  Implement the response according to best practices.


#### View registered webhook

- Request

```
Method: GET
Path: /energy/v1/notifications/{id}
```

{id} is the ID for the webhook registration

The response is similar to the POST request body, but further includes the ID assigned by the server upon adding the webhook.


Body (Exemplary message based on schema):
```
{
"webhook_id": "OIdksUDwveiwe",
"url": "https://localhost:8080/client/",
"country": "NOR",
"calls": 5
}
```

#### View all registered webhooks

Request:
```
Method: GET
Path: /energy/v1/notifications/
```

The response is a collection of all registered webhooks.

Content type: application/json


Body (Exemplary message based on schema):
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
When a webhook is triggered, it should send information as follows. Where multiple webhooks are triggered, the information should be sent separately (i.e., one notification per triggered webhook). Note that for testing purposes, this will require you to set up another service that is able to receive the invocation. During the development, consider using https://webhook.site/ initially.

```
Method: POST
Path: <url specified in the corresponding webhook registration>
```

Body (Exemplary message based on schema):
```
{
"webhook_id": "OIdksUDwveiwe",
"country": "Norway",
"calls": 10
}
```

### Status
/energy/v1/status/

## Dependencies

## Contributers

## Deployment
