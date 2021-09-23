package upcloud

import (
	"demoapp/context"
	"demoapp/db"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type AccountAPIResponse struct {
	Account AccountResponse `json:"account"`
}

type AccountResponse struct {
	Credits        float64        `json:"credits"`
	Username       string         `json:"username"`
	ResourceLimits ResourceLimits `json:"resource_limits"`
}

type ResourceLimits struct {
	Cores          int32 `json:"cores"`
	Memory         int32 `json:"memory"`
	Networks       int32 `json:"networks"`
	PublicIpv4     int32 `json:"public_ipv4"`
	PublicIpv6     int32 `json:"public_ipv6"`
	StorageHdd     int32 `json:"storage_hdd"`
	StorageMaxiops int32 `json:"storage_maxiops"`
	StorageSsd     int32 `json:"storage_ssd"`
}

func GetAccount(ctx *context.Context, account *db.Account) (*db.Account, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.upcloud.com/1.3/account", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(account.Username, *account.Password)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	accountAPIResponse := &AccountAPIResponse{}
	err = json.Unmarshal(body, accountAPIResponse)
	if err != nil {
		log.Fatal(err)
	}

	if accountAPIResponse.Account.Credits == 0 {
		credits := 0.0
		account.Credits = &credits
	}
	if *account.Credits != accountAPIResponse.Account.Credits {
		account.Credits = &accountAPIResponse.Account.Credits

	}
	return account, nil
}
