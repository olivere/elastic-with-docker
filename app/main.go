package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "showenv":
			log.Printf("Environment")
			env := os.Environ()
			sort.Strings(env)
			for _, e := range env {
				log.Printf("- %s", e)
			}
			os.Exit(0)
		}
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ips, err := net.LookupIP("elasticsearch")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Lookup for elasticsearch returns the following IPs:")
	for _, ip := range ips {
		log.Printf("%v", ip)
	}

	log.Printf("Retrieving http://elasticsearch:9200:")
	res, err := http.Get("http://elasticsearch:9200")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", string(body))

	url := "http://elasticsearch:9200"
	log.Printf("connecting to %v", url)
	client, err := elastic.NewClient(elastic.SetURL(url))
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
	ctx := context.Background()
	info, err := client.NodesInfo().Do(ctx)
	if err != nil {
		return err
	}
	log.Printf("Cluster %q with %d node(s)", info.ClusterName, len(info.Nodes))
	for id, node := range info.Nodes {
		log.Printf("- Node %s with IP %s", id, node.IP)
	}
	return nil
}
