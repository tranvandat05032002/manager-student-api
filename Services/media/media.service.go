package media

import "gin-gonic-gom/Models"

type MediaService interface {
	Upload(string) error
	UploadExcelDataUser([]Models.UserModel) error
}
