package tools

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

// GreetingHelloArgsのテスト
func TestGreetingHelloArgs(t *testing.T) {
	// 名前が空の場合
	argsEmpty := GreetingHelloArgs{
		Name: "",
	}
	assert.Equal(t, "", argsEmpty.Name)
	
	// 名前が設定されている場合
	argsWithName := GreetingHelloArgs{
		Name: "テスト太郎",
	}
	assert.Equal(t, "テスト太郎", argsWithName.Name)
}
