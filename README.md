**Run**
```shell
docker-compose up
```
**Tests**
```shell
go test ./... -v
```

**Geo**
- main module
- import *geo data (data_dump.csv) with Importer and save it in data store using Storer
- forward structured data to storer service

**Datastore**
- responsible to store *geo data into database
- implemented with postgres

**Gateway**
- runs on 8000 port (configurable)
- expose GET /geo?ip= endpoint to get geo data based on ip
- uses geo.Search api to search go 