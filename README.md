[![CircleCI](https://circleci.com/gh/tkeech1/githubkey.svg?style=svg)](https://circleci.com/gh/tkeech1/githubkey)
[![codecov](https://codecov.io/gh/tkeech1/githubkey/branch/master/graph/badge.svg)](https://codecov.io/gh/tkeech1/githubkey)
[![Go Report Card](https://goreportcard.com/badge/github.com/tkeech1/githubkey)](https://goreportcard.com/report/github.com/tkeech1/githubkey)

A library to create and delete Github SSH deploy keys.

Usage

```package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/tkeech1/githubkey"
)

func main() {
	var githubUsername, githubPassword, repo, keyTitle, newKey string
	var readOnly bool
	flag.StringVar(&githubUsername, "githubUsername", "", "Github username")
	flag.StringVar(&githubPassword, "githubPassword", "", "Github password")
	flag.StringVar(&repo, "repo", "", "Github repo")
	flag.StringVar(&keyTitle, "keyTitle", "", "Github key name")
	flag.BoolVar(&readOnly, "readOnly", true, "Specifies if the Github key is read-only")
	flag.StringVar(&newKey, "newKey", "", "Github key")
	flag.Parse()

	client := &http.Client{}

	key, err := githubkey.GetDeployKey(client, githubUsername, githubPassword, repo, keyTitle)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if key.ID >= 0 {
		err := githubkey.DeleteDeployKey(client, githubUsername, githubPassword, repo, key.ID)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	githubKey, err := githubkey.CreateDeployKey(client, githubUsername, githubPassword, repo, keyTitle, newKey, readOnly)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("Created new key in repository: %s\n", repo)
	log.Printf("Key Name: %s\n", githubKey.Title)
	log.Printf("New Key ID: %d\n", githubKey.ID)
	log.Printf("Read-only: %t\n", githubKey.ReadOnly)
}```