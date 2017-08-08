package main

import (
	shorten "../.."
	service "../../http"
	store "../../store"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	listenAddr := flag.String("listen", ":9000", "Address to listen on in [host]:port format")
	redisAddr := flag.String("redis", ":6379", "Address of the Redis server to connect to")
	flag.Parse()

	store := store.NewRedis(&store.RedisOptions{
		RedisAddr: *redisAddr,
	})

	s := shorten.New(&short.Options{
		Store: store,
	})

	expandHandler := service.NewExpander(s)

	http.Handle("/", expandHandler)
	err := http.ListenAndServe(*listenAddr, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening on %s: %s\n", *listenAddr, err)
		os.Exit(-1)
	}
}
