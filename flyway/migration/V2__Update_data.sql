USE api;
-- ADD admin User
INSERT IGNORE INTO users (first_name, last_name, user_name, password_hash, role)
VALUES ('admin', 'admin', 'admin', '$2a$10$fuf.PUvm.I5ScIACGmWO8u57LgzlM/aB9FLQ8NNIXIwMpCeJJd8B2', 'ADMIN');

-- ADD school
INSERT IGNORE INTO schools (name)
VALUES
    ('Bouve College of Health Sciences'), -- 1
    ('School of Law'), -- 2
    ('Office of Provost'), -- 3
    ('College of Engineering'), -- 4
    ('College of Arts, Media, and Design'); -- 5

-- ADD Departments
INSERT IGNORE INTO departments (name, school_id) 
VALUES
    ('Physician Assistant', 1),
    ('Pharmacy Health Systems', 1),
    ('Health Informatics', 1),
    ('Applied Psychology', 1);

-- ADD Terms
INSERT IGNORE INTO terms (name)
VALUES
    ('Fall 2020'), -- 1
    ('Spring 2021'), -- 2
    ('Summer 2021'); -- 3

