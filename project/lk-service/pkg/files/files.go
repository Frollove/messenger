package files

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/nickalie/go-webpbin"
	"image"
	"image/jpeg"
	"image/png"
	"lk-service/pkg/custom_errors"
	"mime/multipart"
	"os"
	"strings"
)

func DownloadAvatars(uid int64, file *multipart.FileHeader) (path string, err error) {
	fileNameArr := strings.Split(file.Filename, ".")
	var isJPEG, isPNG bool
	if fileNameArr[1] == "jpeg" || fileNameArr[1] == "jpg" {
		isJPEG = true
	} else if fileNameArr[1] == "png" {
		isPNG = true
	}

	pathFilename := fileNameArr[0] + ".webp"

	err = os.MkdirAll(fmt.Sprintf("C:/NewServer/domains/byTrip/chat-service/src/user/%d", uid), 0750)
	if err != nil {
		return "", fmt.Errorf(fmt.Errorf("mkdir: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	path = fmt.Sprintf("C:/NewServer/domains/byTrip/chat-service/src/user/%d/%s", uid, strings.ReplaceAll(pathFilename, " ", "_"))

	fileOpened, err := file.Open()
	if err != nil {
		return "", fmt.Errorf(fmt.Errorf("open file: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}
	defer fileOpened.Close()

	newFile, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf(fmt.Errorf("create file: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}
	defer newFile.Close()

	var input image.Image
	if isPNG {
		input, err = png.Decode(fileOpened)
		if err != nil {
			return "", fmt.Errorf(fmt.Errorf("jpeg decode: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}
	}

	if isJPEG {
		input, err = jpeg.Decode(fileOpened)
		if err != nil {
			return "", fmt.Errorf(fmt.Errorf("jpeg decode: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}
	}

	input = resize.Resize(608, 608, input, resize.Bicubic)

	err = webpbin.NewCWebP().Quality(100).InputImage(input).Output(newFile).Run()
	if err != nil {
		return "", fmt.Errorf(fmt.Errorf("new CWebP: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	path = strings.ReplaceAll(path, "C:/NewServer/domains/byTrip/chat-service/src/user", "https://img.web-gen.ru/user")

	return
}
