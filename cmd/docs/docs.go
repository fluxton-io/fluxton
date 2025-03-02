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
            "name": "API Support",
            "url": "http://github.com/fluxton-io/fluxton",
            "email": "chief@fluxton.io"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/organizations/{organization_id}/users": {
            "get": {
                "description": "Get all users in an organization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "List all users in an organization",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Organization ID",
                        "name": "organization_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "content": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/resources.UserResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid input"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "422": {
                        "description": "Unprocessable entity"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            },
            "post": {
                "description": "Add a new user to an organization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Create a user in an organization",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Organization ID",
                        "name": "organization_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User ID JSON",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/organization_requests.MemberCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "content": {
                                            "$ref": "#/definitions/resources.UserResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "content": {
                                            "type": "object"
                                        },
                                        "success": {
                                            "type": "boolean"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "content": {
                                            "type": "object"
                                        },
                                        "success": {
                                            "type": "boolean"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "422": {
                        "description": "Unprocessable entity",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "content": {
                                            "type": "object"
                                        },
                                        "success": {
                                            "type": "boolean"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "content": {
                                            "type": "object"
                                        },
                                        "success": {
                                            "type": "boolean"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/organizations/{organization_id}/users/{user_id}": {
            "delete": {
                "description": "Remove a user from an organization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Delete a user from an organization",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Organization ID",
                        "name": "organization_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "User deleted"
                    },
                    "400": {
                        "description": "Invalid input"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "422": {
                        "description": "Unprocessable entity"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        }
    },
    "definitions": {
        "organization_requests.MemberCreateRequest": {
            "type": "object",
            "properties": {
                "user_id": {
                    "type": "string"
                }
            }
        },
        "resources.UserResponse": {
            "type": "object",
            "properties": {
                "bio": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "roleId": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "responses.Response": {
            "type": "object",
            "properties": {
                "content": {},
                "errors": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "success": {
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "fluxton.io/api",
	BasePath:         "/v2",
	Schemes:          []string{},
	Title:            "Fluxton API",
	Description:      "Fluxton is backend as-a-service platform that allows you to build, deploy, and scale applications without managing infrastructure.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
