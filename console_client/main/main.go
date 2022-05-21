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
	var departmentID int32 = 7
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
	log.Println(studentResponse)

	studentsResponse, err := client.GetStudents(context.TODO(), &university_management.GetStudentsRequest{DepartmentId: departmentID})
	if err != nil {
		log.Fatalf("Error occured while fetching student for department id %d,err: %+v", departmentID, err)
	}
	log.Println(studentsResponse)

	var rollNo int32 = 34567
	staffsForStudentResponse, err := client.GetStaffsForStudent(context.TODO(), &university_management.GetStaffsForStudentRequest{RollNo: rollNo})
	if err != nil {
		log.Fatalf("Error occured while fetching staffs for student rollNo %d,err: %+v", rollNo, err)
	}
	log.Println(staffsForStudentResponse)

	var studentId int32 = 3
	loginResponse, err := client.GetLoginForStudent(context.TODO(), &university_management.GetRequestForLogin{StudentId: studentId})
	if err != nil {
		log.Fatalf("Error occured while fetching login for student %d,err: %+v", 5, err)
	}
	log.Println(loginResponse)

	attendance := university_management.Attendance{
		Id:        7,
		StudentId: studentId,
		LoginTime: loginResponse.GetLoginTime(),
		// LoginTime:  "2022-05-20T11:28:49.388646Z",
		LogoutTime: "",
	}
	logoutResponse, err := client.GetLogoutForStudent(context.TODO(), &university_management.GetRequestForLogout{Attendance: &attendance})
	if err != nil {
		log.Fatalf("Error occured while fetching logout for student %d,err: %+v", attendance.GetStudentId(), err)
	}
	log.Println(logoutResponse)

	notifyResponse, err := client.Notify(context.TODO(), &university_management.GetLoginNotifyRequest{StudentIds: []int32{}})
	if err != nil {
		log.Fatalf("Error occured while getting login notification for student,err: %+v", err)
	}
	log.Println(notifyResponse)

	attendanceResponse, err := client.GetAttendances(context.TODO(), &university_management.GetRequestForAttendance{})
	log.Println("new", attendanceResponse)
}
