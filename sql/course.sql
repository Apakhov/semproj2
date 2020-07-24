DROP TRIGGER IF EXISTS id_update_trigger ON course;
DROP TRIGGER IF EXISTS delete_tg ON course;
DROP TABLE IF EXISTS course;


CREATE TABLE course
(
    id          serial PRIMARY KEY NOT NULL,
    short_name  varchar(50),
    full_name   varchar(50)        NOT NULL,
    semester    integer            NOT NULL
        CONSTRAINT proper_semester check ( 1 <= semester AND semester <= 16 ),
    begin_d     date               NOT NULL,
    end_d       date               NOT NULL
        CONSTRAINT proper_bounds check ( begin_d < end_d ),
    subject_id  integer            NOT NULL REFERENCES subject (id),
    st_group_id integer            NOT NULL REFERENCES st_group (id),
    active      BOOLEAN            NOT NULL default TRUE,
    unique (semester, subject_id, st_group_id)
);

CREATE TRIGGER id_update_trigger
    BEFORE UPDATE
    ON course
    FOR EACH ROW
EXECUTE PROCEDURE check_update_id();
CREATE TRIGGER delete_tg
    BEFORE DELETE
    ON course
    FOR EACH ROW
EXECUTE PROCEDURE change_active();

SELECT * from st_group;
SELECT * from subject;
INSERT INTO course (short_name, full_name, semester, begin_d, end_d, subject_id, st_group_id)
VALUES ('asd', 'ads', 1, '1212-12-12', '1213-12-12', 3, 2);
SELECT * FROM course;