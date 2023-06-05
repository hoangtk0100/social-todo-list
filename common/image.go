package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Image struct {
	Id        int    `json:"id" gorm:"column:id;"`
	Name      string `json:"name" gorm:"column:name;"`
	Url       string `json:"url" gorm:"column:url;"`
	Width     int    `json:"width" gorm:"column:width;"`
	Height    int    `json:"height" gorm:"column:height;"`
	CloudName string `json:"cloud_name,omitempty" gorm:"cloud_name"`
	Extension string `json:"extension,omitempty" gorm:"extension"`
}

func (Image) TableName() string { return "images" }

func (i *Image) Fulfill(domain string) {
	i.Url = fmt.Sprintf("%s/%s", domain, i.Url)
}

func (i *Image) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal data from DB:", value))
	}

	var img Image
	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}

	*i = img
	return nil
}

func (i *Image) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	return json.Marshal(i)
}

type Images []Image

func (i *Images) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal data from DB:", value))
	}

	var img []Image
	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}

	*i = img
	return nil
}

func (i *Images) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	return json.Marshal(i)
}
