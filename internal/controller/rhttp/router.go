package rhttp

import (
	"Requester/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(handler *gin.Engine, r usecase.Record) {

	//handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	h := handler.Group("/api")
	{
		newRecordRoutes(h, r)
	}
}

type requestRoutes struct {
	t usecase.Record
}

func LimitMiddleware() gin.HandlerFunc {
	semaphore := make(chan bool, 3)

	return func(c *gin.Context) {
		select {
		case semaphore <- true:
			c.Next()
			<-semaphore
		default:
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Connection refused. Too many connection"})
			c.Abort()
			return
		}
	}
}

func newRecordRoutes(handler *gin.RouterGroup, t usecase.Record) {
	r := &requestRoutes{t}

	h := handler.Group("/v1")
	{
		h.Use(LimitMiddleware())
		h.POST("/process", r.doRequest)
	}
}

func (r *requestRoutes) doRequest(c *gin.Context) {
	var requestURLs []string
	if err := c.ShouldBindJSON(&requestURLs); err != nil {
		c.Abort()
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		fmt.Println("error, ShouldBindJSON", err)
		return
	}
	if len(requestURLs) > 100 {
		c.Abort()
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Count of urls > 100"})
		return
	}

	resultArr, err := r.t.Requester(c.Request.Context(), requestURLs)
	//fmt.Println(resultArr, err)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": resultArr})
	}

}
