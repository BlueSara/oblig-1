#  Oblig_1 - Country Info API  

Welcome to the **Oblig_1 Country Info API!** This API helps you fetch interesting details about countries, including general info, population statistics, and API health status. No passport required! ✈️  

---

##  Endpoints  

###  Get Country Information  
**GET** `/info/{ISO2-country_code}?limit=integer`  
This endpoint returns general details about a country, such as its name, capital, and more!  

🔹 **Limit** is optional – if provided, it specifies the number of cities included in the response.  

####  Example Request:  
```bash
/info/no?limit=10
```

### Get Population Data

GET /population/{ISO2-country_code}?limit="startYear-endYear"
Curious about how a country's population has changed over time? This endpoint provides population statistics for the given time range!

- Limit is optional – specify a start year and an end year to filter the results.
- Example Request:

####  Example Request: 
```bash
/population/no?limit=2010-2020
```
This will return population data for Norway 📈 from 2010 to 2020.


### Check API Status

GET /status
Wondering if everything is running smoothly? This endpoint lets you check the health of the API we are fetching from!
Example Request:

```bash
/status
```
If all systems are go, you'll receive a response confirming the API is operational.



