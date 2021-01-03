# O-fever Server
Provides server functionality for fetching information about orienteering competitions in a series. 

## API

### Scrape
For scraping results off of third party result sites. Currently, only scraping of Geoform results is supported
(https://ilgeoform.no/rankinglop/). 

* **scrape/[season]** - Example https://server/scrape/2020 - Scrapes the whole season and stores it internally.

```
/scrape/<season> - Example https://server/scrape/2020 - Scrapes the whole season and stores it internally.
```

### Athletes

* **athletes** - Example https://server/athletes - Fetches short identifying information about athletes across all series and seasons. Example payload:
```
[
    {
        "id":"3f9e764b-a214-ca5c-75b2-9193db855971",
        "sha":"9j_EhSfisYQqMaL3vxODsaJeczw=",
        "name":"Lastname, Firtsname",
        "club":"O Club"
    }
]
```

### Competitions

* **competitions/[series]/[season]** - Example https://server/competitions/geoform/2020 - Fetches identifying information about all competitions for the given series and season. Example payload:
```
[
    {"series":"geoform","season":"2020","number":1,"name":"Rankingløp 1 (Sørkedalskarusellen 1)"},
    {"series":"geoform","season":"2020","number":11,"name":"OSI/GeoForm rankingløp 11"}
]
```

* **competitions/[series]/[season]/[number]** - Example http://server/competitions/geoform/2020/1 - Fetches detailed information about the competition. Example payload:
```
{
    "series": "geoform",
    "season": "2020",
    "number": 1,
    "name": "Rankingløp 1 (Sørkedalskarusellen 1)",
    "url_live_lox": "https://www.livelox.com/Events/Show/49429",
    "courses": [
        {
            "name": "Resultater Lang (5.1 km)",
            "info": "",
            "results": [
                {
                    "placement": 1,
                    "disqualified": false,
                    "athlete_id": "79bc9a8a-a56f-b95c-47e8-9f42a8080e00",
                    "name": "Flågen, Bjørn Anders",
                    "club": "NTNUI",
                    "elapsed_time_seconds": 2411,
                    "elapsed_time_display": "0:40:11",
                    "missing_controls": 0,
                    "points": 151.99
                },
                {
                    "placement": 2,
                    "disqualified": false,
                    "athlete_id": "d972bdad-587c-b7d4-584b-55e1d3ac0f72",
                    "name": "Jacobsen, Kristoffer",
                    "club": "Tyrving IL",
                    "elapsed_time_seconds": 2457,
                    "elapsed_time_display": "0:40:57",
                    "missing_controls": 0,
                    "points": 151.55
                }
        ]
    }
]
}
```
