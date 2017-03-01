# Elastic on Docker

This is a small example of how to use [Elastic](https://github.com/olivere/elastic),
the Elasticsearch client for Go, with [Docker](https://docs.docker.com/) and
[Docker Compose](https://docs.docker.com/compose/).

Docker Compose orchestrates two containers. It requires Docker/Compose 1.12+
as the `docker-compose.yml` configuration uses version 2.1
(see [here](https://docs.docker.com/compose/compose-file/compose-versioning/) for Docker Compile file versioning matrix).

First, there is the
[Elasticsearch container](https://hub.docker.com/_/elasticsearch/). We use
the default configuration, but we mount the data directory on `/tmp/esdata`.
[See here](https://github.com/docker-library/elasticsearch/issues/74)
why we do this.

Second, we host a simple app in a container that will connect to the
Elasticsearch container and will periodically call the Nodes Info API.

To start everything up, run `docker-compose up`, then watch `docker-compose logs`
in another shell.

The application might restart several times until the Elasticsearch container
is up and running. It might print something like this eventually:

```
app_1            | 2017/03/01 08:21:39 main.go:38: Lookup for elasticsearch returns the following IPs:
app_1            | 2017/03/01 08:21:39 main.go:40: 172.19.0.2
app_1            | 2017/03/01 08:21:39 main.go:43: Retrieving http://elasticsearch:9200:
app_1            | 2017/03/01 08:21:39 main.go:53: {
app_1            |   "name" : "YyRepA3",
app_1            |   "cluster_name" : "elasticsearch",
app_1            |   "cluster_uuid" : "V46AhaVLSJm4ERTgMEz_KA",
app_1            |   "version" : {
app_1            |     "number" : "5.0.0",
app_1            |     "build_hash" : "253032b",
app_1            |     "build_date" : "2016-10-26T05:11:34.737Z",
app_1            |     "build_snapshot" : false,
app_1            |     "lucene_version" : "6.2.0"
app_1            |   },
app_1            |   "tagline" : "You Know, for Search"
app_1            | }
app_1            | 2017/03/01 08:21:39 main.go:56: connecting to http://elasticsearch:9200
app_1            | 2017/03/01 08:21:39 main.go:61: connected to http://elasticsearch:9200
app_1            | 2017/03/01 08:21:39 main.go:101: Cluster "elasticsearch" with 1 node(s)
app_1            | 2017/03/01 08:21:39 main.go:103: - Node YyRepA3vRseVL_DRNoIE7A with IP 172.19.0.2
```

Notice that Elastic has a [dedicated page on Installing Elasticsearch with Docker](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html).

## LICENSE

MIT
