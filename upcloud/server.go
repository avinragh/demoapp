package upcloud

import (
	"demoapp/db"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ServerAPIResponse struct {
	Servers ServersResponse `json:"servers"`
}

type ServersResponse struct {
	Server []ServerResponse `json:"server"`
}

type ServerResponse struct {
	CoreNumber    string `json:"core_number"`
	Created       int32  `json:"created"`
	Hostname      string `json:"hostname"`
	License       int32  `json:"license"`
	MemoryAmount  string `json:"memory_amount"`
	Plan          string `json:"plan"`
	PlanIpv4Bytes string `json:"plan_ipv4_bytes"`
	PlanIpv6Bytes string `json:"plan_ipv6_bytes"`
	SimpleBackup  string `json:"simple_backup"`
	State         string `json:"state"`
	Tags          struct {
		Tag []string `json:"tag"`
	} `json:"tags"`
	Title string `json:"title"`
	UUID  string `json:"uuid"`
	Zone  string `json:"zone"`
}

func GetServers(account *db.Account) ([]*db.Server, error) {
	servers := []*db.Server{}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.upcloud.com/1.3/server", nil)
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

	serverAPIResponse := ServerAPIResponse{}
	err = json.Unmarshal(body, &serverAPIResponse)
	if err != nil {
		log.Fatal(err)
	}
	if len(serverAPIResponse.Servers.Server) > 0 {
		for _, serverResponse := range serverAPIResponse.Servers.Server {
			server := db.Server{}
			server.Uuid = serverResponse.UUID
			server.CoreNumber = &serverResponse.CoreNumber
			server.Hostname = &serverResponse.Hostname
			server.License = &serverResponse.License
			server.MemoryAmount = &serverResponse.MemoryAmount
			server.Plan = &serverResponse.Plan
			server.PlanIpV4Bytes = &serverResponse.PlanIpv4Bytes
			server.PlanIpV6Bytes = &serverResponse.PlanIpv6Bytes
			server.Zone = &serverResponse.Zone
			server.ServerCreationTime = &serverResponse.Created
			server.State = &serverResponse.State
			tags := []string{}
			tags = append(tags, serverResponse.Tags.Tag...)
			tagsWrapper := &db.Tags{}
			tagsWrapper.Tag = &tags
			server.Tags = tagsWrapper
			server.AccountId = *account.Id
			servers = append(servers, &server)

		}
	}

	return servers, nil
}
