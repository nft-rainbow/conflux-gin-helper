package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/pkg/errors"
)

func Md5File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return Md5(file)
}

func Md5(src io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetFileType(fileName string, fileReader io.Reader) (*types.Type, error) {
	var fileContent bytes.Buffer
	if _, err := io.CopyN(&fileContent, fileReader, 100); err != nil {
		return nil, err
	}

	fileType, err := filetype.Match(fileContent.Bytes())
	if err != nil {
		return nil, errors.WithMessage(err, "failed get file type from content")
	}

	if fileType != filetype.Unknown {
		return &fileType, nil
	}

	tmp := strings.Split(fileName, ".")
	return &types.Type{
		Extension: tmp[len(tmp)-1],
	}, nil
}
