package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"time"

	elastic "gopkg.in/olivere/elastic.v3"
)

func main() {
	log.Printf("Environment")
	env := os.Environ()
	sort.Strings(env)
	for _, e := range env {
		log.Printf("- %s", e)
	}

	// Give Elasticsearch some time to startup
	time.Sleep(10 * time.Second)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	url := fmt.Sprintf("http://%s:%s",
		os.Getenv("ELASTICSEARCH_PORT_9200_TCP_ADDR"),
		os.Getenv("ELASTICSEARCH_PORT_9200_TCP_PORT"))
	log.Printf("connecting to %v", url)
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(url))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to %v", url)

	errc := make(chan error)

	go func() {
		err := showNodes(client)
		if err != nil {
			log.Printf("nodes info failed: %v", err)
		}

		t := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-t.C:
				err := showNodes(client)
				if err != nil {
					log.Printf("nodes info failed: %v", err)
				}
			}
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		log.Printf("existing with signal %v", fmt.Sprint(<-c))
		errc <- nil
	}()

	if err := <-errc; err != nil {
		os.Exit(1)
	}
}

func showNodes(client *elastic.Client) error {
	info, err := client.NodesInfo().Do()
	if err != nil {
		return err
	}
	log.Printf("Cluster %q with %d node(s)", info.ClusterName, len(info.Nodes))
	for id, node := range info.Nodes {
		log.Printf("- Node %s with IP %s", id, node.IP)
	}
	return nil
}
