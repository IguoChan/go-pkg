package mysqlx

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

var (
	ErrEmptyConfig = errors.New("empty config")
)

type Client struct {
	*gorm.DB
}

func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, ErrEmptyConfig
	}
	config.setDefault()

	logger := NewLogger(&gLogger.Config{
		SlowThreshold:             config.SlowThreshold,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  config.LogLevel,
	})

	// client
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
				config.Username, config.Password, config.Addr, config.DBName), // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}),
		&gorm.Config{
			SkipDefaultTransaction: true, // 开启提高性能，https://gorm.io/docs/transactions.html
			Logger:                 logger,
		},
	)
	if err != nil {
		return nil, err
	} // get sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxOpenConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour * 4)

	return &Client{db}, nil
}

func (c *Client) Close() {
	if c != nil {
		db, err := c.DB.DB()
		if err == nil {
			_ = db.Close()
		}
	}
}
