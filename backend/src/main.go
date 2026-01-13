package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- Models ---

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username     string    `gorm:"uniqueIndex"`
	PasswordHash string
	Language     string `gorm:"default:'en'"`
	CreatedAt    time.Time
}

type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	Description string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Memo struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID
	Title      string
	Content    string
	LinkedDate *time.Time // Nullable
	ThemeColor string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// --- Database ---

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=db user=user password=password dbname=calendar_db port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	}

	var err error
	// Retry connection
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	// Auto Migrate
	err = db.AutoMigrate(&User{}, &Event{}, &Memo{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
}

// --- Handlers ---

func getEvents(c *gin.Context) {
	var events []Event
	// Simple fetch all for now, ideally filter by date range
	db.Find(&events)
	c.JSON(http.StatusOK, events)
}

func createEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	// user_id hardcoded for now or from context
	// event.UserID = ...
	db.Create(&event)
	c.JSON(http.StatusOK, event)
}

func updateEvent(c *gin.Context) {
	id := c.Param("id")
	var event Event
	if result := db.First(&event, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event.UpdatedAt = time.Now()
	db.Save(&event)
	c.JSON(http.StatusOK, event)
}

func deleteEvent(c *gin.Context) {
	id := c.Param("id")
	db.Delete(&Event{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func getMemos(c *gin.Context) {
	var memos []Memo
	db.Order("updated_at desc").Find(&memos)
	c.JSON(http.StatusOK, memos)
}

func createMemo(c *gin.Context) {
	var memo Memo
	if err := c.ShouldBindJSON(&memo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	memo.CreatedAt = time.Now()
	memo.UpdatedAt = time.Now()
	db.Create(&memo)
	c.JSON(http.StatusOK, memo)
}

func updateMemo(c *gin.Context) {
	id := c.Param("id")
	var memo Memo
	if result := db.First(&memo, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Memo not found"})
		return
	}
	if err := c.ShouldBindJSON(&memo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	memo.UpdatedAt = time.Now()
	db.Save(&memo)
	c.JSON(http.StatusOK, memo)
}

func deleteMemo(c *gin.Context) {
	id := c.Param("id")
	db.Delete(&Memo{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// --- Main ---

func main() {
	initDB()

	r := gin.Default()
	
	// Cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // For dev
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		api.GET("/events", getEvents)
		api.POST("/events", createEvent)
		api.PUT("/events/:id", updateEvent)
		api.DELETE("/events/:id", deleteEvent)

		api.GET("/memos", getMemos)
		api.POST("/memos", createMemo)
		api.PUT("/memos/:id", updateMemo)
		api.DELETE("/memos/:id", deleteMemo)

		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	r.Run(":8080")
}
