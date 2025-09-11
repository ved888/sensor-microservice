package grpcclient

import (
	pb "microservice-a/pb/shared-proto"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestNewGenerator(t *testing.T) {
	addr := "localhost:50051"
	freq := 1 * time.Second

	gen := NewGenerator(addr, freq)

	assert.NotNil(t, gen)
	assert.Equal(t, addr, gen.addr)
	assert.NotNil(t, gen.freqCh)
	assert.NotNil(t, gen.dataCh)
	assert.NotNil(t, gen.stop)
}

func TestGenerator_UpdateFrequency(t *testing.T) {
	gen := NewGenerator("localhost:50051", 1*time.Second)
	newFreq := 2 * time.Second

	gen.UpdateFrequency(newFreq)

	select {
	case receivedFreq := <-gen.freqCh:
		assert.Equal(t, newFreq, receivedFreq)
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected frequency update to be sent to channel")
	}
}

func TestGenerator_Stop(t *testing.T) {
	gen := NewGenerator("localhost:50051", 1*time.Second)

	gen.Stop()

	select {
	case _, ok := <-gen.stop:
		assert.False(t, ok, "Stop channel should be closed")
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected stop channel to be closed")
	}
}

// Server that notifies when it receives data
type FakeSensorServer struct {
	pb.UnimplementedSensorServiceServer
	Received []*pb.SensorData
	Done     chan struct{}
}

func (s *FakeSensorServer) SendSensorData(stream pb.SensorService_SendSensorDataServer) error {
	for {
		data, err := stream.Recv()
		if err != nil {
			return stream.SendAndClose(&pb.Ack{Ok: true, Message: "All data received"})
		}
		s.Received = append(s.Received, data)
		select {
		case s.Done <- struct{}{}:
		default:
		}
	}
}

func TestGenerator_Start_DataGeneration_WithRealServer(t *testing.T) {
	// Listen on free port
	lis, _ := net.Listen("tcp", "127.0.0.1:0")
	server := grpc.NewServer()
	fakeServer := &FakeSensorServer{Done: make(chan struct{}, 1)}
	pb.RegisterSensorServiceServer(server, fakeServer)
	go server.Serve(lis)
	defer server.Stop()

	addr := lis.Addr().String()
	gen := NewGenerator(addr, 50*time.Millisecond)
	gen.Start("Temperature", "A", "1")

	// Wait until server receives at least one message
	select {
	case <-fakeServer.Done:
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for server to receive data")
	}

	gen.Stop()

	assert.Greater(t, len(fakeServer.Received), 0)
	for _, d := range fakeServer.Received {
		assert.Equal(t, "Temperature", d.SensorType)
		assert.Equal(t, "A", d.Id1)
		assert.Equal(t, "1", d.Id2)
		assert.NotNil(t, d.Timestamp)
	}
}

func TestGenerator_UpdateFrequency_DuringGeneration_WithServer(t *testing.T) {
	// Start server on free port
	lis, _ := net.Listen("tcp", "127.0.0.1:0")
	server := grpc.NewServer()
	fakeServer := &FakeSensorServer{Done: make(chan struct{}, 1)}
	pb.RegisterSensorServiceServer(server, fakeServer)
	go server.Serve(lis)
	defer server.Stop()

	addr := lis.Addr().String()
	gen := NewGenerator(addr, 50*time.Millisecond)

	// Start generator
	gen.Start("Temperature", "A", "1")

	time.Sleep(50 * time.Millisecond)

	// Update frequency
	gen.UpdateFrequency(200 * time.Millisecond)

	// Wait until at least one message is received
	select {
	case <-fakeServer.Done:
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for server to receive data")
	}

	gen.Stop()

	assert.Greater(t, len(fakeServer.Received), 0)
	for _, d := range fakeServer.Received {
		assert.Equal(t, "Temperature", d.SensorType)
		assert.Equal(t, "A", d.Id1)
		assert.Equal(t, "1", d.Id2)
		assert.NotNil(t, d.Timestamp)
	}
}

func TestGenerator_ConcurrentOperations(t *testing.T) {
	gen := NewGenerator("localhost:50051", 50*time.Millisecond)
	sensorType := "Temperature"
	id1 := "A"
	id2 := "1"

	gen.Start(sensorType, id1, id2)
	defer gen.Stop() // ensure cleanup

	// Concurrent frequency updates
	go func() {
		for i := 0; i < 5; i++ {
			gen.UpdateFrequency(time.Duration(50+i*10) * time.Millisecond)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Wait for at least one data item, max 1 second
	select {
	case d := <-gen.dataCh:
		assert.Equal(t, sensorType, d.SensorType)
		assert.Equal(t, id1, d.Id1)
		assert.Equal(t, id2, d.Id2)
	case <-time.After(1 * time.Second):
		t.Error("Expected data generated during concurrent updates")
	}
}

func TestGenerator_MultipleDataGeneration_WithServer(t *testing.T) {
	// Start server on free port
	lis, _ := net.Listen("tcp", "127.0.0.1:0")
	server := grpc.NewServer()
	fakeServer := &FakeSensorServer{Done: make(chan struct{}, 1)}
	pb.RegisterSensorServiceServer(server, fakeServer)
	go server.Serve(lis)
	defer server.Stop()

	addr := lis.Addr().String()
	gen := NewGenerator(addr, 20*time.Millisecond)

	// Start generator
	gen.Start("Pressure", "B", "2")

	// Wait until at least one message arrives
	select {
	case <-fakeServer.Done:
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for server to receive data")
	}

	// Stop generator
	gen.Stop()

	// Assert multiple data items
	count := len(fakeServer.Received)
	assert.Greater(t, count, 0, "Expected multiple data items generated")
	for _, d := range fakeServer.Received {
		assert.Equal(t, "Pressure", d.SensorType)
		assert.Equal(t, "B", d.Id1)
		assert.Equal(t, "2", d.Id2)
		assert.NotNil(t, d.Timestamp)
	}
}
