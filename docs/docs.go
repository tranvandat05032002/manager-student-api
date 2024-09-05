// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "tags": [
        {
            "name": "Users",
            "description": "Quản lý Tài Khoản"
        }
    ],
    "paths": {
        "/user/create":  {
            "post": {
                "tags": ["Users"],
                "description": "Tạo mới một người dùng",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Thêm mới một người dùng",
                "parameters": [
                    {
                        "name": "user",
                        "in": "body",
                        "description": "Thông tin để tạo một người dùng",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tạo mới người dùng thành công",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "number",
                                    "example": 200
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Đăng ký thành công"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Lỗi dữ liệu đầu vào"
                    },
                    "404": {
                        "description": "Đã xảy ra trong quá trình chuyển đổi dữ liệu",
                    },
                    "500": {
                        "description": "Lỗi phía hệ thống"
                    }
                }
            }
        },
        "/user/login":  {
            "post": {
                "tags": ["Users"],
                "description": "Đăng nhập vào hệ thống",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Đăng nhập vào hệ thống",
                "parameters": [
                    {
                        "name": "user",
                        "in": "body",
                        "description": "Thông tin đăng nhập",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tạo mới người dùng thành công",
                        "schema": {
                            "$ref": "#/definitions/LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Lỗi dữ liệu đầu vào"
                    },
                    "404": {
                        "description": "Đã xảy ra trong quá trình chuyển đổi dữ liệu",
                    },
                    "500": {
                        "description": "Lỗi phía hệ thống"
                    }
                }
            }
        },
        "/user/me":  {
            "get": {
                "tags": ["Users"],
                "description": "Lấy thông tin cá nhân",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Lấy thông tin cá nhân",
                "parameters": [
                    {
                        "name": "Authorization",
                        "in": "header",
                        "description": "Bearer access token",
                        "required": true,
                        "type": "string",
                        "example": "Bearer <your-access-token>"
                     }
                ],
                "responses": {
                    "200": {
                        "description": "",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "number",
                                    "example": 200
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Get me thành công"
                                },
                                "data": {
                                    "$ref": "#/definitions/GetMeResponse"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Lỗi dữ liệu đầu vào"
                    },
                    "404": {
                        "description": "Đã xảy ra trong quá trình chuyển đổi dữ liệu",
                    },
                    "500": {
                        "description": "Lỗi phía hệ thống"
                    }
                }
            }
        },
        "/user/{user_id}":  {
            "get": {
                "tags": ["Users"],
                "description": "Lấy thông tin người dùng theo ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Lấy thông tin người dùng theo ID",
                "parameters": [
                    {
                        "name": "user_id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Unique ID of the user"
                    },
                    {
                        "name": "Authorization",
                        "in": "header",
                        "description": "Bearer access token",
                        "required": true,
                        "type": "string",
                        "example": "Bearer <your-access-token>"
                     }
                ],
                "responses": {
                    "200": {
                        "description": "",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "number",
                                    "example": 200
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Get user thành công"
                                },
                                "data": {
                                    "$ref": "#/definitions/GetMeResponse"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Lỗi dữ liệu đầu vào"
                    },
                    "404": {
                        "description": "Đã xảy ra trong quá trình chuyển đổi dữ liệu",
                    },
                    "500": {
                        "description": "Lỗi phía hệ thống"
                    }
                }
            }
        },
         "/user/me/update":  {
            "patch": {
                "tags": ["Users"],
                "description": "Cập nhật thông tin cá nhân",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Cập nhật thông tin cá nhân",
                "parameters": [
                    {
                        "name": "user",
                        "in": "body",
                        "description": "Thông tin cập nhật",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UserUpdate"
                        }
                    },
                    {
                        "name": "Authorization",
                        "in": "header",
                        "description": "Bearer access token",
                        "required": true,
                        "type": "string",
                        "example": "Bearer <your-access-token>"
                     }
                ],
                "responses": {
                    "200": {
                        "description": "Cập nhật thông tin thành công",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "number",
                                    "example": 200
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Cập nhật thông tin thành công"
                                },
                                "data": {
                                    "$ref": "#/definitions/GetMeResponse"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Lỗi dữ liệu đầu vào"
                    },
                    "404": {
                        "description": "Đã xảy ra trong quá trình chuyển đổi dữ liệu",
                    },
                    "500": {
                        "description": "Lỗi phía hệ thống"
                    }
                }
            }
        }
    },
    "definitions": {
        "LoginResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "integer",
                    "example": 200
                },
                "message": {
                    "type": "string",
                    "example": "Đăng nhập thành công!"
                },
                "data": {
                    "type": "object",
                    "properties": {
                        "access_token": {
                            "type": "string",
                            "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUzNzcyMDgsImlhdCI6MTcyNTM3NTQwOCwicm9sZSI6InN0dWRlbnQiLCJ1c2VySUQiOiI2NmNjNGJjZGQ5N2EyYjliMGE3YzUxOWIifQ.OMUJhGatDjcbp7Q4c2be3olz4Mvq4XxVXfz5LqmBD3w"
                        },
                        "refresh_token": {
                            "type": "string",
                            "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU2MzQ2MDgsImlhdCI6MTcyNTM3NTQwOCwicm9sZSI6InN0dWRlbnQiLCJ1c2VySUQiOiI2NmNjNGJjZGQ5N2EyYjliMGE3YzUxOWIifQ.fZPqkg7yhs2pPNI6YbTUCBsOE2tqtgejH95S_O9UhdU"
                        }
                    },
                    "required": ["access_token", "refresh_token"]
                }
            },
            "required": ["status", "message", "data"]
        },
        "RegisterRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "example": "Trần Văn Đạt"
                },
                "email": {
                    "type": "string",
                    "example": "tranvandatevondev0503@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "password123"
                },
                "phone": {
                    "type": "string",
                    "example": "0987654321"
                },
                "role_type": {
                    "type": "string",
                    "enum": ["student", "teacher", "admin"],
                    "description": "Role của tài khoản",
                    "example": "student"
                },
                "avatar": {
                    "type": "string",
                    "example": "https://images2.thanhnien.vn/528068263637045248/2024/1/25/c3c8177f2e6142e8c4885dbff89eb92a-65a11aeea03da880-1706156293184503262817.jpg"
                }
            },
            "required": ["email", "password", "name"]
        },
        "LoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "tranvandatevondev0503@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "password123"
                }
            },
            "required": ["email", "password"]
        },
        "GetMeResponse": {
            "type": "object",
            "properties": {
                "Id": {
                    "type": "string",
                    "example": "66d1858ecae5906ce4f4df98"
                },
                "MajorId": {
                    "type": "string",
                    "nullable": true,
                    "example": null
                },
                "email": {
                    "type": "string",
                    "example": "admin@meteor.com"
                },
                "password": {
                    "type": "string",
                    "example": ""
                },
                "role_type": {
                    "type": "string",
                    "example": "admin"
                },
                "phone": {
                    "type": "string",
                    "example": "076852312"
                },
                "name": {
                    "type": "string",
                    "example": "Admin"
                },
                "avatar": {
                    "type": "string",
                    "format": "uri",
                    "example": "http://localhost:3000/static/images/f8d0eabf-0c2c-4e28-ba77-4c5ba77fbe61.jpg"
                },
                "gender": {
                    "type": "integer",
                    "example": 1
                },
                "department": {
                    "type": "string",
                    "example": ""
                },
                "date_of_birth": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2024-08-11T17:00:00Z"
                },
                "enrollment_date": {
                    "type": "string",
                    "format": "date-time",
                    "example": "0001-01-01T00:00:00Z"
                },
                "hire_date": {
                    "type": "string",
                    "format": "date-time",
                    "example": "0001-01-01T00:00:00Z"
                },
                "address": {
                    "type": "string",
                    "example": ""
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2024-08-30T08:40:46.206Z"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2024-08-30T09:58:37.356Z"
                }
            }
        },
        "UserUpdate": {
            "type": "object",
            "properties": {
                "major_name": {
                    "type": "string",
                    "example": "Computer Science"
                },
                "email": {
                    "type": "string",
                    "example": "admin@meteor.com"
                },
                "phone": {
                    "type": "string",
                    "example": "0768523123"
                },
                "name": {
                    "type": "string",
                    "example": "Admin"
                },
                "avatar": {
                    "type": "string",
                    "format": "uri",
                    "example": "https://images2.thanhnien.vn/528068263637045248/2024/1/25/c3c8177f2e6142e8c4885dbff89eb92a-65a11aeea03da880-1706156293184503262817.jpg"
                },
                "gender": {
                    "type": "integer",
                    "example": 1
                },
                "department": {
                    "type": "string",
                    "example": "Engineering"
                },
                "date_of_birth": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2024-08-30T19:13:03.972Z"
                },
                "enrollment_date": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2024-08-30T19:13:03.972Z"
                },
                "hire_date": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2024-08-30T19:13:03.972Z"
                },
                "address": {
                    "type": "string",
                    "example": "Lê Lợi"
                }
            },
            "required": [
                "email",
                "role_type",
                "name",
                "avatar",
                "gender",
                "date_of_birth",
                "address"
            ]
        }
    }
}`


// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:4000",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Manager Student Service API",
	Description:      "Manager Student service API in Go using Gin Framework",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
