//    Copyright 2018 cclin
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	ca = `-----BEGIN CERTIFICATE-----
MIID4zCCAsugAwIBAgIJAJNowMdJfs8UMA0GCSqGSIb3DQEBCwUAMIGHMQswCQYD
VQQGEwJUVzEPMA0GA1UECAwGVGFpd2FuMQ8wDQYDVQQHDAZUYWlwZWkxDjAMBgNV
BAoMBWNjbGluMQ4wDAYDVQQLDAVjY2xpbjERMA8GA1UEAwwIY2EuY2NsaW4xIzAh
BgkqhkiG9w0BCQEWFGNjbGluODE5MjJAZ21haWwuY29tMB4XDTE4MDIyNzA3NTI0
MloXDTI4MDEwNjA3NTI0MlowgYcxCzAJBgNVBAYTAlRXMQ8wDQYDVQQIDAZUYWl3
YW4xDzANBgNVBAcMBlRhaXBlaTEOMAwGA1UECgwFY2NsaW4xDjAMBgNVBAsMBWNj
bGluMREwDwYDVQQDDAhjYS5jY2xpbjEjMCEGCSqGSIb3DQEJARYUY2NsaW44MTky
MkBnbWFpbC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDKveco
eF+ZlJDgiX/Vtlg492DC/wrF4ZFzdZVpIxuD0kJxS9w8i2Xhf5tD72Uu8KAHR28W
u3NjjDtzH/KQWJd9XNeN6lkzTTTQUzFUfP3uYEELnmaTRLzdO5UDJzC/n98YOJwE
HK9K9M7C++EEjbzNYEBdhIaLlbU/p4mZSOqIUhOLjg5EzlaCgvgTuvS9YEh5hXix
5by3j9GCRIj39E8R+gdXkr4XqTVIhg4xUF81iJk1yFYEgoO2gJd0KEXCPW8NqKjZ
lUcYRWw3LvKcuhqjWfOxQVIMZeRw6nM2syJfU6umO0mFVwl19ajcq6Ic2QxZTBEf
3osMjh9966/xefkBAgMBAAGjUDBOMB0GA1UdDgQWBBRiUWcH93pfph/aBKJcsvgq
ObVJpzAfBgNVHSMEGDAWgBRiUWcH93pfph/aBKJcsvgqObVJpzAMBgNVHRMEBTAD
AQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAzMbnRbXLYOFawr1N3i2NOBwqKjyare+TD
JSYnDvIQjnPm1TVLhCRN10XH98nqbT8tx1pWVmmH0XtINVYq9KaWm4089oiopYc5
MG4Ru3a3vdNADJBvh+EUtO3pAYbHExIfBCP0Vo/gp3n1LLcUn49sIkHXKmbbzPcP
scWwlL72mtmtcrbL4HGjX632xpvuyc1ZzIiHcdwKLnxrUtZcl6oQCMlGVGJHLJyH
RWQ8gXZ71sbYOHghHrIcK4XG+ChaZJxYskd3112RIeC/5/QK8U8FKNo2KhcOOv1n
cUupgdgewtbVnl1p09PRnnBkEHhUZeXItoCUAeAC36Lhb981v5h1
-----END CERTIFICATE-----
`
	key = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAzwaCAE57wEdKOA3VlOd1i/WECflvxFuCtxTqYSYocGQoWzV0
petHzD4CPDmKmEUDXSEr7jhY1cQhbrsVxbSdMWyAGWaLTYHP3e0Jb6rCfeEm1OXQ
M8Op7IOeqBhURe6EDo7cm8maCIquBaNLeGXadCr1lSGqnvDjOHhqzr0+WdLoZmq7
rBpzyGNhj+jl17a44SM7Jab4gevJ7kF6+iXn497PMuSpw71JZO3VSfy/GUrfgEpi
d+IlYZynmneF3MHimkupUVXELVWIp9dLX3DJ0KdPF+7D34InSEq33A6adbMozAe6
O9PjhSk9bVsxltl4jdpF5QzCzizqxslDk1o6XQIDAQABAoIBAAD97P2HYrxnGwnt
twpBmaSUBo/trAVq0tOBvCW/Aw6WzEKznN12pR9rZKNZOzrDieKWWBmKF5Len7Ji
HxaVaPNlq69zeFNkvdQ4YIUyckAcJg17FGZLF7NUZhw0EaNMI0HWmWP9sUk9MNxC
DgiDpGpOEmMmiCS+zrfki851xai0gcV0lgHTmgCCUN3tqphbYZb9HlL2wi4wt7DB
mL9OlgapHibW3648ZrOi4aV8D9oJqxgsIyuDgjNpQjkGeUIzfSk9i/+Ypuio8uwC
9bSQdebiksvC3wQwtSyClifSq4tLwEUkRUeeoEBEimDsbjpHsub6tQbcl7bwKyP9
bRmp+RkCgYEA/C5bMuZmpN07aSWfj5zmQK4bDEJ3QLGC3VBym2tFyE6BrZHGTj6s
F5pJ2dA8t3Jx5YDEwMmzrAUWB+Yz2XWpvouAYfC6pDXoJQOEwMBS0zKSmbVR6HkB
dBZV1+ofTFKHI9x7EdN2dEQrkBxaX0NUgZNNa9ii+oKx7BwcvXix+zMCgYEA0ikY
IBRJ9fZPnOdn5BIRvGG9vK4gHCLWRRiuqqSd0EyvoNlaQ8l8DEIDI/EFDhe/mEsn
mC5xj3CIjYlNObZiCV++2tPYrMG7/2TTlmxvHKd6EGyq5fL0/TVpu70VUw6Mk7LH
A88q/gQDrB9tIaUbpOWr0+CysDgCRKbzYVe+dC8CgYEAy/1fjkProc7HYR2q/Yuo
gYeUn40gU/eDaSzLGEdk8kv3AAUcSWzO3mTS+ltE0gvEcCaCgYRnT23pzUf8hxpz
zYugtRj6kRx+BXrcJuMr3GVbSvTuJcPEVjg/BmH/IUjcwjh2YQwSFKiUKIWW4Nph
AFO8W9GovEV+UQTIhsecCRsCgYAg7HUuGV+Y29SPFSWOclI6++j4lSLMpZyByKMc
cpuSlWDyRvrAIeGAHhtV1x1entPSLPvv+F6sBQovejIR94OWSlyg9Y09S0CDey02
pJgnmgkiZ5PCYHSG8oY09iNQFrhpLxnEfAEVOFXG8klreu1AwQZRNCNqPewFC06X
kmJw2wKBgAmyjpKaWooIFZOs6pe9q5FI1q2jkiWAXPkZh25Sm9U8QFpYyIeVq57H
lQmnIKzwJhkVHAqegKYBfWimE6OD/slT7NosjR3+Dod9v8Xa59hupQnVZDvIZqqP
3CBKWLJppjCXSpqCo1OwwvWTdQg8DgH6//lPIQIqJ6/fZBD3FcRc
-----END RSA PRIVATE KEY-----
`
	cert = `-----BEGIN CERTIFICATE-----
MIIDmTCCAoECCQDHvcJjcXEkZjANBgkqhkiG9w0BAQUFADCBhzELMAkGA1UEBhMC
VFcxDzANBgNVBAgMBlRhaXdhbjEPMA0GA1UEBwwGVGFpcGVpMQ4wDAYDVQQKDAVj
Y2xpbjEOMAwGA1UECwwFY2NsaW4xETAPBgNVBAMMCGNhLmNjbGluMSMwIQYJKoZI
hvcNAQkBFhRjY2xpbjgxOTIyQGdtYWlsLmNvbTAeFw0xODAyMjgwODQ4NDlaFw0x
OTAyMjgwODQ4NDlaMIGUMQswCQYDVQQGEwJUVzEPMA0GA1UECBMGVGFpd2FuMQ8w
DQYDVQQHEwZUYWlwZWkxDjAMBgNVBAoTBWNjbGluMQ4wDAYDVQQLEwVjY2xpbjEe
MBwGA1UEAxMVbG9jYWxob3N0LmxvY2FsZG9tYWluMSMwIQYJKoZIhvcNAQkBFhRj
Y2xpbjgxOTIyQGdtYWlsLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBAM8GggBOe8BHSjgN1ZTndYv1hAn5b8RbgrcU6mEmKHBkKFs1dKXrR8w+Ajw5
iphFA10hK+44WNXEIW67FcW0nTFsgBlmi02Bz93tCW+qwn3hJtTl0DPDqeyDnqgY
VEXuhA6O3JvJmgiKrgWjS3hl2nQq9ZUhqp7w4zh4as69PlnS6GZqu6wac8hjYY/o
5de2uOEjOyWm+IHrye5Bevol5+PezzLkqcO9SWTt1Un8vxlK34BKYnfiJWGcp5p3
hdzB4ppLqVFVxC1ViKfXS19wydCnTxfuw9+CJ0hKt9wOmnWzKMwHujvT44UpPW1b
MZbZeI3aReUMws4s6sbJQ5NaOl0CAwEAATANBgkqhkiG9w0BAQUFAAOCAQEALSl4
7plFGrQ22IIQJ9uBU+MOtBPatzGUY3vVYcMVsI75s/FcUNZ64j2zygOXoKrD6Yys
XYqsQf7D0bUXusdTGhg7uYwJH6VmjfS2kOt22T1IdWJlTWtgu885tIJTj9R5pfhY
S4sogqc/yH2ZefOyWv17t3QU2RYbqKe6qtfUZJqLEzsuP2KVwGkjF0I8XnOICVj4
fkJXm4nEg3ZgpHohB1u5JqI6P1ru8qo7MNICKIg0rGAHOVEgReYpxiks4tdIUlNa
LzHO60WSht0a0czSido2ozU1NyRB5UiM16X6pAhjMNNd0fDOX2sCSZtWyPbUO2Tr
K4mhPL/qnKeks8nz9A==
-----END CERTIFICATE-----
`
)

func echo(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

var (
	flagPort int
)

func init() {
	flag.IntVar(&flagPort, "port", 443, "HTTPS PORT")
	flag.Parse()
}

func main() {
	log.Printf("Base url is https://localhost.localdomain:%d/", flagPort)

	// register http routes
	http.HandleFunc("/echo", echo)

	// make a ca cert pool object
	caCert := []byte(ca)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// make a tls config object
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// make a server object
	addr := fmt.Sprintf(":%d", flagPort)
	server := &http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
	}

	// make a temp server key file
	serverKey := []byte(key)
	serverKeyFile, err := ioutil.TempFile("", "server.key")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(serverKeyFile.Name())
	if _, err := serverKeyFile.Write(serverKey); err != nil {
		log.Fatal(err)
	}
	if err := serverKeyFile.Close(); err != nil {
		log.Fatal(err)
	}

	// make a temp server cert file
	serverCert := []byte(cert)
	serverCertFile, err := ioutil.TempFile("", "server.cert")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(serverCertFile.Name())
	if _, err := serverCertFile.Write(serverCert); err != nil {
		log.Fatal(err)
	}
	if err := serverCertFile.Close(); err != nil {
		log.Fatal(err)
	}

	// start listening
	log.Fatal(server.ListenAndServeTLS(serverCertFile.Name(), serverKeyFile.Name()))
}
