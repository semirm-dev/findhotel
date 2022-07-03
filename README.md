**Run**
```shell
docker-compose up
```
**Tests**
```shell
go test ./... -v
```

**Geo**
- import data from given -path (data_dump.csv)
- validate and encode csv data to go structure
- forward structured data to storer service

**Storer**
- runs on 8001 port (configurable)
- listen for data to be stored
- store data into database
- implemented with postgres

**Gateway**
- runs on 8000 port (configurable)
- expose GET /geo?ip= endpoint to get geo data based on ip

Tasks
- [ ] storer service - pg implementation by default
- [ ] api gateway - expose GET /geo endpoint to return geo data based on ip address
- [ ] geo library - import geo data, csv implementation by default
- [ ] geo model should match csv data
- [ ] statistics after importing csv data (total time elapsed, imported, discarded)
- [ ] validate imported data (duplicates, missing data, corrupted)