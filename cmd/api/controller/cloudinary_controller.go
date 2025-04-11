package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type imageController struct {
	imageSvc service.ImageService
}

func NewImageController(di do.Injector) (api.Controller, error) {
	return &imageController{
		imageSvc: do.MustInvoke[service.ImageService](di),
	}, nil
}

func (c *imageController) Register(g *echo.Group) {
	g = g.Group("/image")
	g.POST("/upload", c.Upload) // POST /image/upload
}

func (c *imageController) Upload(e echo.Context) error {
	// <form-data name="file" type="file">
	fileHeader, err := e.FormFile("file")
	if err != nil {
		return e.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	folder := e.FormValue("folder") // optional – e.g. "products"

	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.imageSvc.UploadImage(ctx, fileHeader, folder)
	})
}
