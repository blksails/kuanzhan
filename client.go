package kuanzhan

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type Client struct {
	BaseURL   string
	AppKey    string
	AppSecret string
	debug     bool
	impls     *impls
}

// NewClient creates a new client
func NewClient(appKey, appSecret string) *Client {
	return &Client{
		BaseURL:   "https://cloud.kuaizhan.com/api/v1",
		AppKey:    appKey,
		AppSecret: appSecret,
		impls: &impls{
			CreateSite: &methodImpl[SiteResponse, SiteRequest]{
				Path:   "/tbk/createSite",
				Method: "POST",
			},
			GetSiteIds: &methodImpl[GetSiteIdsResponse, GetSiteIdsRequest]{
				Path:   "/tbk/getSiteIds",
				Method: "POST",
			},
			GetPageIds: &methodImpl[GetPageIdsResponse, GetPageIdsRequest]{
				Path:   "/tbk/getPageIds",
				Method: "POST",
			},
			CreateSitePage: &methodImpl[CreateSitePageResponse, CreateSitePageRequest]{
				Path:   "/tbk/createSitePage",
				Method: "POST",
			},
			PublishSite: &methodImpl[PublishSiteResponse, PublishSiteRequest]{
				Path:   "/tbk/publishSite",
				Method: "POST",
			},
			PublishPage: &methodImpl[PublishPageResponse, PublishPageRequest]{
				Path:   "/tbk/publishPage",
				Method: "POST",
			},
			DeleteSitePage: &methodImpl[DeleteSitePageResponse, DeleteSitePageRequest]{
				Path:   "/tbk/deleteSitePage",
				Method: "POST",
			},
			UpdatePageName: &methodImpl[UpdatePageNameResponse, UpdatePageNameRequest]{
				Path:   "/tbk/updatePageName",
				Method: "POST",
			},
			GetPageName: &methodImpl[GetPageNameResponse, GetPageNameRequest]{
				Path:   "/tbk/getPageName",
				Method: "GET",
			},
			GetSiteInfo: &methodImpl[GetSiteInfoResponse, GetSiteInfoRequest]{
				Path:   "/tbk/getSiteInfo",
				Method: "POST",
			},
			ModifyPageJs: &methodImpl[ModifyPageJsResponse, ModifyPageJsRequest]{
				Path:   "/tbk/modifyPageJs",
				Method: "POST",
			},
			BatchModifyPagePublishPageJs: &methodImpl[BatchModifyPagePublishPageJsResponse, BatchModifyPagePublishPageJsRequest]{
				Path:   "/tbk/batchModifyPublishPageJs",
				Method: "POST",
			},
			OpenBusinessPackage: &methodImpl[OpenBusinessPackageResponse, OpenBusinessPackageRequest]{
				Path:   "/agent/openBusinessPackage",
				Method: "POST",
			},
			ChangeDomain: &methodImpl[ChangeDomainResponse, ChangeDomainRequest]{
				Path:   "/tbk/changeDomain",
				Method: "POST",
			},
			UpdateSiteInfo: &methodImpl[UpdateSiteInfoResponse, UpdateSiteInfoRequest]{
				Path:   "/tbk/updateSiteSetting",
				Method: "POST",
			},
		},
	}
}

// SetDebug 设置调试模式
func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

// SignMethod 生成API请求签名
// params: 包含所有请求参数的map（不包含sign参数）
// 返回: MD5签名字符串（32位十六进制）
func (c *Client) SignMethod(params map[string]string) string {
	// 1. 获取所有参数名并按ASCII顺序排序
	keys := make([]string, 0, len(params))
	for key := range params {
		if key != "sign" { // 排除sign参数
			keys = append(keys, key)
		}
	}
	sort.Strings(keys) // 按ASCII顺序排序

	// 2. 拼装参数名和参数值
	var builder strings.Builder
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString(params[key])
	}
	paramString := builder.String()

	// 3. 在参数字符串前后加上AppSecret
	signString := c.AppSecret + paramString + c.AppSecret

	// 4. 计算MD5摘要
	h := md5.New()
	h.Write([]byte(signString))
	// 5. 转换为32位十六进制字符串
	return hex.EncodeToString(h.Sum(nil))
}

// BuildSignedParams 构建带签名的参数map
// params: 原始请求参数
// 返回: 包含签名的完整参数map
func (c *Client) BuildSignedParams(params map[string]interface{}) map[string]string {
	// 复制原始参数
	signedParams := make(map[string]string)
	for k, v := range params {
		signedParams[k] = fmt.Sprintf("%v", v)
	}

	// 添加签名
	signedParams["sign"] = c.SignMethod(signedParams)
	signedParams["appKey"] = c.AppKey
	return signedParams
}

// 使用示例:
// client := NewClient("your_app_key", "your_app_secret")
// params := map[string]string{
//     "foo": "1",
//     "bar": "2",
//     "foo_bar": "3",
//     "foobar": "4",
// }
// signedParams := client.BuildSignedParams(params)
// // signedParams 现在包含了sign字段

// CreateSite
func (c *Client) CreateSite(siteName string, domain string, siteType string, httpsForward bool) (*SiteResponse, error) {

	return c.impls.CreateSite.Do(c, SiteRequest{
		SiteName:     siteName,
		Domain:       domain,
		SiteType:     siteType,
		HTTPSForward: httpsForward,
	})
}

// CreateSitePage
func (c *Client) CreateSitePage(siteId int, tpl string) (*CreateSitePageResponse, error) {
	return c.impls.CreateSitePage.Do(c, CreateSitePageRequest{
		SiteId: siteId,
		Tpl:    tpl,
	})
}

// GetSiteIds
func (c *Client) GetSiteIds() (*GetSiteIdsResponse, error) {
	return c.impls.GetSiteIds.Do(c, GetSiteIdsRequest{})
}

// GetPageIds
func (c *Client) GetPageIds(siteId int) (*GetPageIdsResponse, error) {
	return c.impls.GetPageIds.Do(c, GetPageIdsRequest{
		SiteId: siteId,
	})
}

// PublishSite
func (c *Client) PublishSite(siteId int) (*PublishSiteResponse, error) {
	return c.impls.PublishSite.Do(c, PublishSiteRequest{
		SiteId: siteId,
	})
}

// PublishPage
func (c *Client) PublishPage(siteId int, pageId int) (*PublishPageResponse, error) {
	return c.impls.PublishPage.Do(c, PublishPageRequest{
		SiteId: siteId,
		PageId: pageId,
	})
}

// UpdatePageName
func (c *Client) UpdatePageName(pageId int, pageName string) (*UpdatePageNameResponse, error) {
	return c.impls.UpdatePageName.PostJSON(c, UpdatePageNameRequest{
		PageId:   pageId,
		PageName: pageName,
	})
}

// DeleteSitePage
func (c *Client) DeleteSitePage(pageId int) (*DeleteSitePageResponse, error) {
	return c.impls.DeleteSitePage.Do(c, DeleteSitePageRequest{
		PageId: pageId,
	})
}

// GetPageName
func (c *Client) GetPageName(siteId int) (*GetPageNameResponse, error) {
	return c.impls.GetPageName.Do(c, GetPageNameRequest{
		SiteId: siteId,
	})
}

// GetSiteInfo
func (c *Client) GetSiteInfo(siteId int) (*GetSiteInfoResponse, error) {
	return c.impls.GetSiteInfo.Do(c, GetSiteInfoRequest{
		SiteId: siteId,
	})
}

// ModifyPageJs
func (c *Client) ModifyPageJs(siteId int, pageId string, content string, isEncryptContent bool) (*ModifyPageJsResponse, error) {
	return c.impls.ModifyPageJs.Do(c, ModifyPageJsRequest{
		SiteId:           siteId,
		PageId:           pageId,
		Content:          content,
		IsEncryptContent: isEncryptContent,
	})
}

// BatchModifyPagePublishPageJs
func (c *Client) BatchModifyPagePublishPageJs(siteIds []int, pageIds []int, content string, isSecure bool, taskId string) (*BatchModifyPagePublishPageJsResponse, error) {
	return c.impls.BatchModifyPagePublishPageJs.PostJSON(c, BatchModifyPagePublishPageJsRequest{
		SiteIds:  siteIds,
		PageIds:  pageIds,
		Content:  content,
		IsSecure: isSecure,
		TaskId:   taskId,
	})
}

// OpenBusinessPackage
func (c *Client) OpenBusinessPackage(businessType string, siteId int64, appId string, phoneNo string) (*OpenBusinessPackageResponse, error) {
	return c.impls.OpenBusinessPackage.Do(c, OpenBusinessPackageRequest{
		BusinessType: businessType,
		SiteId:       siteId,
		AppId:        appId,
		PhoneNo:      phoneNo,
	})
}

// ChangeDomain
func (c *Client) ChangeDomain(siteId int64, domain string, httpsForward bool) (*ChangeDomainResponse, error) {
	return c.impls.ChangeDomain.Do(c, ChangeDomainRequest{
		SiteId:       siteId,
		Domain:       domain,
		HTTPSForward: httpsForward,
	})
}

// UpdateSiteInfo
func (c *Client) UpdateSiteInfo(siteId int64, siteName string) (*UpdateSiteInfoResponse, error) {
	return c.impls.UpdateSiteInfo.Do(c, UpdateSiteInfoRequest{
		SiteId:   siteId,
		SiteName: siteName,
	})
}
