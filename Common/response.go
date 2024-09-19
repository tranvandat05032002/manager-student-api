package common

/*
Cấu trúc của một response:

	{
		status: 200,
		message: "Success!",
		data: [
				{
					id: .....
					name: ....
					.......
				}
			  ]
	}
*/
type successRes struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Page    interface{} `json:"page,omitempty"`
	Limit   interface{} `json:"limit,omitempty"`
	Total   interface{} `json:"total,omitempty"`
}

func NewSuccessResponse(status int, message string, data interface{}, total int, page interface{}, limit interface{}) *successRes {
	return &successRes{
		Status:  status,
		Message: message,
		Data:    data,
		Page:    page,
		Limit:   limit,
		Total:   total,
	}
}
func SimpleSuccessResponse(status int, message string, data interface{}) *successRes {
	return &successRes{
		Status:  status,
		Message: message,
		Data:    data,
		Total:   nil,
		Page:    nil,
		Limit:   nil,
	}
}
