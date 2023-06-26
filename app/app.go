package app

import (
	"log"

	"gorm.io/gorm"
)

type MarketPlaceAPIs struct {
	DB *gorm.DB
}

func NewMarketPlaceAPIs(db *gorm.DB) MarketPlaceAPIs {
	return MarketPlaceAPIs{
		DB: db,
	}

}

func apiErr(apiName string, errMsg string, err error) {
	if err != nil {
		log.Printf("ApiName: %s, %s: %v", apiName, errMsg, err)
	}
}
