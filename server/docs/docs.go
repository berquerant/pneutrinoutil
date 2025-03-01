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
        "/debug": {
            "get": {
                "description": "debug info",
                "produces": [
                    "application/json"
                ],
                "summary": "debug info",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.SuccessResponse-handler_DebugResponseData"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "health check",
                "produces": [
                    "application/json"
                ],
                "summary": "health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.SuccessResponse-string"
                        }
                    }
                }
            }
        },
        "/proc": {
            "get": {
                "description": "list results of processes",
                "produces": [
                    "application/json"
                ],
                "summary": "list results",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.SuccessResponse-handler_ListResponseData"
                        }
                    }
                }
            },
            "post": {
                "description": "start a pneutrinoutil process with given arguments",
                "produces": [
                    "application/json"
                ],
                "summary": "start a process",
                "parameters": [
                    {
                        "type": "file",
                        "description": "musicxml",
                        "name": "score",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "[0, 100]%, default: 0",
                        "name": "enhanceBreathiness",
                        "in": "formData"
                    },
                    {
                        "type": "number",
                        "description": "default: 1.0",
                        "name": "formantShift",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "[2, 3, 4], default: 2",
                        "name": "inference",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "default: MERROW",
                        "name": "model",
                        "in": "formData"
                    },
                    {
                        "type": "number",
                        "description": "default: 0",
                        "name": "pitchShiftNsf",
                        "in": "formData"
                    },
                    {
                        "type": "number",
                        "description": "default: 0",
                        "name": "pitchShiftWorld",
                        "in": "formData"
                    },
                    {
                        "type": "number",
                        "description": "[0, 100]%, default: 0",
                        "name": "smoothFormant",
                        "in": "formData"
                    },
                    {
                        "type": "number",
                        "description": "[0, 100]%, default: 0",
                        "name": "smoothPitch",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "default: 0",
                        "name": "styleShift",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "202": {
                        "description": "new process started",
                        "schema": {
                            "$ref": "#/definitions/handler.SuccessResponse-string"
                        },
                        "headers": {
                            "string x-request-id": {
                                "type": "string",
                                "description": "request id, or just id"
                            }
                        }
                    },
                    "400": {
                        "description": "bad score",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "413": {
                        "description": "too big score",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/proc/{id}/config": {
            "get": {
                "description": "download pneutrinoutil config as json",
                "produces": [
                    "application/json"
                ],
                "summary": "download config",
                "parameters": [
                    {
                        "type": "string",
                        "description": "request id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.SuccessResponse-ctl_Config"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/proc/{id}/detail": {
            "get": {
                "description": "get process info",
                "produces": [
                    "application/json"
                ],
                "summary": "get process info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "request id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.SuccessResponse-handler_GetDetailResponseData"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/proc/{id}/log": {
            "get": {
                "description": "download process log file",
                "summary": "download log",
                "parameters": [
                    {
                        "type": "string",
                        "description": "request id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/proc/{id}/musicxml": {
            "get": {
                "description": "download musicxml file",
                "summary": "download musicxml",
                "parameters": [
                    {
                        "type": "string",
                        "description": "request id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/proc/{id}/wav": {
            "get": {
                "description": "download wav file generated by pneutrinoutil",
                "summary": "download wav",
                "parameters": [
                    {
                        "type": "string",
                        "description": "request id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/proc/{id}/world_wav": {
            "get": {
                "description": "download world wav file generated by pneutrinoutil",
                "summary": "download world wav",
                "parameters": [
                    {
                        "type": "string",
                        "description": "request id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ctl.Config": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "enhanceBreathiness": {
                    "type": "number"
                },
                "formantShift": {
                    "type": "number",
                    "default": 1
                },
                "inference": {
                    "type": "integer",
                    "default": 3
                },
                "model": {
                    "description": "NEUTRINO",
                    "type": "string",
                    "default": "MERROW"
                },
                "parallel": {
                    "type": "integer",
                    "default": 1
                },
                "pitchShiftNsf": {
                    "description": "NSF",
                    "type": "number"
                },
                "pitchShiftWorld": {
                    "description": "WORLD",
                    "type": "number"
                },
                "randomSeed": {
                    "type": "integer",
                    "default": 1234
                },
                "score": {
                    "description": "musicXML_to_label\nSuffix string ` + "`" + `yaml:\"suffix\"` + "`" + `\nProject settings",
                    "type": "string"
                },
                "smoothFormant": {
                    "type": "number"
                },
                "smoothPitch": {
                    "type": "number"
                },
                "styleShift": {
                    "type": "integer"
                },
                "thread": {
                    "type": "integer",
                    "default": 4
                }
            }
        },
        "handler.DebugResponseData": {
            "type": "object",
            "properties": {
                "routes": {}
            }
        },
        "handler.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "ok": {
                    "description": "false",
                    "type": "boolean"
                }
            }
        },
        "handler.GetDetailResponseData": {
            "type": "object",
            "properties": {
                "basename": {
                    "description": "original musicxml file name except extension",
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "rid": {
                    "description": "request id, or just id",
                    "type": "string"
                },
                "salt": {
                    "type": "integer"
                }
            }
        },
        "handler.ListResponseDataElement": {
            "type": "object",
            "properties": {
                "basename": {
                    "description": "original musicxml file name except extension",
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "rid": {
                    "description": "request id, or just id",
                    "type": "string"
                },
                "salt": {
                    "type": "integer"
                }
            }
        },
        "handler.SuccessResponse-ctl_Config": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/ctl.Config"
                },
                "ok": {
                    "description": "true",
                    "type": "boolean"
                }
            }
        },
        "handler.SuccessResponse-handler_DebugResponseData": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/handler.DebugResponseData"
                },
                "ok": {
                    "description": "true",
                    "type": "boolean"
                }
            }
        },
        "handler.SuccessResponse-handler_GetDetailResponseData": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/handler.GetDetailResponseData"
                },
                "ok": {
                    "description": "true",
                    "type": "boolean"
                }
            }
        },
        "handler.SuccessResponse-handler_ListResponseData": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.ListResponseDataElement"
                    }
                },
                "ok": {
                    "description": "true",
                    "type": "boolean"
                }
            }
        },
        "handler.SuccessResponse-string": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                },
                "ok": {
                    "description": "true",
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9101",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "pneutrinoutil API",
	Description:      "pneutrinoutil http server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
