# client-s3

S3-compatible object storage client for omcrgnt apps (AWS S3 / MinIO): ecfg catalog resource,
AWS SDK Go v2 under the hood, OpenTelemetry traces via `otelaws` + smithy adapters.

One client instance binds to **one bucket**. Credentials are a separate SDI port
([CredentialsProvider](credentials.go)), wired like `srv-http.Server[T]` / `logger.Output`.

## Catalog

```go
import clients3 "github.com/omcrgnt/client-s3"

type catalog struct {
	S3     *clients3.Client[*clients3.CredentialsStatic] `ecfg:"S3"`
	S3Cred *clients3.CredentialsStatic                   `ecfg:"S3_CREDENTIALS_STATIC"`
}
```

`Client[C]` depends on concrete credentials type `C` (SDI `Deps` / `Inject`).
Multiple clients with separate credential sets — later (distinct `C` / tags).

## Environment

Prefix comes from the app (`EnvPrefix`); fields below are relative to each slot.

### `S3` (client)

| Variable | Description |
|----------|-------------|
| `S3_LABEL` | Resource label (otel) |
| `S3_ENDPOINT` | S3 API base URL (`http.v1.URL`, required http/https) |
| `S3_REGION` | Region (default `us-east-1`) |
| `S3_BUCKET` | Bucket name (required) |
| `S3_USE_PATH_STYLE` | Reserved (path-style is always on for custom endpoints) |

### `S3_CREDENTIALS_STATIC`

| Variable | Description |
|----------|-------------|
| `S3_CREDENTIALS_STATIC_ACCESS_KEY` | Access key (required) |
| `S3_CREDENTIALS_STATIC_SECRET_KEY` | Secret key (required) |

## API

| Method | Role |
|--------|------|
| `Put` | Upload object (`PutOptions`: ContentType, ContentLength) |
| `Get` | Download body + `Head` metadata |
| `Head` | Metadata only |
| `Delete` | Remove object |
| `PresignGet` | Time-limited GET URL |
| `Ready` | `HeadBucket` readiness probe |

`ErrNotFound` / `IsNotFound` map missing-object errors.

SDK client is created in `Inject` after SDI resolves credentials (not in `Build`).

## Telemetry

On `Inject`, the package attaches OpenTelemetry middleware so S3 calls emit spans to the process
`TracerProvider` (set by `omcrgnt/telemetry` before `runner.Run`).

## MinIO (local)

```bash
# example env (app prefix CLIENT_S3_EXAMPLE)
CLIENT_S3_EXAMPLE_S3_LABEL=assets
CLIENT_S3_EXAMPLE_S3_ENDPOINT=http://127.0.0.1:9000
CLIENT_S3_EXAMPLE_S3_REGION=us-east-1
CLIENT_S3_EXAMPLE_S3_BUCKET=bemvpgame-assets

CLIENT_S3_EXAMPLE_S3_CREDENTIALS_STATIC_ACCESS_KEY=minioadmin
CLIENT_S3_EXAMPLE_S3_CREDENTIALS_STATIC_SECRET_KEY=minioadmin
```

See `example/app` for a minimal `app.Run` catalog.
