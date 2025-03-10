package server

import (
	"testing"

	"github.com/cnosuke/mcp-greeting/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// TestNewGreetingServer - GreetingServerの初期化テスト
func TestNewGreetingServer(t *testing.T) {
	// テスト用のロガーを設定
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// テスト用の設定
	cfg := &config.Config{}
	cfg.Greeting.DefaultMessage = "テスト用挨拶"

	// サーバーの作成
	server, err := NewGreetingServer(cfg)

	// アサーション
	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.Equal(t, "テスト用挨拶", server.DefaultMessage)
}

// TestSetupServerComponents - サーバーのセットアップロジックをテスト
func TestSetupServerComponents(t *testing.T) {
	// テスト用のロガーを設定
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// テスト用の設定
	cfg := &config.Config{}
	cfg.Greeting.DefaultMessage = "テスト用挨拶"

	// サーバーを作成してテスト
	greetingServer, err := NewGreetingServer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, greetingServer)

	// 挨拶生成機能のテスト
	greeting, err := greetingServer.GenerateGreeting("テストユーザー")
	assert.NoError(t, err)
	assert.Equal(t, "テスト用挨拶 テストユーザーさん！", greeting)
}
