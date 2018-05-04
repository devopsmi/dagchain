package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/465583030/dagchain"
)

var laddr = flag.String("l", "", "")
var saddr = flag.String("s", "", "")
var messaging = flag.Bool("m", false, "")
var id = flag.String("i", "", "")

func main() {
	flag.Parse()

	send := make(chan interface{}, 1)
	recv := make(chan interface{}, 1)
	// start the p2p node
	go func() {
		err := dotray.StartNode(*laddr, *saddr, send, recv)
		if err != nil {
			panic("node start panic:" + err.Error())
		}
	}()

	// wait 1 second for p2p node started
	time.Sleep(1 * time.Second)

	// query 10 nodes address from p2p network
	addrs := dotray.QueryNodes(10)

	// excute some actions with the nodes address,like downloads blockchain from these nodes
	// all depends on yourself
	fmt.Println("query nodes:", addrs)

	// send message to all the other nodes
	if *messaging {
		data := "hello-" + *id
		go func() {
			for {
				send <- data
				time.Sleep(5 * time.Second)
				fmt.Println("send message：", data)
			}
		}()
	}

	// receive message from other nodes
	for {
		select {
		case r := <-recv:
			res := r.(*dotray.Request)
			fmt.Printf("receive message: %v from other node: \"%s\" \n", res.Data, res.From)
		}
	}

}
