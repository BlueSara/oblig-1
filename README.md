# Country Info API - README

### Overview

This API provides country-related data, including general information, population statistics, and service status. It is built using Go and serves HTTP endpoints to fetch data from external APIs.

Features

Retrieve general information about a country.

Fetch population data for a specific country within a given year range.

Check the status of external APIs used by the service.

Endpoints

## 1. Get Country Information

GET /countryinfo/v1/info/{two_letter_country_code}

Retrieves general country information using a two-letter country code (ISO2).

Example Request:

```curl -X GET http://localhost:8080/countryinfo/v1/info/NO


## 2. Get Population Data

GET /countryinfo/v1/population/{two_letter_country_code}?limit=XXXX-YYYY

Fetches population data for a country within a given year range (optional limit).

Example Request:

curl -X GET "http://localhost:8080/countryinfo/v1/population/NO?limit=2000-2020"

## 3. Check API Status

GET /countryinfo/v1/status/

Returns the status of external APIs used by the service and uptime.

Example Request: