CREATE DATABASE IF NOT EXISTS task;
USE task;
DROP TABLE IF EXISTS task;
CREATE TABLE task (
  ID         INT AUTO_INCREMENT NOT NULL,
  name      VARCHAR(128) NOT NULL,
  completed     VARCHAR(255) NOT NULL,
  PRIMARY KEY (`ID`)
);

INSERT INTO task
  (name, completed)
VALUES
    ('Luana Dantas', 'yes'),
    ('Marina Silva', 'no'),
    ('Sofia Soares', 'no'),
    ('Lais Loreto', 'yes');
 