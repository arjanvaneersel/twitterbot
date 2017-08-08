package main

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/arjanvaneersel/twitterbot/log"
)

var StreamActive = false

func StartTwitterStream(cfg *Config, api *anaconda.TwitterApi, log *log.Logger) {
	if cfg.HasSubscriptions() {
		stream := api.PublicStreamFilter(url.Values{
			"track": cfg.Subscriptions,
		})
		StreamActive = true

		defer func() {
			stream.Stop()
			StreamActive = false
		}()

		for v := range stream.C {
			t, ok := v.(anaconda.Tweet)
			if !ok {
				log.Warningf("Received a value of %T instead of anaconda.Tweet", v)
				continue
			}

			if t.RetweetedStatus != nil {
				continue
			}

			_, err := api.Retweet(t.Id, false)
			if err != nil {
				log.Errorf("Could not retweet %d: %v", t.Id, err)
				continue
			}
			log.Infof("Retweeted %d", t.Id)
		}
	}
}
