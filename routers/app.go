package routers

import (
	"astroauth-api/database"
	"astroauth-api/middleware"
	"astroauth-api/models"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func AppRouter(router *gin.Engine) {
	appuser := router.Group("/site")

	appuser.Use(middleware.CheckSession())
	{
		appuser.POST("/app", middleware.AppCreateValidation(), middleware.CheckAppSite(), CreateApp)
	}
}

func CreateApp(c *gin.Context) {
	var rApp models.App
	c.ShouldBindBodyWith(&rApp, binding.JSON)

	// Get max apps a user can create
	var maxapp uint
	database.DBB.QueryRow(context.Background(), "SELECT max_app FROM site_users WHERE id = $1", c.MustGet("userID")).Scan(&maxapp)

	//get number of apps the user has
	var appcount uint
	database.DBB.QueryRow(context.Background(), "SELECT COUNT(*) FROM apps WHERE owned_by = $1", c.MustGet("userID")).Scan(&appcount)

	//check if number apps is less than or equal to max apps
	if appcount >= maxapp {
		c.JSON(200, models.Error{Message: "Max apps reached"})
		return
	}

	FindName, err := database.DBB.Exec(context.Background(), "SELECT name FROM apps WHERE name = $1", rApp.Name)
	if err == nil {
		c.JSON(500, models.Error{Message: "Internal server error"})
		return
	}
	if FindName.RowsAffected() != 0 {
		c.JSON(200, models.Error{Message: "Name not available"})
		return
	}

	rApp.OwnedBy = c.MustGet("userID").(uint)
	database.DB.Create(&rApp)
	c.JSON(200, rApp)
}
