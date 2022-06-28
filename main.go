package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

//type StringInterfaceMap map[string]interface{}

func HealthCheck(c *gin.Context) {
	message := "pong"
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
	return
}

func InitApplication() {
	log.Println("Initialize application ...")
	r := gin.Default()
	r.GET("/ping", HealthCheck)
	r.GET("/", GetPayload)
	r.Run()
}

const (
	STOCKPRICE = "stock_price"
	BUY        = "buy"
	SELL       = "sell"
)

type Rates struct {
	Id          uint           `gorm:"primaryKey" json:"id"`
	Origin      string         `json:"origin"`      // PAXOS, Binance
	Transaction string         `json:"transaction"` // constantes
	Payload     datatypes.JSON `json:"payload"`     // Metadata do resultado retornado pelo serviço
	CreatedAt   time.Time      `json:"created_at"`  // Data hora criação do rate
}

var Db *gorm.DB

var mockStockBTCBuy = datatypes.JSON([]byte(`{"id": "366a26d2-3098-4226-a520-4bb43ae4d921","market": "BTCUSD","side": "BUY","price": "6001.2","base_asset": "BTC","quote_asset": "USD","created_at": "2020-01-17T18:36:08Z","expires_at": "2020-01-17T18:36:38Z"}`))
var mockStockEthBuy = datatypes.JSON([]byte(`{"id": "366a26d2-3098-4226-a520-4bb43ae4d922","market": "ETHUSD","side": "BUY","price": "6001.2","base_asset": "ETC","quote_asset": "USD","created_at": "2020-01-17T18:36:08Z","expires_at": "2020-01-17T18:36:38Z"}`))
var mockStockBTCSell = datatypes.JSON([]byte(`{"id": "366a26d2-3098-4226-a520-4bb43ae4d923","market": "BTCUSD","side": "SELL","price": "6001.2","base_asset": "BTC","quote_asset": "USD","created_at": "2020-01-17T18:36:08Z","expires_at": "2020-01-17T18:36:38Z"}`))
var mockStockEthSell = datatypes.JSON([]byte(`{"id": "366a26d2-3098-4226-a520-4bb43ae4d924","market": "ETHUSD","side": "SELL","price": "6001.2","base_asset": "ETC","quote_asset": "USD","created_at": "2020-01-17T18:36:08Z","expires_at": "2020-01-17T18:36:38Z"}`))

func InitDatabase() {
	log.Println("Initialize database ...")
	dsn := "root:root@tcp(127.0.0.1:3306)/poc?charset=utf8&parseTime=true"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Panicf("Error in database %s", err)
	}

	conn.AutoMigrate(&Rates{})

	conn.Exec("delete from rates;")

	conn.Create(&Rates{
		Origin:      "Paxos",
		Transaction: STOCKPRICE,
		Payload:     mockStockBTCBuy,
		CreatedAt:   time.Now().UTC(),
	})

	conn.Create(&Rates{
		Origin:      "Paxos",
		Transaction: BUY,
		Payload:     mockStockEthBuy,
		CreatedAt:   time.Now().UTC(),
	})

	conn.Create(&Rates{
		Origin:      "Paxos",
		Transaction: SELL,
		Payload:     mockStockEthSell,
		CreatedAt:   time.Now().UTC(),
	})

	conn.Create(&Rates{
		Origin:      "Paxos",
		Transaction: SELL,
		Payload:     mockStockBTCSell,
		CreatedAt:   time.Now().UTC(),
	})

	Db = conn
	log.Println("Database completed ..")
	return
}

func main() {
	InitDatabase()
	InitApplication()
}

func GetPayload(c *gin.Context) {
	rates := []Rates{}
	stock := c.Query("stock")

	if stock == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "use params stock = BTCUSD or ETHUSD,",
		})
		return
	}

	values := Db.Find(&rates, datatypes.JSONQuery("payload").Equals(stock, "market"))

	if values.Error != nil {
		log.Printf("ERROR IN GET DATABASE %s", values.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprint(values.Error),
		})
		return
	}

	c.JSON(http.StatusOK, rates)
}
