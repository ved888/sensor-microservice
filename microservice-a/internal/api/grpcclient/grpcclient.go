package grpcclient

import (
	"context"
	"log"
	"math/rand"
	pb "microservice-a/pb/shared-proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Generator struct {
	addr   string              // gRPC server address
	freqCh chan time.Duration  // channel to dynamically update frequency
	dataCh chan *pb.SensorData // internal channel to buffer data
	stop   chan struct{}       // to stop the generator
	client pb.SensorServiceClient
}

// NewGenerator creates a new generator
func NewGenerator(addr string, freq time.Duration) *Generator {
	return &Generator{
		addr:   addr,
		freqCh: make(chan time.Duration, 1),
		dataCh: make(chan *pb.SensorData, 100),
		stop:   make(chan struct{}),
	}
}

// Start the generator: sends data to gRPC server and handles reconnections
func (g *Generator) Start(sensorType, id1, id2 string) {
	go g.generateDataLoop(sensorType, id1, id2) // continuously generate data
	go g.sendDataLoop()                         // continuously send data via gRPC
}

// generateDataLoop produces sensor data at the current frequency
func (g *Generator) generateDataLoop(sensorType, id1, id2 string) {
	freq := 1 * time.Second
	ticker := time.NewTicker(freq)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			data := &pb.SensorData{
				Value:      rand.Float64() * 100,
				SensorType: sensorType,
				Id1:        id1,
				Id2:        id2,
				Timestamp:  timestamppb.Now(),
			}
			select {
			case g.dataCh <- data:
			default:
				log.Println("dataCh full, dropping data")
			}
		case newFreq := <-g.freqCh:
			ticker.Stop()
			ticker = time.NewTicker(newFreq)
			log.Printf("Frequency updated to %v\n", newFreq)
		case <-g.stop:
			return
		}
	}
}

// sendDataLoop handles gRPC connection and reconnection
func (g *Generator) sendDataLoop() {
	for {
		select {
		case <-g.stop:
			return
		default:
		}

		conn, err := grpc.Dial(g.addr, grpc.WithInsecure())
		if err != nil {
			log.Println("Failed to connect gRPC, retrying in 1s:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		client := pb.NewSensorServiceClient(conn)
		stream, err := client.SendSensorData(context.Background())
		if err != nil {
			log.Println("Failed to start gRPC stream, retrying:", err)
			conn.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		// send buffered data
		for data := range g.dataCh {
			err := stream.Send(data)
			if err != nil {
				log.Println("gRPC send failed, reconnecting:", err)
				conn.Close()
				break // break inner loop for reconnect
			}
		}

		conn.Close()
	}
}

// UpdateFrequency dynamically updates data generation frequency
func (g *Generator) UpdateFrequency(freq time.Duration) {
	select {
	case g.freqCh <- freq:
	default:
		log.Println("freqCh full, skipping update")
	}
}

// Stop the generator
func (g *Generator) Stop() {
	close(g.stop)
}
