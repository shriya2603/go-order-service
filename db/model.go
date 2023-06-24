package db

import "time"

type Order struct {
	ID         uint `gorm:"primarykey"`
	CustomerID uint
	Status     string `gorm:"not null;type:text;default:placed"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Product struct {
	ID          uint `gorm:"primarykey"`
	ProductName string
	EanBarcode  string
	Category    string
	SubCategory string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderProduct struct {
	ID         uint `gorm:"primarykey"`
	OrderID    uint `gorm:"foreignKey:OrderID"`
	ProductID  uint `gorm:"foreignKey:ProductID"`
	CustomerID uint `gorm:"foreignKey:CustomerID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Customer struct {
	ID           uint `gorm:"primarykey"`
	CustomerName string
	Address      string
	Phone        int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
