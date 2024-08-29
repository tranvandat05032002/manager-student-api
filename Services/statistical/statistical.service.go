package statiscal

import (
	"gin-gonic-gom/Models"
)

type StatisticalService interface {
	StatisticalOfTerm(int, int) ([]Models.StatisticalOfTermRes, error)
	ExportStatisticalOfTerm([]Models.StatisticalExportInput) error
}
