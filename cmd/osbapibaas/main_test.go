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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func newClient() *http.Client {
	// Load client cert
	cert, err := tls.LoadX509KeyPair("../../pki/client.cert.pem", "../../pki/client.key.pem")
	if err != nil {
		panic(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile("../../pki/ca.cert.pem")
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	return client
}

func TestEcho(t *testing.T) {
	// Bring up a local https server
	port := 8443
	api := fmt.Sprintf("https://localhost.localdomain:%d/echo", port)
	setup(port)
	go listenAndServeTLS()

	// Make a https client
	client := newClient()

	testcases := []struct {
		name    string
		message string
	}{
		{"hi", "hi"},
		{"hello", "hello"},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Post(api, "Content-Type: application/x-www-form-urlencoded", strings.NewReader(tc.message))

			if err != nil {
				t.Fatal(err)
			}

			reply, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				t.Fatal(err)
			}

			if string(reply) != tc.message {
				t.Fatalf("reply error | expected %s | got %s", tc.message, string(reply))
			}
		})
	}
}
