package bmw

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/oauth"
	"github.com/evcc-io/evcc/util/request"
	"github.com/jpfielding/go-http-digest/pkg/digest"
	"golang.org/x/oauth2"
)

const AuthURI = "https://customer.bmwgroup.com/gcdm/oauth/authenticate"

type Identity struct {
	*request.Helper
	oauth2.TokenSource
	user, password string
}

// NewIdentity creates BMW identity
func NewIdentity(log *util.Logger) *Identity {
	v := &Identity{
		Helper: request.NewHelper(log),
	}

	return v
}

func (v *Identity) Login(user, password string) error {
	v.user = user
	v.password = password

	token, err := v.RefreshToken(nil)

	if err == nil {
		v.TokenSource = oauth.RefreshTokenSource(token, v)
	}

	return err
}

func (v *Identity) RefreshToken(_ *oauth2.Token) (*oauth2.Token, error) {
	data := url.Values{
		"client_id":             []string{"31c357a0-7a1d-4590-aa99-33b97244d048"},
		"response_type":         []string{"code"},
		"scope":                 []string{"authenticate_user vehicle_data remote_services"},
		"redirect_uri":          []string{"com.bmw.connected://oauth"},
		"state":                 []string{"cwU-gIE27j67poy2UcL3KQ"},
		"nonce":                 []string{"login_nonce"},
		"code_challenge_method": []string{"S256"},
		"username":              []string{v.user},
		"password":              []string{v.password},
		"grant_type":            []string{"authorization_code"},
	}

	v.Client.Transport = digest.NewTransport(v.user, v.password, v.Client.Transport)

	req, err := request.New(http.MethodPost, AuthURI, strings.NewReader(data.Encode()), map[string]string{
		"Accept":       "application/json, text/plain, */*",
		"Content-Type": "application/x-www-form-urlencoded",
	})
	if err != nil {
		return nil, err
	}

	// don't follow redirects
	v.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	defer func() { v.Client.CheckRedirect = nil }()

	resp, err := v.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		var desc struct {
			Error            string
			ErrorDescription string `json:"error_description"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&desc); err != nil {
			return nil, fmt.Errorf("could not obtain token")
		}

		return nil, fmt.Errorf("%s:%s", desc.Error, desc.ErrorDescription)
	}

	uri := resp.Header.Get("Location")

	query, err := url.ParseQuery(strings.TrimPrefix(uri, "com.bmw.connected://oauth#"))
	if err != nil {
		return nil, err
	}

	at := query.Get("access_token")
	expires, err := strconv.Atoi(query.Get("expires_in"))
	if err != nil || at == "" || expires == 0 {
		return nil, errors.New("could not obtain token")
	}

	token := &oauth2.Token{
		AccessToken: at,
		Expiry:      time.Now().Add(time.Duration(expires) * time.Second),
	}

	return token, nil
}
