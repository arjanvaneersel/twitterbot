package main

import (
	"compress/gzip"
	"encoding/gob"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/BalkanTech/kit/stringutils"

	"golang.org/x/crypto/bcrypt"
)

type AccessKey struct {
	Domain string `gob:"domain"`
	Key    string `gob:"key"`
}

type Config struct {
	ConsumerKey       string `gob:"consumer_key"`
	ConsumerSecret    string `gob:"consumer_key"`
	AccessToken       string `gob:"access_token"`
	AccessTokenSecret string `gob:"access_token_secret"`

	ServerAddress    string `gob:"server_address"`
	ServerAddressTLS string `gob:"server_address_tls"`
	ServerDomain     string `gob:"server_domain"`

	AccessKeys    []AccessKey `gob:"access_keys"`
	Subscriptions []string    `gob:"subscriptions"`
}

func (c *Config) AddSubscription(o string) {
	c.Subscriptions = append(c.Subscriptions, o)
}

func (c *Config) RemoveSubscription(i int) {
	c.Subscriptions = append(c.Subscriptions[:i], c.Subscriptions[i+1:]...)
}

func (c *Config) HasSubscriptions() bool {
	return len(c.Subscriptions) != 0
}

func (c *Config) CreateKey() (string, error) {
	key := stringutils.RandomString(25)
	h, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	c.AccessKeys = append(c.AccessKeys, AccessKey{Key: string(h)})
	return key, nil
}

func (c *Config) ValidKey(s string) bool {
	for i := range c.AccessKeys {
		err := bcrypt.CompareHashAndPassword([]byte(c.AccessKeys[i].Key), []byte(s))
		if err == nil {
			return true
		}
	}
	return false
}

func LoadConfig(fn string) (*Config, error) {
	c := Config{}
	file, err := os.Open(fn)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if os.IsNotExist(err) {
		key, err := c.CreateKey()
		if err != nil {
			return nil, err
		}
		if err := SaveConfig(&c, fn); err != nil {
			return nil, err
		}
		logrus.Infof("Created a new config file. Your access key is: %q.", key)
		return &c, nil
	}
	defer file.Close()

	zfile, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer zfile.Close()

	if err := gob.NewDecoder(zfile).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func SaveConfig(c *Config, fn string) error {
	file, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer file.Close()

	zfile := gzip.NewWriter(file)
	defer zfile.Close()

	if err := gob.NewEncoder(zfile).Encode(c); err != nil {
		return err
	}
	return nil
}

func WithValidKey(c *Config, n http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if !c.ValidKey(key) {
			http.Error(w, "Invalid key", http.StatusForbidden)
			return
		}
		n(w, r)
		return
	}
}
