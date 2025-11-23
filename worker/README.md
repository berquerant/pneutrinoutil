# Worker

```
pneutrinoutil-worker -- pneutrinoutil process worker

e.g.
pneutrinoutil-worker --neutrinoDir /path/to/NEUTRINO --workDir /path/to/workingDirectory --pneutrinoutil /path/to/pneutrinoutil --mysqlDSN DSN --redisDSN DSN

Flags:
  -c, --concurrency int                   pneutrinoutil process concurrency (default 1)
      --debug                             enable debug logs
      --mysqlConnMaxLifetimeSeconds int   max amount of time a connection may be reused (default 300)
      --mysqlDSN string                   format: USER:PASS@tcp(HOST:PORT)/DB
      --mysqlMaxIdleConns int             maximum number of connections in the idle connection pool (default 3)
      --mysqlMaxOpenConns int             maximum number of open connections to the database (default 3)
  -n, --neutrinoDir string                NEUTRINO directory (default "./dist/NEUTRINO")
  -x, --pneutrinoutil string              pneutrinoutil executable (default "./dist/pneutrinoutil")
      --redisDSN string                   format: redis://HOST:PORT/DB
  -s, --shell string                      shell command to execute (default "bash")
      --shutdownPeriodSeconds int         duration the server needs to shut down gracefully (default 10)
      --storageBucket string              storage bucket (default "pneutrinoutil-worker")
      --storageDir string                 local storage directory; $HOME/.pneutrinoutil-worker/storage or .pneutrinoutil-worker/storage if no $HOME
      --storagePath string                storage base path
      --storageS3                         use s3 as the object storage; if set, storageDir is ignored
      --version                           print pneutrinoutil-worker version
      --webhook string                    webhook endpoint to notify task completion
      --webhookTimeoutSeconds int         duration webhook timeout (default 10)
  -w, --workDir string                    working directory; $HOME/.pneutrinoutil-worker/workspace or .pneutrinoutil-worker/workspace if no $HOME
```
