# town-diff
Calculates Geometric Difference between two town

## Start all

```bash
docker-compose up
```

## Docker build

```bash
docker build --tag tdiff .
```

## Docker Run

```bash
docker run -p 8080:8080 --rm --name tdiff tdiff
```


## Get db container ip

```bash
docker container list | grep town-diff_db_1 | awk '{ print $1 }'
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' container_name_or_id
```
