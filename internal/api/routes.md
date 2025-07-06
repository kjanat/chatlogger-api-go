# [routes.go](routes.go)

Package `api` implements the REST API `router` and `routes` for the ChatLogger API.

This file defines the route setup for different API endpoints, including public
routes for chat plugins, authenticated routes for dashboard users, and admin routes.

```mermaid
classDiagram
    class Routes {
        +setupSwaggerRoutes(router *gin.Engine, version string, host string, port string)
        +addRoutes(router *gin.Engine, services *AppServices, jwtSecret string)
    }

    class SwaggerRoutes {
        +StaticFile("/docs/api.json", "./docs/OpenAPI_swagger.json")
        +StaticFile("/docs/api.yaml", "./docs/OpenAPI_swagger.yaml")
        +GET("/openapi/*any")
        +GET("/swagger/*any")
    }

    class PublicRoutes {
        +GET("/health")
        +GET("/version")
    }

    class AuthRoutes {
        +POST("/auth/login")
        +POST("/auth/register")
        +POST("/auth/logout")
    }

    class DashboardRoutes {
        +Group("/api/v1")
        +JWTAuth(jwtSecret)
    }

    class UserRoutes {
        +GET("/users/me")
        +PATCH("/users/me")
        +POST("/users/me/password")
    }

    class OrganizationRoutes {
        +GET("/orgs/me/apikeys")
        +POST("/orgs/me/apikeys")
        +DELETE("/orgs/me/apikeys/:id")
    }

    class ChatRoutes {
        +POST("/chats")
        +GET("/chats")
        +GET("/chats/:chatID")
        +PATCH("/chats/:chatID")
        +DELETE("/chats/:chatID")
    }

    class MessageRoutes {
        +GET("/chats/:chatID/messages")
        +GET("/analytics/messages")
    }

    class ExportRoutes {
        +POST("/exports")
        +GET("/exports")
        +GET("/exports/:id")
        +GET("/exports/:id/download")
        +POST("/exports/sync")
    }

    class PublicAPIRoutes {
        +POST("/api/v1/orgs/:slug/chats")
        +POST("/api/v1/orgs/:slug/chats/:chatID/messages")
    }

    Routes --> SwaggerRoutes
    Routes --> PublicRoutes
    Routes --> AuthRoutes
    Routes --> DashboardRoutes
    DashboardRoutes --> UserRoutes
    DashboardRoutes --> OrganizationRoutes
    DashboardRoutes --> ChatRoutes
    DashboardRoutes --> MessageRoutes
    DashboardRoutes --> ExportRoutes
    Routes --> PublicAPIRoutes
```
