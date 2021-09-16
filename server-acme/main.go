package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"go.step.sm/go-grpc-example/health"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type helloServer struct {
	helloworld.UnimplementedGreeterServer
}

func (*helloServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("method=SayHello name=%s", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	var mtls bool
	var address, domain, directory, caCert string
	hostname, _ := os.Hostname()

	flag.StringVar(&address, "address", ":443", "The address to listen to.")
	flag.StringVar(&domain, "domain", hostname, "The domain to use.")
	flag.StringVar(&directory, "directory", "", "Url of the ACME server.")
	flag.StringVar(&caCert, "cacert", "", "CA certificate to verify peer against.")
	flag.BoolVar(&mtls, "mtls", false, "Enables mTLS mode, usually used with --cacert.")
	flag.Parse()

	if domain == "" {
		log.Fatalln("flag --domain is required")
	}
	if !strings.Contains(strings.Trim(domain, "."), ".") {
		log.Fatalln("acme/autocert does not support domains without a .")
	}

	var pool *x509.CertPool
	if caCert != "" {
		b, err := os.ReadFile(caCert)
		if err != nil {
			log.Fatalf("error reading %s: %v", caCert, err)
		}
		pool = x509.NewCertPool()
		pool.AppendCertsFromPEM(b)
	}

	var client *acme.Client
	if directory != "" {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Fatalf("error generating key: %v", err)
		}

		var httpClient *http.Client
		if pool != nil {
			transport := http.DefaultTransport.(*http.Transport).Clone()
			transport.TLSClientConfig = &tls.Config{
				RootCAs: pool,
			}
			httpClient = &http.Client{
				Transport: transport,
			}
		}

		client = &acme.Client{
			Key:          key,
			DirectoryURL: directory,
			HTTPClient:   httpClient,
		}
	}

	m := &autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		HostPolicy:  autocert.HostWhitelist(domain),
		Client:      client,
		RenewBefore: 8 * time.Hour,
	}

	tlsConfig := m.TLSConfig()
	if mtls {
		// Require client certificate.
		// It will use system pool if --cacert is not given.
		tlsConfig.ClientCAs = pool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	tlsConfig.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cert, err := m.GetCertificate(hello)
		if err != nil {
			log.Println(err)
		}
		return cert, err
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer(opts...)
	helloworld.RegisterGreeterServer(srv, &helloServer{})
	grpc_health_v1.RegisterHealthServer(srv, health.New())
	reflection.Register(srv)

	log.Printf("server listening at %v with domain %s", lis.Addr(), domain)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
