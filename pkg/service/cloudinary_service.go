package service

import (
	"SangXanh/pkg/common/api"
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
	"mime/multipart"
)

type ImageService interface {
	UploadImages(ctx context.Context, files []*multipart.FileHeader, folder string) (api.Response, error)
}

type imageService struct {
	cld *cloudinary.Cloudinary
}

func NewImageService(di do.Injector) (ImageService, error) {
	cld, err := do.Invoke[*cloudinary.Cloudinary](di)
	if err != nil {
		return nil, fmt.Errorf("init ImageService: %w", err)
	}
	return &imageService{cld: cld}, nil
}

func (s *imageService) UploadImages(ctx context.Context, files []*multipart.FileHeader, folder string) (api.Response, error) {
	type imgMeta struct {
		URL      string `json:"url"`
		PublicID string `json:"public_id"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
	}

	meta := make([]imgMeta, len(files))

	// one goroutine per file, fail fast on the first error
	g, gctx := errgroup.WithContext(ctx)

	for i, fh := range files {
		i, fh := i, fh // capture loop vars
		g.Go(func() error {
			src, err := fh.Open()
			if err != nil {
				return fmt.Errorf("open file %q: %w", fh.Filename, err)
			}
			defer src.Close()

			upParams := uploader.UploadParams{
				Folder:   folder,
				PublicID: uuid.NewString(), // random => no collisions
			}

			res, err := s.cld.Upload.Upload(gctx, src, upParams)
			if err != nil {
				return fmt.Errorf("cloudinary upload %q: %w", fh.Filename, err)
			}

			meta[i] = imgMeta{
				URL:      res.SecureURL,
				PublicID: res.PublicID,
				Width:    res.Width,
				Height:   res.Height,
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err // bubble up the first failure
	}

	return api.Success(meta), nil
}
