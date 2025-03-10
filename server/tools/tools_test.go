package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// GreetingHelloToolのテスト
func TestGreeterFunctionality(t *testing.T) {
	// このテストは、Greeterインターフェースの基本機能をテストします
	// 実際のMCPサーバー統合テストはここでは行いません

	// モックGreeterインスタンス
	mockGreeter := &TestGreeter{
		defaultMessage: "こんにちは！",
	}

	// テスト
	greeting1, err := mockGreeter.GenerateGreeting("")
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは！", greeting1)

	greeting2, err := mockGreeter.GenerateGreeting("田中")
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは！ 田中さん！", greeting2)
}

// テスト用のGreeter実装
type TestGreeter struct {
	defaultMessage string
}

func (g *TestGreeter) GenerateGreeting(name string) (string, error) {
	if name == "" {
		return g.defaultMessage, nil
	}
	return g.defaultMessage + " " + name + "さん！", nil
}
