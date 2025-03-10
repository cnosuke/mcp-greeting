package server

import (
	"testing"

	"github.com/cnosuke/mcp-greeting/config"
	"github.com/stretchr/testify/assert"
)

func TestGreetingServer_GenerateGreeting(t *testing.T) {
	// テストケース
	testCases := []struct {
		name           string
		defaultMessage string
		inputName      string
		expected       string
	}{
		{
			name:           "デフォルトメッセージのみ",
			defaultMessage: "こんにちは！",
			inputName:      "",
			expected:       "こんにちは！",
		},
		{
			name:           "名前付きの挨拶",
			defaultMessage: "こんにちは！",
			inputName:      "田中",
			expected:       "こんにちは！ 田中さん！",
		},
		{
			name:           "別のデフォルトメッセージ",
			defaultMessage: "Hello!",
			inputName:      "Smith",
			expected:       "Hello! Smithさん！",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用の設定を作成
			cfg := &config.Config{}
			cfg.Greeting.DefaultMessage = tc.defaultMessage

			// GreetingServerを初期化
			server, err := NewGreetingServer(cfg)
			assert.NoError(t, err)

			// 挨拶を生成
			greeting, err := server.GenerateGreeting(tc.inputName)
			
			// アサーション
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, greeting)
		})
	}
}
