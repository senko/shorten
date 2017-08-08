package main

import (
	shorten "../.."
	store "../../store"
	"flag"
	"fmt"
	"os"
)

func main() {
	shortUrl := flag.String("expand", "", "URL to expand (mutually-exclusive with -shorten)")
	fullUrl := flag.String("shorten", "", "URL to shorten (mutually-exclusive with -expand)")
	redisAddr := flag.String("redis", ":6379", "Address of the Redis server to connect to")
	flag.Parse()

	if *shortUrl != "" && *fullUrl != "" {
		fmt.Fprintf(os.Stderr, "Can't use both -expand and -shorten at the same time\n")
		os.Exit(-1)
	}

	if *shortUrl == "" && *fullUrl == "" {
		fmt.Fprintf(os.Stderr, "At least one of -expand, -shorten must be used (see -help for usage details)\n")
		os.Exit(-1)
	}

	store := store.NewRedis(&store.RedisOptions{
		RedisAddr: *redisAddr,
	})
	s := shorten.New(&short.Options{
		Store: store,
	})

	var result string
	var err error

	if *shortUrl != "" {
		result = s.Expand(*shortUrl)
		if result == "" {
			fmt.Fprintf(os.Stderr, "Short URL not found: %s\n", *shortUrl)
			os.Exit(-1)
		}
	}

	if *fullUrl != "" {
		result, err = s.Shorten(*fullUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't shorten the URL: %s\n", err)
			os.Exit(-1)
		}
	}

	fmt.Fprintln(os.Stderr, result)
}
