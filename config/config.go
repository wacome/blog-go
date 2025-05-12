package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"blog-go/ent"
	"blog-go/ent/migrate"

	_ "github.com/lib/pq"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

// Config 应用配置
type Config struct {
	Environment string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
	JWTSecret   string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	return &Config{
		Environment: getEnv("ENV", "development"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		ServerPort:  os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
}

// InitDB 初始化数据库连接
func InitDB() (*sql.DB, error) {
	cfg := LoadConfig()
	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("数据库连接成功")
	return db, nil
}

// InitEntClient 初始化ent客户端
func InitEntClient() (*ent.Client, error) {
	cfg := LoadConfig()
	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 创建ent驱动
	drv := entsql.OpenDB(dialect.Postgres, db)

	// 创建ent客户端
	client := ent.NewClient(ent.Driver(drv))

	// 运行迁移
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		return nil, fmt.Errorf("创建schema失败: %w", err)
	}

	log.Println("ent客户端初始化成功")
	return client, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
