package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// --- Models (データモデル) ---

// User: ユーザー情報を格納する構造体
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username     string    `gorm:"uniqueIndex"`
	PasswordHash string
	Language     string `gorm:"default:'en'"`
	CreatedAt    time.Time
}

// Event: カレンダーの予定情報を格納する構造体
type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID
	Title       string    // 予定のタイトル
	StartTime   time.Time // 開始日時
	EndTime     time.Time // 終了日時
	Description string    // 詳細説明
	Color       string    // 表示色
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Memo: 日々のメモ（日記やフリーメモ）を格納する構造体
type Memo struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID
	Title      string     // メモのタイトル
	Content    string     // メモの内容
	LinkedDate *time.Time // カレンダーの日付とリンクする場合の日付 (NULL可)
	ThemeColor string     // 本棚に表示する際の色
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// --- Database (データベース設定) ---

var db *gorm.DB

// initDB: データベース接続の初期化とマイグレーションを実行
func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// 環境変数が設定されていない場合のデフォルト接続文字列
		dsn = "host=db user=user password=password dbname=calendar_db port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	}

	var err error
	// データベース接続のリトライ処理 (最大10回)
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

	// Auto Migrate: 定義した構造体を元にデータベースのテーブルを自動作成・更新
	err = db.AutoMigrate(&User{}, &Event{}, &Memo{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
}

// --- Handlers (APIハンドラー) ---

// getEvents: 全ての予定を取得する
func getEvents(c *gin.Context) {
	var events []Event
	// 現状は全ての予定を取得していますが、理想的には日付範囲でフィルタリングすべきです
	db.Find(&events)
	c.JSON(http.StatusOK, events)
}

// createEvent: 新しい予定を作成する
func createEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	// user_id は認証機能実装後にコンテキストから取得するように修正が必要
	// event.UserID = ...
	db.Create(&event)
	c.JSON(http.StatusOK, event)
}

// updateEvent: 既存の予定を更新する
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

// deleteEvent: 予定を削除する
func deleteEvent(c *gin.Context) {
	id := c.Param("id")
	db.Delete(&Event{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// getMemos: 全てのメモを取得する (更新日時順)
func getMemos(c *gin.Context) {
	var memos []Memo
	db.Order("updated_at desc").Find(&memos)
	c.JSON(http.StatusOK, memos)
}

// createMemo: 新しいメモを作成する
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

// updateMemo: 既存のメモを更新する
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

// deleteMemo: メモを削除する
func deleteMemo(c *gin.Context) {
	id := c.Param("id")
	db.Delete(&Memo{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// --- AI Advisor (AIアドバイザー機能) ---

type AdviceRequest struct {
	WeatherCode int     `json:"weather_code"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Language    string  `json:"language"` // "en" or "ja"
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// getAdvice: 天気情報に基づいてAIからのアドバイスを取得する
func getAdvice(c *gin.Context) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Error: GEMINI_API_KEY is not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error: API Key missing"})
		return
	}
	log.Printf("Using API Key: %s...%s", apiKey[:4], apiKey[len(apiKey)-4:])

	var req AdviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prompt string
	if req.Language == "ja" {
		prompt = fmt.Sprintf("天気は%s (%.1f°C)です。今日の服装を具体的に提案してください。100文字以内で簡潔に答えてください。", req.Description, req.Temperature)
	} else {
		prompt = fmt.Sprintf("The weather is %s (%.1f°C). Suggest a specific outfit for today. Keep it concise (under 100 characters).", req.Description, req.Temperature)
	}
	log.Printf("Sending prompt to Gemini: %s", prompt)

	geminiReq := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(geminiReq)
	// Gemini APIのエンドポイント (gemini-flash-latest を使用)
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent?key=" + apiKey

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error calling Gemini API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call AI service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Gemini API Error: Status=%d Body=%s", resp.StatusCode, string(bodyBytes))
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("AI service returned error: %d", resp.StatusCode)})
		return
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		log.Printf("Error decoding response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
		return
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		advice := geminiResp.Candidates[0].Content.Parts[0].Text
		log.Printf("Received advice: %s", advice)
		c.JSON(http.StatusOK, gin.H{"advice": advice})
	} else {
		log.Println("No advice candidates returned")
		c.JSON(http.StatusOK, gin.H{"advice": "No advice available."})
	}
}

// --- Main (メイン関数) ---

func main() {
	// データベース接続初期化
	initDB()

	// Ginルーターの初期化
	r := gin.Default()
	
	// CORS設定: 異なるオリジン（フロントエンドなど）からのアクセスを許可
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 開発環境用のため全て許可。本番では制限すべき
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// APIエンドポイントグループ
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

		api.POST("/advice", getAdvice) // AIアドバイス取得

		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	// サーバーをポート8080で起動
	r.Run(":8080")
}
