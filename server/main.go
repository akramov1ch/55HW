package main

import (
  "log"
  "net"
  "google.golang.org/grpc"
  pb "55HW/proto"
)

func main() {
  lis, err := net.Listen("tcp", ":9001")
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  s := grpc.NewServer()
  pb.RegisterTaskServiceServer(s, &server{})
  log.Println("Server is running on port :9001")
  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}
