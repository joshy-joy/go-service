package businessV1

import (
	"context"
	"fmt"
	"time"

	"go-service/externals"
	pdf "go-service/utils/pdf"
)

func GetStudentPDFReport(ctx context.Context, id string) ([]byte, string, error) {
	data, err := externals.GetStudentReport(ctx, id)
	if err != nil {
		return nil, "", err
	}
	b, err := pdf.GeneratePDF(ctx, "studentReport", data)
	if err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("%s_report_%s.pdf", id, time.Now().Format("20060102_150405"))
	return b, filename, nil
}
