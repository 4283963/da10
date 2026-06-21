package routes

import (
	"transfer-tracker/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	vehicleHandler := handlers.NewVehicleHandler()
	nodeHandler := handlers.NewTransferNodeHandler()
	expenseHandler := handlers.NewExpenseHandler()
	uploadHandler := handlers.NewUploadHandler()

	r.Static(uploadHandler.PublicURL(), uploadHandler.UploadDir())

	api := r.Group("/api")
	{
		upload := api.Group("/upload")
		{
			upload.POST("", uploadHandler.Upload)
		}

		vehicles := api.Group("/vehicles")
		{
			vehicles.POST("", vehicleHandler.Create)
			vehicles.GET("", vehicleHandler.List)
			vehicles.GET("/vin/:vin", vehicleHandler.GetByVIN)
			vehicles.DELETE("/:id", vehicleHandler.Delete)
		}

		nodes := api.Group("/nodes")
		{
			nodes.PUT("/:id", nodeHandler.Update)
			nodes.GET("/vehicle/:vehicleId/progress", nodeHandler.GetProgress)
		}

		expenses := api.Group("/expenses")
		{
			expenses.POST("", expenseHandler.Create)
			expenses.GET("/vehicle/:vehicleId/stats", expenseHandler.GetStatistics)
			expenses.DELETE("/:id", expenseHandler.Delete)
		}
	}

	return r
}
