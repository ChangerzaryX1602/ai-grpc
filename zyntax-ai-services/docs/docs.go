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
    "paths": {
        "/api/v1/ask/": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Ask a question to the AI service",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Ask"
                ],
                "summary": "Ask a question",
                "parameters": [
                    {
                        "description": "Question to ask",
                        "name": "question",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ask.Ask"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/ask/history": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all histories",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Get all histories",
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a history",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Create a history",
                "parameters": [
                    {
                        "description": "History to create",
                        "name": "history",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ask.History"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/ask/history/messages/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get history messages by history ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Get history messages by history ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "History ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/ask/history/{id}": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update a history",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Update a history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "History ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "History to update",
                        "name": "history",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ask.History"
                        }
                    }
                ],
                "responses": {}
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete a history",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Delete a history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "History ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/auth/": {
            "post": {
                "description": "https://oauth.kku.ac.th/authorize?response_type=code\u0026client_id=e8fdb4894be17a3a\u0026redirect_uri=http://localhost:8080/api/v1/swagger/index.html",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Code from oauth",
                        "name": "code",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Oauth"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/files": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all files",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Get all files",
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Create a new file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/files/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a file by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Get a file by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "File ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update a file",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Update a file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "File ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "File object",
                        "name": "file",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.File"
                        }
                    }
                ],
                "responses": {}
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete a file",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Delete a file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "File ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "ask.Ask": {
            "type": "object",
            "properties": {
                "history_id": {
                    "type": "integer"
                },
                "question": {
                    "type": "string"
                }
            }
        },
        "ask.History": {
            "type": "object",
            "properties": {
                "place_holder": {
                    "type": "string"
                }
            }
        },
        "models.File": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                }
            }
        },
        "models.Oauth": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{"http", "https"},
	Title:            "Zyntax ai API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
