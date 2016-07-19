CREATE DATABASE sample_router character set utf8;
GRANT ALL ON sample_router.* TO 'gtest'@'172.30.23.39' IDENTIFIED BY 'mypwd123';
FLUSH PRIVILEGES;
USE sample_router;
CREATE TABLE rules (
 id INT auto_increment,
 rule_id VARCHAR(36) NOT NULL,
 deleted INT NOT NULL,
 host VARCHAR(128) NOT NULL,
 band VARCHAR(64) NOT NULL,
 expire INT NOT NULL,
 product VARCHAR(64) NOT NULL,
 action_type VARCHAR(64) NOT NULL,
 action_value VARCHAR(64) NOT NULL,
 result INT not NULL,
 type VARCHAR(32) not NULL,
 primary key (id)
);

CREATE TABLE servers (
  id INT AUTO_INCREMENT,
  addr VARCHAR(32) DEFAULT NULL,
  state INT not NULL,
  product VARCHAR(64) DEFAULT NULL,
  PRIMARY KEY (id)
);
