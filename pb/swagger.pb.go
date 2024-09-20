package pb 

const (
  Swagger = `{
  "swagger": "2.0",
  "info": {
    "title": "fxgateway.proto",
    "version": "version not set"
  },
  "tags": [
    {
      "name": "FxGateway"
    }
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/healthz": {
      "get": {
        "operationId": "FxGateway_HealthCheck",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/function-log/{FunctionName}": {
      "get": {
        "operationId": "FxGateway_GetLog",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "FunctionName",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/function/{FunctionName}": {
      "get": {
        "operationId": "FxGateway_GetMeta",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbFunction"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "FunctionName",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "FxGateway"
        ]
      },
      "delete": {
        "operationId": "FxGateway_Delete",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "FunctionName",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/function/{Service}": {
      "post": {
        "operationId": "FxGateway_Invoke",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "Service",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "Input": {
                  "type": "string",
                  "format": "byte"
                }
              }
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/functions": {
      "get": {
        "operationId": "FxGateway_List",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbFunctions"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "tags": [
          "FxGateway"
        ]
      },
      "post": {
        "operationId": "FxGateway_Deploy",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbCreateFunctionRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      },
      "put": {
        "operationId": "FxGateway_Update",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbCreateFunctionRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/info": {
      "get": {
        "operationId": "FxGateway_Info",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/scale-function": {
      "put": {
        "operationId": "FxGateway_ReplicaUpdate",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbScaleServiceRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/vm/{VMName}": {
      "delete": {
        "operationId": "FxGateway_DeleteVM",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "VMName",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/vms": {
      "post": {
        "operationId": "FxGateway_DeployVM",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbCreateVMRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    }
  },
  "definitions": {
    "pbCreateFunctionRequest": {
      "type": "object",
      "properties": {
        "Service": {
          "type": "string"
        },
        "Image": {
          "type": "string"
        },
        "EnvVars": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "Labels": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "Annotations": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "Constraints": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "Secrets": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "RegistryAuth": {
          "type": "string"
        },
        "Limits": {
          "$ref": "#/definitions/pbFunctionResources"
        },
        "Requests": {
          "$ref": "#/definitions/pbFunctionResources"
        },
        "MinReplicas": {
          "type": "integer",
          "format": "int32"
        },
        "MaxReplicas": {
          "type": "integer",
          "format": "int32"
        },
        "NodeName": {
          "type": "string"
        }
      }
    },
    "pbCreateVMRequest": {
      "type": "object",
      "properties": {
        "Instance": {
          "type": "string"
        },
        "Domain": {
          "type": "string"
        },
        "UserData": {
          "type": "string"
        },
        "Requests": {
          "$ref": "#/definitions/pbFunctionResources"
        }
      }
    },
    "pbFunction": {
      "type": "object",
      "properties": {
        "Name": {
          "type": "string"
        },
        "Image": {
          "type": "string"
        },
        "InvocationCount": {
          "type": "string",
          "format": "uint64"
        },
        "Replicas": {
          "type": "string",
          "format": "uint64"
        },
        "AvailableReplicas": {
          "type": "string",
          "format": "uint64"
        },
        "Annotations": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "Labels": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "pbFunctionResources": {
      "type": "object",
      "properties": {
        "Memory": {
          "type": "string"
        },
        "CPU": {
          "type": "string"
        },
        "GPU": {
          "type": "string"
        }
      }
    },
    "pbFunctions": {
      "type": "object",
      "properties": {
        "Functions": {
          "type": "array",
          "items": {
            "type": "object",
            "$ref": "#/definitions/pbFunction"
          }
        }
      }
    },
    "pbMessage": {
      "type": "object",
      "properties": {
        "Msg": {
          "type": "string"
        }
      }
    },
    "pbScaleServiceRequest": {
      "type": "object",
      "properties": {
        "ServiceName": {
          "type": "string"
        },
        "Replicas": {
          "type": "string",
          "format": "uint64"
        }
      }
    },
    "protobufAny": {
      "type": "object",
      "properties": {
        "@type": {
          "type": "string"
        }
      },
      "additionalProperties": {}
    },
    "rpcStatus": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        },
        "details": {
          "type": "array",
          "items": {
            "type": "object",
            "$ref": "#/definitions/protobufAny"
          }
        }
      }
    }
  }
}
`
)
