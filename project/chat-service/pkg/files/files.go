package files

import (
	"chat-service/pkg/custom_errors"
	"fmt"
	"github.com/nickalie/go-webpbin"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
)

func DownloadFiles(chatID int64, files []*multipart.FileHeader) (map[string]string, error) {
	paths := make(map[string]string)

	for _, file := range files {
		messageType := "file"

		fileNameArr := strings.Split(file.Filename, ".")
		matchedJPEG, err := regexp.Match("jpeg|jpg", []byte(fileNameArr[1]))
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("match: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}
		matchedPNG, err := regexp.Match("png", []byte(fileNameArr[1]))
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("match: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}

		pathFilename := file.Filename
		if matchedJPEG || matchedPNG {
			pathFilename = fileNameArr[0] + ".webp"
			messageType = "picture"
		}

		err = os.MkdirAll(fmt.Sprintf("src/chat/files/%d", chatID), 0750)
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("mkdir: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}

		path := fmt.Sprintf("src/chat/files/%d/%s", chatID, strings.ReplaceAll(pathFilename, " ", "_"))
		paths[path] = messageType

		fileOpened, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("open file: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}

		newFile, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("create file: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}

		if matchedPNG || matchedJPEG {
			var input image.Image
			if matchedPNG {
				input, err = png.Decode(fileOpened)
				if err != nil {
					return nil, fmt.Errorf(fmt.Errorf("jpeg decode: %w", err).Error()+": %w", custom_errors.ErrInternal)
				}
			}

			if matchedJPEG {
				input, err = jpeg.Decode(fileOpened)
				if err != nil {
					return nil, fmt.Errorf(fmt.Errorf("jpeg decode: %w", err).Error()+": %w", custom_errors.ErrInternal)
				}
			}

			err = webpbin.NewCWebP().Quality(100).InputImage(input).Output(newFile).Run()
			if err != nil {
				return nil, fmt.Errorf(fmt.Errorf("new CWebP: %w", err).Error()+": %w", custom_errors.ErrInternal)
			}

		} else {
			_, err = io.Copy(newFile, fileOpened)
			if err != nil {
				return nil, fmt.Errorf(fmt.Errorf("copy file: %w", err).Error()+": %w", custom_errors.ErrInternal)
			}
		}

		fileOpened.Close()
		newFile.Close()
	}

	return paths, nil
}
