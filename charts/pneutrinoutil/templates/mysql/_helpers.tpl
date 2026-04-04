{{- define "pneutrinoutil.mysql.myCnf" -}}
[mysqld]
innodb_redo_log_capacity=512M
host-cache-size=0
skip-name-resolve
datadir=/var/lib/mysql
socket=/var/run/mysqld/mysqld.sock
secure-file-priv=/var/lib/mysql-files
user=mysql
character-set-server=utf8mb4
collation-server=utf8mb4_unicode_ci
default-time-zone=Asia/Tokyo
pid-file=/var/run/mysqld/mysqld.pid
[client]
socket=/var/run/mysqld/mysqld.sock
{{- end}}

{{- define "pneutrinoutil.mysql.dbSQL" -}}
CREATE DATABASE IF NOT EXISTS pneutrinoutil;

DROP DATABASE IF EXISTS test;
CREATE DATABASE IF NOT EXISTS test;
{{- end }}

{{- define "pneutrinoutil.mysql.usersSQL" -}}
DROP USER IF EXISTS pneutrinoutil, test;

CREATE USER IF NOT EXISTS pneutrinoutil IDENTIFIED BY 'userpass';
GRANT SELECT, INSERT, UPDATE, DELETE ON pneutrinoutil.* TO `pneutrinoutil`@`%`;

CREATE USER IF NOT EXISTS test IDENTIFIED BY 'test';
GRANT ALL PRIVILEGES ON test.* TO `test`@`%`;
{{- end }}

{{- define "pneutrinoutil.mysql.tablesSQL" -}}
CREATE TABLE IF NOT EXISTS master_statuses (
  id INT PRIMARY KEY,
  name VARCHAR(255)
);

INSERT IGNORE INTO master_statuses (id, name) VALUES
(1, 'pending'),
(2, 'running'),
(3, 'succeed'),
(4, 'failed');

CREATE TABLE IF NOT EXISTS master_object_types (
  id INT PRIMARY KEY,
  name VARCHAR(255)
);

INSERT IGNORE INTO master_object_types (id, name) VALUES
(1, 'file'),
(2, 'dir');

CREATE TABLE IF NOT EXISTS objects (
  id INT AUTO_INCREMENT PRIMARY KEY,
  type_id INT NOT NULL,
  bucket VARCHAR(4096) NOT NULL,
  path VARCHAR(4096) NOT NULL,
  size_bytes BIGINT UNSIGNED,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  -- bucket_path_sha256 VARCHAR(64) GENERATED ALWAYS AS (SHA2(CONCAT(bucket, '___', path), 256)) STORED,
  -- UNIQUE INDEX bucket_path_sha256_idx (bucket_path_sha256),
  INDEX bucket_prefix_idx (bucket(64)),
  INDEX path_prefix_idx (path(256)),
  CONSTRAINT fk_type_id FOREIGN KEY (type_id) REFERENCES master_object_types(id)
);

CREATE TABLE IF NOT EXISTS process_details (
  id INT AUTO_INCREMENT PRIMARY KEY,
  command VARCHAR(4096),
  title VARCHAR(4096) NOT NULL,
  score_object_id INT NOT NULL,
  log_object_id INT,
  result_object_id INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX title_prefix_idx (title(64)),
  CONSTRAINT fk_score_object_ids FOREIGN KEY (score_object_id) REFERENCES objects(id),
  CONSTRAINT fk_log_object_ids FOREIGN KEY (log_object_id) REFERENCES objects(id),
  CONSTRAINT fk_result_object_ids FOREIGN KEY (result_object_id) REFERENCES objects(id)
);

CREATE TABLE IF NOT EXISTS processes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  request_id VARCHAR(255),
  status_id INT NOT NULL,
  details_id INT NOT NULL,
  started_at TIMESTAMP,
  completed_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX status_id_idx (status_id),
  INDEX created_at_idx (created_at),
  UNIQUE INDEX request_id_idx (request_id),
  UNIQUE INDEX details_id_idx (details_id),
  CONSTRAINT fk_status_id FOREIGN KEY (status_id) REFERENCES master_statuses(id),
  CONSTRAINT fk_details_id FOREIGN KEY (details_id) REFERENCES process_details(id)
);
{{- end }}

{{- define "pneutrinoutil.mysql.setupSh" -}}
#!/bin/bash

readonly mysql_host="$MYSQL_HOST"
readonly mysql_user="$MYSQL_USER"
readonly mysql_pass="$MYSQL_PASSWORD"
readonly mysql_db="$MYSQL_DATABASE"
readonly root="$ROOT_DIR"

echo "USER=${mysql_user}"
echo "DATABASE=${mysql_db}"

client() {
  mysql -h "$mysql_host" -u"$mysql_user" -p"$mysql_pass" "$@"
}

run() {
  local -r sql="$1"
  echo >&2 "Run $*"
  shift
  client "$@" < "$sql"
}

set -e
set -o pipefail
run "${root}/db.sql"
run "${root}/users.sql"
run "${root}/tables.sql" "$mysql_db"

client <<< "SHOW DATABASES;"
client -D "$mysql_db" <<< "SHOW TABLES;"
{{- end }}

{{- define "pneutrinoutil.mysql.config" -}}
rootPassword: {{ .Values.mysql.rootPassword }}
my.cnf: |
  {{- include "pneutrinoutil.mysql.myCnf" . | nindent 4 }}
{{- end }}

{{- define "pneutrinoutil.mysql.configSetup" -}}
db.sql: |
  {{- include "pneutrinoutil.mysql.dbSQL" . | nindent 4 }}
users.sql: |
  {{- include "pneutrinoutil.mysql.usersSQL" . | nindent 4 }}
tables.sql: |
  {{- include "pneutrinoutil.mysql.tablesSQL" . | nindent 4 }}
setup.sh: |
  {{- include "pneutrinoutil.mysql.setupSh" . | nindent 4 }}
{{- end }}
