package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/sirupsen/logrus"
)

func GetProducts(search geo.Search) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _ := c.GetQuery("ip")

		geoData, err := search.ByIp(ip)
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(
			http.StatusOK,
			geoData,
		)
	}
}
