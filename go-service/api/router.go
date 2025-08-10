package api

import (
	"context"

	apiV1 "go-service/api/v1"
	"go-service/constant"

	"github.com/gin-gonic/gin"
)

func GetRouter(ctx context.Context) *gin.Engine {
	router := gin.New()

	// add the routes
	v1 := router.Group(constant.V1Route)
	{
		addStudentRoutes(v1)
	}
	return router
}

func addStudentRoutes(r *gin.RouterGroup) {
	r.GET(constant.StudentReportRoute, apiV1.GetStudentPDFReport)
}
