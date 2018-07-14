# Elastic on Docker (v6)

This is a small example of how to use [Elastic](https://github.com/olivere/elastic) for Elasticsearch 6.x with [Docker](https://docs.docker.com/) and [Docker Compose](https://docs.docker.com/compose/).

We use Docker Compose to orchestrates the containers. I've been using it successfully with Docker Compose 1.19+.  The Docker compose configuration files use version 3 (see [here](https://docs.docker.com/compose/compose-file/compose-versioning/) for Docker Compile file versioning matrix).

Notice there are two official Elasticsearch 6.x images: The OSS image at `docker.elastic.co/elasticsearch/elasticsearch-oss:6.3.1` (which we use here as default) and the full image `docker.elastic.co/elasticsearch/elasticsearch:6.3.1`. You need a license for the latter after a grace period. For more information, read the [official documentation on Installing Elasticsearch with Docker](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html).

Make sure to create a `./data` directory locally and uncomment the `volumnes` section in Docker Compose file(s) if you want your data to be persistent.

The examples are split into two configurations. The first is for running Elasticsearch and a test application inside containers with Docker Compose. The second is for running Elasticsearch inside Docker only, but running the test application on the host.

## Containers

This section describes how to run both Elasticsearch and the application inside Docker containers.

Use the `docker-compose.containers.yml` file for starting the setup via Docker Compose. Notice that the test application might restart several times before Elasticsearch becomes available.

```
$ docker-compose -f docker-compose.containers.yml up --build
...
app_1            | 2018/07/14 11:52:12 main.go:50: Looking up hostname "elasticsearch"
app_1            | 2018/07/14 11:52:12 main.go:55: Lookup for hostname "elasticsearch" returns the following IPs:
app_1            | 2018/07/14 11:52:12 main.go:57: 172.18.0.2
app_1            | 2018/07/14 11:52:12 main.go:62: Retrieving http://elasticsearch:9200:
app_1            | 2018/07/14 11:52:13 main.go:72: {
app_1            |   "name" : "5dd3if1",
app_1            |   "cluster_name" : "escontainers",
app_1            |   "cluster_uuid" : "8SdPtWC_SVm2supbip2Bwg",
app_1            |   "version" : {
app_1            |     "number" : "6.3.1",
app_1            |     "build_flavor" : "oss",
app_1            |     "build_type" : "tar",
app_1            |     "build_hash" : "eb782d0",
app_1            |     "build_date" : "2018-06-29T21:59:26.107521Z",
app_1            |     "build_snapshot" : false,
app_1            |     "lucene_version" : "7.3.1",
app_1            |     "minimum_wire_compatibility_version" : "5.6.0",
app_1            |     "minimum_index_compatibility_version" : "5.0.0"
app_1            |   },
app_1            |   "tagline" : "You Know, for Search"
app_1            | }
app_1            | 2018/07/14 11:52:13 main.go:77: Retrieving http://elasticsearch:9200/_nodes/http?pretty=true:
app_1            | 2018/07/14 11:52:13 main.go:87: {
app_1            |   "_nodes" : {
app_1            |     "total" : 1,
app_1            |     "successful" : 1,
app_1            |     "failed" : 0
app_1            |   },
app_1            |   "cluster_name" : "escontainers",
app_1            |   "nodes" : {
app_1            |     "5dd3if1pThS5LAY6VX1rbQ" : {
app_1            |       "name" : "5dd3if1",
app_1            |       "transport_address" : "172.18.0.2:9300",
app_1            |       "host" : "172.18.0.2",
app_1            |       "ip" : "172.18.0.2",
app_1            |       "version" : "6.3.1",
app_1            |       "build_flavor" : "oss",
app_1            |       "build_type" : "tar",
app_1            |       "build_hash" : "eb782d0",
app_1            |       "roles" : [
app_1            |         "master",
app_1            |         "data",
app_1            |         "ingest"
app_1            |       ],
app_1            |       "http" : {
app_1            |         "bound_address" : [
app_1            |           "0.0.0.0:9200"
app_1            |         ],
app_1            |         "publish_address" : "172.18.0.2:9200",
app_1            |         "max_content_length_in_bytes" : 104857600
app_1            |       }
app_1            |     }
app_1            |   }
app_1            | }
app_1            | 2018/07/14 11:52:13 main.go:90: Connecting to http://elasticsearch:9200
app_1            | 2018/07/14 11:52:13 main.go:95: Connected to http://elasticsearch:9200
app_1            | 2018/07/14 11:52:13 main.go:135: Cluster "escontainers" with 1 node(s)
app_1            | 2018/07/14 11:52:13 main.go:137: - Node 5dd3if1pThS5LAY6VX1rbQ with IP 172.18.0.2
app_1            | 2018/07/14 11:52:23 main.go:135: Cluster "escontainers" with 1 node(s)
...
```

To stop everything, run:

```
$ docker-compose -f docker-compose.containers.yml down
```

## Local

This section describes how to run Elasticsearch inside Docker only, and running the test application on the host.

Use the `docker-compose.local.yml` file for starting the setup via Docker Compose. Again, notice that this only starts Elasticsearch and makes it accessible locally on the host.

```
$ docker-compose -f docker-compose.local.yml up
```

Once Elasticsearch is successfully started, switch to a second console and run the test application. Notice that we need to pass the `-url=http://localhost:9200` parameter as Elasticsearch is configured to listen on that IP/port.

```
$ cd app
$ go run main.go -url http://localhost:9200
2018/07/14 13:59:36 main.go:50: Looking up hostname "localhost"
2018/07/14 13:59:36 main.go:55: Lookup for hostname "localhost" returns the following IPs:
2018/07/14 13:59:36 main.go:57: ::1
2018/07/14 13:59:36 main.go:57: 127.0.0.1
2018/07/14 13:59:36 main.go:62: Retrieving http://localhost:9200:
2018/07/14 13:59:36 main.go:72: {
  "name" : "t-Hcpbj",
  "cluster_name" : "docker-cluster",
  "cluster_uuid" : "zOFxQHUCTy2V3y45JiO2yg",
  "version" : {
    "number" : "6.3.1",
    "build_flavor" : "oss",
    "build_type" : "tar",
    "build_hash" : "eb782d0",
    "build_date" : "2018-06-29T21:59:26.107521Z",
    "build_snapshot" : false,
    "lucene_version" : "7.3.1",
    "minimum_wire_compatibility_version" : "5.6.0",
    "minimum_index_compatibility_version" : "5.0.0"
  },
  "tagline" : "You Know, for Search"
}
2018/07/14 13:59:36 main.go:77: Retrieving http://localhost:9200/_nodes/http?pretty=true:
2018/07/14 13:59:36 main.go:87: {
  "_nodes" : {
    "total" : 1,
    "successful" : 1,
    "failed" : 0
  },
  "cluster_name" : "docker-cluster",
  "nodes" : {
    "t-HcpbjwQX-h8KW8wwtv0w" : {
      "name" : "t-Hcpbj",
      "transport_address" : "127.0.0.1:9300",
      "host" : "127.0.0.1",
      "ip" : "127.0.0.1",
      "version" : "6.3.1",
      "build_flavor" : "oss",
      "build_type" : "tar",
      "build_hash" : "eb782d0",
      "roles" : [
        "master",
        "data",
        "ingest"
      ],
      "http" : {
        "bound_address" : [
          "172.18.0.2:9200",
          "127.0.0.1:9200"
        ],
        "publish_address" : "127.0.0.1:9200",
        "max_content_length_in_bytes" : 104857600
      }
    }
  }
}
2018/07/14 13:59:36 main.go:90: Connecting to http://localhost:9200
2018/07/14 13:59:37 main.go:95: Connected to http://localhost:9200
2018/07/14 13:59:37 main.go:135: Cluster "docker-cluster" with 1 node(s)
2018/07/14 13:59:37 main.go:137: - Node t-HcpbjwQX-h8KW8wwtv0w with IP 127.0.0.1
...
```

To stop everything, quit the test application, then run:

```
$ docker-compose -f docker-compose.local.yml down
```

## LICENSE

MIT
