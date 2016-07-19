CREATE DATABASE nginxsecurity character set utf8;
GRANT ALL ON nginxsecurity.* TO 'gtest'@'172.30.23.39' IDENTIFIED BY 'mypwd123';
FLUSH PRIVILEGES;
USE nginxsecurity;
CREATE TABLE rules (
 id INT auto_increment,
 rule_id VARCHAR(36) NOT NULL,
 deleted INT NOT NULL,
 ip VARCHAR(36) DEFAULT NULL,
 uid VARCHAR(36) DEFAULT NULL,
 uuid VARCHAR(36) DEFAULT NULL,
 expire INT NOT NULL,
 product VARCHAR(64) NOT NULL,
 action VARCHAR(64) NOT NULL,
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
