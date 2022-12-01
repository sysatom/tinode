package pocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tinode/chat/server/extra/cache"
	"net/http"
	"time"
)

const (
	ID          = "pocket"
	ClientIdKey = "consumer_key"
)

type CodeResponse struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}

type ListResponse struct {
	Status int             `json:"status"`
	List   map[string]Item `json:"list"`
}

type Item struct {
	Id            string `json:"item_id"`
	ResolvedId    string `json:"resolved_id"`
	GivenUrl      string `json:"given_url"`
	GivenTitle    string `json:"given_title"`
	Favorite      string `json:"favorite"`
	Status        string `json:"status"`
	TimeAdded     string `json:"time_added"`
	TimeUpdated   string `json:"time_updated"`
	TimeRead      string `json:"time_read"`
	TimeFavorited string `json:"time_favorited"`
	ResolvedTitle string `json:"resolved_title"`
	ResolvedUrl   string `json:"resolved_url"`
	Excerpt       string `json:"excerpt"`
	IsArticle     string `json:"is_article"`
	IsIndex       string `json:"is_index"`
	HasVideo      string `json:"has_video"`
	HasImage      string `json:"has_image"`
	WordCount     string `json:"word_count"`
}

type ItemResponse struct {
	Status int  `json:"status"`
	Item   Item `json:"item"`
}

type Pocket struct {
	c            *resty.Client
	clientId     string // ConsumerKey
	clientSecret string
	redirectURI  string
	accessToken  string
	code         string
}

func NewPocket(clientId, clientSecret, redirectURI, accessToken string) *Pocket {
	v := &Pocket{clientId: clientId, clientSecret: clientSecret, redirectURI: redirectURI, accessToken: accessToken}

	v.c = resty.New()
	v.c.SetBaseURL("https://getpocket.com")
	v.c.SetTimeout(time.Minute)

	return v
}

func (v *Pocket) GetCode(state string) (*CodeResponse, error) {
	resp, err := v.c.R().
		SetResult(&CodeResponse{}).
		SetHeader("X-Accept", "application/json").
		SetBody(map[string]interface{}{"consumer_key": v.clientId, "redirect_uri": v.redirectURI, "state": state}).
		Post("/v3/oauth/request")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusOK {
		result := resp.Result().(*CodeResponse)
		v.code = result.Code

		_ = cache.DB.Set([]byte("pocket:code"), []byte(v.code)) // todo

		return result, nil
	} else {
		return nil, fmt.Errorf("%d, %s (%s)", resp.StatusCode(), resp.Header().Get("X-Error-Code"), resp.Header().Get("X-Error"))
	}
}

func (v *Pocket) AuthorizeURL() string {
	return fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", v.code, v.redirectURI)
}

func (v *Pocket) GetAccessToken(code string) (interface{}, error) {
	resp, err := v.c.R().
		SetResult(&TokenResponse{}).
		SetHeader("X-Accept", "application/json").
		SetBody(map[string]interface{}{"consumer_key": v.clientId, "code": code}).
		Post("/v3/oauth/authorize")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusOK {
		result := resp.Result().(*TokenResponse)
		v.accessToken = result.AccessToken
		return result, nil
	} else {
		return nil, fmt.Errorf("%d, %s (%s)", resp.StatusCode(), resp.Header().Get("X-Error-Code"), resp.Header().Get("X-Error"))
	}
}

func (v *Pocket) Redirect(_ *http.Request) (string, error) {
	_ = cache.DB.Set([]byte("pocket:code"), []byte(v.code)) // fixme uid key

	appRedirectURI := v.AuthorizeURL()
	return appRedirectURI, nil
}

func (v *Pocket) StoreAccessToken(_ *http.Request) (map[string]interface{}, error) {
	data, err := cache.DB.Get([]byte("pocket:code")) // fixme uid key
	if err != nil {
		return nil, err
	}
	code := string(data)

	if code != "" {
		tokenResp, err := v.GetAccessToken(code)
		if err != nil {
			return nil, err
		}

		extra, err := json.Marshal(&tokenResp)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"name":  ID,
			"type":  ID,
			"token": v.accessToken,
			"extra": extra,
		}, nil
	}
	return nil, errors.New("error")
}

func (v *Pocket) Retrieve(count int) (*ListResponse, error) {
	resp, err := v.c.R().
		SetResult(&ListResponse{}).
		SetBody(map[string]interface{}{
			"consumer_key": v.clientId,
			"access_token": v.accessToken,
			"count":        count,
			"detailType":   "simple",
			"state":        "all",
			"sort":         "newest",
		}).
		Post("/v3/get")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusOK {
		return resp.Result().(*ListResponse), nil
	} else {
		return nil, fmt.Errorf("%d, %s (%s)", resp.StatusCode(), resp.Header().Get("X-Error-Code"), resp.Header().Get("X-Error"))
	}
}

func (v *Pocket) Add(url string) (int, error) {
	resp, err := v.c.R().
		SetResult(&ItemResponse{}).
		SetBody(map[string]interface{}{
			"consumer_key": v.clientId,
			"access_token": v.accessToken,
			"url":          url,
		}).
		Post("/v3/add")
	if err != nil {
		return 0, err
	}

	if resp.StatusCode() == http.StatusOK {
		return resp.Result().(*ItemResponse).Status, nil
	} else {
		return 0, fmt.Errorf("%d, %s (%s)", resp.StatusCode(), resp.Header().Get("X-Error-Code"), resp.Header().Get("X-Error"))
	}
}
