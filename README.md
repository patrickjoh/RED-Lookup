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

Example request: /energy/v1/renewables/current/nor, /energy/v1/renewables/current/nor?neighbours=true

Body (Exemplary message based on schema) - with country code:
```
{
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2021",
    "percentage": 71.558365
}
```

### Renewables history
/energy/v1/renewables/history

### Notifications
/energy/v1/notifications/

### Status
/energy/v1/status/

## Dependencies

## Contributers

## Deployment
