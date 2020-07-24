DROP TRIGGER IF EXISTS id_update_trigger ON student;
DROP TRIGGER IF EXISTS delete_tg ON student;
DROP TABLE IF EXISTS student;


CREATE TABLE student
(
    id         serial PRIMARY KEY,
    study_id   varchar(50) unique NOT NULL, -- зачетка
    first_name VARCHAR(50)         NOT NULL,
    last_name  VARCHAR(50)         NOT NULL,
    email      VARCHAR(355) UNIQUE NOT NULL
        CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    active     BOOLEAN NOT NULL default TRUE
);
-- create ок
-- update не id
-- delete active


CREATE TRIGGER id_update_trigger BEFORE UPDATE ON student FOR EACH ROW
      EXECUTE PROCEDURE check_update_id();
CREATE TRIGGER delete_tg BEFORE DELETE ON student FOR EACH ROW
      EXECUTE PROCEDURE change_active();

INSERT INTO student (study_id, first_name, last_name, email) VALUES ('asd', 'm', 'a', 'asd@m.a')