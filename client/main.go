package main

import (
	pb "55HW/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
  conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()
  c := pb.NewTaskServiceClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  res, err := c.CreateTask(ctx, &pb.TaskRequest{TaskDescription: "New Task"})
  if err != nil {
    log.Fatalf("could not create task: %v", err)
  }
  log.Printf("Task created: %v", res)

  tasks, err := c.ListTasks(context.Background(), &pb.Empty{})
  if err != nil {
    log.Fatalf("could not list tasks: %v", err)
  }
  log.Printf("Tasks: %v", tasks)

  cancelRes, err := c.CancelTask(context.Background(), &pb.CancelRequest{TaskId: res.TaskId})
  if err != nil {
    log.Fatalf("could not cancel task: %v", err)
  }
  log.Printf("Cancel task response: %v", cancelRes)
}
