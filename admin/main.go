package main

import (
	"context"
	"fmt"
	"io"
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

	doServerStreaming(client)
}

func doServerStreaming(c university_management.UniversityManagementServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &university_management.GetRequestForStreamingAttendance{}

	resStream, err := c.StreamAttendanceResponse(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling StreamAttendanceResponse RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		att := msg.GetAttendance()
		for _, att := range att {
			fmt.Printf("Response from StreamAttendanceResponse: %v\n", att)
		}
		// fmt.Printf("Response from StreamAttendanceResponse: %v\n", msg.GetAttendance())
	}
}
