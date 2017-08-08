package main

import (
	"crypto/tls"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"golang.org/x/crypto/acme/autocert"
)

func StartTLSServer(cfg *Config, r *mux.Router) {
	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		//HostPolicy: autocert.HostWhitelist(cfg.ServerDomain),
	}
	var addr string = cfg.ServerAddressTLS
	if cfg.ServerAddressTLS == "" {
		addr = ":https"
	}
	s := &http.Server{
		Addr:      addr,
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   r,
	}
	logrus.Infof("Starting TLS server on %q", s.Addr)
	logrus.Fatal(s.ListenAndServeTLS("", ""))
}

func StartServer(cfg *Config, r *mux.Router) {
	var addr string = cfg.ServerAddress
	if cfg.ServerAddress == "" {
		addr = ":http"
	}
	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	logrus.Infof("Starting server on %q", s.Addr)
	logrus.Fatal(s.ListenAndServe())
}
