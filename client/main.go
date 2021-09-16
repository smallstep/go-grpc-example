package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	var address, caCert, certFile, keyFile string
	hostname, _ := os.Hostname()
	if hostname != "" {
		hostname += ":443"
	}

	flag.StringVar(&address, "address", hostname, "The server url.")
	flag.StringVar(&certFile, "cert", "", "Server certificate file.")
	flag.StringVar(&keyFile, "key", "", "Private key file.")
	flag.StringVar(&caCert, "cacert", "", "CA certificate to verify peer against.")
	flag.Parse()

	switch {
	case address == "":
		log.Fatalln("flag --address is required")
	case certFile != "" && keyFile == "":
		log.Fatalln("flag --cert requires the --key flag")
	case keyFile != "" && certFile == "":
		log.Fatalln("flag --key requires the --cert flag")
	}

	// Configure TLS
	tlsConfig := &tls.Config{}
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("error loading X.509 key pair: %v", err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}
	if caCert != "" {
		b, err := os.ReadFile(caCert)
		if err != nil {
			log.Fatalf("error reading %s: %v", caCert, err)
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(b)
		tlsConfig.RootCAs = pool
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	fmt.Printf("What's your name? ")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &helloworld.HelloRequest{
		Name: strings.TrimSpace(name),
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Printf("Greeting: %s\n", r.GetMessage())
}
