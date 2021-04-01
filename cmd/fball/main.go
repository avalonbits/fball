package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"git.cana.pw/avalonbits/fball"
	"github.com/kr/pretty"
	"go.uber.org/ratelimit"
)

var (
	key = flag.String("key", "", "API key for football-api")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stderr, "fball - ", log.LstdFlags|log.Lshortfile)
	limit := ratelimit.New(10, ratelimit.Per(time.Minute))
	c := fball.NewClient(*key, limit, &http.Client{Timeout: 10 * time.Second}, logger)

	tr, err := c.Timezone()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(tr))
}
