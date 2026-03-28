CREATE TABLE IF NOT EXISTS Groups(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Subjects(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Teachers(
    id SERIAL PRIMARY KEY,
    fullname TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Classrooms(
    id SERIAL PRIMARY KEY,
    num TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Schedule(
    id SERIAL PRIMARY KEY,

    group_id INT NOT NULL REFERENCES Groups(id),
    subject_id INT NOT NULL REFERENCES Subjects(id),
    teacher_id INT NOT NULL REFERENCES Teachers(id),
    classroom_id INT NOT NULL REFERENCES Classrooms(id),

    weekday INT CHECK (weekday >= 1 AND weekday <= 7),
    lesson_number INT CHECK (lesson_number >= 1 AND lesson_number <= 10),
    week_type INT CHECK (week_type >= 1 AND week_type <= 2), -- 1 - для первой недели, 2 - для второй недели, NULL - для всех недель

    subgroup INT CHECK (subgroup IN(1, 2)), -- 1 - для первой подгруппы, 2 - для второй подгруппы, NULL - для всех студентов

    UNIQUE(group_id, weekday, lesson_number, week_type, subgroup)
);