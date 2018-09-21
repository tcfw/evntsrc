package passport

import (
	"encoding/json"
	fmt "fmt"
	"net/http"
	"time"

	pb "github.com/tcfw/evntsrc/pkg/passport/protos"
)

type socialUserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func validateGoogleLogin(request *pb.SocialRequest) (*socialUserInfo, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token="+request.GetIdpTokens().GetToken(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	userInfo := &socialUserInfo{}
	json.NewDecoder(resp.Body).Decode(userInfo)

	return userInfo, nil
}

func validateFacebookLogin(request *pb.SocialRequest) (*socialUserInfo, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", "https://graph.facebook.com/me?fields=name,email,picture", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+request.GetIdpTokens().GetToken())
	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	userInfo := &socialUserInfo{}
	json.NewDecoder(resp.Body).Decode(userInfo)

	return userInfo, nil
}

func getSocialInfo(request *pb.SocialRequest) (*socialUserInfo, error) {
	switch provider := request.GetProvider(); provider {
	case "google":
		return validateGoogleLogin(request)
	case "facebook":
		return validateFacebookLogin(request)
	default:
		return nil, fmt.Errorf("Unknown social provider %s", provider)
	}
}
