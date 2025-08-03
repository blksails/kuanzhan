package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"
)

func TestFindCompatibleAsset(t *testing.T) {
	// 创建测试资产，使用当前系统的架构
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	assets := []Asset{
		{
			Name:               fmt.Sprintf("kuanzhan-%s-%s", currentOS, currentArch),
			BrowserDownloadURL: fmt.Sprintf("https://example.com/kuanzhan-%s-%s", currentOS, currentArch),
		},
		{
			Name:               "kuanzhan-linux-amd64",
			BrowserDownloadURL: "https://example.com/kuanzhan-linux-amd64",
		},
		{
			Name:               "kuanzhan-windows-amd64.exe",
			BrowserDownloadURL: "https://example.com/kuanzhan-windows-amd64.exe",
		},
		{
			Name:               fmt.Sprintf("kuanzhan-%s-%s.tar.gz", currentOS, currentArch),
			BrowserDownloadURL: fmt.Sprintf("https://example.com/kuanzhan-%s-%s.tar.gz", currentOS, currentArch),
		},
	}

	// 测试查找兼容资产
	asset := findCompatibleAsset(assets)
	if asset == nil {
		t.Fatal("Expected to find compatible asset, but got nil")
	}

	// 验证找到的资产是否匹配当前系统
	expectedOS := runtime.GOOS
	expectedArch := runtime.GOARCH

	if !contains(asset.Name, expectedOS) {
		t.Errorf("Expected asset to contain OS '%s', but got '%s'", expectedOS, asset.Name)
	}

	if !contains(asset.Name, expectedArch) {
		t.Errorf("Expected asset to contain arch '%s', but got '%s'", expectedArch, asset.Name)
	}

	t.Logf("Found compatible asset: %s", asset.Name)
}

func TestFindCompatibleAssetNoMatch(t *testing.T) {
	// 创建不匹配的测试资产
	assets := []Asset{
		{
			Name:               "kuanzhan-linux-arm64",
			BrowserDownloadURL: "https://example.com/kuanzhan-linux-arm64",
		},
		{
			Name:               "kuanzhan-windows-amd64.exe",
			BrowserDownloadURL: "https://example.com/kuanzhan-windows-amd64.exe",
		},
	}

	// 测试查找兼容资产（应该返回 nil）
	asset := findCompatibleAsset(assets)
	if asset != nil {
		t.Errorf("Expected no compatible asset, but got: %s", asset.Name)
	}
}

func TestExtractTarGz(t *testing.T) {
	// 创建测试 tar.gz 文件
	tempFile, err := os.CreateTemp("", "test-*.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 创建 gzip writer
	gw := gzip.NewWriter(tempFile)
	defer gw.Close()

	// 创建 tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 添加一个测试文件
	content := []byte("test binary content")
	header := &tar.Header{
		Name: "kuanzhan-test",
		Mode: 0755,
		Size: int64(len(content)),
	}
	err = tw.WriteHeader(header)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tw.Write(content)
	if err != nil {
		t.Fatal(err)
	}

	// 关闭 writers
	tw.Close()
	gw.Close()
	tempFile.Close()

	// 重新打开文件进行测试
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// 测试解压
	reader, err := extractTarGz(file)
	if err != nil {
		t.Fatal(err)
	}

	// 读取解压后的内容
	extractedContent, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	if string(extractedContent) != string(content) {
		t.Errorf("Expected content '%s', got '%s'", string(content), string(extractedContent))
	}
}

func TestExtractZip(t *testing.T) {
	// 创建测试 zip 文件
	tempFile, err := os.CreateTemp("", "test-*.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 创建 zip writer
	zw := zip.NewWriter(tempFile)
	defer zw.Close()

	// 添加一个测试文件
	content := []byte("test binary content")
	file, err := zw.Create("kuanzhan-test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write(content)
	if err != nil {
		t.Fatal(err)
	}

	// 关闭 writer
	zw.Close()
	tempFile.Close()

	// 重新打开文件进行测试
	zipFile, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer zipFile.Close()

	// 测试解压
	reader, err := extractZip(zipFile)
	if err != nil {
		t.Fatal(err)
	}

	// 读取解压后的内容
	extractedContent, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	if string(extractedContent) != string(content) {
		t.Errorf("Expected content '%s', got '%s'", string(content), string(extractedContent))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
