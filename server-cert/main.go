package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net"
	"os"

	"go.step.sm/crypto/tlsutil"
	"go.step.sm/go-grpc-example/health"

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
	var address, certFile, keyFile, caCert string
	flag.StringVar(&address, "address", ":443", "The address to listen to.")
	flag.StringVar(&certFile, "cert", "", "Server certificate file.")
	flag.StringVar(&keyFile, "key", "", "Private key file.")
	flag.StringVar(&caCert, "cacert", "", "CA certificate to verify peer against.")
	flag.BoolVar(&mtls, "mtls", false, "Enables mTLS mode, usually used with --cacert.")
	flag.Parse()

	if certFile == "" {
		log.Fatalln("flag --cert is required")
	}
	if keyFile == "" {
		log.Fatalln("flag --key is required")
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

	m, err := tlsutil.NewServerCredentialsFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("error creating server credentials: %v", err)
	}

	tlsConfig := m.TLSConfig()
	if mtls {
		// Require client certificate.
		// It will use system pool if --cacert is not given.
		tlsConfig.ClientCAs = pool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.GetConfigForClient = nil
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

	log.Printf("server listening at %v", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
