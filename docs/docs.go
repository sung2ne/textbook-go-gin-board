package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "게시판 API 서버입니다.",
        "title": "Go Board API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api/v1",
    "paths": {}
}`

// SwaggerInfo holds exported Swagger Info
var SwaggerInfo = &swag.Spec{
    Version:          "1.0",
    Host:             "localhost:8080",
    BasePath:         "/api/v1",
    // ...
}

func init() {
    swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
