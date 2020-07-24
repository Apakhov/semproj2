DROP TRIGGER IF EXISTS id_update_trigger ON teacher;
DROP TRIGGER IF EXISTS delete_tg ON teacher;
DROP TABLE IF EXISTS teacher;



CREATE TABLE teacher
(
    id          serial PRIMARY KEY,
    contract_id varchar(50) unique  NOT NULL,
    first_name  VARCHAR(50)         NOT NULL,
    last_name   VARCHAR(50)         NOT NULL,
    email       VARCHAR(355) UNIQUE NOT NULL
        CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    cathedra_id integer             NOT NULL REFERENCES cathedra (id),
    active      BOOLEAN NOT NULL default TRUE
);
-- create ок
-- update не id
-- delete active

CREATE TRIGGER id_update_trigger BEFORE UPDATE ON teacher FOR EACH ROW
      EXECUTE PROCEDURE check_update_id();
CREATE TRIGGER delete_tg BEFORE DELETE ON teacher FOR EACH ROW
      EXECUTE PROCEDURE change_active();

INSERT INTO teacher (contract_id, first_name, last_name, email, cathedra_id) VALUES ('112', 'asd', 'ads', 'das', 2)
SELECT id, contract_id, first_name, last_name, email, cathedra_id, active FROM teacher WHERE active = false;