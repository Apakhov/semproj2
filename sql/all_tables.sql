CREATE TABLE faculty
(
    id         serial PRIMARY KEY NOT NULL,
    short_name varchar(50) unique NOT NULL,
    full_name  varchar(50) unique NOT NULL,
    deleted    boolean            NOT NULL DEFAULT false
);

CREATE TABLE cathedra
(
    id         serial PRIMARY KEY NOT NULL,
    short_name varchar(50) unique NOT NULL,
    full_name  varchar(50) unique NOT NULL,
    faculty_id integer            NOT NULL REFERENCES faculty (id),
    deleted    boolean            NOT NULL DEFAULT false
);

CREATE TABLE subject
(
    id          serial PRIMARY KEY NOT NULL,
    short_name  varchar(50) unique NOT NULL,
    full_name   varchar(50) unique NOT NULL,
    cathedra_id integer            NOT NULL REFERENCES cathedra (id),
    deleted     boolean            NOT NULL DEFAULT false
);

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

CREATE TABLE st_group
(
    id         serial PRIMARY KEY,
    group_id   varchar(50) unique NOT NULL, -- название группы
    begin_d    date               NOT NULL,
    end_d      date               NOT NULL CONSTRAINT proper_bounds check ( begin_d < end_d ) ,
    teacher_id integer            NOT NULL REFERENCES teacher (id),
    active     BOOLEAN NOT NULL default TRUE
);

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

CREATE TABLE student_st_group
(
    id          serial PRIMARY KEY,
    active      BOOLEAN default TRUE,
    student_id  integer NOT NULL REFERENCES student (id),
    st_group_id integer NOT NULL REFERENCES st_group (id),
    unique (student_id, st_group_id)
);

CREATE TABLE course
(
    id          serial PRIMARY KEY NOT NULL,
    short_name  varchar(50),
    full_name   varchar(50)        NOT NULL,
    semester    integer            NOT NULL,
    begin_d     date               NOT NULL,
    end_d       date               NOT NULL,
    subject_id  integer            NOT NULL REFERENCES subject (id),
    st_group_id integer            NOT NULL REFERENCES st_group (id),
    unique (semester,subject_id, st_group_id)
);

CREATE TABLE teacher_course
(
    id         serial PRIMARY KEY NOT NULL,
    teacher_id integer            NOT NULL REFERENCES teacher (id),
    course_id  integer            NOT NULL REFERENCES course (id),
    unique (teacher_id, course_id)
);

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

CREATE TYPE exam_type AS ENUM ('point', 'non point');
CREATE TABLE exam
(
    id         serial PRIMARY KEY NOT NULL,
    date       date               NOT NULL DEFAULT   CURRENT_TIMESTAMP,
    type       exam_type          NOT NULL,
    points     integer            NULL,
    student_id integer            NOT NULL REFERENCES student (id),
    teacher_id integer            NOT NULL REFERENCES teacher (id),
    course_id  integer            NOT NULL REFERENCES course (id)
);