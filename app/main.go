package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/olivere/elastic"
)

func main() {
	var (
		esURL = flag.String("url", "http://elasticsearch:9200", "Elasticsearch connection string")
	)
	flag.Parse()

	if *esURL == "" {
		log.Fatal("missing -url flag")
	}
	url, err := url.Parse(*esURL)
	if err != nil {
		log.Fatalf("invalid -url flag: %v", err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "showenv":
			log.Println("Environment")
			env := os.Environ()
			sort.Strings(env)
			for _, e := range env {
				log.Printf("- %s", e)
			}
			os.Exit(0)
		}
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Looking up hostname %q", url.Hostname())
	ips, err := net.LookupIP(url.Hostname())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Lookup for hostname %q returns the following IPs:", url.Hostname())
	for _, ip := range ips {
		log.Printf("%v", ip)
	}

	// Check ES version and status
	{
		log.Printf("Retrieving %s:", *esURL)
		res, err := http.Get(*esURL)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v", string(body))
	}

	// Check ES nodes configuration
	{
		log.Printf("Retrieving %s:", *esURL+"/_nodes/http?pretty=true")
		res, err := http.Get(*esURL + "/_nodes/http?pretty=true")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v", string(body))
	}

	log.Printf("Connecting to %s", *esURL)
	client, err := elastic.NewClient(elastic.SetURL(*esURL))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to %s", *esURL)

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
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
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
