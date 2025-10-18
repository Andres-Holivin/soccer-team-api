package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var (
	cld        *cloudinary.Cloudinary
	once       sync.Once
	ctx        = context.Background()
	rootFolder string
)

func InitCloudinary() *cloudinary.Cloudinary {
	once.Do(func() {
		env := LoadEnv()
		rootFolder = env.CloudinaryRootFolder

		if env.CloudinaryURL == "" {
			log.Fatal("❌ CLOUDINARY_URL is not set in environment variables")
		}
		if rootFolder == "" {
			log.Fatal("❌ CLOUDINARY_ROOT_FOLDER is not set in environment variables")
		}

		client, err := cloudinary.NewFromURL(env.CloudinaryURL)
		if err != nil {
			log.Fatalf("Failed to initialize Cloudinary: %v", err)
		}

		cld = client
	})
	return cld
}

func UploadImage(fileHeader *multipart.FileHeader, subFolder string) (string, string, error) {
	cld := InitCloudinary()

	file, err := fileHeader.Open()
	if err != nil {
		return "", "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fullFolder := path.Join(rootFolder, subFolder)

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: fullFolder,
	})
	if err != nil {
		return "", "", fmt.Errorf("upload failed: %w", err)
	}

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

func DeleteImage(publicID string) error {
	cld := InitCloudinary()

	_, err := cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
		PublicIDs: []string{publicID},
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}
