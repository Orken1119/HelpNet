package main

import (
	"fmt"
	"log"

	_ "github.com/Orken1119/HelpNet/docs"
	controller "github.com/Orken1119/HelpNet/internal/controller"
	pkg "github.com/Orken1119/HelpNet/pkg"
	"github.com/gin-gonic/gin"
)

// @title           HelpNet API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      helpnet-1.onrender.com

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	app, err := pkg.App()

	if err != nil {
		log.Fatal(err)
	}
	defer app.CloseDBConnection()

	ginRouter := gin.Default()

	controller.Setup(app, ginRouter)

	ginRouter.Run(fmt.Sprintf(":%d", 2000))
}
