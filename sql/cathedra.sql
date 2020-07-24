DROP TRIGGER IF EXISTS delete_tg ON cathedra;
DROP TRIGGER IF EXISTS id_update_trigger ON cathedra;
DROP TABLE IF EXISTS cathedra;

CREATE TABLE cathedra
(
    id         serial PRIMARY KEY NOT NULL,
    short_name varchar(50) unique NOT NULL,
    full_name  varchar(50) unique NOT NULL,
    faculty_id integer            NOT NULL REFERENCES faculty (id),
    deleted    boolean            NOT NULL DEFAULT false
);
-- create ок
-- update не id
-- delete deleted

CREATE TRIGGER id_update_trigger BEFORE UPDATE ON cathedra FOR EACH ROW
      EXECUTE PROCEDURE check_update_id();
CREATE TRIGGER delete_tg BEFORE DELETE ON cathedra FOR EACH ROW
      EXECUTE PROCEDURE change_deleted();

SELECT * FROM cathedra;