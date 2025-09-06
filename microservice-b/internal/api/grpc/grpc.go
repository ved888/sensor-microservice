package grpc

import (
	"io"
	"log"
	"microservice-b/internal/repository"
	"net"

	pb "microservice-b/pb/shared-proto"

	"google.golang.org/grpc"
)

type SensorServer struct {
	pb.UnimplementedSensorServiceServer
	Repo *repository.SensorRepository
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

		if err := s.Repo.Save(data); err != nil {
			log.Printf("DB error: %v", err)
		}
		log.Printf("Sent data: %v", data)
	}
}

func StartGRPCServer(repo *repository.SensorRepository, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSensorServiceServer(grpcServer, &SensorServer{Repo: repo})

	log.Printf("gRPC server running on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
