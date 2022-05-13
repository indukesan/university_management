package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

	_, err = json.Marshal(&department)
	fmt.Println("response", &department)
	if err != nil {
		log.Fatalf("Error while marshaling %+v", err)
	}

	return &um.GetDepartmentResponse{Department: &um.Department{
		Id:   department.Id,
		Name: department.Name,
	}}, nil
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
	_, err = json.Marshal(&student)
	if err != nil {
		log.Fatalf("Error while marshaling %+v", err)
	}

	return &um.GetStudentResponse{Student: &um.Student{
		Id:     student.Id,
		Name:   student.Name,
		RollNo: student.RollNo,
	}}, nil
}

func (u *universityManagementServer) GetStudents(ctx context.Context, request *um.GetStudentsRequest) (*um.GetStudentsResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var tempStudents []um.Student
	fmt.Println("stud dpmnt id", request.GetDepartmentId())
	connection.GetSession().Select("id", "name", "roll_no").From("students").Where("department_id = ?", request.GetDepartmentId()).Load(&tempStudents)
	fmt.Println("debug", &tempStudents)
	_, err = json.Marshal(&tempStudents)
	if err != nil {
		log.Fatalf("Error while marshaling %+v", err)
	}

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

func NewUniversityManagementHandler(connectionmanager connection.DatabaseConnectionManager) um.UniversityManagementServiceServer {
	return &universityManagementServer{
		connectionManager: connectionmanager,
	}
}
