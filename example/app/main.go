package main

import (
	"log"

	clients3 "github.com/omcrgnt/client-s3"
	"github.com/omcrgnt/app"
	"github.com/omcrgnt/res/unique"
)

const envPrefix = "CLIENT_S3_EXAMPLE"

type catalog struct {
	S3     *clients3.Client[*clients3.CredentialsStatic] `ecfg:"S3"`
	S3Cred *clients3.CredentialsStatic                   `ecfg:"S3_CREDENTIALS_STATIC"`
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
