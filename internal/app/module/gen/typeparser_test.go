package gen

import (
	"testing"
)

func Test_parseLineTag(t *testing.T) {
	value := `gorm:"not null;PRIMARY_KEY;comment:'用户uid'" json:"uid"`
	tags := parseLineTag(value)
	t.Log(tags)
}
