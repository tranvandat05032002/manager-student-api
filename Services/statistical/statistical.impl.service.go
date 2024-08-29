package statiscal

import (
	"context"
	"fmt"
	"gin-gonic-gom/Models"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatisticalImplementService struct {
	termcollection *mongo.Collection
	ctx            context.Context
}

func NewStatisticalService(termcolecttion *mongo.Collection, ctx context.Context) StatisticalService {
	return &StatisticalImplementService{
		termcollection: termcolecttion,
		ctx:            ctx,
	}
}
func (a *StatisticalImplementService) StatisticalOfTerm(page, limit int) ([]Models.StatisticalOfTermRes, error) {
	var statisticalOfTermRes []Models.StatisticalOfTermRes
	skip := (page - 1) * limit
	pipeline := bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "Subjects",
				"localField":   "_id",
				"foreignField": "term_id",
				"as":           "subjects",
			},
		},
		bson.M{
			"$project": bson.M{
				"created_at": 0,
				"updated_at": 0,
			},
		},
		bson.M{
			"$addFields": bson.M{
				"total_credits":  bson.M{"$sum": "$subjects.credits"},
				"total_subjects": bson.M{"$size": "$subjects"},
			},
		},
		bson.M{
			"$skip": skip,
		},
		bson.M{
			"$limit": limit,
		},
	}
	cursor, err := a.termcollection.Aggregate(a.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(a.ctx)

	if err = cursor.All(a.ctx, &statisticalOfTermRes); err != nil {
		return nil, err
	}
	//// export file PDF
	//pdf := gofpdf.New("P", "mm", "A4", "")
	//pdf.AddPage()
	//pdf.SetFont("Arial", "B", 16)
	//pdf.Cell(40, 10, "Term Statistics")
	//
	//pdf.SetFont("Arial", "", 12)
	//for _, termStat := range statisticalOfTermRes {
	//	termID := termStat.ID
	//	totalCredits := termStat.TotalCredits
	//	totalSubjects := termStat.TotalSubjects
	//	pdf.Ln(10)
	//	pdf.Cell(40, 10, "Mã học kỳ: "+termID.Hex())
	//	pdf.Ln(6)
	//	pdf.Cell(40, 10, "Tổng số tín chỉ: "+strconv.Itoa(totalCredits))
	//	pdf.Ln(6)
	//	pdf.Cell(40, 10, "Tổng số môn học: "+strconv.Itoa(totalSubjects))
	//	pdf.Ln(10)
	//}
	//fileName := fmt.Sprintf("%s.pdf", statisticalOfTermRes[0].ID.Hex()) // Sử dụng ID của tài liệu đầu tiên
	//c.Header("Content-Type", "application/pdf")
	//c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	//err = pdf.Output(c.Writer)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return nil, err
	//}
	return statisticalOfTermRes, nil
}
func (a *StatisticalImplementService) ExportStatisticalOfTerm(statisticalExportInput []Models.StatisticalExportInput) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return err
	}
	prefixSheet := "Sheet1"
	f.SetCellValue(prefixSheet, "A1", "Năm học")
	f.SetCellValue(prefixSheet, "B1", "Học kỳ")
	f.SetCellValue(prefixSheet, "C1", "Ngày bắt đầu")
	f.SetCellValue(prefixSheet, "D1", "Ngày kết thúc")
	f.SetCellValue(prefixSheet, "E1", "Tổng số môn học")
	f.SetCellValue(prefixSheet, "F1", "Tổng số tín chỉ")
	for i, termStat := range statisticalExportInput {
		row := i + 2
		//startDateConvert, _ := utils.ConvertISOToDate(termStat.StartDate.String())
		//endDateConvert, _ := utils.ConvertISOToDate(termStat.EndDate.String())
		f.SetCellValue(prefixSheet, fmt.Sprintf("A%d", row), termStat.TermAcademicYear)
		f.SetCellValue(prefixSheet, fmt.Sprintf("B%d", row), termStat.TermSemester)
		f.SetCellValue(prefixSheet, fmt.Sprintf("C%d", row), termStat.StartDate)
		f.SetCellValue(prefixSheet, fmt.Sprintf("D%d", row), termStat.EndDate)
		f.SetCellValue(prefixSheet, fmt.Sprintf("E%d", row), termStat.TotalSubjects)
		f.SetCellValue(prefixSheet, fmt.Sprintf("F%d", row), termStat.TotalCredits)
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Thống kê theo học kỳ.xlsx"); err != nil {
		return err
	}
	// return file download
	//newFileName := uuid.New().String()
	//newFileNamePDFWithExt := newFileName + ".pdf"
	//path := filepath.Join("uploads/files", newFileNamePDFWithExt)
	// lưu file vào trong uploads/images
	//if ctx.SaveUploadedFile(file, path); err != nil {
	//	ctx.String(http.StatusInternalServerError, fmt.Sprintln("Upload image error: %s", err.Error()))
	//	return
	//}
	//// Tạo URL cho hình ảnh đã upload
	//imageURL := GeneratorURLImage(ctx, newFileNameWithExt)
	//err = mediaController.MediaService.Upload(imageURL)
	//if err != nil {
	//	common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, err.Error())
	//	return
	//}
	return nil
}
