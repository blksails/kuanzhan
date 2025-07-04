package kuanzhan

import (
	_ "embed"
	"os"
	"strconv"
	"testing"
)

func TestClient_SignMethod(t *testing.T) {
	var client = NewClient("adde13Efcse", "helloWord")
	var params = map[string]string{
		"appKey":  "adde13Efcse",
		"url":     "https://www.baidu.com",
		"urlType": "w.url.cn",
	}

	sign := client.SignMethod(params)
	t.Log(sign)
}

// GetSiteIds
func TestClient_GetSiteIds(t *testing.T) {
	var client = NewClient(os.Getenv("KUAIZHAN_APP_KEY"), os.Getenv("KUAIZHAN_APP_SECRET"))
	client.SetDebug(true)
	resp, err := client.GetSiteIds()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Data.SiteIds)
}

// GetSiteInfo
func TestClient_GetSiteInfo(t *testing.T) {
	var client = NewClient(os.Getenv("KUAIZHAN_APP_KEY"), os.Getenv("KUAIZHAN_APP_SECRET"))

	client.SetDebug(true)
	var siteId = os.Getenv("KUAIZHAN_SITE_ID")
	siteIdInt, err := strconv.Atoi(siteId)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.GetSiteInfo(siteIdInt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

//go:embed test.html
var testHtml string

// ModifyPageJs
func TestClient_ModifyPageJs(t *testing.T) {
	var client = NewClient(os.Getenv("KUAIZHAN_APP_KEY"), os.Getenv("KUAIZHAN_APP_SECRET"))
	client.SetDebug(true)
	var siteId = os.Getenv("KUAIZHAN_SITE_ID")
	var pageId = os.Getenv("KUAIZHAN_PAGE_ID")
	siteIdInt, err := strconv.Atoi(siteId)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.ModifyPageJs(siteIdInt, pageId, testHtml, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)

	pageIdInt, err := strconv.Atoi(pageId)
	if err != nil {
		t.Fatal(err)
	}
	resp1, err := client.PublishPage(siteIdInt, pageIdInt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp1)
}

// GetPageName
func TestClient_GetPageName(t *testing.T) {
	var client = NewClient(os.Getenv("KUAIZHAN_APP_KEY"), os.Getenv("KUAIZHAN_APP_SECRET"))
	client.SetDebug(true)
	var siteId = os.Getenv("KUAIZHAN_SITE_ID")
	siteIdInt, err := strconv.Atoi(siteId)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.GetPageName(siteIdInt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

// BatchModifyPagePublishPageJs
func TestClient_BatchModifyPagePublishPageJs(t *testing.T) {
	var client = NewClient(os.Getenv("KUAIZHAN_APP_KEY"), os.Getenv("KUAIZHAN_APP_SECRET"))
	client.SetDebug(true)
	var siteId = os.Getenv("KUAIZHAN_SITE_ID")
	var pageId = os.Getenv("KUAIZHAN_PAGE_ID")
	siteIdInt, err := strconv.Atoi(siteId)
	if err != nil {
		t.Fatal(err)
	}
	pageIdInt, err := strconv.Atoi(pageId)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.BatchModifyPagePublishPageJs([]int{siteIdInt}, []int{pageIdInt}, testHtml, true, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
