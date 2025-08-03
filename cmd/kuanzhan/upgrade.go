/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "升级kuanzhan",
	Long:  `自动检查并升级kuanzhan到最新版本`,
	Run: func(cmd *cobra.Command, args []string) {
		release, err := getLatestVersion()
		if err != nil {
			fmt.Println("Failed to get latest version:", err)
			return
		}

		if Version == "dev" {
			fmt.Println("You are using the development version, skipping upgrade")
			return
		}

		log.Println("Latest version:", release.TagName)
		if semver.Compare(release.TagName, Version) > 0 {
			fmt.Println("New version found:", release.TagName)
			log.Printf("release % #v", pretty.Formatter(release))

			// 查找适合当前系统的资产
			asset := findCompatibleAsset(release.Assets)
			if asset == nil {
				fmt.Println("No compatible binary found for your system")
				return
			}

			fmt.Printf("Found compatible binary: %s\n", asset.Name)

			// 下载并更新
			if err := downloadAndUpdate(asset.BrowserDownloadURL); err != nil {
				fmt.Printf("Failed to update: %v\n", err)
				return
			}

			fmt.Println("Successfully updated to version:", release.TagName)
		} else {
			fmt.Println("You are using the latest version")
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only work when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getLatestVersion() (*Release, error) {
	resp, err := http.Get("https://api.github.com/repos/blksails/kuanzhan/releases")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dec = json.NewDecoder(resp.Body)

	var releases []Release
	err = dec.Decode(&releases)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found")
	}

	return &releases[0], nil
}

// findCompatibleAsset 根据当前系统和 CPU 架构查找兼容的二进制文件
func findCompatibleAsset(assets []Asset) *Asset {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// 构建目标文件名模式
	var targetPatterns []string

	// 常见的二进制文件命名模式
	patterns := []string{
		fmt.Sprintf("kuanzhan-%s-%s", goos, goarch),
		fmt.Sprintf("kuanzhan_%s_%s", goos, goarch),
		fmt.Sprintf("kuanzhan-%s-%s.exe", goos, goarch), // Windows
		fmt.Sprintf("kuanzhan_%s_%s.exe", goos, goarch), // Windows
	}

	// 添加带版本号的模式
	for _, pattern := range patterns {
		targetPatterns = append(targetPatterns, pattern)
		targetPatterns = append(targetPatterns, pattern+".tar.gz")
		targetPatterns = append(targetPatterns, pattern+".zip")
	}

	log.Printf("Looking for assets matching patterns: %v", targetPatterns)
	log.Printf("Current system: %s/%s", goos, goarch)

	for _, asset := range assets {
		log.Printf("Checking asset: %s", asset.Name)
		for _, pattern := range targetPatterns {
			if strings.Contains(strings.ToLower(asset.Name), strings.ToLower(pattern)) {
				log.Printf("Found compatible asset: %s", asset.Name)
				return &asset
			}
		}
	}

	return nil
}

// downloadAndUpdate 下载并更新二进制文件
func downloadAndUpdate(downloadURL string) error {
	fmt.Printf("Downloading from: %s\n", downloadURL)

	// 下载文件
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	fmt.Printf("Updating executable: %s\n", execPath)

	// 根据文件扩展名决定是否需要解压
	var reader io.Reader = resp.Body
	fileName := filepath.Base(downloadURL)

	if strings.HasSuffix(fileName, ".tar.gz") || strings.HasSuffix(fileName, ".tgz") {
		fmt.Println("Detected tar.gz file, extracting...")
		reader, err = extractTarGz(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to extract tar.gz: %w", err)
		}
	} else if strings.HasSuffix(fileName, ".zip") {
		fmt.Println("Detected zip file, extracting...")
		reader, err = extractZip(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to extract zip: %w", err)
		}
	}

	// 使用 selfupdate 进行更新
	err = selfupdate.Apply(reader, selfupdate.Options{
		TargetPath: execPath,
	})
	if err != nil {
		return fmt.Errorf("failed to apply update: %w", err)
	}

	return nil
}

// extractTarGz 从 tar.gz 文件中提取二进制文件
func extractTarGz(body io.Reader) (io.Reader, error) {
	// 创建 gzip 读取器
	gzr, err := gzip.NewReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	// 创建 tar 读取器
	tr := tar.NewReader(gzr)

	// 查找第一个可执行文件
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar header: %w", err)
		}

		// 检查是否是文件且名称包含 "kuanzhan"
		if header.Typeflag == tar.TypeReg && strings.Contains(header.Name, "kuanzhan") {
			fmt.Printf("Found binary in tar.gz: %s\n", header.Name)
			return tr, nil
		}
	}

	return nil, fmt.Errorf("no kuanzhan binary found in tar.gz")
}

// extractZip 从 zip 文件中提取二进制文件
func extractZip(body io.Reader) (io.Reader, error) {
	// 创建临时文件来存储 zip 内容
	tempFile, err := os.CreateTemp("", "kuanzhan-update-*.zip")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// 将响应体写入临时文件
	_, err = io.Copy(tempFile, body)
	if err != nil {
		tempFile.Close()
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tempFile.Close()

	// 重新打开文件用于 zip 读取
	zipReader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %w", err)
	}

	// 查找第一个可执行文件
	for _, file := range zipReader.File {
		if strings.Contains(file.Name, "kuanzhan") {
			fmt.Printf("Found binary in zip: %s\n", file.Name)
			rc, err := file.Open()
			if err != nil {
				zipReader.Close()
				return nil, fmt.Errorf("failed to open file in zip: %w", err)
			}
			// 返回一个包装的 reader，确保在读取完成后关闭 zipReader
			return &zipReaderWrapper{
				reader:    rc,
				zipReader: zipReader,
				closeOnce: false,
			}, nil
		}
	}

	zipReader.Close()
	return nil, fmt.Errorf("no kuanzhan binary found in zip")
}

// zipReaderWrapper 包装 zip 文件的 reader，确保正确关闭资源
type zipReaderWrapper struct {
	reader    io.ReadCloser
	zipReader *zip.ReadCloser
	closeOnce bool
}

func (w *zipReaderWrapper) Read(p []byte) (n int, err error) {
	return w.reader.Read(p)
}

func (w *zipReaderWrapper) Close() error {
	if !w.closeOnce {
		w.closeOnce = true
		w.reader.Close()
		return w.zipReader.Close()
	}
	return nil
}

type Release struct {
	URL             string    `json:"url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	HTMLURL         string    `json:"html_url"`
	ID              int64     `json:"id"`
	Author          Author    `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Immutable       bool      `json:"immutable"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []Asset   `json:"assets"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
	MentionsCount   int64     `json:"mentions_count"`
}

type Asset struct {
	URL                string    `json:"url"`
	ID                 int64     `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	Uploader           Author    `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int64     `json:"size"`
	Digest             string    `json:"digest"`
	DownloadCount      int64     `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}

type Author struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}
