# town-diff
Calculates Geometric Difference between two towns

## Prerequisits

### Mysql setup

To run the application it is required to have `mysql` locally.
In the command line execute:

```bash
docker volume create tdiff-vol
docker run --rm --name tdiff-mysql -v tdiff-vol:/var/lib/mysql -p 3306:3306 -e MYSQL_DATABASE=tdiff -e MYSQL_USER=tdiff -e MYSQL_PASSWORD=tdiff -e MYSQL_ROOT_PASSWORD=root -d mysql:latest
```

This will initialize a mysql docker container.

In order to get the docker Ip of the MySql container to connect to, execute the following:

```bash
docker container list | grep tdiff-mysql | awk '{ print $1 }'

# replace the container id in the insepct command
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' <container_id>
```

## Build

Before building we need to make sure that the configuration is updated.

### TDiff configuration

TDiff requires a configuration file in order to connect to the db.

You can update the [config.docker.json](./config.docker.json) with the necessary information.

> **_NOTE:_**  Password should be base64 encoded

---

Once the config is done, to build the application execute:

```bash
docker build --tag tdiff .
```

## Run

Before running the application locally, make sure that there is a mysql instance is already running.

To start tdiff execute

```bash
docker run --rm --name tdiff -p 8080:8080 tdiff
```

## Postman Collection

Attached you can find a postman collection that contains the 5 operations of this app.

To find the Lat, Lon coordinates of a town, visit [https://www.latlong.net/](https://www.latlong.net/)

To validate the distance between two towns, visit [https://www.distancefromto.net/](https://www.distancefromto.net/)