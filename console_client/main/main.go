package main

import (
	"context"
	"fmt"
	"log"
	"university-management-golang/protoclient/university_management"

	"google.golang.org/grpc"
)

const (
	host = "localhost"
	port = "2345"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error occured %+v", err)
	}
	client := university_management.NewUniversityManagementServiceClient(conn)
	var departmentID int32 = 1
	departmentResponse, err := client.GetDepartment(context.TODO(), &university_management.GetDepartmentRequest{Id: departmentID})
	if err != nil {
		log.Fatalf("Error occured while fetching department for id %d,err: %+v", departmentID, err)
	}
	log.Println(departmentResponse)

	// var dpmntID int32 = 1
	studentResponse, err := client.GetStudent(context.TODO(), &university_management.GetStudentRequest{DepartmentId: departmentID})
	if err != nil {
		log.Fatalf("Error occured while fetching student for department id %d,err: %+v", departmentID, err)
	}
	log.Println("main stud", studentResponse)

	studentsResponse, err := client.GetStudents(context.TODO(), &university_management.GetStudentsRequest{DepartmentId: departmentID})
	if err != nil {
		log.Fatalf("Error occured while fetching student for department id %d,err: %+v", departmentID, err)
	}
	log.Println("main stud", studentsResponse)
}
