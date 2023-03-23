package utils

import "github.com/tinode/chat/server/extra/store/model"

func ConvertJSON(content interface{}) (model.JSON, error) {
	var v model.JSON
	err := v.Scan(content)
	if err != nil {
		return nil, err
	}
	return v, nil
}
