DROP TRIGGER IF EXISTS id_update_trigger ON faculty;
DROP TRIGGER IF EXISTS delete_tg ON faculty;
DROP TABLE IF EXISTS faculty;

CREATE TABLE faculty
(
    id         serial PRIMARY KEY NOT NULL,
    short_name varchar(50) unique NOT NULL,
    full_name  varchar(50) unique NOT NULL,
    deleted    boolean            NOT NULL DEFAULT false
);

CREATE TRIGGER id_update_trigger
    BEFORE UPDATE
    ON faculty
    FOR EACH ROW
EXECUTE PROCEDURE check_update_id();
CREATE TRIGGER delete_tg
    BEFORE DELETE
    ON faculty
    FOR EACH ROW
EXECUTE PROCEDURE change_deleted();

INSERT INTO faculty (short_name, full_name)
VALUES ('kkaa', 'ddaa'),
       ('asda', 'sda')
RETURNING id, deleted;
SELECT *
FROM faculty AS f;

SELECT (id, short_name, full_name, deleted)
FROM faculty
WHERE id >= 2
  AND short_name = 'dsass'
LIMIT 2;

UPDATE faculty
SET short_name = 'kakakakkaak',
    full_name  = 'kakakakkaak'
WHERE id = 2;

DELETE FROM faculty where id = 11;

SELECT id, short_name, full_name, deleted
FROM faculty
WHERE id >= -1
ORDER BY id LIMIT 2;