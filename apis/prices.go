package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/KyberNetwork/reserve-data/market"
)

func AllPrices(c *gin.Context) {
	fmt.Printf("Getting all prices \n")
	data, err := market.GetAllPrice()
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success":   true,
				"version":   data.Version,
				"timestamp": data.Timestamp,
				"data":      data.AllPairData,
			},
		)
	}
}

func Price(c *gin.Context) {
	base := c.Param("base")
	quote := c.Param("quote")
	fmt.Printf("Getting price for %s - %s \n", base, quote)
	data, err := market.GetPrice(base, quote)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success":   true,
				"version":   data.Version,
				"timestamp": data.Timestamp,
				"exchanges": data.ExchangeData,
			},
		)
	}
}
