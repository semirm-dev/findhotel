**Run**
* Please start database and redis first because they take some time to boot. 
> I didnt implement healthchecks to make sure services are not started before databases are ready
```shell
docker-compose up db redis
```
* Then start loader to load all dump_data.csv into database and gateway to expose an api to search for geolocations
* It takes about 30-35 seconds to import all data
```shell
docker-compose up loader gateway
```

**Tests**
```shell
go test ./... -v
```

**Geo**
- main module
- import *geo data (data_dump.csv) with Importer and save it in database using Storer

**Gateway**
- runs on 8000 port (configurable)
- expose GET /geo?ip= endpoint to get *geo data based on ip
- uses geo.Search api to search for *geo data

**Todo**
- [ ] implement re-try logic if insert into database fails! Really important!! Right now data loss is possible.
- [ ] improve importer, split data_dump.csv into smaller files and then process each file concurrently
- [ ] improve error handling when storing data into database, some inserts can fail and no feedback is provided
- [ ] improve data load/import time to be less than 20s