package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"context"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type videoController struct {
	imageSvc service.ImageService
}

func NewVideoController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &videoController{
		imageSvc: do.MustInvoke[service.ImageService](di),
	}, nil
}

func (c *videoController) Register(g *echo.Group) {
	g = g.Group("/file")
	g.POST("/upload", c.Upload)
}

func (c *videoController) Upload(e echo.Context) error {
	// Parse the whole multipart form once. Echo re‑uses http.Request.ParseMultipartForm
	// so this is cached and cheap even for many fields.
	form, err := e.MultipartForm()
	if err != nil {
		return e.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid multipart form: " + err.Error(),
		})
	}

	// Preferred: <input type="file" name="files" multiple>
	files := form.File["files"]

	// Fallback: single <input type="file" name="file">
	if len(files) == 0 {
		if fh, err := e.FormFile("file"); err == nil {
			files = []*multipart.FileHeader{fh}
		}
	}

	if len(files) == 0 {
		return e.JSON(http.StatusBadRequest, echo.Map{
			"error": "no file uploaded",
		})
	}

	folder := e.FormValue("folder") // may be ""

	// Pass the slice straight through to the service
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.imageSvc.UploadImages(ctx, files, folder)
	})
}
