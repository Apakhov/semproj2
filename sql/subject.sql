DROP TRIGGER IF EXISTS id_update_trigger ON subject;
DROP TRIGGER IF EXISTS delete_tg ON subject;
DROP TABLE IF EXISTS subject;

CREATE TABLE subject
(
    id          serial PRIMARY KEY NOT NULL,
    short_name  varchar(50) unique NOT NULL,
    full_name   varchar(50) unique NOT NULL,
    cathedra_id integer            NOT NULL REFERENCES cathedra (id),
    deleted     boolean            NOT NULL DEFAULT false
);
-- create ок
-- update не id
-- delete deleted

CREATE TRIGGER id_update_trigger BEFORE UPDATE ON subject FOR EACH ROW
      EXECUTE PROCEDURE check_update_id();
CREATE TRIGGER delete_tg BEFORE DELETE ON subject FOR EACH ROW
      EXECUTE PROCEDURE change_deleted();