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

> Optimisation of -w (number of database store workers) in loader is important, 5 should be optimal.
> Having too many workers can lead to data loss during the insert into database. 
> It depends on your laptop and local postgres setup/performance - make sure to not overkill with pg queries.

**Todo**
- [ ] implement re-try logic if insert into database fails! Really important!! Right now data loss is possible.
- [ ] improve importer, split data_dump.csv into smaller files and then process each file concurrently
- [ ] improve error handling when storing data into database, some entries can fail insert and no feedback is provided
- [ ] improve data load time to be less than 20s, right now average time is ~33s