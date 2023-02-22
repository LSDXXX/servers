package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	"github.com/LSDXXX/libs/pkg/log"
	http "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type OpenAIAuth struct {
	Email    string
	Password string
	Proxy    string

	client        tlsClient.HttpClient
	userAgent     string
	sessionToken  string
	accessToken   string
	sessionCookie *http.Cookie
}

func NewOpenAIAuth(email, password, proxy string) *OpenAIAuth {
	out := &OpenAIAuth{
		Email:     email,
		Password:  password,
		Proxy:     proxy,
		userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	}
	jar := tlsClient.NewCookieJar()
	options := []tlsClient.HttpClientOption{
		tlsClient.WithTimeoutSeconds(30),
		tlsClient.WithClientProfile(tlsClient.Chrome_105),
		tlsClient.WithNotFollowRedirects(),
		tlsClient.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		//tls_client.WithInsecureSkipVerify(),
	}
	if len(out.Proxy) != 0 {
		options = append(options, tlsClient.WithProxyUrl(out.Proxy))
	}

	client, err := tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), options...)
	if err != nil {
		panic("create tls client err: " + err.Error())
	}
	out.client = client
	return out
}

func (o *OpenAIAuth) Login() error {
	if len(o.Email) == 0 || len(o.Password) == 0 {
		return errors.New("email and password is required")
	}
	logrus.Debugf("start login ---")
	url := "https://explorer.api.openai.com/"
	headers := map[string][]string{
		"Host":            {"ask.openai.com"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"User-Agent":      {o.userAgent},
		"Accept-Language": {"en-GB,en-US;q=0.9,en;q=0.8"},
		"Accept-Encoding": {"gzip, deflate, br"},
		"Connection":      {"keep-alive"},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	req.Header = headers
	res, err := o.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Errorf("step 1: status code: %d, content: %s", res.StatusCode, string(data))
		return errors.New("step 1 error")
	}
	return o.step2()
}

func (o *OpenAIAuth) step2() error {
	logrus.Debugf("step2 begin ----")
	url := "https://explorer.api.openai.com/api/auth/csrf"
	headers := map[string][]string{
		"Host":            {"ask.openai.com"},
		"Accept":          {"*/*"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {o.userAgent},
		"Accept-Language": {"en-GB,en-US;q=0.9,en;q=0.8"},
		"Referer":         {"https://explorer.api.openai.com/auth/login"},
		"Accept-Encoding": {"gzip, deflate, br"},
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	req.Header = headers
	res, err := o.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Errorf("step 2: status code: %d, content: %s", res.StatusCode, string(data))
		return errors.New("step 2 error")
	}
	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}
	token, ok := m["csrfToken"]
	if !ok {
		return errors.New("empty csrf token")
	}
	return o.step3(token.(string))
}

func (o *OpenAIAuth) step3(token string) error {
	logrus.Debugf("step3 begin ----")
	logrus.Debugf("step3 csrf token: %s", token)
	form := url.Values{}
	url := "https://explorer.api.openai.com/api/auth/signin/auth0?prompt=login"

	form.Add("callbackUrl", "/")
	form.Add("csrfToken", token)
	form.Add("json", "true")

	logrus.Debugf("form data: %s", form.Encode())
	payload := bytes.NewBufferString(form.Encode())
	headers := map[string][]string{
		"Host":            {"explorer.api.openai.com"},
		"User-Agent":      {o.userAgent},
		"Content-Type":    {"application/x-www-form-urlencoded"},
		"Accept":          {"*/*"},
		"Sec-Gpc":         {"1"},
		"Accept-Language": {"en-US,en;q=0.8"},
		"Origin":          {"https://explorer.api.openai.com"},
		"Sec-Fetch-Site":  {"same-origin"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Dest":  {"empty"},
		"Referer":         {"https://explorer.api.openai.com/auth/login"},
		"Accept-Encoding": {"gzip, deflate"},
	}
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	req.Header = headers
	res, err := o.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Errorf("step 3: status code: %d, content: %s", res.StatusCode, string(data))
		return errors.New("step 3 error")
	}
	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}
	url = m["url"].(string)
	if url == "https://explorer.api.openai.com/api/auth/error?error=OAuthSignin" || strings.Contains(url, "error") {
		return errors.New("you have been rate limited")
	}
	return o.step4(url)
}

func (o *OpenAIAuth) step4(url string) error {
	logrus.Debugf("begin step4 ----")
	logrus.Debugf("step4 url: %s", url)
	headers := map[string][]string{
		"Host":            {"auth0.openai.com"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {o.userAgent},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Referer":         {"https://explorer.api.openai.com/"},
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	req.Header = headers
	res, err := o.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 302 && res.StatusCode != 200 {
		logrus.Errorf("step 4: status code: %d, content: %s", res.StatusCode, string(data))
		return errors.New("step 4 error")
	}
	logrus.Debugf("step 4 response: %s", string(data))
	re := regexp.MustCompile("state=(.*)")
	allMatch := re.FindAllString(string(data), -1)
	if len(allMatch) == 0 {
		return errors.New("rate limit")
	}
	state := strings.TrimLeft(strings.Split(allMatch[0], `"`)[0], "state=")
	return o.step5(state)
}

func (o *OpenAIAuth) step5(state string) error {
	logrus.Debugf("begin step5 ----")
	url := fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state)
	logrus.Debugf("step5 state: %s, url: %s", state, url)

	headers := map[string][]string{
		"Host":            {"auth0.openai.com"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {o.userAgent},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Referer":         {"https://explorer.api.openai.com/"},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	req.Header = headers
	res, err := o.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Errorf("step 5: status code: %d, content: %s", res.StatusCode, string(data))
		return errors.New("step 5 error")
	}

	return o.step6(state)
}

func (o *OpenAIAuth) step6(state string) error {
	logrus.Debugf("begin step6")
	form := url.Values{}
	url := fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state)

	form.Add("state", state)
	form.Add("username", o.Email)
	form.Add("js-available", "false")
	form.Add("webauthn-available", "true")
	form.Add("is-brave", "false")
	form.Add("webauthn-platform-available", "true")
	form.Add("action", "default")

	headers := map[string][]string{
		"Host":            {"auth0.openai.com"},
		"Origin":          {"https://auth0.openai.com"},
		"Connection":      {"keep-alive"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"User-Agent":      {o.userAgent},
		"Referer":         {"https://auth0.openai.com/u/login/identifier?state={state}"},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Content-Type":    {"application/x-www-form-urlencoded"},
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	req.Header = headers
	res, err := o.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 302 {
		logrus.Errorf("step 6: status code: %d, content: %s", res.StatusCode, string(data))
		return errors.New("step 6 error")
	}
	return o.step7(state)
}

func (o *OpenAIAuth) step7(state string) error {
	logrus.Debugf("begin step7 ----")
	form := url.Values{}
	url := fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state)

	form.Add("state", state)
	form.Add("username", o.Email)
	form.Add("password", o.Password)
	form.Add("action", "default")

	headers := map[string][]string{
		"Host":            {"auth0.openai.com"},
		"Origin":          {"https://auth0.openai.com"},
		"Connection":      {"keep-alive"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"User-Agent":      {o.userAgent},
		"Referer":         {"https://auth0.openai.com/u/login/password?state={state}"},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Content-Type":    {"application/x-www-form-urlencoded"},
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	req.Header = headers
	res, err := o.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 302 && res.StatusCode != 200 {
		logrus.Errorf("step 7: status code: %d, content: %s", res.StatusCode, string(data))
		return errors.New("step 7 error")
	}
	logrus.Debugf("step 7: status code: %d, content: %s", res.StatusCode, string(data))

	re := regexp.MustCompile("state=(.*)")
	allMatch := re.FindAllString(string(data), -1)
	if len(allMatch) == 0 {
		return errors.New("rate limit")
	}
	newState := strings.TrimLeft(strings.Split(allMatch[0], `"`)[0], "state=")
	return o.step8(state, newState)
}

func (o *OpenAIAuth) step8(oldState, newState string) error {
	logrus.Debugf("begin step8 ----")
	logrus.Debugf("step8 newState: %s", newState)
	url := fmt.Sprintf("https://auth0.openai.com/authorize/resume?state=%s", newState)

	headers := map[string][]string{
		"Host":            {"auth0.openai.com"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {o.userAgent},
		"Accept-Language": {"en-GB,en-US;q=0.9,en;q=0.8"},
		"Referer":         {fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", oldState)},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	req.Header = headers
	res, err := o.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 302 {
		logrus.Errorf("step 8: status code: %d, content: %s", res.StatusCode, string(data))
		return errors.New("step 8 error")
	}
	logrus.Debugf("step 8: status code: %d, content: %s, location: %s", res.StatusCode, string(data), res.Header.Get("location"))
	return o.step9(res.Header.Get("location"), url)
}

func (o *OpenAIAuth) step9(redirectUrl, previousUrl string) error {
	log.WithContext(context.Background()).Debugf("begin step9 ----")
	log.WithContext(context.Background()).Debugf("redirect url: %s, previous url: %s", redirectUrl, previousUrl)

	headers := map[string][]string{
		"Host":            {"explorer.api.openai.com"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {o.userAgent},
		"Accept-Language": {"en-GB,en-US;q=0.9,en;q=0.8"},
		"Referer":         {previousUrl},
	}

	req, err := http.NewRequest(http.MethodGet, redirectUrl, nil)
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	req.Header = headers
	res, err := o.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 302 {
		logrus.Errorf("step 9: status code: %d, content: %s", res.StatusCode, string(data))
		return errors.New("step 9 error")
	}

	cookies := res.Cookies()
	var sessionToken string
	for _, cookie := range cookies {
		if cookie.Name == "__Secure-next-auth.session-token" {
			sessionToken = cookie.Value
			o.sessionCookie = cookie
			break
		}
	}
	o.sessionToken = sessionToken
	return o.getAccessToken()
}

func (o *OpenAIAuth) getAccessToken() error {
	o.client.SetCookies(&url.URL{
		Host: "auth0.openai.com",
	}, []*http.Cookie{o.sessionCookie})

	res, err := o.client.Get("https://explorer.api.openai.com/api/auth/session")
	if err != nil {
		return errors.Wrap(err, "get access token")
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		if res.StatusCode != 200 {
			logrus.Errorf("getAccessToken status code: %d, content: %s", res.StatusCode, string(data))
			return errors.New("getAccessToken error")
		}
	}
	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}
	token, ok := m["accessToken"]
	if !ok {
		return errors.New("empty csrf token")
	}
	o.accessToken = token.(string)
	logrus.Debugf("access token: %s", o.accessToken)
	return nil
}

func (o *OpenAIAuth) AccessToken() string {
	return o.accessToken
}

func (o *OpenAIAuth) SessionToken() string {
	return o.sessionToken
}
