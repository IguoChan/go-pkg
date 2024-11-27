package mysqlx

import (
	"gorm.io/gorm"
)

type Option func(*gorm.DB) *gorm.DB

// DB 与其他 Option 一起使用时务必放在第一位，因为这里会换一个db对象
func DB(db *gorm.DB) Option {
	return func(*gorm.DB) *gorm.DB {
		return db
	}
}

func Select(fields ...string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fields)
	}
}

func TableName(tableName string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(tableName)
	}
}

func ID(id uint32) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func IDs(ids ...uint32) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in (?)", ids)
	}
}

func Limit(limit uint32) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(int(limit))
	}
}

func Offset(offset uint32) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset))
	}
}

func OrderBy(orderBy string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(orderBy)
	}
}
