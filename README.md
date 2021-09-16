# go-grpc-example

An example of using TLS with gRPC in Go.

## Prerequisites

The following examples requires [`step`](https://github.com/smallstep/cli/) and [`step-ca`](https://github.com/smallstep/certificates/).

Once installed initialize `step-ca` with:

```shell
step ca init
```

Add a new ACME provisioner running:

```shell
step ca provisioner add --type acme acme
```

And start step-ca:

```shell
step-ca $(step path)/config/ca.json
```

## Test

Before continuing compile the example:

```shell
make
```

### Using ACME

Run the ACME server using your private CA, let's say it's in https://localhost:9000:

```shell
bin/server-acme --directory https://localhost:9000/acme/acme/directory \
--cacert $(step path)/certs/root_ca.crt
```

And test it with `grpcurl`:

```shell
$ grpcurl -cacert $(step path)/certs/root_ca.crt -d '{"name":"Smallstep"}' $(hostname):443 helloworld.Greeter/SayHello
{
  "message": "Hello Smallstep"
}
```

Or the client

```shell
$ bin/client --cacert ~/.step/certs/root_ca.crt
What's your name? Smallstep
Greeting: Hello Smallstep
```

### Using a certificate

First create a certificate running:

```shell
step ca certificate $(hostname) local.crt local.key
```

And run `server-cert` with:

```shell
bin/server-cert --cert local.crt --key local.key 
```

And you can test it in the same way as before.

### mTLS

To enable mTLS to `server-acme` or `server-cert` just add the `--mtls` flag to
the previous commands. And if you haven't installed step's root certificate in
your truststore, make sure to add `--cacert $(step path)/certs/root_ca.crt` too.

Run `bin/server-acme` or `bin/server-cert`

```shell
bin/server-acme --directory https://localhost:9000/acme/acme/directory 
--cacert $(step path)/certs/root_ca.crt \
--cert local.crt --key local.key
```

```shell
bin/server-cert --cacert $(step path)/certs/root_ca.crt \
--cert local.crt --key local.key
```

And test it with the same or a different certificate from `step-ca`:

```shell
bin/client --cacert $(step path)/certs/root_ca.crt \
--cert local.crt --key local.key
```
