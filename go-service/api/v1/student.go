package apiV1

import (
	"net/http"

	businessV1 "go-service/business/v1"
	"go-service/constant"

	"github.com/gin-gonic/gin"
)

func GetStudentPDFReport(ctx *gin.Context) {
	id := ctx.Param("id")
	pdfBytes, filename, err := businessV1.GetStudentPDFReport(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Header(constant.HeaderContentDisposition, "attachment; filename=\""+filename+"\"")
	ctx.Data(http.StatusOK, constant.ContentTypePDF, pdfBytes)
}
