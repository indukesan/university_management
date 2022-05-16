CREATE table IF NOT EXISTS departments_staffs(
       id SERIAL PRIMARY KEY ,
       department_id integer REFERENCES departments(id),
       staff_id integer REFERENCES staffs(id)
);