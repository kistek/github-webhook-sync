package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/chrishiestand/github-webhook-sync/repotrack"
	"github.com/chrishiestand/github-webhook-sync/webhook"
	"github.com/joho/godotenv"
	"golang.org/x/sys/unix"
	yaml "gopkg.in/yaml.v2"
)

//When you create a new webhook, we'll send you a simple ping event to let you know you've set up the webhook correctly. This event isn't stored so it isn't retrievable via the Events API. You can trigger a ping again by calling the ping endpoint.

// Ping Event Payload
// Key	Value
// zen	Random string of GitHub zen
// hook_id	The ID of the webhook that triggered the ping
// hook	The webhook configuration

//
//Header	Description
// X-GitHub-Event	Name of the event type that triggered the delivery.
// X-GitHub-Delivery	A GUID to identify the delivery.
// X-Hub-Signature	The HMAC hex digest of the response body. This header will be sent if the webhook is configured with a secret. The HMAC hex digest is generated using the sha1 hash function and the secret as the HMAC key.

// We only care about the `push` event
//
// environment variables:
// port
// message signing keys
// path to repo root
// TODO git credentials
// TODO whitelist git repos or regexes?
func main() {

	dotenvPath := os.Getenv("DOTENV_PATH")

	if dotenvPath == "" {
		dotenvPath = "./.env"
	}

	err := godotenv.Load(dotenvPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	repoSourcePath := os.Getenv("REPO_SOURCE_PATH")
	repoTargetPath := os.Getenv("REPO_TARGET_PATH")
	port := os.Getenv("PORT")

	rts, err := ConfigPathToRepoTracks(repoSourcePath)

	if err != nil {
		log.Fatalln(err)
	}

	start(port, rts, repoTargetPath)
}

func start(port string, rts []repotrack.RepoTrack, repoRootPath string) {

	mainMux := http.NewServeMux()

	mainMux.HandleFunc("/webhook", webhookHandler(rts, repoRootPath)) // TODO allow main path to be configurable
	mainMux.HandleFunc("/_ready", readyHandler(repoRootPath))
	mainMux.HandleFunc("/_alive", aliveHandler())

	log.Printf("Starting on port %s...", port)

	if err := http.ListenAndServe(":"+port, mainMux); err != nil {
		panic(err)
	}
}

func webhookHandler(rts []repotrack.RepoTrack, repoRootPath string) func(http.ResponseWriter, *http.Request) {

	var keys = []string{"TO", "DO"}

	return func(w http.ResponseWriter, req *http.Request) {

		contentType := req.Header.Get("content-type")

		if contentType != "application/json" {

			errorString := fmt.Sprintf("request http header content-type must be applicaton/json but is: %s", contentType)

			log.Println(errorString)

			w.WriteHeader(http.StatusBadRequest)
			_, err := io.WriteString(w, errorString)

			if err != nil {
				log.Panicln("failed to write bad content-type response body in webhookHandler")
			}
			return
		}

		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			log.Println("Failed to read body of webhook request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var wh webhook.Webhook

		if err := json.Unmarshal(body, &wh); err != nil {

			log.Println(fmt.Sprintf("could not unmarshal json: %s", string(body)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println(fmt.Sprintf("webhook: %s", wh))

		signatureHeader := req.Header.Get("X-Hub-Signature")
		signature := strings.Split(signatureHeader, "=")[1]
		computedSignatures := make([]string, 0, len(keys))

		log.Println(fmt.Sprintf("blah request %s %s\ncontent-type: %s, repo name: %s", req.Method, req.URL.Path, req.Header.Get("content-type"), wh.Repository.FullName))

		if verifyHubSignature(signature, body, wh.Repository, &rts, &computedSignatures) != true {
			log.Println(fmt.Sprintf("invalid signature hash %s, computed signatures:\n%s", signatureHeader, strings.Join(computedSignatures, "\n")))
			w.WriteHeader(http.StatusUnauthorized)
			return

		}

	}
}

// TODO - return 200 only when *all* repos are cloned successfully
// pass in repos struct
func readyHandler(repoRootPath string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {

		log.Println(fmt.Sprintf("request %s %s", req.Method, req.URL.Path))

		if err := unix.Access(repoRootPath, unix.W_OK); err != nil {
			log.Println(fmt.Sprintf("path %s not available or writable: %s", repoRootPath, err))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func aliveHandler() func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func verifyHubSignature(signature string, bodyBytes []byte, whr webhook.Repository, rts *[]repotrack.RepoTrack, computedSignatures *[]string) bool {

	// TODO in order to allow no signatures we must be able to map repositories to known-valid keys, otherwise anyone could push to a secured repo. So this requires app specific configuration files instead of .env for this data.
	// webhooks not configured with a secret will not contain a signature
	// if signature == "" {
	// 	return true
	// }
	//
	//
	rt := repotrack.NewRepoTrack()

	for _, testRt := range *rts {
		if testRt.URL == whr.HTTPURL || testRt.URL == whr.SSHURL {
			rt = testRt
		}
	}

	if rt.URL == "" {
		// could not find a matching repository
		log.Printf("Could not find matching repository for webhook: %v", whr)
		return false
	}

	if rt.WebhookSecretRequired == false {
		return true
	}

	return false
	// keyBytes := make([][]byte, 0)
	//
	// for _, str := range keys {
	// 	key := []byte(str)
	// 	keyBytes = append(keyBytes, key)
	//
	// }
	//
	// signatureBytes, err := hex.DecodeString(signature)
	//
	// if err != nil {
	// 	log.Println("got bad signature: " + signature)
	// 	return false
	// }
	//
	// for _, key := range keyBytes {
	//
	// 	mac := hmac.New(sha1.New, key)
	// 	mac.Write(bodyBytes)
	// 	computedMac := mac.Sum(nil)
	// 	match := hmac.Equal(signatureBytes, computedMac)
	//
	// 	log.Printf("computedMac: %v ", hex.EncodeToString(computedMac))
	// 	*computedSignatures = append(*computedSignatures, hex.EncodeToString(computedMac))
	//
	// 	if match == true {
	// 		return true
	// 	}
	//
	// }
	// return false
}

// TODO add a config path watch function
// ConfigPathToRepoTracks reads repo metadata from folder and converts each file into a struct
func ConfigPathToRepoTracks(repoSourcePath string) ([]repotrack.RepoTrack, error) {
	repoFiles, err := ioutil.ReadDir(repoSourcePath)

	if err != nil {
		return nil, err
	}

	rts := make([]repotrack.RepoTrack, 0)

	for _, file := range repoFiles {
		filename := file.Name()
		repoPath := path.Join(repoSourcePath, filename)
		fmt.Println("filename: " + filename)

		content, err := ioutil.ReadFile(repoPath)
		if err != nil {
			return nil, err
		}

		rt := repotrack.NewRepoTrack()

		err2 := yaml.Unmarshal(content, &rt)
		if err2 != nil {
			return nil, err2
		}

		// TODO
		// rt.Populate()
		fmt.Println("RT: ", rt)
	}

	return rts, nil
}
