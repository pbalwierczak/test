package events

import (
	"context"

	"scootin-aboot/internal/models"
	"scootin-aboot/internal/services"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTripService struct {
	mock.Mock
}

func (m *MockTripService) StartTrip(ctx context.Context, scooterID, userID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	args := m.Called(ctx, scooterID, userID, lat, lng)
	return args.Get(0).(*models.Trip), args.Error(1)
}

func (m *MockTripService) EndTrip(ctx context.Context, scooterID uuid.UUID, lat, lng float64) (*models.Trip, error) {
	args := m.Called(ctx, scooterID, lat, lng)
	return args.Get(0).(*models.Trip), args.Error(1)
}

func (m *MockTripService) CancelTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, scooterID)
	return args.Get(0).(*models.Trip), args.Error(1)
}

func (m *MockTripService) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	args := m.Called(ctx, scooterID, lat, lng)
	return args.Error(0)
}

func (m *MockTripService) GetActiveTrip(ctx context.Context, scooterID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, scooterID)
	return args.Get(0).(*models.Trip), args.Error(1)
}

func (m *MockTripService) GetActiveTripByUser(ctx context.Context, userID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.Trip), args.Error(1)
}

func (m *MockTripService) GetTrip(ctx context.Context, tripID uuid.UUID) (*models.Trip, error) {
	args := m.Called(ctx, tripID)
	return args.Get(0).(*models.Trip), args.Error(1)
}

type MockScooterService struct {
	mock.Mock
}

func (m *MockScooterService) GetScooters(ctx context.Context, params services.ScooterQueryParams) (*services.ScooterListResult, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*services.ScooterListResult), args.Error(1)
}

func (m *MockScooterService) GetScooter(ctx context.Context, id uuid.UUID) (*services.ScooterDetailsResult, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*services.ScooterDetailsResult), args.Error(1)
}

func (m *MockScooterService) GetClosestScooters(ctx context.Context, params services.ClosestScootersQueryParams) (*services.ClosestScootersResult, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*services.ClosestScootersResult), args.Error(1)
}

func (m *MockScooterService) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	args := m.Called(ctx, scooterID, lat, lng)
	return args.Error(0)
}

type MockConsumerGroupSession struct {
	mock.Mock
}

func (m *MockConsumerGroupSession) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *MockConsumerGroupSession) Claims() map[string][]int32 {
	args := m.Called()
	return args.Get(0).(map[string][]int32)
}

func (m *MockConsumerGroupSession) MemberID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConsumerGroupSession) GenerationID() int32 {
	args := m.Called()
	return args.Get(0).(int32)
}

func (m *MockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *MockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *MockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.Called(msg, metadata)
}

func (m *MockConsumerGroupSession) Commit() {
	m.Called()
}

type MockConsumerGroupClaim struct {
	mock.Mock
	messages chan *sarama.ConsumerMessage
}

func NewMockConsumerGroupClaim() *MockConsumerGroupClaim {
	return &MockConsumerGroupClaim{
		messages: make(chan *sarama.ConsumerMessage, 10),
	}
}

func (m *MockConsumerGroupClaim) Topic() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConsumerGroupClaim) Partition() int32 {
	args := m.Called()
	return args.Get(0).(int32)
}

func (m *MockConsumerGroupClaim) InitialOffset() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockConsumerGroupClaim) HighWaterMarkOffset() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	args := m.Called()
	if args.Get(0) == nil {
		return m.messages
	}
	return args.Get(0).(<-chan *sarama.ConsumerMessage)
}

func (m *MockConsumerGroupClaim) SendMessage(msg *sarama.ConsumerMessage) {
	m.messages <- msg
}

func (m *MockConsumerGroupClaim) Close() {
	close(m.messages)
}

type MockConsumerGroup struct {
	mock.Mock
}

func (m *MockConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	args := m.Called(ctx, topics, handler)
	return args.Error(0)
}

func (m *MockConsumerGroup) Errors() <-chan error {
	args := m.Called()
	return args.Get(0).(<-chan error)
}

func (m *MockConsumerGroup) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConsumerGroup) Pause(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *MockConsumerGroup) Resume(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *MockConsumerGroup) PauseAll() {
	m.Called()
}

func (m *MockConsumerGroup) ResumeAll() {
	m.Called()
}
