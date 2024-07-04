package domain

import (
	"github.com/labstack/echo/v4"
)

// @Title Get example
// @Description Get example.
// @Param  exampleID  path  int  true  "Id of example" "120"
// @Success  200  object  string  "example JSON"
// @Failure  400  object  string  "example JSON"
// @Resource example
// @Route /api/example/{exampleID}/example [get]
func RegisterHandlers(e *echo.Echo, baseURL string) {
	//session.RegisterHandlersWithBaseURL(e, &session.ServerInterfaceImpl{}, baseURL)
}
