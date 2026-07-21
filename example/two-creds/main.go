// Two S3 clients with two different static credential sets (distinct Tags).
//
// Env prefix: CLIENT_S3_TWO_CREDS
package main

import (
	"log"

	clients3 "github.com/omcrgnt/client-s3"
	"github.com/omcrgnt/app"
	"github.com/omcrgnt/res/unique"
)

const envPrefix = "CLIENT_S3_TWO_CREDS"

type (
	assets  struct{}
	backups struct{}
)

type (
	assetsCred  = clients3.CredentialsStatic[assets]
	backupsCred = clients3.CredentialsStatic[backups]
)

type catalog struct {
	Assets      *clients3.Client[*assetsCred]  `ecfg:"S3_ASSETS"`
	Backups     *clients3.Client[*backupsCred] `ecfg:"S3_BACKUPS"`
	AssetsCred  *assetsCred                    `ecfg:"S3_ASSETS_CRED"`
	BackupsCred *backupsCred                   `ecfg:"S3_BACKUPS_CRED"`
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
