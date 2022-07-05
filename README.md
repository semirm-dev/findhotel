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
- import *geo data (data_dump.csv) with Importer and save it in database using Storer

**Gateway**
- runs on 8000 port (configurable)
- expose GET /geo?ip= endpoint to get geo data based on ip
- uses geo.Search api to search go 

**Todo**
- [ ] improve importer, split data_dump.csv into smaller files and then process each file concurrently
- [ ] improve error handling when storing data into database, some entries can fail insert and no feedback is provided