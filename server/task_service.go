package main

import (
  "context"
  "sync"
  "time"
  "github.com/google/uuid"
  pb "55HW/proto"
)

type server struct {
  pb.UnimplementedTaskServiceServer
  tasks map[string]*pb.TaskResponse
  mu    sync.Mutex
}

func (s *server) CreateTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
  taskID := generateUniqueID()
  task := &pb.TaskResponse{
    TaskId: taskID,
    Status: "In Progress",
  }

  s.mu.Lock()
  if s.tasks == nil {
    s.tasks = make(map[string]*pb.TaskResponse)
  }
  s.tasks[taskID] = task
  s.mu.Unlock()

  go func() {
    select {
    case <-time.After(10 * time.Second):
      s.mu.Lock()
      if s.tasks[taskID].Status != "Cancelled" {
        s.tasks[taskID].Status = "Completed"
      }
      s.mu.Unlock()
    case <-ctx.Done():
      s.mu.Lock()
      s.tasks[taskID].Status = "Cancelled"
      s.mu.Unlock()
    }
  }()

  return task, nil
}

func (s *server) ListTasks(ctx context.Context, req *pb.Empty) (*pb.TaskList, error) {
  s.mu.Lock()
  defer s.mu.Unlock()
  tasks := []*pb.TaskResponse{}
  for _, task := range s.tasks {
    tasks = append(tasks, task)
  }
  return &pb.TaskList{Tasks: tasks}, nil
}

func (s *server) CancelTask(ctx context.Context, req *pb.CancelRequest) (*pb.CancelResponse, error) {
  s.mu.Lock()
  defer s.mu.Unlock()
  if task, exists := s.tasks[req.TaskId]; exists {
    task.Status = "Cancelled"
    return &pb.CancelResponse{Status: "Cancelled"}, nil
  }
  return &pb.CancelResponse{Status: "Task not found"}, nil
}

func generateUniqueID() string {
  return uuid.New().String()
}
