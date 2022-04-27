package main

import (
	"Requester/internal/controller/rhttp"
	"Requester/internal/usecase"
	"Requester/internal/usecase/repo"
	"Requester/internal/usecase/webapi"
	"github.com/gin-gonic/gin"
)

func main() {
	//Написать приложение, которое принимает http запросы от клиентов.
	//
	//- [ ] - Клиент передает список url, приложение отвечает клиенту с информацией о каждом из url (размер страницы) в формате json
	//- [ ] - Приложение, за один запрос, может обработать до 100 url.
	//- [ ] - Одновременно обрабатываются 3 url. (3 goroutine)
	//- [ ] - Если страница недоступна, то завершается работа приложения с ошибкой для клиента.
	//- [ ] - Если клиент отключается не дождавшись ответа, то вся работа по проверке url  прерывается.  (context)
	//- [ ] - Приложение отдает ошибку если уже одновременно обрабатываются 3 запроса от клиентов. (count in work process =< 3)
	//- [ ] - Если запрашивается информация о странице, которая уже была получена не более 10  секунд назад, то отдается старая информация (Cache with map and mutex)
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
