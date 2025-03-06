#  Oblig_1 - Country Info API  

Welcome to my **Oblig_1**
This API helps you fetch details about countries, including general info, population statistics, and API health status.

---


##  Endpoints  

###  Get Country Information  
**GET** `/info/{ISO2-country_code}?limit=integer`  
This endpoint returns general details about a country, such as its name, capital, and more!  

**Limit** is optional – if provided, it specifies the number of cities included in the response.  

####  Example Request:  
```bash
/info/no?limit=10
```

This will give the following JSON:
```bash
"name": {
        "common": "Norway"
    },
    "continents": [
        "Europe"
    ],
    "population": 5379475,
    "languages": {
        "nno": "Norwegian Nynorsk",
        "nob": "Norwegian Bokmål",
        "smi": "Sami"
    },
    "borders": [
        "FIN",
        "SWE",
        "RUS"
    ],
    "flag": "🇳🇴",
    "capital": [
        "Oslo"
    ],
    "cities": [
        "Abelvaer",
        "Adalsbruk",
        "Adland",
        "Agdenes",
        "Agotnes",
        "Agskardet",
        "Aker",
        "Akkarfjord",
        "Akrehamn",
        "Al"
    ]
}
```


### Get Population Data

**GET** `/population/{ISO2-country_code}?limit="startYear-endYear"`
Curious about how a country's population has changed over time? This endpoint provides population statistics for the given time range!

- Limit is optional – specify a start year and an end year to filter the results.
- Example Request:

####  Example Request: 
```bash
/population/no?limit=2010-2020
```

This will give the following JSON:
```bash
{
    "mean": 5121086,
    "values": [
        {
            "year": 2010,
            "value": 4889252
        },
        {
            "year": 2011,
            "value": 4953088
        },
        {
            "year": 2012,
            "value": 5018573
        },
        {
            "year": 2013,
            "value": 5079623
        },
        {
            "year": 2014,
            "value": 5137232
        },
        {
            "year": 2015,
            "value": 5188607
        },
        {
            "year": 2016,
            "value": 5234519
        },
        {
            "year": 2017,
            "value": 5276968
        },
        {
            "year": 2018,
            "value": 5311916
        }
    ]
}
```

This will return population data for Norway from 2010 to 2020.



### Check API Status

**GET** `/status`
Wondering if everything is running smoothly? This endpoint lets you check the health of the API we are fetching from!
Example Request:

```bash
/status
```

This will give the following JSON:
```bash
{
    "countriesnowapi": "200 OK",
    "restcountriesapi": "200 OK",
    "version": "v1",
    "uptime": 697
}
```
If all systems are good, you'll receive a response confirming the API is operational.



