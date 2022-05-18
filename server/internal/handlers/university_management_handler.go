package handlers

import (
	"context"
	"fmt"
	"log"
	"time"
	"university-management-golang/db/connection"
	um "university-management-golang/protoclient/university_management"
)

type universityManagementServer struct {
	um.UniversityManagementServiceServer
	connectionManager connection.DatabaseConnectionManager
}

func (u *universityManagementServer) GetDepartment(ctx context.Context, request *um.GetDepartmentRequest) (*um.GetDepartmentResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	// defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var department um.Department
	connection.GetSession().Select("id", "name").From("departments").Where("id = ?", request.GetId()).LoadOne(&department)

	// _, err = json.Marshal(&department)
	// fmt.Println("response", &department)
	// if err != nil {
	// 	log.Fatalf("Error while marshaling %+v", err)
	// }

	return &um.GetDepartmentResponse{Department: &department}, nil
}

func (u *universityManagementServer) GetStudent(ctx context.Context, request *um.GetStudentRequest) (*um.GetStudentResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	// defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var student um.Student
	fmt.Println("stud dpmnt id", request.GetDepartmentId())
	connection.GetSession().Select("id", "name", "roll_no").From("students").Where("department_id = ?", request.GetDepartmentId()).LoadOne(&student)
	fmt.Println("debug", &student)
	// _, err = json.Marshal(&student)
	// if err != nil {
	// 	log.Fatalf("Error while marshaling %+v", err)
	// }

	return &um.GetStudentResponse{Student: &student}, nil
}

func (u *universityManagementServer) GetStudents(ctx context.Context, request *um.GetStudentsRequest) (*um.GetStudentsResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	// defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var tempStudents []um.Student
	fmt.Println("stud dpmnt id", request.GetDepartmentId())
	connection.GetSession().Select("id", "name", "roll_no").From("students").Where("department_id = ?", request.GetDepartmentId()).Load(&tempStudents)
	fmt.Println("debug", &tempStudents)
	// _, err = json.Marshal(&tempStudents)
	// if err != nil {
	// 	log.Fatalf("Error while marshaling %+v", err)
	// }

	var students []*um.Student
	for _, student := range tempStudents {
		newStudent := um.Student{
			Id:           student.Id,
			Name:         student.Name,
			RollNo:       student.RollNo,
			DepartmentId: student.DepartmentId,
		}
		students = append(students, &newStudent)
	}
	return &um.GetStudentsResponse{Student: students}, nil
}

func (u *universityManagementServer) GetStaffsForStudent(ctx context.Context, request *um.GetStaffsForStudentRequest) (*um.GetStaffsForStudentResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	// defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var tempStaffs []um.Staff
	connection.GetSession().Select("staffs.name").From("students").Join("departments", "students.department_id = departments.id").Join("departments_staffs", "departments.id = departments_staffs.department_id").Join("staffs", "departments_staffs.staff_id = staffs.id").Where("students.roll_no = ?", request.GetRollNo()).Load(&tempStaffs)
	// _, err = json.Marshal(&tempStaffs)
	// if err != nil {
	// 	log.Fatalf("Error while marshaling %+v", err)
	// }

	var staffs []*um.Staff
	for _, staff := range tempStaffs {
		newStaff := um.Staff{
			Id:   staff.Id,
			Name: staff.Name,
		}
		staffs = append(staffs, &newStaff)
	}
	return &um.GetStaffsForStudentResponse{Staff: staffs}, nil
}

func (u *universityManagementServer) GetLoginForStudent(ctx context.Context, request *um.GetRequestForLogin) (*um.GetResponseForLogin, error) {
	connection, err := u.connectionManager.GetConnection()
	// defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var loginTime string
	rows, err := connection.GetSession().InsertInto("attendances").Columns("student_id", "login_time", "logout_time").Values(request.GetStudentId(), time.Now().UTC(), nil).Exec()
	if err != nil {
		log.Fatalf("Error while inserting into attendances table %+v", err)
	}
	log.Println(rows, "login record created!")
	connection.GetSession().Select("min(login_time)").From("attendances").Where("login_time >= current_date AND login_time < current_date + interval '1 day' AND student_id = ?", request.GetStudentId()).LoadOne(&loginTime)

	// fmt.Println(&loginTime)
	return &um.GetResponseForLogin{LoginTime: loginTime}, nil
}

func (u *universityManagementServer) GetLogoutForStudent(ctx context.Context, request *um.GetRequestForLogout) (*um.GetResponseForLogout, error) {
	connection, err := u.connectionManager.GetConnection()
	// defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var logoutTime string
	rows, err := connection.GetSession().Update("attendances").Set("logout_time", time.Now().UTC()).Where("id = ? AND student_id = ?", request.GetAttendance().Id, request.GetAttendance().StudentId).Exec()
	if err != nil {
		log.Fatalf("Error while inserting into attendances table for logout %+v", err)
	}
	log.Println(rows, "logout record created!")
	connection.GetSession().Select("max(logout_time)").From("attendances").Where("logout_time >= current_date AND logout_time < current_date + interval '1 day' AND student_id = ?", request.GetAttendance().StudentId).LoadOne(&logoutTime)

	// fmt.Println(&logoutTime)
	return &um.GetResponseForLogout{LogoutTime: logoutTime}, nil
}

func NewUniversityManagementHandler(connectionmanager connection.DatabaseConnectionManager) um.UniversityManagementServiceServer {
	return &universityManagementServer{
		connectionManager: connectionmanager,
	}
}
