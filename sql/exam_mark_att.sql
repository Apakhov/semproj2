DROP TRIGGER IF EXISTS check_date_bounds_trigger ON mark;
DROP TRIGGER IF EXISTS check_date_bounds_trigger ON attendance;
DROP TRIGGER IF EXISTS check_date_bounds_trigger ON exam;
DROP FUNCTION IF EXISTS check_date_bounds;
DROP TABLE IF EXISTS mark;
DROP TABLE IF EXISTS attendance;
DROP TABLE IF EXISTS exam;
DROP TYPE IF EXISTS exam_type;

CREATE TABLE mark
(
    id         serial PRIMARY KEY NOT NULL,
    date       date               NOT NULL DEFAULT   CURRENT_TIMESTAMP,
    points     integer            NOT NULL,
    student_id integer            NOT NULL REFERENCES student (id),
    teacher_id integer            NOT NULL REFERENCES teacher (id),
    course_id  integer            NOT NULL REFERENCES course (id)
);

CREATE TABLE attendance
(
    id         serial PRIMARY KEY NOT NULL,
    date       date               NOT NULL DEFAULT   CURRENT_TIMESTAMP,
    student_id integer            NOT NULL REFERENCES student (id),
    teacher_id integer            NOT NULL REFERENCES teacher (id),
    course_id  integer            NOT NULL REFERENCES course (id)
);

CREATE TYPE exam_type AS ENUM ('p', 'np');
CREATE TABLE exam
(
    id         serial PRIMARY KEY NOT NULL,
    date       date               NOT NULL DEFAULT   CURRENT_TIMESTAMP,
    type       exam_type          NOT NULL,
    points     integer            NOT NULL
        CONSTRAINT proper_points check ( points >= 0 AND (type = 'p' OR points = 0 OR points = 1) ),
    student_id integer            NOT NULL REFERENCES student (id),
    teacher_id integer            NOT NULL REFERENCES teacher (id),
    course_id  integer            NOT NULL REFERENCES course (id),
    unique (student_id, course_id)
);

CREATE FUNCTION check_date_bounds() RETURNS TRIGGER AS
$$
DECLARE
    begin_d date;
    end_d date;
BEGIN
    SELECT INTO begin_d c.begin_d FROM course c WHERE id = new.course_id;
    SELECT INTO end_d c.end_d FROM course c WHERE id = new.course_id;
    IF begin_d is Null OR end_d is Null THEN
        return NEW;
    END IF;
    IF NEW.date < begin_d OR NEW.date > end_d THEN
        RAISE EXCEPTION 'bad date';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_date_bounds_trigger
    BEFORE INSERT
    ON mark
    FOR EACH ROW
EXECUTE PROCEDURE check_date_bounds();
CREATE TRIGGER check_date_bounds_trigger
    BEFORE INSERT
    ON attendance
    FOR EACH ROW
EXECUTE PROCEDURE check_date_bounds();
CREATE TRIGGER check_date_bounds_trigger
    BEFORE INSERT
    ON exam
    FOR EACH ROW
EXECUTE PROCEDURE check_date_bounds();
