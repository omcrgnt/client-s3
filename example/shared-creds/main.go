// Two S3 clients share one static credentials resource (same Tag / one env block).
//
// Env prefix: CLIENT_S3_SHARED_CREDS
package main

import (
	"log"

	clients3 "github.com/omcrgnt/client-s3"
	"github.com/omcrgnt/app"
	"github.com/omcrgnt/res/unique"
)

const envPrefix = "CLIENT_S3_SHARED_CREDS"

type static = clients3.CredentialsStatic[clients3.Default]

type catalog struct {
	Assets  *clients3.Client[*static] `ecfg:"S3_ASSETS"`
	Backups *clients3.Client[*static] `ecfg:"S3_BACKUPS"`
	Cred    *static                   `ecfg:"S3_CREDENTIALS_STATIC"`
}

func main() {
	c := catalog{}
	pipeline := app.Pipeline{
		Registry:  unique.Global(),
		EnvPrefix: envPrefix,
	}
	if err := app.Run(&c, pipeline); err != nil {
		log.Fatal(err)
	}
}
