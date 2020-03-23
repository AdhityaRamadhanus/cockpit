package keygenerator

import "fmt"

func NationalLevelRedisKey() string {
	return fmt.Sprintf("indonesia:national-level")
}

func ProvincialLevelRedisKey(province string) string {
	return fmt.Sprintf("%s:provincial-level", province)
}
