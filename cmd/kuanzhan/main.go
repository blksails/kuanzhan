package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/html"
	"pkg.blksails.net/kuanzhan"
)

var rootCmd = &cobra.Command{
	Use:   "kuanzhan",
	Short: "kuanzhan",
	Long:  "kuanzhan",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createSiteCmd = &cobra.Command{
	Use:   "create-site",
	Short: "创建站点",
	Long:  "快站快速创建站点",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()
		for i := 0; i < siteSize; i++ {
			uniqueDomain := randomUniqueDomain()
			resp, err := client.CreateSite(createSiteName, uniqueDomain, createSiteType, true)
			if err != nil {
				log.Fatal(err)
			}

			siteId, err := strconv.ParseInt(resp.Data.SiteID, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			_, err = client.OpenBusinessPackage(businessType, siteId, "", "")
			if err != nil {
				log.Fatal(err)
			}

			log.Println("create site", resp.Data.SiteID)
		}
	},
}

var siteListCmd = &cobra.Command{
	Use:   "list",
	Short: "站点列表",
	Long:  "快站快速获取站点列表",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()
		resp, err := client.GetSiteIds()
		if err != nil {
			log.Fatal(err)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("站点ID", "站点名称", "套餐类型", "站点URL", "站点状态", "页面ID", "页面名称", "页面URL")

		rows := [][]string{}
		for _, siteId := range resp.Data.SiteIds {
			siteInfo, err := client.GetSiteInfo(siteId)
			if err != nil {
				log.Fatal(err)
			}

			prerow := []string{
				fmt.Sprintf("%d", siteId),
				siteInfo.Data.SiteName,
				siteInfo.Data.SiteDomain,
				siteInfo.Data.PackageName,
				siteInfo.Data.SiteStatus,
			}

			if !onlySite {
				pageNames, err := client.GetPageName(siteId)
				if err != nil {
					rows = append(rows, prerow)
					continue
				}

				for _, pageName := range pageNames.Data {
					row := append(prerow, fmt.Sprintf("%d", pageName.PageId), pageName.Title, siteInfo.Data.SiteDomain+"/"+strconv.Itoa(pageName.PageId))
					rows = append(rows, row)
				}
			} else {
				rows = append(rows, prerow)
			}

		}
		table.Bulk(rows)
		table.Render()
	},
}

var uploadSiteCmd = &cobra.Command{
	Use:   "upload",
	Short: "上传站点",
	Long:  "快站快速上传站点",
	Run: func(cmd *cobra.Command, args []string) {
		pagehtml, err := downloadPage(sourceUrl)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("pagehtml", pagehtml)
		log.Println("upload site ", sourceUrl, " to site ", siteIds, " page ", pageSize, " pageIds ", pageIds)

		var (
			allPageIds []int
		)

		client := newClient()
		if taskId != "" {
			resp, err := client.BatchModifyPagePublishPageJs(siteIds, allPageIds, pagehtml, true, taskId)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("task", resp.Data)
			return
		}
		if len(pageIds) > 0 {
			for _, pageId := range pageIds {
				_, err := client.UpdatePageName(pageId, pageName)
				if err != nil {
					log.Fatal(err)
				}
			}
			allPageIds = pageIds
		}

		for _, siteId := range siteIds {
			_, err := client.PublishSite(siteId)
			if err != nil {
				log.Fatal(err)
			}

			if len(pageIds) == 0 {
				sitePageIds := []int{}
				for i := 0; i < pageSize; i++ {
					resp, err := client.CreateSitePage(siteId, createPateTpl)
					if err != nil {
						log.Fatal(err)
					}
					_, err = client.UpdatePageName(resp.Data.PageId, pageName)
					if err != nil {
						log.Fatal(err)
					}
					sitePageIds = append(sitePageIds, resp.Data.PageId)
				}

				allPageIds = append(allPageIds, sitePageIds...)
			}
		}

		if len(allPageIds) == 0 || len(siteIds) == 0 {
			log.Fatal("no page ids or site ids")
		}

		resp, err := client.BatchModifyPagePublishPageJs(siteIds, allPageIds, pagehtml, true, "")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("taskId", resp.Data.TaskId)
	},
}

var updatePageCmd = &cobra.Command{
	Use:   "update",
	Short: "更新页面",
	Long:  "更新页面",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()
		for _, pageId := range pageIds {
			_, err := client.UpdatePageName(pageId, pageName)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var deletePageCmd = &cobra.Command{
	Use:   "delete-page",
	Short: "删除页面",
	Long:  "快站快速删除页面",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()
		for _, pageId := range pageIds {
			_, err := client.DeleteSitePage(pageId)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var upgradeSiteCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "升级站点",
	Long:  "升级站点",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()
		for _, siteId := range siteIds {
			_, err := client.OpenBusinessPackage(businessType, int64(siteId), "", "")
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var changeDomainCmd = &cobra.Command{
	Use:   "change-domain",
	Short: "更换域名",
	Long:  "更换域名",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()
		for _, siteId := range siteIds {
			domain := randomUniqueDomain()
			_, err := client.ChangeDomain(int64(siteId), domain, true)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var updateSiteInfoCmd = &cobra.Command{
	Use:   "update-site",
	Short: "更新站点信息",
	Long:  "更新站点信息",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()
		for _, siteId := range siteIds {
			_, err := client.UpdateSiteInfo(int64(siteId), siteName)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var publishPageCmd = &cobra.Command{
	Use:   "publish-page",
	Short: "发布页面",
	Long:  "发布页面",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()

		siteResp, err := client.PublishSite(siteId)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("siteResp", siteResp.Data.Url)

		pageResp, err := client.PublishPage(siteId, pageId)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("pageResp", pageResp.Data.Url)
	},
}

func Execute() {
	rootCmd.AddCommand(createSiteCmd)
	rootCmd.AddCommand(uploadSiteCmd)
	rootCmd.AddCommand(siteListCmd)
	rootCmd.AddCommand(updatePageCmd)
	rootCmd.AddCommand(deletePageCmd)
	rootCmd.AddCommand(upgradeSiteCmd)
	rootCmd.AddCommand(changeDomainCmd)
	rootCmd.AddCommand(updateSiteInfoCmd)
	rootCmd.AddCommand(publishPageCmd)
	rootCmd.Execute()
}

var (
	sourceUrl      string // 上传站点源站URL
	appKey         string // 快站APPKEY
	appSecret      string // 快站APPSECRET
	siteIds        []int  // 站点ID
	siteSize       int    // 创建站点数量
	siteName       string // 站点名称
	pageSize       int    // 创建页面数量
	pageName       string // 创建页面名称
	businessType   string // 创建站点套餐类型
	createPateTpl  string // 创建页面模板
	createSiteName string // 创建站点名称
	createSiteType string // 创建站点类型
	debug          bool   // 是否debug
	pageIds        []int  // 页面ID
	onlySite       bool   // 是否只显示站点
	taskId         string // 任务ID
	pageId         int    // 页面ID
	siteId         int    // 站点ID
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug")

	siteListCmd.PersistentFlags().BoolVarP(&onlySite, "only-site", "o", false, "是否只显示站点")

	pflags := uploadSiteCmd.PersistentFlags()

	pflags.StringVarP(&sourceUrl, "source-url", "s", "", "源站URL") // required
	uploadSiteCmd.MarkPersistentFlagRequired("source-url")

	pflags.IntSliceVarP(&siteIds, "site-ids", "i", []int{}, "站点ID") // required
	uploadSiteCmd.MarkPersistentFlagRequired("site-ids")

	createSiteCmd.PersistentFlags().IntVarP(&siteSize, "size", "s", 5, "创建站点数量")
	createSiteCmd.MarkPersistentFlagRequired("size")
	createSiteCmd.PersistentFlags().StringVarP(&createSiteName, "name", "n", "", "创建站点域名")
	createSiteCmd.MarkPersistentFlagRequired("name")
	createSiteCmd.PersistentFlags().StringVarP(&createSiteType, "type", "t", "FAST", "创建站点类型")
	createSiteCmd.PersistentFlags().StringVarP(&businessType, "business-type", "b", "SITE_EXCLUSIVE_YEAR", "创建站点套餐类型")

	uploadSiteCmd.PersistentFlags().IntVarP(&pageSize, "page", "p", 1, "创建页面数量")
	uploadSiteCmd.PersistentFlags().StringVarP(&createPateTpl, "tpl", "t", "WHITE", "创建页面模板")
	uploadSiteCmd.PersistentFlags().StringVarP(&pageName, "name", "n", "", "创建页面名称")
	uploadSiteCmd.PersistentFlags().IntSliceVarP(&pageIds, "page-ids", "g", []int{}, "指定页面ID")
	uploadSiteCmd.PersistentFlags().StringVarP(&taskId, "task-id", "a", "", "任务ID")

	uploadSiteCmd.MarkPersistentFlagRequired("name")

	updatePageCmd.PersistentFlags().StringVarP(&pageName, "name", "n", "", "更新页面名称")
	updatePageCmd.MarkPersistentFlagRequired("name")
	updatePageCmd.PersistentFlags().IntSliceVarP(&pageIds, "page-ids", "i", []int{}, "页面ID")

	deletePageCmd.PersistentFlags().IntSliceVarP(&pageIds, "page-ids", "i", []int{}, "页面ID")
	deletePageCmd.MarkPersistentFlagRequired("page-ids")

	upgradeSiteCmd.PersistentFlags().StringVarP(&businessType, "business-type", "b", "SITE_EXCLUSIVE_YEAR", "升级站点套餐类型")
	upgradeSiteCmd.PersistentFlags().IntSliceVarP(&siteIds, "site-ids", "i", []int{}, "站点ID")
	upgradeSiteCmd.MarkPersistentFlagRequired("site-ids")

	changeDomainCmd.PersistentFlags().IntSliceVarP(&siteIds, "site-ids", "i", []int{}, "站点ID")
	changeDomainCmd.MarkPersistentFlagRequired("site-ids")

	updateSiteInfoCmd.PersistentFlags().StringVarP(&siteName, "name", "n", "", "更新站点名称")
	updateSiteInfoCmd.MarkPersistentFlagRequired("name")
	updateSiteInfoCmd.PersistentFlags().IntSliceVarP(&siteIds, "site-ids", "i", []int{}, "站点ID")
	updateSiteInfoCmd.MarkPersistentFlagRequired("site-ids")

	publishPageCmd.PersistentFlags().IntVarP(&siteId, "site-id", "i", 0, "站点ID")
	publishPageCmd.MarkPersistentFlagRequired("site-id")
	publishPageCmd.PersistentFlags().IntVarP(&pageId, "page-id", "p", 0, "页面ID")
	publishPageCmd.MarkPersistentFlagRequired("page-id")

	initConfig()
}

func initConfig() {
	viper.SetConfigName("kuanzhan")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	appKey = viper.GetString("app_key")
	appSecret = viper.GetString("app_secret")

}

func main() {
	Execute()
}

func newClient() *kuanzhan.Client {
	client := kuanzhan.NewClient(appKey, appSecret)
	client.SetDebug(debug)
	return client
}

func downloadPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	// 查找body节点
	bodyNode := findBodyNode(doc)
	if bodyNode == nil {
		return "", fmt.Errorf("body node not found")
	}

	// 将body节点内部内容转换为HTML字符串
	body := renderBodyContent(bodyNode)
	return body, nil
}

// 查找body节点
func findBodyNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "body" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findBodyNode(c); result != nil {
			return result
		}
	}
	return nil
}

// 将节点渲染为HTML字符串
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	html.Render(&buf, n)
	return buf.String()
}

// 将body节点内部内容渲染为HTML字符串（不包含body标签本身）
func renderBodyContent(bodyNode *html.Node) string {
	var buf bytes.Buffer
	for c := bodyNode.FirstChild; c != nil; c = c.NextSibling {
		html.Render(&buf, c)
	}
	return buf.String()
}

func randomUniqueDomain() string {
	// 字符数字混合字符集
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	// 获取当前时间戳（纳秒级）
	now := time.Now().UnixNano()

	// 将时间戳转换为36进制字符串（使用字符集）
	timeStr := strconv.FormatInt(now%1000000, 36) // 取时间戳的后6位并转换为36进制

	// 生成随机字符部分（4位）
	randomLength := 4
	randomPart := make([]byte, randomLength)
	for i := range randomPart {
		randomPart[i] = charset[rand.Intn(len(charset))]
	}

	// 组合时间因子和随机字符
	return timeStr + string(randomPart)
}

func init() {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())
}
