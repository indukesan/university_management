CREATE TABLE IF NOT EXISTS students (
    id serial PRIMARY KEY,
    name text,
    roll_no bigint,
    department_id integer REFERENCES departments(id)
);