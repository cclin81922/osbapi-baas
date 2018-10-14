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
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.Println("Try `curl --cert pki/client.cert.pem --key pki/client.key.pem --cacert pki/ca.cert.pem https://localhost.localdomain:8443`")

	if caCert, err := ioutil.ReadFile("pki/ca.cert.pem"); err != nil {
		log.Fatal(err)
	} else {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		tlsConfig.BuildNameToCertificate()

		server := &http.Server{
			Addr:      ":8443",
			TLSConfig: tlsConfig,
		}

		log.Fatal(server.ListenAndServeTLS("pki/server.cert.pem", "pki/server.key.pem"))
	}
}
