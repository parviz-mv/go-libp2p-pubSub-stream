package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	// "log"
	"os"
	"sync"

	// "time"

	libp2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

var (
	topicNameFlag = flag.String("topicName", "topic", "Name of topic to join")
	listenAddrs = flag.String("listenAddrs", "/ip4/127.0.0.1/tcp/9009", "Addrs for listen")
)

func main() {
	// topicStr := "alice_bob_topic"
	flag.Parse()
	ctx := context.Background()

	host, err := libp2p.New(libp2p.ListenAddrStrings(*listenAddrs))
	if err != nil {
		panic(err)
	}
	
	go discoverPeers(ctx, host)

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
	
	topic, err := ps.Join(*topicNameFlag)
	if err != nil {
		panic(err)
	}

	go streamConsoleTo(ctx, topic)

	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	sub.Topic()

	setInterval(2, func() {
		topics := ps.GetTopics()
		fmt.Println("Topic list::",topics )
		peers := ps.ListPeers(*topicNameFlag)
		fmt.Printf("Peers on %s: %s \n", *topicNameFlag, peers)
		publishMessage(ctx, topic)
	})

	printMessagesFrom(ctx, sub)

}

	func initDHT(ctx context.Context, h host.Host) *dht.IpfsDHT {
		// Start a DHT, for use in peer discovery. We can't just make a new DHT
		// client because we want each peer to maintain its own local copy of the
		// DHT, so that the bootstrapping node of the DHT can go down without
		// inhibiting future peer discovery.
		kademliaDHT, err := dht.New(ctx, h)
		if err != nil {
			panic(err)
		}
		if err = kademliaDHT.Bootstrap(ctx); err != nil {
			panic(err)
		}
		var wg sync.WaitGroup
		for _, peerAddr := range dht.DefaultBootstrapPeers {
			peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := h.Connect(ctx, *peerinfo); err != nil {
					fmt.Println("Bootstrap warning:", err)
				}
			}()
		}
		wg.Wait()
	
		return kademliaDHT
	}
	
	func discoverPeers(ctx context.Context, h host.Host) {
		kademliaDHT := initDHT(ctx, h)
		// topicStr := "alice_bob_topic"
		routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
		dutil.Advertise(ctx, routingDiscovery, *topicNameFlag)
	
		// Look for others who have announced and attempt to connect to them
		anyConnected := false
		for !anyConnected {
			fmt.Println("Searching for peers...")
			peerChan, err := routingDiscovery.FindPeers(ctx, *topicNameFlag)
			if err != nil {
				panic(err)
			}
			for peer := range peerChan {
				if peer.ID == h.ID() {
					continue // No self connection
				}
				err := h.Connect(ctx, peer)
				if err != nil {
					fmt.Println("Failed connecting to ", peer.ID.Pretty(), ", error:", err)
				} else {
					fmt.Println("Connected to:", peer.ID.Pretty())
					anyConnected = true
				}
			}
		}
		fmt.Println("Peer discovery complete")
	}
	
	func streamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			if err := topic.Publish(ctx, []byte(s)); err != nil {
				fmt.Println("### Publish error:", err)
			}
		}
	}
	
	func printMessagesFrom(ctx context.Context, sub *pubsub.Subscription) {
		for {
			m, err := sub.Next(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
		}
	}

	func setInterval(tick time.Duration, callback func()) {
		ticker := time.NewTicker(tick * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					callback()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}

	func publishMessage(ctx context.Context, topic *pubsub.Topic) {
		m:= "Hello Alice!!"
		msgBytes, err := json.Marshal(m)
		if err != nil {
			 fmt.Println("err::", err)
		}
		 topic.Publish(ctx, msgBytes)
	}