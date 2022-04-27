package rhttp

import (
	"Requester/internal/entity"
	"Requester/internal/usecase"
	"fmt"
	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
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

func newRecordRoutes(handler *gin.RouterGroup, t usecase.Record) {
	r := &requestRoutes{t}

	h := handler.Group("/v1")
	{
		h.Use(limit.MaxAllowed(3))
		h.POST("/process", r.doRequest)
	}
}

type returnRecord struct {
	mu   sync.RWMutex
	Mapa map[string]entity.Record
}

func (r *requestRoutes) doRequest(c *gin.Context) {
	var requestURLs []string

	//Одновременно обрабатываются 3 url. (3 goroutine)
	guard := make(chan struct{}, 3)
	result := make(chan entity.Record, 100)
	if err := c.ShouldBindJSON(&requestURLs); err != nil {
		fmt.Println("error, ShouldBindJSON", err)
		//TODO
	}
	resultArr := make([]entity.Record, 0)
	counter := 0
	wg := sync.WaitGroup{}
	for _, val := range requestURLs {
		guard <- struct{}{}
		wg.Add(1)
		if counter <= 100 {

			go func(val string) {
				defer wg.Done()
				record, ok, err := r.t.Get(c.Request.Context(), val)
				if err != nil {
					c.Abort()
					//TODO err
				}
				temp := time.Now().Unix() - record.Ttl
				fmt.Println(temp)
				if !ok || time.Now().Unix()-record.Ttl > 10 {
					res, err := http.Get(val)
					if err != nil {
						c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "BadRequest"})
						return
						//TODO err
					}
					f, err := ioutil.ReadAll(res.Body)
					if err != nil {
						c.Abort()
						//TODO err
					}
					rec := entity.Record{URL: val, Size: len(f), Ttl: time.Now().Unix()}
					err = r.t.Add(c.Request.Context(), rec)
					if err != nil {
						c.Abort()
						//TODO err
					}
					result <- rec
				} else {
					result <- record
				}

				//return record
				fmt.Println(record)
				<-guard
			}(val)
			//wg.Wait()

			select {
			case t := <-result:
				resultArr = append(resultArr, t)
				fmt.Println("t", t)
			case <-c.Request.Context().Done():
				fmt.Println("ABORTED11111")
				return
			}
		}

		counter++

	}

	c.JSON(http.StatusOK, gin.H{"data": resultArr})
}
