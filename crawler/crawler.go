package crawler

import (
	"compress/gzip"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"schannel-qt5/parser"
	"schannel-qt5/urls"
)

// GetCsrfToken 返回本次回话所使用的csrftoken
func GetCsrfToken(proxy string) (string, []*http.Cookie, error) {
	client, err := genClientWithProxy(proxy)
	if err != nil {
		return "", nil, err
	}

	request, err := http.NewRequest("GET", urls.AuthPath, nil)
	if err != nil {
		return "", nil, err
	}
	setRequestHeader(request, nil, urls.RootPath)

	resp, err := client.Do(request)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	htmlReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", nil, err
	}
	defer htmlReader.Close()

	dom, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", nil, err
	}

	csrfToken, exists := dom.Find("input[type='hidden'][name='token']").Eq(0).Attr("value")
	if !exists {
		return "", nil, errors.New("csrftoken doesn't exist")
	}

	u2, _ := url.Parse(urls.RootPath)
	return csrfToken, client.Jar.Cookies(u2), nil
}

// GetAuth 登录schannel并返回登陆成功后获得的cookies
// 这些cookies在后续的页面访问中需要使用
func GetAuth(user, passwd, proxy string) ([]*http.Cookie, error) {
	client, err := genClientWithProxy(proxy)
	if err != nil {
		return nil, err
	}

	csrfToken, session, err := GetCsrfToken(proxy)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("token", csrfToken)
	form.Set("username", user)
	form.Set("password", passwd)
	getlogin, err := http.NewRequest("POST", urls.LoginPath, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	setRequestHeader(getlogin, session, urls.AuthPath)
	getlogin.Header.Set("content-type", "application/x-www-form-urlencoded")

	resp, err := client.Do(getlogin)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 验证登录是否成功，如果incorrect的值为“true”，则登录失败
	incorrect := resp.Request.FormValue("incorrect")
	if incorrect == "true" {
		return nil, errors.New("登录验证失败")
	}

	u2, _ := url.Parse(urls.RootPath)
	cookies := make([]*http.Cookie, 0, 2)
	cookies = append(cookies, client.Jar.Cookies(u2)...)
	// 添加cfuid会话cookie
	for _, c := range session {
		if c.Name == "__cfduid" {
			cookies = append(cookies, c)
		}
	}

	return cookies, nil
}

// GetServiceHTML 获取所有已购买服务的状态信息，包括详细页面的地址
func GetServiceHTML(cookies []*http.Cookie, proxy string) (string, error) {
	return getAccountPage(urls.ServiceListPath, cookies, proxy)
}

// GetInvoiceHTML 获取账单页面的HTML,包含未付款和已付款账单
// 未付款账单显示在最前列
// 现只支持获取第一页
func GetInvoiceHTML(cookies []*http.Cookie, proxy string) (string, error) {
	return getAccountPage(urls.InvoicePath, cookies, proxy)
}

// GetSSRInfoHTML 获取服务详细信息页面的HTML，包含使用情况和节点信息
func GetSSRInfoHTML(service *parser.Service, cookies []*http.Cookie, proxy string) (string, error) {
	return getAccountPage(service.Link, cookies, proxy)
}
