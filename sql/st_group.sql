DROP TRIGGER IF EXISTS id_update_trigger ON st_group;
DROP TRIGGER IF EXISTS delete_tg ON st_group;
DROP TABLE IF EXISTS st_group;

CREATE TABLE st_group
(
    id         serial PRIMARY KEY,
    group_id   varchar(50) unique NOT NULL, -- название группы
    begin_d    date               NOT NULL,
    end_d      date               NOT NULL
        CONSTRAINT proper_bounds check ( begin_d < end_d ),
    teacher_id integer            NOT NULL REFERENCES teacher (id),
    active     BOOLEAN            NOT NULL default TRUE
);
-- create ок
-- update не id
-- delete active

CREATE TRIGGER id_update_trigger
    BEFORE UPDATE
    ON st_group
    FOR EACH ROW
EXECUTE PROCEDURE check_update_id();
CREATE TRIGGER delete_tg
    BEFORE DELETE
    ON st_group
    FOR EACH ROW
EXECUTE PROCEDURE change_active();


INSERT INTO st_group (group_id, begin_d, end_d, teacher_id)
VALUES ('QW-12', '1999-2-1', '2003-2-9', 3)
SELECT id, group_id, begin_d, end_d, teacher_id, active
FROM st_group