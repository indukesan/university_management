CREATE TABLE IF NOT EXISTS attendances (
	id SERIAL PRIMARY KEY,
    student_id INT,
	login_time TIMESTAMP,
	logout_time TIMESTAMP,
    CONSTRAINT fk_student FOREIGN KEY(student_id) REFERENCES students(id)
);