package main

import (
	"Requester/internal/controller/rhttp"
	"Requester/internal/usecase"
	"Requester/internal/usecase/repo"
	"Requester/internal/usecase/webapi"
	"github.com/gin-gonic/gin"
)

func main() {

	//Initialize Repo
	repo := repo.New()

	//Initialize WebApi
	webapi := webapi.New()
	//Initialize UseCase
	recordUseCase := usecase.New(repo, webapi)
	//Initialize Router
	handler := gin.Default()
	rhttp.NewRouter(handler, recordUseCase)
	//v := handler.Group("v1")
	//{
	//	v.POST("/requester", Requester)
	//}

	handler.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
