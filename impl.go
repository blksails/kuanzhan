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
	CreateSitePage               *methodImpl[CreateSitePageResponse, CreateSitePageRequest]
	GetSiteInfo                  *methodImpl[GetSiteInfoResponse, GetSiteInfoRequest]
	ModifyPageJs                 *methodImpl[ModifyPageJsResponse, ModifyPageJsRequest]
	BatchModifyPagePublishPageJs *methodImpl[BatchModifyPagePublishPageJsResponse, BatchModifyPagePublishPageJsRequest]
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

	req, err = http.NewRequest(m.Method, apiUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
	Msg  string `json:"message"`
	Data struct {
		SiteID     string `json:"site_id"`
		SiteDomain string `json:"site_domain"`
		SiteStatus string `json:"site_status"`
	} `json:"data"`
}

type SiteRequest struct {
	SiteName     string `json:"siteName"`
	Domain       string `json:"domain"`
	SiteType     string `json:"siteType"`
	HTTPSForward bool   `json:"httpsForward"`
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

type CreateSitePageResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Url string `json:"url"`
	} `json:"data"`
}

type CreateSitePageRequest struct {
	SiteId int    `json:"siteId"`
	Tpl    string `json:"tpl"`
}

type GetSiteInfoResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		SiteId               string `json:"siteId"`
		SiteName             string `json:"siteName"`
		SiteDomain           string `json:"siteType"`
		SiteStatus           string `json:"domain"`
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
	Code int                              `json:"code"`
	Msg  string                           `json:"msg"`
	Data BatchModifyPagePublishPageJsData `json:"data"`
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

func jsonToMap(v any, m *map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()

	return dec.Decode(m)
}
