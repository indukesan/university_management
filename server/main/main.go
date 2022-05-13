package main

import (
	"log"
	"net"
	migrations "university-management-golang/db"
	"university-management-golang/db/connection"
	um "university-management-golang/protoclient/university_management"
	"university-management-golang/server/internal/handlers"

	"google.golang.org/grpc"
)

const port = "2345"

type Department struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type Student struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	RollNo       string `db:"roll_no"`
	DepartmentId int64  `db:"department_id"`
}

//db
const (
	username = "postgres"
	password = "admin"
	host     = "localhost"
	dbPort   = "5436"
	dbName   = "postgres"
	schema   = "public"
)

func main() {
	err := migrations.MigrationsUp(username, password, host, dbPort, dbName, schema)
	if err != nil {
		log.Fatalf("Failed to migrate, err: %+v\n", err)
	}

	connectionmanager := &connection.DatabaseConnectionManagerImpl{
		&connection.DBConfig{
			host, dbPort, username, password, dbName, schema,
		},
		nil,
	}

	// insertSeedData(connectionmanager)

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen to port: %s, err: %+v\n", port, err)
	}
	log.Printf("Starting to listen on port: %s\n", port)

	um.RegisterUniversityManagementServiceServer(grpcServer, handlers.NewUniversityManagementHandler(connectionmanager))
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to start GRPC Server: %+v\n", err)
	}
}

func insertSeedData(connectionManager connection.DatabaseConnectionManager) {
	connection, err := connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Println("Cleaning up department table")
	_, err = connection.GetSession().DeleteFrom("departments").Exec()
	if err != nil {
		log.Fatalf("Could not delete from department table. Err: %+v", err)
	}

	log.Println("Cleaning up students table")
	_, err = connection.GetSession().DeleteFrom("students").Exec()
	if err != nil {
		log.Fatalf("Could not delete from students table. Err: %+v", err)
	}

	log.Println("Inserting into departments table")
	departmentData := &Department{
		ID:   2,
		Name: "Computer Science",
	}
	_, err = connection.GetSession().InsertInto("departments").Columns("id", "name").Record(departmentData).Exec()

	log.Println("Inserting into students table")
	studentData := &Student{
		Name:         "Mathi",
		RollNo:       "12345",
		DepartmentId: 2,
		ID:           2,
	}
	_, err = connection.GetSession().InsertInto("students").Columns("id", "name", "roll_no", "department_id").Record(studentData).Exec()

	if err != nil {
		log.Fatalf("Could not insert into departments table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}
