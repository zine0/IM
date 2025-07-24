package initialization

import (
	"context"
	"fmt"
	"github/zine0/IM/internal/repository"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func InitDB() *repository.Queries {
	// 获取配置
	dbConfig := viper.GetStringMapString("db")

	// 转换端口
	port, err := strconv.Atoi(dbConfig["port"])
	if err != nil {
		panic(fmt.Errorf("invalid port number: %w", err))
	}

	// 构建连接字符串
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		dbConfig["user"],
		dbConfig["password"],
		dbConfig["host"],
		port,
		dbConfig["database"],
	)

	// 配置连接池
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(fmt.Errorf("failed to parse connection config: %w", err))
	}

	// 自定义连接池设置
	config.MaxConns = 50                               // 最大连接数
	config.MinConns = 5                                // 最小保持连接数
	config.MaxConnLifetime = time.Hour                 // 连接最大生存时间
	config.MaxConnIdleTime = 30 * time.Minute          // 连接最大空闲时间
	config.HealthCheckPeriod = time.Minute             // 健康检查间隔
	config.ConnConfig.ConnectTimeout = 5 * time.Second // 连接超时时间

	// 创建连接池
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(fmt.Errorf("failed to create connection pool: %w", err))
	}

	// 测试连接
	if err := dbpool.Ping(ctx); err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}

	// 返回生成的查询接口
	return repository.New(dbpool)
}
