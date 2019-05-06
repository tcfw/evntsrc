package evntsrc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gogo/protobuf/proto"

	passportPB "github.com/tcfw/evntsrc/internal/passport/protos"
)

func (api *APIClient) getToken() error {
	reqObj := &passportPB.AuthRequest{Creds: &passportPB.AuthRequest_OAuthCodeCreds{OAuthCodeCreds: &passportPB.OAuthCodeCreds{Code: api.auth}}}

	pbBytes, err := reqObj.Marshal()
	if err != nil {
		return err
	}

	reqBody := bytes.NewReader(pbBytes)

	req, err := http.NewRequest(http.MethodPost, api.formatURL(endpointAPI, "auth/login"), reqBody)
	if err != nil {
		return err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	authResp := &passportPB.AuthResponse{}
	if err := proto.Unmarshal(body, authResp); err != nil {
		return err
	}

	if !authResp.Success {
		return fmt.Errorf("Failed to get token")
	}

	if authResp.MFAResponse != nil {
		return fmt.Errorf("Token requires MFA which is not supported in client libraries")
	}

	api.token = authResp.Tokens.Token

	return nil
}
