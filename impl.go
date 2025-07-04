package kuanzhan

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

type methodImpl[R, Q any] struct {
	Path   string
	Method string
}

type impls struct {
	CreateSite                   *methodImpl[SiteResponse, SiteRequest]
	GetSiteIds                   *methodImpl[GetSiteIdsResponse, GetSiteIdsRequest]
	GetPageIds                   *methodImpl[GetPageIdsResponse, GetPageIdsRequest]
	PublishSite                  *methodImpl[PublishSiteResponse, PublishSiteRequest]
	PublishPage                  *methodImpl[PublishPageResponse, PublishPageRequest]
	UpdatePageName               *methodImpl[UpdatePageNameResponse, UpdatePageNameRequest]
	DeleteSitePage               *methodImpl[DeleteSitePageResponse, DeleteSitePageRequest]
	GetPageName                  *methodImpl[GetPageNameResponse, GetPageNameRequest]
	CreateSitePage               *methodImpl[CreateSitePageResponse, CreateSitePageRequest]
	GetSiteInfo                  *methodImpl[GetSiteInfoResponse, GetSiteInfoRequest]
	ModifyPageJs                 *methodImpl[ModifyPageJsResponse, ModifyPageJsRequest]
	BatchModifyPagePublishPageJs *methodImpl[BatchModifyPagePublishPageJsResponse, BatchModifyPagePublishPageJsRequest]
	OpenBusinessPackage          *methodImpl[OpenBusinessPackageResponse, OpenBusinessPackageRequest]
	ChangeDomain                 *methodImpl[ChangeDomainResponse, ChangeDomainRequest]
	UpdateSiteInfo               *methodImpl[UpdateSiteInfoResponse, UpdateSiteInfoRequest]
}

func (m *methodImpl[R, Q]) Do(client *Client, params Q) (resp *R, err error) {
	var (
		apiUrl  = client.BaseURL + m.Path
		pParams = make(map[string]interface{})
		req     *http.Request
	)

	if err = jsonToMap(params, &pParams); err != nil {
		return
	}

	signedParams := client.BuildSignedParams(pParams)
	var form = url.Values{}
	for k, v := range signedParams {
		form.Add(k, v)
	}

	req, err = m.request(client, apiUrl, form)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var httpClient = &http.Client{}
	if client.debug {
		dump, _ := httputil.DumpRequest(req, true)
		log.Printf("request: %s", string(dump))
	}

	resp1, err := httpClient.Do(req)
	if err != nil {
		return
	}

	if client.debug {
		dump, _ := httputil.DumpResponse(resp1, true)
		log.Printf("response: %s", string(dump))
	}

	defer resp1.Body.Close()

	resp = new(R)
	var decoder = json.NewDecoder(resp1.Body)
	err = decoder.Decode(resp)
	if err != nil {
		return
	}

	getErr := getError(resp)
	if getErr != nil {
		return nil, getErr
	}

	return
}

// request
func (m *methodImpl[R, Q]) request(client *Client, apiUrl string, form url.Values) (req *http.Request, err error) {
	if m.Method == "GET" {
		apiUrl = apiUrl + "?" + form.Encode()
		return http.NewRequest(m.Method, apiUrl, nil)
	} else {
		return http.NewRequest(m.Method, apiUrl, strings.NewReader(form.Encode()))
	}
}

// PostJSON
func (m *methodImpl[R, Q]) PostJSON(client *Client, params Q) (resp *R, err error) {
	var (
		apiUrl  = client.BaseURL + m.Path
		pParams = make(map[string]interface{})
		req     *http.Request
	)

	if err = mapstructure.Decode(params, &pParams); err != nil {
		return
	}

	signedParams := client.BuildSignedParams(pParams)

	var form = url.Values{}
	form.Add("appKey", client.AppKey)
	form.Add("sign", signedParams["sign"])

	au, err := url.Parse(apiUrl)
	if err != nil {
		return
	}
	au.RawQuery = form.Encode()
	apiUrl = au.String()

	b, err := json.Marshal(params)
	if err != nil {
		return
	}

	req, err = http.NewRequest(m.Method, apiUrl, bytes.NewReader(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	var httpClient = &http.Client{}

	resp1, err := httpClient.Do(req)
	if err != nil {
		return
	}

	if client.debug {
		dump, _ := httputil.DumpResponse(resp1, true)
		log.Printf("response: %s", string(dump))
	}

	defer resp1.Body.Close()

	resp = new(R)
	var decoder = json.NewDecoder(resp1.Body)
	err = decoder.Decode(resp)
	if err != nil {
		return
	}

	return
}

type SiteResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		SiteID     string `json:"siteId"`
		SiteDomain string `json:"siteDomain"`
		SiteStatus string `json:"siteStatus"`
	} `json:"data"`
}

type SiteRequest struct {
	SiteName     string `json:"siteName"`
	Domain       string `json:"domain,omitempty"`
	SiteType     string `json:"siteType,omitempty"`
	HTTPSForward bool   `json:"httpsForward,omitempty"`
}

type GetSiteIdsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		SiteIds []int `json:"siteIds"`
	} `json:"data"`
}

type GetSiteIdsRequest struct {
}

type GetPageIdsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		PageIds []int `json:"pageIds"`
	} `json:"data"`
}

type GetPageIdsRequest struct {
	SiteId int `json:"siteId"`
}

type PublishSiteResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Url string `json:"url"`
	} `json:"data"`
}

type PublishSiteRequest struct {
	SiteId int `json:"siteId"`
}

type PublishPageResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Url string `json:"url"`
	} `json:"data"`
}

type PublishPageRequest struct {
	SiteId int `json:"siteId"`
	PageId int `json:"pageId"`
}

type UpdatePageNameResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data,omitempty"`
}

type UpdatePageNameRequest struct {
	PageId   int    `json:"pageId"`
	PageName string `json:"pageName"`
}

type DeleteSitePageResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type DeleteSitePageRequest struct {
	PageId int `json:"pageId"`
}

type CreateSitePageResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		PageId int `json:"pageId"`
	} `json:"data"`
}

type CreateSitePageRequest struct {
	SiteId int    `json:"siteId"`
	Tpl    string `json:"tpl"`
}

type GetPageNameResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		PageId int    `json:"pageId"`
		Title  string `json:"title"`
	} `json:"data,omitempty"`
}

type GetPageNameRequest struct {
	SiteId int `json:"siteId"`
}

type GetSiteInfoResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		SiteId               string `json:"siteId"`
		SiteName             string `json:"siteName"`
		SiteDomain           string `json:"siteDomain"`
		SiteStatus           string `json:"siteStatus"`
		PackageName          string `json:"packageName"`
		PackageRemainingDays int    `json:"packageRemainingDays"`
	} `json:"data"`
}

type GetSiteInfoRequest struct {
	SiteId int `json:"siteId"`
}

type ModifyPageJsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
}

type ModifyPageJsRequest struct {
	SiteId           int    `json:"siteId"`
	PageId           string `json:"pageId"`
	Content          string `json:"content"`
	IsEncryptContent bool   `json:"isEncryptContent"`
}

type BatchModifyPagePublishPageJsResponse struct {
	Code int                               `json:"code"`
	Msg  string                            `json:"msg"`
	Data *BatchModifyPagePublishPageJsData `json:"data,omitempty"`
}

type BatchModifyPagePublishPageJsData struct {
	TaskId string `json:"taskId"`
	Task   Task   `json:"task"`
}

type Task struct {
	TaskCreateTime int64         `json:"taskCreateTime"`
	FailedPages    []SucceedPage `json:"failedPages"`
	WaitingPages   []SucceedPage `json:"waitingPages"`
	SucceedPages   []SucceedPage `json:"succeedPages"`
	TaskStatus     string        `json:"taskStatus"`
}

// UnmarshalJSON
func (m *BatchModifyPagePublishPageJsData) UnmarshalJSON(data []byte) error {
	var taskId string
	if err := json.Unmarshal(data, &taskId); err != nil {
		var task Task
		if err := json.Unmarshal(data, &task); err != nil {
			return err
		}
		m.Task = task
		return err
	}
	m.TaskId = taskId
	return nil
}

type SucceedPage struct {
	PageID   int64  `json:"pageId"`
	Status   string `json:"status"`
	ErrorMsg string `json:"errorMsg"`
}

type BatchModifyPagePublishPageJsRequest struct {
	SiteIds  []int  `json:"siteIds"`
	PageIds  []int  `json:"pageIds"`
	Content  string `json:"content"`
	IsSecure bool   `json:"isSecure"`
	TaskId   string `json:"taskId"`
}

// 套餐类型枚举值
const (
	BusinessTypeSiteAdvancedYear         = "SITE_ADVANCED_YEAR"           // 站点高级包年套餐
	BusinessTypeSiteExclusiveYear        = "SITE_EXCLUSIVE_YEAR"          // 站点尊享包年套餐
	BusinessTypeSiteExclusiveLifetime    = "SITE_EXCLUSIVE_LIFETIME"      // 站点尊享终身套餐
	BusinessTypeVoteAdvancedYear         = "VOTE_ADVANCED_YEAR"           // 投票高级包年套餐
	BusinessTypeVoteAdvancedLifetime     = "VOTE_ADVANCED_LIFETIME"       // 投票高级终身套餐
	BusinessTypeVoteAdvancedMonth        = "VOTE_ADVANCED_MONTH"          // 投票高级包年套餐
	BusinessTypeVoteDiamondOneYear       = "VOTE_DIAMOND_ONE_YEAR"        // 投票钻石版包年套餐
	BusinessTypeAppAdvancedYear          = "APP_ADVANCED_YEAR"            // 小程序高级包年套餐
	BusinessTypeAppAdvancedLifetime      = "APP_ADVANCED_LIFETIME"        // 小程序高级终身套餐
	BusinessTypeSiteExclusiveAndApp      = "SITE_EXCLUSIVE_AND_APP"       // 站点尊享小程序联合套餐
	BusinessTypeKmApiBasicYear           = "KM_API_BASIC_YEAR"            // 快码短链API基础版包年套餐
	BusinessTypeKmApiAdvancedYear        = "KM_API_ADVANCED_YEAR"         // 快码短链API进阶版包年套餐
	BusinessTypeKmApiExclusiveYear       = "KM_API_EXCLUSIVE_YEAR"        // 快码短链API高级版包年套餐
	BusinessTypeKmShortLinkAdvancedYear  = "KM_SHORT_LINK_ADVANCED_YEAR"  // 快码短链进阶版包年套餐
	BusinessTypeKmShortLinkExclusiveYear = "KM_SHORT_LINK_EXCLUSIVE_YEAR" // 快码短链高级版包年套餐
)

type OpenBusinessPackageResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
	} `json:"data,omitempty"`
}

type OpenBusinessPackageRequest struct {
	BusinessType string `json:"businessType"`      // 枚举，见套餐类型枚举值，必填字段
	SiteId       int64  `json:"siteId,omitempty"`  // 站点类型套餐的站点id，非必填
	AppId        string `json:"appId,omitempty"`   // 小程序类型套餐的小程序id，非必填
	PhoneNo      string `json:"phoneNo,omitempty"` // 投票类型套餐、快码短链、快码短链api的使用用户手机号，非必填
}

type ChangeDomainResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		NewDomain string `json:"newDomain"`
	} `json:"data,omitempty"`
}

type ChangeDomainRequest struct {
	SiteId       int64  `json:"siteId"`
	Domain       string `json:"domain"`
	HTTPSForward bool   `json:"httpsForward"`
}

type UpdateSiteInfoResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data,omitempty"`
}

type UpdateSiteInfoRequest struct {
	SiteId   int64  `json:"siteId"`
	SiteName string `json:"siteName"`
}
