# Server

```
pneutrinoutil-server -- pneutrinoutil http server

e.g.
pneutrinoutil-server --mysqlDSN DSN --redisDSN DSN

Flags:
      --accessLogFile string              access log file; stdout, stderr are available (default "stderr")
      --debug                             enable debug logs
      --host string                       server host
      --mysqlConnMaxLifetimeSeconds int   max amount of time a connection may be reused (default 300)
      --mysqlDSN string                   format: USER:PASS@tcp(HOST:PORT)/DB
      --mysqlMaxIdleConns int             maximum number of connections in the idle connection pool (default 3)
      --mysqlMaxOpenConns int             maximum number of open connections to the database (default 3)
  -p, --port uint                         server port (default 9101)
      --processTimeoutSeconds int         duration pneutrinoutil timeout (default 1200)
      --redisDSN string                   format: redis://HOST:PORT/DB
      --shutdownPeriodSeconds int         duration the server needs to shut down gracefully (default 10)
      --storageBucket string              storage bucket (default "pneutrinoutil-worker")
      --storageDir string                 local storage directory; $HOME/.pneutrinoutil-worker/storage or .pneutrinoutil-worker/storage if no $HOME
      --storagePath string                storage base path
      --storageS3                         use s3 as the object storage; if set, storageDir is ignored
      --version                           print pneutrinoutil-server version
```
