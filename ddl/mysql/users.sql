DROP USER IF EXISTS pneutrinoutil, test;

CREATE USER IF NOT EXISTS pneutrinoutil IDENTIFIED BY 'userpass';
GRANT SELECT, INSERT, UPDATE, DELETE ON pneutrinoutil.* TO `pneutrinoutil`@`%`;

CREATE USER IF NOT EXISTS test IDENTIFIED BY 'test';
GRANT ALL PRIVILEGES ON test.* TO `test`@`%`;
