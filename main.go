package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Sirupsen/logrus"
	"github.com/arjanvaneersel/twitterbot/log"
	"github.com/gorilla/mux"
)

func routes(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello gopher!")
	})
}

func main() {
	cfgFile := flag.String("config", "twiibt.cfg", "Configuration file")
	cfg, err := LoadConfig(*cfgFile)
	if err != nil {
		panic(err)
	}
	anaconda.SetConsumerKey(cfg.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.ConsumerSecret)
	api := anaconda.NewTwitterApi(cfg.AccessToken, cfg.AccessTokenSecret)

	log := &log.Logger{logrus.New()}
	api.SetLogger(log)

	r := mux.NewRouter()

	StartServer(cfg, r)
}
