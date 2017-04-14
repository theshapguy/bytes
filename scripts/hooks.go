package main

/*

Copyright 2017 Shapath Neupane (@theshapguy)

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

------------------------------------------------

Listner - written in Go because it's native web server is much more robust than Python. Plus its fun to write Go!

NOTE: Make sure you are using the right domain for travis [.com] or [.org]

*/

import (
    "bytes"
    "crypto"
    "crypto/rsa"
    "crypto/sha1"
    "crypto/x509"
    "encoding/base64"
    "encoding/json"
    "encoding/pem"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
)

var logPrint = log.Println

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

func PayloadSignature(r *http.Request) ([]byte, error) {

    signature := r.Header.Get("Signature")
    b64, err := base64.StdEncoding.DecodeString(signature)
    if err != nil {
        return nil, errors.New("cannot decode signature")
    }

    return b64, nil
}

func parsePublicKey(key string) (*rsa.PublicKey, error) {

    // https://golang.org/pkg/encoding/pem/#Block
    block, _ := pem.Decode([]byte(key))

    if block == nil || block.Type != "PUBLIC KEY" {
        return nil, errors.New("invalid public key")
    }

    publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, errors.New("invalid public key")
    }

    return publicKey.(*rsa.PublicKey), nil

}

func TravisPublicKey() (*rsa.PublicKey, error) {
    // NOTE: Use """https://api.travis-ci.com/config""" for private repos.
    response, err := http.Get("https://api.travis-ci.org/config")

    if err != nil {
        return nil, errors.New("cannot fetch travis public key")
    }
    defer response.Body.Close()

    decoder := json.NewDecoder(response.Body)
    var t ConfigKey
    err = decoder.Decode(&t)
    if err != nil {
        return nil, errors.New("cannot decode travis public key")
    }

    key, err := parsePublicKey(t.Config.Notifications.Webhook.PublicKey)
    if err != nil {
        return nil, err
    }

    return key, nil

}

func PayloadDigest(payload string) []byte {
    hash := sha1.New()
    hash.Write([]byte(payload))
    return hash.Sum(nil)

}

func RespondWithError(w http.ResponseWriter, m string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(401)
    message := fmt.Sprintf("{\"message\": \"%s\"}", m)
    w.Write([]byte(message))
}

func RespondWithSuccess(w http.ResponseWriter, m string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    message := fmt.Sprintf("{\"message\": \"%s\"}", m)
    w.Write([]byte(message))
}

func DeployHandler(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        RespondWithError("only travis ci allowed here")
    }

    key, err := TravisPublicKey()
    if err != nil {
        RespondWithError(w, err.Error())
        return
    }
    signature, err := PayloadSignature(r)
    if err != nil {
        RespondWithError(w, err.Error())
        return
    }
    payload := PayloadDigest(r.FormValue("payload"))

    err = rsa.VerifyPKCS1v15(key, crypto.SHA1, payload, signature)

    if err != nil {
        RespondWithError(w, errors.New("unauthorized payload").Error())
        return
    }

    RespondWithSuccess(w, "payload verified")
    Deploy()
}

func SlackSend(payload string) {
    // notify = Get Slack URL

    notify := os.Getenv("SLACK_HOOK")

    if notify == "" {
        logPrint("Slack wehhook not found")
        return
    }

    _, err := http.Post(notify, "application/json", bytes.NewBufferString(payload))

    if err != nil {
        // Simple Retry for Now
        http.Post(notify, "application/json", bytes.NewBufferString(payload))
    }

}

func Deploy() {

    cmd, err := exec.Command("sh", "./scripts/deploy.sh").Output()
    if err != nil {

        payload := `{"channel": "#blog", "username": "deploybot", "text": " _deployment failed_: errors with deploy shell script. Check supervisord logs", "icon_emoji": ":warning:"}`
        logPrint(string(cmd))
        SlackSend(payload)
        return
    }

    payload := `{"channel": "#blog", "username": "deploybot", "text": "New posts are now _live_!! :pray:", "icon_emoji": ":tada:"}'`
    SlackSend(payload)
}

func main() {
    http.HandleFunc("/", DeployHandler)
    log.Fatal(http.ListenAndServe(":5000", nil))
}
