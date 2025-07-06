package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/testutils"
	"github.com/stretchr/testify/require"
)

// BenchmarkChatRepository_Create benchmarks chat creation performance.
func BenchmarkChatRepository_Create(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(b, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			chat := fixtures.CreateTestChat(user.ID)
			chat.OrganizationID = org.ID

			err := repo.Create(chat)
			if err != nil {
				b.Fatalf("Failed to create chat: %v", err)
			}
		}
	})
}

// BenchmarkChatRepository_FindByID benchmarks chat retrieval by ID.
func BenchmarkChatRepository_FindByID(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(b, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(b, err)

	// Create test chats for lookup
	chatIDs := make([]uint64, 100)
	for i := 0; i < 100; i++ {
		chat := fixtures.CreateTestChat(user.ID)
		chat.OrganizationID = org.ID
		err := repo.Create(chat)
		require.NoError(b, err)
		chatIDs[i] = chat.ID
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, err := repo.FindByID(chatIDs[i%len(chatIDs)])
			if err != nil {
				b.Fatalf("Failed to find chat: %v", err)
			}
			i++
		}
	})
}

// BenchmarkChatRepository_FindByOrganizationID benchmarks paginated organization queries.
func BenchmarkChatRepository_FindByOrganizationID(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(b, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(b, err)

	// Create many chats for the organization
	for i := 0; i < 1000; i++ {
		chat := fixtures.CreateTestChat(user.ID)
		chat.OrganizationID = org.ID
		chat.Title = fmt.Sprintf("Chat %d", i)
		err := repo.Create(chat)
		require.NoError(b, err)
	}

	b.ResetTimer()
	b.Run("Limit10", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			offset := 0
			for pb.Next() {
				_, err := repo.FindByOrganizationID(org.ID, 10, offset%900)
				if err != nil {
					b.Fatalf("Failed to find chats: %v", err)
				}
				offset += 10
			}
		})
	})

	b.Run("Limit50", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			offset := 0
			for pb.Next() {
				_, err := repo.FindByOrganizationID(org.ID, 50, offset%950)
				if err != nil {
					b.Fatalf("Failed to find chats: %v", err)
				}
				offset += 50
			}
		})
	})

	b.Run("Limit100", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			offset := 0
			for pb.Next() {
				_, err := repo.FindByOrganizationID(org.ID, 100, offset%900)
				if err != nil {
					b.Fatalf("Failed to find chats: %v", err)
				}
				offset += 100
			}
		})
	})
}

// BenchmarkChatRepository_Update benchmarks chat updates.
func BenchmarkChatRepository_Update(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	dbWrapper := &Database{DB: db}
	repo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(b, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(b, err)

	// Create test chats for updating
	chats := make([]*domain.Chat, 100)
	for i := 0; i < 100; i++ {
		chat := fixtures.CreateTestChat(user.ID)
		chat.OrganizationID = org.ID
		err := repo.Create(chat)
		require.NoError(b, err)
		chats[i] = chat
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			chat := chats[i%len(chats)]
			chat.Title = fmt.Sprintf("Updated Chat %d", time.Now().UnixNano())
			chat.UpdatedAt = time.Now()

			err := repo.Update(chat)
			if err != nil {
				b.Fatalf("Failed to update chat: %v", err)
			}
			i++
		}
	})
}

// BenchmarkMessageRepository_Create benchmarks message creation.
func BenchmarkMessageRepository_Create(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	dbWrapper := &Database{DB: db}
	messageRepo := NewMessageRepository(dbWrapper)
	chatRepo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(b, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(b, err)

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = chatRepo.Create(chat)
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			message := fixtures.CreateTestMessage(chat.ID)

			err := messageRepo.Create(message)
			if err != nil {
				b.Fatalf("Failed to create message: %v", err)
			}
		}
	})
}

// BenchmarkMessageRepository_FindByChatID benchmarks message retrieval by chat.
func BenchmarkMessageRepository_FindByChatID(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	dbWrapper := &Database{DB: db}
	messageRepo := NewMessageRepository(dbWrapper)
	chatRepo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(b, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(b, err)

	chat := fixtures.CreateTestChat(user.ID)
	chat.OrganizationID = org.ID
	err = chatRepo.Create(chat)
	require.NoError(b, err)

	// Create many messages for the chat
	for i := 0; i < 1000; i++ {
		message := fixtures.CreateTestMessage(chat.ID)
		message.Content = fmt.Sprintf("Message %d content", i)
		err := messageRepo.Create(message)
		require.NoError(b, err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := messageRepo.FindByChatID(chat.ID)
			if err != nil {
				b.Fatalf("Failed to find messages: %v", err)
			}
		}
	})
}

// BenchmarkDatabase_ConcurrentOperations benchmarks concurrent database operations.
func BenchmarkDatabase_ConcurrentOperations(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	dbWrapper := &Database{DB: db}
	chatRepo := NewChatRepository(dbWrapper)
	messageRepo := NewMessageRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(b, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(b, err)

	// Pre-create some chats for read operations
	existingChats := make([]*domain.Chat, 50)
	for i := 0; i < 50; i++ {
		chat := fixtures.CreateTestChat(user.ID)
		chat.OrganizationID = org.ID
		err := chatRepo.Create(chat)
		require.NoError(b, err)
		existingChats[i] = chat
	}

	b.ResetTimer()
	b.Run("ReadWriteMix", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				switch i % 4 {
				case 0: // Create chat
					chat := fixtures.CreateTestChat(user.ID)
					chat.OrganizationID = org.ID
					_ = chatRepo.Create(chat)
				case 1: // Read chat
					_, _ = chatRepo.FindByID(existingChats[i%len(existingChats)].ID)
				case 2: // Create message
					message := fixtures.CreateTestMessage(existingChats[i%len(existingChats)].ID)
					_ = messageRepo.Create(message)
				case 3: // List chats
					_, _ = chatRepo.FindByOrganizationID(org.ID, 10, 0)
				}
				i++
			}
		})
	})
}

// BenchmarkQueryPerformance benchmarks complex query performance.
func BenchmarkQueryPerformance(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	dbWrapper := &Database{DB: db}
	chatRepo := NewChatRepository(dbWrapper)

	// Setup test data with multiple organizations
	for orgIdx := 0; orgIdx < 10; orgIdx++ {
		org := fixtures.CreateTestOrganization()
		err := db.Create(org).Error
		require.NoError(b, err)

		user := fixtures.CreateTestUser(org.ID)
		err = db.Create(user).Error
		require.NoError(b, err)

		// Create chats with various tags
		for i := 0; i < 100; i++ {
			chat := fixtures.CreateTestChat(user.ID)
			chat.OrganizationID = org.ID
			chat.Tags = fmt.Sprintf(`["tag-%d", "category-%d"]`, i%5, i%3)
			err := chatRepo.Create(chat)
			require.NoError(b, err)
		}
	}

	b.ResetTimer()
	b.Run("CountByDateRange", func(b *testing.B) {
		start := time.Now().AddDate(0, 0, -7)
		end := time.Now()

		b.RunParallel(func(pb *testing.PB) {
			orgID := uint64(1)
			for pb.Next() {
				_, err := chatRepo.CountByOrgIDAndDateRange(orgID, start, end)
				if err != nil {
					b.Fatalf("Failed to count chats: %v", err)
				}
				orgID = (orgID % 10) + 1
			}
		})
	})

	b.Run("GetTagStats", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			orgID := uint64(1)
			for pb.Next() {
				_, err := chatRepo.GetTagStats(orgID)
				if err != nil {
					b.Fatalf("Failed to get tag stats: %v", err)
				}
				orgID = (orgID % 10) + 1
			}
		})
	})
}

// BenchmarkMemoryUsage benchmarks memory usage patterns.
func BenchmarkMemoryUsage(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	dbWrapper := &Database{DB: db}
	chatRepo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err := db.Create(org).Error
	require.NoError(b, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(b, err)

	// Create large dataset
	for i := 0; i < 1000; i++ {
		chat := fixtures.CreateTestChat(user.ID)
		chat.OrganizationID = org.ID
		chat.Metadata = fmt.Sprintf(
			`{"large_field": "%s"}`,
			string(make([]byte, 1024)),
		) // 1KB metadata
		err := chatRepo.Create(chat)
		require.NoError(b, err)
	}

	b.ResetTimer()
	b.Run("LargeBatchQuery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			chats, err := chatRepo.FindByOrganizationID(org.ID, 500, 0)
			if err != nil {
				b.Fatalf("Failed to find chats: %v", err)
			}
			_ = chats // Prevent optimization
		}
	})
}

// BenchmarkConnectionPooling benchmarks database connection efficiency.
func BenchmarkConnectionPooling(b *testing.B) {
	db := testutils.SetupBenchmarkDB(b)
	defer testutils.CleanupBenchmarkDB()

	// Configure connection pool
	sqlDB, err := db.DB()
	require.NoError(b, err)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbWrapper := &Database{DB: db}
	chatRepo := NewChatRepository(dbWrapper)

	// Setup test data
	org := fixtures.CreateTestOrganization()
	err = db.Create(org).Error
	require.NoError(b, err)

	user := fixtures.CreateTestUser(org.ID)
	err = db.Create(user).Error
	require.NoError(b, err)

	b.ResetTimer()
	b.Run("HighConcurrency", func(b *testing.B) {
		b.SetParallelism(20) // More goroutines than connections
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				chat := fixtures.CreateTestChat(user.ID)
				chat.OrganizationID = org.ID
				err := chatRepo.Create(chat)
				if err != nil {
					b.Fatalf("Failed to create chat: %v", err)
				}
			}
		})
	})
}
