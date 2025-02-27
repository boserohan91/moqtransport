package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/mengelbart/moqtransport"
)

func main() {
	certFile := flag.String("cert", "localhost.pem", "TLS certificate file")
	keyFile := flag.String("key", "localhost-key.pem", "TLS key file")
	addr := flag.String("addr", "localhost:8080", "listen address")
	wt := flag.Bool("webtransport", false, "Use webtransport instead of QUIC")
	flag.Parse()

	if err := run(*addr, *wt, *certFile, *keyFile); err != nil {
		log.Fatal(err)
	}
}

func run(addr string, wt bool, certFile, keyFile string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tlsConfig, err := generateTLSConfigWithCertAndKey(certFile, keyFile)
	if err != nil {
		log.Printf("failed to generate TLS config from cert file and key, generating in memory certs: %v", err)
		tlsConfig = generateTLSConfig()
	}

	s := moqtransport.Server{
		Handler:   moqtransport.SessionHandlerFunc(handler()),
		TLSConfig: tlsConfig,
	}
	if wt {
		return s.ListenWebTransport(ctx, addr)
	}
	return s.ListenQUIC(ctx, addr)
}

func handler() moqtransport.SessionHandlerFunc {
	return func(p *moqtransport.Session) {
		go func() {
			s, err := p.ReadSubscription(context.Background())
			if err != nil {
				panic(err)
			}
			if fmt.Sprintf("%v/%v", s.Namespace(), s.Trackname()) != "clock/second" {
				s.Reject(errors.New("unknown namespace/trackname"))
			}
			t := s.Accept()
			go func() {
				ticker := time.NewTicker(time.Second)
				for ts := range ticker.C {
					w, err := t.StartCancellableObject()
					if err != nil {
						log.Println(err)
						return
					}
					if _, err := fmt.Fprintf(w, "%v", ts); err != nil {
						log.Println(err)
						return
					}
				}
			}()
		}()
		if err := p.Announce(context.Background(), "clock"); err != nil {
			panic(err)
		}
	}
}

func generateTLSConfigWithCertAndKey(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"moq-00"},
	}, nil
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"moq-00"},
	}
}
