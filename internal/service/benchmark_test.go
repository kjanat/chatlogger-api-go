package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/mocks"
	"github.com/stretchr/testify/mock"
)

// BenchmarkChatService_CreateChat benchmarks chat creation service performance.
func BenchmarkChatService_CreateChat(b *testing.B) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	// Setup mock to always succeed quickly
	mockRepo.On("Create", mock.AnythingOfType("context.Context"), mock.AnythingOfType("*domain.Chat")).
		Return(nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			chat := fixtures.CreateTestChat(1)
			err := service.CreateChat(context.Background(), chat)
			if err != nil {
				b.Fatalf("Failed to create chat: %v", err)
			}
		}
	})
}

// BenchmarkChatService_GetByID benchmarks chat retrieval performance.
func BenchmarkChatService_GetByID(b *testing.B) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	testChat := fixtures.CreateTestChat(1)
	testChat.ID = 1

	// Setup mock to return pre-created chat
	mockRepo.On("FindByID", mock.AnythingOfType("context.Context"), mock.AnythingOfType("uint64")).
		Return(testChat, nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := service.GetByID(context.Background(), 1)
			if err != nil {
				b.Fatalf("Failed to get chat: %v", err)
			}
		}
	})
}

// BenchmarkChatService_GetByOrganizationID benchmarks organization listing performance.
func BenchmarkChatService_GetByOrganizationID(b *testing.B) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	// Create test data of various sizes
	testChats := []domain.Chat{
		*fixtures.CreateTestChat(1),
		*fixtures.CreateTestChat(1),
		*fixtures.CreateTestChat(1),
	}

	b.Run("SmallResult", func(b *testing.B) {
		smallChats := testChats[:1]
		mockRepo.On("FindByOrganizationID", mock.AnythingOfType("context.Context"), uint64(1), 10, 0).
			Return(smallChats, nil)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := service.GetByOrganizationID(context.Background(), 1, 10, 0)
				if err != nil {
					b.Fatalf("Failed to get chats: %v", err)
				}
			}
		})
	})

	b.Run("MediumResult", func(b *testing.B) {
		mediumChats := make([]domain.Chat, 50)
		for i := range mediumChats {
			mediumChats[i] = *fixtures.CreateTestChat(1)
		}

		mockRepo.On("FindByOrganizationID", mock.AnythingOfType("context.Context"), uint64(2), 50, 0).
			Return(mediumChats, nil)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := service.GetByOrganizationID(context.Background(), 2, 50, 0)
				if err != nil {
					b.Fatalf("Failed to get chats: %v", err)
				}
			}
		})
	})

	b.Run("LargeResult", func(b *testing.B) {
		largeChats := make([]domain.Chat, 200)
		for i := range largeChats {
			largeChats[i] = *fixtures.CreateTestChat(1)
		}

		mockRepo.On("FindByOrganizationID", mock.AnythingOfType("context.Context"), uint64(3), 200, 0).
			Return(largeChats, nil)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := service.GetByOrganizationID(context.Background(), 3, 200, 0)
				if err != nil {
					b.Fatalf("Failed to get chats: %v", err)
				}
			}
		})
	})
}

// BenchmarkChatService_UpdateChat benchmarks chat update performance.
func BenchmarkChatService_UpdateChat(b *testing.B) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	existingChat := fixtures.CreateTestChat(1)
	existingChat.ID = 1

	// Setup mocks
	mockRepo.On("FindByID", mock.AnythingOfType("context.Context"), uint64(1)).
		Return(existingChat, nil)
	mockRepo.On("Update", mock.AnythingOfType("context.Context"), mock.AnythingOfType("*domain.Chat")).
		Return(nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			chatToUpdate := &domain.Chat{
				ID:             1,
				OrganizationID: 1,
				UserID:         uint64Ptr(1),
				Title:          "Updated Title",
				Tags:           `["updated"]`,
				Metadata:       `{"updated": true}`,
			}

			err := service.UpdateChat(context.Background(), chatToUpdate)
			if err != nil {
				b.Fatalf("Failed to update chat: %v", err)
			}
		}
	})
}

// BenchmarkChatService_GetChatStats benchmarks statistics aggregation performance.
func BenchmarkChatService_GetChatStats(b *testing.B) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	// Setup mock responses for stats
	mockRepo.On("CountByOrgIDAndDateRange", mock.AnythingOfType("context.Context"), mock.AnythingOfType("uint64"),
		mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
		Return(int64(1000), nil)

	tagStats := map[string]int64{
		"support":   500,
		"technical": 300,
		"general":   200,
	}
	mockRepo.On("GetTagStats", mock.AnythingOfType("context.Context"), mock.AnythingOfType("uint64")).
		Return(tagStats, nil)

	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := service.GetChatStats(context.Background(), 1, start, end)
			if err != nil {
				b.Fatalf("Failed to get chat stats: %v", err)
			}
		}
	})
}

// BenchmarkUserService_Authenticate benchmarks authentication performance.
func BenchmarkUserService_Authenticate(b *testing.B) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-jwt-secret")

	// Create test user with pre-hashed password
	testUser := fixtures.CreateTestUser(1)
	testUser.ID = 1
	testUser.Email = "test@example.com"
	// This is a bcrypt hash of "testpassword"
	testUser.PasswordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

	mockRepo.On("FindByEmail", mock.AnythingOfType("string")).Return(testUser, nil)

	b.ResetTimer()
	b.Run("ValidPassword", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _, err := service.Authenticate("test@example.com", "testpassword")
				if err != nil {
					b.Fatalf("Failed to authenticate: %v", err)
				}
			}
		})
	})

	b.Run("InvalidPassword", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _, err := service.Authenticate("test@example.com", "wrongpassword")
				// We expect this to fail, so only fail benchmark if no error
				if err == nil {
					b.Fatalf("Expected authentication to fail")
				}
			}
		})
	})
}

// BenchmarkUserService_Register benchmarks user registration performance.
func BenchmarkUserService_Register(b *testing.B) {
	mockRepo := &mocks.MockUserRepository{}
	service := NewUserService(mockRepo, "test-jwt-secret")

	// Setup mocks
	mockRepo.On("FindByEmail", mock.AnythingOfType("string")).Return((*domain.User)(nil), nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			user := &domain.User{
				Email:          "testuser@example.com",
				FirstName:      "Test",
				LastName:       "User",
				OrganizationID: 1,
			}

			err := service.Register(user, "securepassword123")
			if err != nil {
				b.Fatalf("Failed to register user: %v", err)
			}
			i++
		}
	})
}

// BenchmarkConcurrentServiceCalls benchmarks concurrent service operations.
func BenchmarkConcurrentServiceCalls(b *testing.B) {
	mockChatRepo := &mocks.MockChatRepository{}
	mockUserRepo := &mocks.MockUserRepository{}

	chatService := NewChatService(mockChatRepo)
	userService := NewUserService(mockUserRepo, "test-jwt-secret")

	// Setup mocks for concurrent operations
	testChat := fixtures.CreateTestChat(1)
	testChat.ID = 1
	testUser := fixtures.CreateTestUser(1)
	testUser.ID = 1

	mockChatRepo.On("Create", mock.AnythingOfType("context.Context"), mock.AnythingOfType("*domain.Chat")).
		Return(nil)
	mockChatRepo.On("FindByID", mock.AnythingOfType("context.Context"), mock.AnythingOfType("uint64")).
		Return(testChat, nil)
	mockChatRepo.On(
		"FindByOrganizationID",
		mock.AnythingOfType("context.Context"),
		mock.AnythingOfType("uint64"),
		mock.AnythingOfType(
			"int",
		),
		mock.AnythingOfType("int"),
	).Return([]domain.Chat{*testChat}, nil)

	mockUserRepo.On("FindByID", mock.AnythingOfType("uint64")).Return(testUser, nil)
	mockUserRepo.On(
		"FindByOrganizationID",
		mock.AnythingOfType("uint64"),
		mock.AnythingOfType(
			"int",
		),
		mock.AnythingOfType("int"),
	).Return([]domain.User{*testUser}, nil)

	b.ResetTimer()
	b.Run("MixedOperations", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				switch i % 5 {
				case 0:
					chat := fixtures.CreateTestChat(1)
					_ = chatService.CreateChat(context.Background(), chat)
				case 1:
					_, _ = chatService.GetByID(context.Background(), 1)
				case 2:
					_, _ = chatService.GetByOrganizationID(context.Background(), 1, 10, 0)
				case 3:
					_, _ = userService.GetByID(1)
				case 4:
					_, _ = userService.GetByOrganizationID(1, 10, 0)
				}
				i++
			}
		})
	})
}

// BenchmarkServiceLatency benchmarks service call latency under load.
func BenchmarkServiceLatency(b *testing.B) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	testChat := fixtures.CreateTestChat(1)
	testChat.ID = 1

	// Add artificial delay to simulate database latency
	mockRepo.On("FindByID", mock.AnythingOfType("context.Context"), mock.AnythingOfType("uint64")).
		Return(testChat, nil).
		Run(func(args mock.Arguments) {
			time.Sleep(time.Microsecond * 100) // 100μs simulated DB latency
		})

	b.ResetTimer()
	b.Run("LowConcurrency", func(b *testing.B) {
		b.SetParallelism(1)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := service.GetByID(context.Background(), 1)
				if err != nil {
					b.Fatalf("Failed to get chat: %v", err)
				}
			}
		})
	})

	b.Run("MediumConcurrency", func(b *testing.B) {
		b.SetParallelism(10)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := service.GetByID(context.Background(), 1)
				if err != nil {
					b.Fatalf("Failed to get chat: %v", err)
				}
			}
		})
	})

	b.Run("HighConcurrency", func(b *testing.B) {
		b.SetParallelism(50)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := service.GetByID(context.Background(), 1)
				if err != nil {
					b.Fatalf("Failed to get chat: %v", err)
				}
			}
		})
	})
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns.
func BenchmarkMemoryAllocation(b *testing.B) {
	mockRepo := &mocks.MockChatRepository{}
	service := NewChatService(mockRepo)

	// Create large slice to test memory allocation
	largeChats := make([]domain.Chat, 1000)
	for i := range largeChats {
		largeChats[i] = *fixtures.CreateTestChat(1)
		largeChats[i].Metadata = `{"large_data": "` + strings.Repeat("x", 1024) + `"}` // 1KB each
	}

	mockRepo.On("FindByOrganizationID", mock.AnythingOfType("context.Context"), uint64(1), 1000, 0).
		Return(largeChats, nil)

	b.ResetTimer()
	b.Run("LargeResultSet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			chats, err := service.GetByOrganizationID(context.Background(), 1, 1000, 0)
			if err != nil {
				b.Fatalf("Failed to get chats: %v", err)
			}
			// Force access to prevent optimization
			_ = len(chats)
		}
	})
}

// Helper function for benchmarks.
func uint64Ptr(i uint64) *uint64 {
	return &i
}
