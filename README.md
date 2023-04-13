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

The initial endpoint focuses on returning historical percentages of renewables in the energy mix, including individual levels, as well as mean values for individual or selections of countries.

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
/energy/v1/notifications/

### Status
/energy/v1/status/

## Dependencies

## Contributers

## Deployment
