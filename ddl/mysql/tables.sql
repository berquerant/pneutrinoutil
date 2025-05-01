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

  INDEX (status_id),
  UNIQUE INDEX request_id_idx (request_id),
  CONSTRAINT fk_status_id FOREIGN KEY (status_id) REFERENCES master_statuses(id),
  CONSTRAINT fk_details_id FOREIGN KEY (details_id) REFERENCES process_details(id)
);
