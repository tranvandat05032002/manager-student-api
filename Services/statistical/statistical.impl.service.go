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
	f.SetCellValue(prefixSheet, "C1", "Tổng số môn học")
	f.SetCellValue(prefixSheet, "D1", "Tổng số tín chỉ")
	for i, termStat := range statisticalExportInput {
		row := i + 2
		f.SetCellValue(prefixSheet, fmt.Sprintf("A%d", row), fmt.Sprintf("%d-%d", termStat.TermFromYear, termStat.TermToYear))
		f.SetCellValue(prefixSheet, fmt.Sprintf("B%d", row), termStat.TermSemester)
		f.SetCellValue(prefixSheet, fmt.Sprintf("C%d", row), termStat.TotalSubjects)
		f.SetCellValue(prefixSheet, fmt.Sprintf("D%d", row), termStat.TotalCredits)
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Thống kê theo học kỳ.xlsx"); err != nil {
		return err
	}
	// Đặt header cho việc tải xuống file
	//c.Header("Content-Disposition", "attachment; filename=example.xlsx")
	//c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	//
	//// Gửi file Excel cho client
	//c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileContent.Bytes())
	return nil
}
