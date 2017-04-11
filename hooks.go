package main

/*

Copyright 2017 Shapath Neupane

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

------------------------------------------------

Listner - written in Golang because it's native web server is much more robust than Python.

NOTE: Make sure you are using the right domain for travis [.com] or [.org]

*/

import (
	// "crypto"
	"crypto/rsa"
	// "crypto/sha1"
	"crypto/x509"
	"encoding/base64"
    "encoding/pem"
	"encoding/json"
	"log"
	"net/http"
)

type Payload struct {
	ID             int         `json:"id"`
	Number         string      `json:"number"`
	Status         interface{} `json:"status"`
	StartedAt      interface{} `json:"started_at"`
	FinishedAt     interface{} `json:"finished_at"`
	StatusMessage  string      `json:"status_message"`
	Commit         string      `json:"commit"`
	Branch         string      `json:"branch"`
	Message        string      `json:"message"`
	CompareURL     string      `json:"compare_url"`
	CommittedAt    string      `json:"committed_at"`
	CommitterName  string      `json:"committer_name"`
	CommitterEmail string      `json:"committer_email"`
	AuthorName     string      `json:"author_name"`
	AuthorEmail    string      `json:"author_email"`
	Type           string      `json:"type"`
	BuildURL       string      `json:"build_url"`
	Repository     struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		OwnerName string `json:"owner_name"`
		URL       string `json:"url"`
	} `json:"repository"`
	Config struct {
		Notifications struct {
			Webhooks []string `json:"webhooks"`
		} `json:"notifications"`
	} `json:"config"`
	Matrix []struct {
		ID           int         `json:"id"`
		RepositoryID int         `json:"repository_id"`
		Number       string      `json:"number"`
		State        string      `json:"state"`
		StartedAt    interface{} `json:"started_at"`
		FinishedAt   interface{} `json:"finished_at"`
		Config       struct {
			Notifications struct {
				Webhooks []string `json:"webhooks"`
			} `json:"notifications"`
		} `json:"config"`
		Status         interface{} `json:"status"`
		Log            string      `json:"log"`
		Result         interface{} `json:"result"`
		ParentID       int         `json:"parent_id"`
		Commit         string      `json:"commit"`
		Branch         string      `json:"branch"`
		Message        string      `json:"message"`
		CommittedAt    string      `json:"committed_at"`
		CommitterName  string      `json:"committer_name"`
		CommitterEmail string      `json:"committer_email"`
		AuthorName     string      `json:"author_name"`
		AuthorEmail    string      `json:"author_email"`
		CompareURL     string      `json:"compare_url"`
	} `json:"matrix"`
}

type ConfigKey struct {
	Config struct {
		Host        string `json:"host"`
		ShortenHost string `json:"shorten_host"`
		Assets      struct {
			Host string `json:"host"`
		} `json:"assets"`
		Pusher struct {
			Key string `json:"key"`
		} `json:"pusher"`
		Github struct {
			APIURL string   `json:"api_url"`
			Scopes []string `json:"scopes"`
		} `json:"github"`
		Notifications struct {
			Webhook struct {
				PublicKey string `json:"public_key"`
			} `json:"webhook"`
		} `json:"notifications"`
	} `json:"config"`
}

var logPrint = log.Println

func PayloadSignature(h *http.Request) string {

	signature := h.Header.Get("HTTP_SIGNATURE")
	b64, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		log.Println("Error Base 64")
	}

	return string(b64)
}

func parsePublicKey(key string) *rsa.PublicKey {

	// https://golang.org/pkg/encoding/pem/#Block
	block, _ := pem.Decode([]byte(key))

	if block == nil || block.Type != "PUBLIC KEY" {
		log.Println("Error Parsing Public Key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        log.Fatal(err)
    }

    return publicKey.(*rsa.PublicKey)

}

func TravisPublicKey() *rsa.PublicKey {
	// NOTE: Use """https://api.travis-ci.com/config""" for private repos.
	response, err := http.Get("https://api.travis-ci.org/config")

	if err != nil {
		log.Println("Error Fetching Travis Config Private Key")
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var t ConfigKey
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
	}

	return parsePublicKey(t.Config.Notifications.Webhook.PublicKey)

}

func DeployHandler(w http.ResponseWriter, h *http.Request) {

    key := TravisPublicKey()
    signature := PayloadSignature(h)

    // decodeBytes := bytes.NewBufferString(h.FormValue("payload"))
    // decoder := json.NewDecoder(decodeBytes)
    // var p Payload
    // err := decoder.Decode(&p)
    // if err != nil {
    //     log.Println(err)
    // }
    logPrint("Key")
    logPrint(key)
    logPrint("\n")

    logPrint("payload")
    logPrint(h.FormValue("payload"))
    logPrint("\n")


    logPrint("signature")
    logPrint(signature)
    logPrint("\n")
}

func main() {
	http.HandleFunc("/", DeployHandler)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
