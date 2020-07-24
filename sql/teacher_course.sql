DROP TRIGGER IF EXISTS update_trigger ON teacher_course;
DROP FUNCTION IF EXISTS ban_update_teacher_course;
DROP TABLE IF EXISTS teacher_course;

CREATE TABLE teacher_course
(
    id         serial PRIMARY KEY NOT NULL,
    teacher_id integer            NOT NULL REFERENCES teacher (id),
    course_id  integer            NOT NULL REFERENCES course (id),
    unique (teacher_id, course_id)
);

CREATE FUNCTION ban_update_teacher_course() RETURNS TRIGGER AS
$$
BEGIN
   RAISE EXCEPTION 'cannot change teacher_course fields';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_trigger
    BEFORE UPDATE
    ON teacher_course
    FOR EACH ROW
EXECUTE PROCEDURE ban_update_teacher_course();

SELECT * FROM teacher;
SELECT * FROM course;
SELECT * FROM teacher_course;
DELETE FROM teacher_course WHERE (teacher_id, course_id) = (3, 3)