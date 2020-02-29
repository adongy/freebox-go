package freebox

import "fmt"

type PortForwardingConfig struct {
	ID           int     `json:"id"`
	Enabled      bool    `json:"enabled"`
	IPProto      string  `json:"ip_proto"`
	WANPortStart int     `json:"wan_port_start"`
	WANPortEnd   int     `json:"wan_port_end"`
	LANIp        string  `json:"lan_ip"`
	LANPort      int     `json:"lan_port"`
	Hostname     string  `json:"hostname"`
	Host         LanHost `json:"host"`
	SrcIP        string  `json:"src_ip"`
	Comment      string  `json:"comment"`
}

type PortForwardingListResponse struct {
	APIResponse
	Result []PortForwardingConfig `json:"result"`
}

type PortForwardingResponse struct {
	APIResponse
	Result PortForwardingConfig `json:"result"`
}

type LanHost struct {
	ID                string `json:"id"`
	PrimaryName       string `json:"primary_name"`
	HostType          string `json:"host_type"`
	PrimaryNameManual bool   `json:"primary_name_manual"`
}

func (c *Client) ListPortForwarding() (*PortForwardingListResponse, error) {
	var output *PortForwardingListResponse
	if err := c.Get("fw/redir/", true, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (c *Client) UpdatePortForwarding(id int, config *PortForwardingConfig) (*PortForwardingResponse, error) {
	var output *PortForwardingResponse
	if err := c.Put(fmt.Sprintf("fw/redir/%d", id), true, config, output); err != nil {
		return nil, err
	}

	return output, nil
}

func (c *Client) AddPortForwarding(config *PortForwardingConfig) (*PortForwardingResponse, error) {
	var output *PortForwardingResponse
	if err := c.Post("fw/redir/", true, config, &output); err != nil {
		return nil, err
	}

	return output, nil
}
