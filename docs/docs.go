// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "zsc",
            "url": "https://github.com/zoushucai/random-images",
            "email": "zscmoyujian@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/random": {
            "get": {
                "description": "随机返回一张图片",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "随机返回一张图片",
                "parameters": [
                    {
                        "type": "string",
                        "description": "图像目录下的子目录",
                        "name": "sub",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "图像目录下的子目录",
                        "name": "sub",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "图像宽度, 默认为 0,则是原始大小",
                        "name": "width",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "图像高度, 默认为 0,则是原始大小",
                        "name": "height",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "图像索引, 默认为 -1,则是随机",
                        "name": "index",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/routers.ResponseImageInfo"
                        }
                    }
                }
            }
        },
        "/v1/upload": {
            "post": {
                "description": "该接口接收一个文件，将其保存到服务器，并将文件信息（包括 MD5 值、文件名、宽高等）存储在内存中",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "上传文件并保存到服务器",
                "parameters": [
                    {
                        "type": "file",
                        "description": "要上传的文件",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/routers.ResponseImageInfo"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "routers.ResponseImageInfo": {
            "type": "object",
            "properties": {
                "file": {
                    "type": "string"
                },
                "msg": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9113",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "API 文档",
	Description:      "API 详细的描述信息 ......",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
