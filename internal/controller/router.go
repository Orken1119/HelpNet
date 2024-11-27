package controller

import (
	"time"

	"github.com/Orken1119/HelpNet/pkg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Orken1119/HelpNet/internal/controller/auth_controller"
	"github.com/Orken1119/HelpNet/internal/controller/auth_controller/middleware"
	"github.com/Orken1119/HelpNet/internal/controller/event_controller"
	organization "github.com/Orken1119/HelpNet/internal/controller/organization_controller"
	"github.com/Orken1119/HelpNet/internal/controller/volunteer_controller"
	repository "github.com/Orken1119/HelpNet/internal/repository"
)

func Setup(app pkg.Application, router *gin.Engine) {
	db := app.DB

	loginController := &auth_controller.AuthController{
		UserRepository: repository.NewUserRepository(db),
	}

	userController := &volunteer_controller.UserController{
		UserRepository: repository.NewUserRepository(db),
	}

	organizationController := &organization.OrganizationController{
		OrganizationRepository: repository.NewOrganizationRepository(db),
	}

	eventController := &event_controller.EventController{
		EventRepository: repository.NewEventRepository(db),
	}

	// CORS settings
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Укажите URL вашего фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authenticationRouter := router.Group("/authentication")
	{
		authenticationRouter.POST("/signup-as-volunteer", loginController.SignupAsVolunteer)
		authenticationRouter.POST("/signin-as-volunteer", loginController.SigninAsVolunteer)
		authenticationRouter.POST("/manual-organization-registration", loginController.Signup)
		authenticationRouter.POST("/signin-as-organization", loginController.Signin)
		authenticationRouter.POST("/forgot-password", loginController.ForgotPassword)
		authenticationRouter.POST("/change-forgotten-password", loginController.ChangeForgottenPassword)

	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(middleware.JWTAuth(`access-secret-key`))

	userRouter := router.Group("/user")
	{
		userRouter.GET("/profile", userController.GetProfile)
		userRouter.PUT("/edit-profile", userController.EditPersonalData)
		userRouter.PUT("/change-password", userController.ChangePassword)
	}

	eventRouter := router.Group("/events")
	{
		eventRouter.GET("/create-ivent", eventController.CreateEvent)
		eventRouter.GET("/delete-ivent/:id", eventController.CreateEvent)
		eventRouter.GET("/finish/:id", eventController.CreateEvent)
		eventRouter.GET("/get-ivent-by-id/:id", eventController.CreateEvent)
		eventRouter.GET("/get-ivents", eventController.CreateEvent)
		eventRouter.GET("/finished-events-by-organization", eventController.CreateEvent)
		eventRouter.GET("/get-organizations-in-process", eventController.CreateEvent)
		eventRouter.GET("/get-user-participating", eventController.CreateEvent)
		eventRouter.GET("/get-volunteer-finished-events", eventController.CreateEvent)
		eventRouter.GET("/participate-event", eventController.CreateEvent)
		eventRouter.GET("/update-ivent", eventController.CreateEvent)
	}

	organizationRouter := router.Group("/organizations")
	{
		organizationRouter.DELETE("/delete-organiztion/:id", organizationController.DeleteOrganization)
		organizationRouter.GET("/profile", organizationController.DeleteOrganization)
		organizationRouter.PUT("/edit-organization-profile", organizationController.EditOrganization)
	}

}
