server: go run ./server --port $PORT --mysqlDSN $MYSQL_DSN --redisDSN $REDIS_DSN --storageBucket $STORAGE_BUCKET --storageS3 --debug
worker: go run ./worker --mysqlDSN $MYSQL_DSN --redisDSN $REDIS_DSN --shell /bin/bash --storageBucket $STORAGE_BUCKET --storageS3 --debug
