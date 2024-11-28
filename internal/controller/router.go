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
		userRouter.GET("/profile/:id", userController.GetProfile)
		userRouter.PUT("/edit-profile", userController.EditPersonalData)
		userRouter.PUT("/change-password", userController.ChangePassword)
		userRouter.PUT("/change-password-for-org", userController.ChangePasswordForOrg)
		userRouter.POST("/add-certificate", userController.AddCertificate)
		userRouter.DELETE("/delete-certificate/:id", userController.DeleteCertificate)
		userRouter.GET("/get-all-volunteers-profile", userController.GetAllVolunteersProfile)
	}

	eventRouter := router.Group("/events")
	{
		eventRouter.POST("/create-event", eventController.CreateEvent)
		eventRouter.DELETE("/delete-event/:id", eventController.DeleteEvent)
		eventRouter.PUT("/finish/:id", eventController.FinishEvent)
		eventRouter.GET("/get-event-by-id/:id", eventController.GetEventById)
		eventRouter.GET("/get-events", eventController.GetAllEvent)
		eventRouter.GET("/finished-events-by-organization", eventController.GetOrgFinishedEvents)
		eventRouter.GET("/get-organizations-in-process", eventController.GetOrgEvents)
		eventRouter.GET("/get-user-participating", eventController.GetVolEvents)
		eventRouter.GET("/get-volunteer-finished-events", eventController.GetVolFinishedEvents)
		eventRouter.POST("/participate-event/:userID/:id", eventController.JoinEvent)
		eventRouter.PUT("/update-event/:id", eventController.UpdateEvent)
		eventRouter.GET("/get-event-by-direction/:direction", eventController.GetEventsByDirection)
	}

	organizationRouter := router.Group("/organizations")
	{
		organizationRouter.DELETE("/delete-organiztion/:id", organizationController.DeleteOrganization)
		organizationRouter.GET("/profile/:id", organizationController.GetOrganizationProfile)
		organizationRouter.PUT("/edit-organization-profile", organizationController.EditOrganization)
		organizationRouter.GET("/get-all-organizations-profile", organizationController.GetAllOrganizationsProfile)
	}
}
