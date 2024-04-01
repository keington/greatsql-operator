package utils

import (
	"errors"
	"strings"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-18 18:21:36
 * @file: type_valid.go
 * @description: greatsql type valid
 */

var greatSqlTypeToAPI = map[string]string{
	"single":                    "Single",
	"replicaofGroupCluster":     "ReplicaofGroupCluster",
	"singlePrimaryGroupCluster": "SinglePrimaryGroupCluster",
	"multiPrimaryGroupCluster":  "MultiPrimaryGroupCluster",
}

// IsValidGreatSqlTypeToApi checks if the given greatSqlType is valid
func IsValidGreatSqlTypeToApi(greatSqlType string) (string, error) {
	api, ok := greatSqlTypeToAPI[strings.ToLower(greatSqlType)]
	if !ok {
		return "", errors.New("no corresponding API found for the given greatSqlType")
	}
	return api, nil
}

func Equal[T comparable](a, b T) bool {
	return a == b
}
