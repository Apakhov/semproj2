DROP FUNCTION IF EXISTS check_update_student_st_group;
DROP FUNCTION IF EXISTS new_student_st_group;
DROP TRIGGER IF EXISTS id_update_trigger ON student_st_group;
DROP TRIGGER IF EXISTS delete_tg ON student_st_group;
DROP TRIGGER IF EXISTS update_trigger_student_st_group ON student_st_group;
DROP TRIGGER IF EXISTS insert_trigger_student_st_group ON student_st_group;
DROP TABLE IF EXISTS student_st_group;

CREATE TABLE student_st_group
(
    id          serial PRIMARY KEY,
    active      BOOLEAN default TRUE,
    student_id  integer NOT NULL REFERENCES student (id),
    st_group_id integer NOT NULL REFERENCES st_group (id),
    unique (student_id, st_group_id)
);

CREATE TRIGGER id_update_trigger
    BEFORE UPDATE
    ON student_st_group
    FOR EACH ROW
EXECUTE PROCEDURE check_update_id();
CREATE TRIGGER delete_tg
    BEFORE DELETE
    ON student_st_group
    FOR EACH ROW
EXECUTE PROCEDURE change_active();

CREATE FUNCTION check_update_student_st_group() RETURNS TRIGGER AS
$$
BEGIN
    IF OLD.student_id <> NEW.student_id THEN
        RAISE EXCEPTION 'cannot change student_id';
    END IF;
    IF OLD.st_group_id <> NEW.st_group_id THEN
        RAISE EXCEPTION 'cannot change student_id';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_trigger_student_st_group
    BEFORE UPDATE
    ON student_st_group
    FOR EACH ROW
EXECUTE PROCEDURE check_update_student_st_group();

CREATE FUNCTION new_student_st_group() RETURNS TRIGGER AS
$$
DECLARE
    old_id INTEGER;
BEGIN
    UPDATE student_st_group SET active= false WHERE student_id = NEW.student_id;
    IF OLD.st_group_id <> NEW.st_group_id THEN
        RAISE EXCEPTION 'cannot change student_id';
    END IF;
    SELECT INTO old_id id FROM student_st_group WHERE student_id = NEW.student_id AND st_group_id = NEW.st_group_id;
    IF old_id IS NOT NULL THEN
        UPDATE student_st_group SET active= true WHERE id = old_id;
        RETURN NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_trigger_student_st_group
    BEFORE INSERT
    ON student_st_group
    FOR EACH ROW
EXECUTE PROCEDURE new_student_st_group();

SELECT *
FROM st_group;
SELECT *
FROM student;

SELECT *
FROM student_st_group;

INSERT INTO student_st_group (student_id, st_group_id)
VALUES (1, 2);
SELECT s.id, s.study_id, s.first_name, s.last_name, s.email, s.active
FROM student s
         JOIN student_st_group sg on s.id = sg.student_id
WHERE sg.st_group_id = 2;

SELECT gr.id, gr.group_id, gr.begin_d, gr.end_d, gr.teacher_id, gr.active
FROM st_group gr
         JOIN student_st_group sg on gr.id = sg.st_group_id
WHERE sg.student_id = 1;