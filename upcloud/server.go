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

	serverAPIResponse := &ServerAPIResponse{}
	err = json.Unmarshal(body, serverAPIResponse)
	if err != nil {
		log.Fatal(err)
	}
	for _, serverResponse := range serverAPIResponse.Servers.Server {
		server := &db.Server{}
		server.Uuid = serverResponse.UUID
		if server.CoreNumber != nil && *server.CoreNumber != serverResponse.CoreNumber && serverResponse.CoreNumber == "" {
			server.CoreNumber = &serverResponse.CoreNumber
		}
		if server.Hostname != nil && *server.Hostname != serverResponse.Hostname && serverResponse.Hostname != "" {
			server.Hostname = &serverResponse.Hostname
		}
		if server.License != nil && *server.License != serverResponse.License && serverResponse.License != 0 {
			server.License = &serverResponse.License
		}
		if server.MemoryAmount != nil && *server.MemoryAmount != serverResponse.MemoryAmount && serverResponse.MemoryAmount != "" {
			server.MemoryAmount = &serverResponse.MemoryAmount
		}
		if server.Plan != nil && *server.Plan != serverResponse.Plan && serverResponse.Plan != "" {
			server.Plan = &serverResponse.Plan
		}
		if server.PlanIpV4Bytes != nil && *server.PlanIpV4Bytes != serverResponse.PlanIpv4Bytes && serverResponse.PlanIpv4Bytes != "" {
			server.PlanIpV4Bytes = &serverResponse.PlanIpv4Bytes
		}
		if server.PlanIpV6Bytes != nil && *server.PlanIpV6Bytes != serverResponse.PlanIpv6Bytes && serverResponse.PlanIpv6Bytes != "" {
			server.PlanIpV6Bytes = &serverResponse.PlanIpv6Bytes
		}
		if server.Zone != nil && *server.Zone != serverResponse.Zone && serverResponse.Zone != "" {
			server.Zone = &serverResponse.Zone
		}
		if server.ServerCreationTime != nil && *server.ServerCreationTime != serverResponse.Created && serverResponse.Created != 0 {
			server.ServerCreationTime = &serverResponse.Created
		}
		if server.State != nil && *server.State != serverResponse.State && serverResponse.State != "" {
			server.State = &serverResponse.State
		}

		if server.Tags != nil && server.Tags.Tag != nil && len(*server.Tags.Tag) == len(serverResponse.Tags.Tag) {
			tags := *server.Tags.Tag
			var change bool
			for index, tag := range serverResponse.Tags.Tag {
				if tag != tags[index] {
					change = true
					break
				}
			}
			if change {
				tags := []string{}
				server.Tags.Tag = &tags
			}

		}

		server.AccountId = *account.Id
		servers = append(servers, server)
	}
	return servers, nil
}
