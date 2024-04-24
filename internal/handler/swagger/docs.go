// Package docs Code generated by swaggo/swag. DO NOT EDIT
package swagger

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
        "/v1/": {
            "get": {
                "summary": "Check server live",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/cake": {
            "put": {
                "summary": "Delete cake by ID",
                "responses": {
                    "200": {
                        "description": "deleted",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/cake/{id}": {
            "get": {
                "summary": "Get cake by ID",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dao.Cake"
                        }
                    }
                }
            }
        },
        "/v1/cakes": {
            "get": {
                "summary": "List all cakes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dao.Cake"
                            }
                        }
                    }
                }
            },
            "post": {
                "summary": "Search cakes by name and yum-factor",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dao.Cake"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dao.Cake": {
            "type": "object",
            "properties": {
                "comment": {
                    "description": "A comment about the cake, max 200 characters.",
                    "type": "string"
                },
                "id": {
                    "description": "Unique identifier for the cake.",
                    "type": "integer"
                },
                "image_url": {
                    "description": "URL to an image of the cake.",
                    "type": "string"
                },
                "name": {
                    "description": "Name of the cake, max 30 characters.",
                    "type": "string"
                },
                "yum_factor": {
                    "description": "Rating from 1 to 5 inclusive.",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
