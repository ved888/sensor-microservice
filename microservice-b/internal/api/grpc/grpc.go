package grpc

import (
	"context"
	"io"
	"log"
	"microservice-b/internal/repository"
	"net"
	"sync"

	pb "microservice-b/pb/shared-proto"

	"google.golang.org/grpc"
)

type SensorServer struct {
	pb.UnimplementedSensorServiceServer
	Repo   *repository.SensorRepository
	dataCh chan *pb.SensorData
	wg     sync.WaitGroup
}

// Init server with worker pool
func NewSensorServer(repo *repository.SensorRepository, workerCount int) *SensorServer {
	server := &SensorServer{
		Repo:   repo,
		dataCh: make(chan *pb.SensorData, 1000),
	}

	// start worker pool
	for i := 0; i < workerCount; i++ {
		server.wg.Add(1)
		go func() {
			defer server.wg.Done()
			for data := range server.dataCh {
				if err := repo.Save(data); err != nil {
					log.Printf("DB error: %v", err)
				}
			}
		}()
	}

	return server
}

func (s *SensorServer) SendSensorData(stream pb.SensorService_SendSensorDataServer) error {
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Ack{Ok: true, Message: "All data received"})
		}
		if err != nil {
			log.Printf("Stream error: %v", err)
			return err
		}

		// Send data to worker pool
		select {
		case s.dataCh <- data:
		default:
			log.Println("dataCh full, dropping data")
		}
		log.Printf("Received data: %v", data)
	}
}

func StartGRPCServer(ctx context.Context, repo *repository.SensorRepository, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := NewSensorServer(repo, 10) //example 10 workers
	pb.RegisterSensorServiceServer(grpcServer, server)

	// Graceful shutdown on context cancellation
	go func() {
		<-ctx.Done()
		log.Println("Stopping gRPC server gracefully...")
		grpcServer.GracefulStop()
		server.Stop()
		log.Println("gRPC server stopped")
	}()

	log.Printf("gRPC server running on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func (s *SensorServer) Stop() {
	log.Println("Closing data channel...")
	close(s.dataCh)
	log.Println("Waiting for workers to finish...")
	s.wg.Wait()
	log.Println("All workers done. Server stopped.")
}
