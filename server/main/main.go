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

	log.Println("Cleaning up students table")
	_, err = connection.GetSession().DeleteFrom("students").Exec()
	if err != nil {
		log.Fatalf("Could not delete from students table. Err: %+v", err)
	}

	log.Println("Cleaning up departments_staffs table")
	_, err = connection.GetSession().DeleteFrom("departments_staffs").Exec()
	if err != nil {
		log.Fatalf("Could not delete from departments_staffs table. Err: %+v", err)
	}

	log.Println("Cleaning up departments table")
	_, err = connection.GetSession().DeleteFrom("departments").Exec()
	if err != nil {
		log.Fatalf("Could not delete from departments table. Err: %+v", err)
	}

	log.Println("Cleaning up staffs table")
	_, err = connection.GetSession().DeleteFrom("staffs").Exec()
	if err != nil {
		log.Fatalf("Could not delete from staffs table. Err: %+v", err)
	}

	log.Println("Inserting into departments table")
	_, err = connection.GetSession().InsertInto("departments").Columns("name").Values("Computer Science").Exec()
	_, err = connection.GetSession().InsertInto("departments").Columns("name").Values("DLC").Exec()
	_, err = connection.GetSession().InsertInto("departments").Columns("name").Values("LIC").Exec()
	if err != nil {
		log.Fatalf("Could not insert into departments table. Err: %+v", err)
	}

	log.Println("Inserting into students table")
	_, err = connection.GetSession().InsertInto("students").Columns("name", "roll_no", "department_id").Values("Vaishnavi", 34567, 1).Exec()
	_, err = connection.GetSession().InsertInto("students").Columns("name", "roll_no", "department_id").Values("Lulu", 12345, 2).Exec()
	_, err = connection.GetSession().InsertInto("students").Columns("name", "roll_no", "department_id").Values("ABC", 67890, 3).Exec()
	if err != nil {
		log.Fatalf("Could not insert into students table. Err: %+v", err)
	}

	log.Println("Inserting into staffs table")
	_, err = connection.GetSession().InsertInto("staffs").Columns("name").Values("Rajesh").Exec()
	_, err = connection.GetSession().InsertInto("staffs").Columns("name").Values("Raj").Exec()
	_, err = connection.GetSession().InsertInto("staffs").Columns("name").Values("Rakesh").Exec()
	if err != nil {
		log.Fatalf("Could not insert into staffs table. Err: %+v", err)
	}

	log.Println("Inserting into departments_staffs table")

	_, err = connection.GetSession().InsertInto("departments_staffs").Columns("department_id", "staff_id").Values(1, 1).Exec()
	_, err = connection.GetSession().InsertInto("departments_staffs").Columns("department_id", "staff_id").Values(1, 2).Exec()
	_, err = connection.GetSession().InsertInto("departments_staffs").Columns("department_id", "staff_id").Values(3, 3).Exec()
	if err != nil {
		log.Fatalf("Could not insert into departments_staffs table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}
