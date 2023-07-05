package entity

import "errors"

var (
	ErrFileTooLarge           = errors.New("file too large")
	ErrFileNotImage           = errors.New("file is not image")
	ErrFileMissing            = errors.New("file is missing")
	ErrCannotReadFile         = errors.New("cannot read file")
	ErrCannotSaveFile         = errors.New("cannot save file")
	ErrCannotGetFileDimension = errors.New("cannot get file dimension")
)
