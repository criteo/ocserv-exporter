package occtl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// StatusMessage is a structure to decode JSON returned by "occtl --json show status"
type StatusMessage struct {
	Status                      string `json:"Status"`
	ServerPID                   int64  `json:"Server PID"`
	SecmodPID                   int64  `json:"Sec-mod PID"`
	RawUpSince                  int64  `json:"raw_up_since"`
	Uptime                      int64  `json:"uptime"`
	ActiveSessions              int64  `json:"Active sessions"`
	HandledSessions             int64  `json:"Total sessions"`
	TotalAuthenticationFailures int64  `json:"Total authentication failures"`
	IPsBanned                   int64  `json:"IPs in ban list"`
	RawLastStatsReset           int64  `json:"raw_last_stats_reset"`
	SessionsHandled             int64  `json:"Sessions handled"`
	TimedOutSessions            int64  `json:"Timed out sessions"`
	TimedOutIdleSessions        int64  `json:"Timed out (idle) sessions"`
	ClosedDueToErrorSessions    int64  `json:"Closed due to error sessions"`
	AuthenticationFailures      int64  `json:"Authentication failures"`
	RawAverageAuthTime          int64  `json:"raw_avg_auth_time"`
	RawMaxAuthTime              int64  `json:"raw_max_auth_time"`
	RawAverageSessionTime       int64  `json:"raw_avg_session_time"`
	RawMaxSessionTime           int64  `json:"raw_max_session_time"`
	RawRX                       int64  `json:"raw_rx"`
	RawTX                       int64  `json:"raw_tx"`
	MinMTU                      int16  `json:"Min MTU"`
	MaxMTU                      int16  `json:"Max MTU"`
}

// UsersMessage is a structure to decode JSON returned by "occtl --json show users"
type UsersMessage struct {
	ID        int64  `json:"ID"`
	Username  string `json:"Username"`
	Vhost     string `json:"vhost"`
	Device    string `json:"Device"`
	MTU       string `json:"MTU"`
	RemoteIP  string `json:"Remote IP"`
	VPNIPv4   string `json:"IPv4"`
	VPNIPv6   string `json:"IPv6"`
	RawRX     int64  `json:"RX,string"`
	RawTX     int64  `json:"TX,string"`
	AverageRX string `json:"Average RX"`
	AverageTX string `json:"Average TX"`
}

// Commander is an interface implementing exec commands
type Commander interface {
	Exists() (bool, error)
	RunCommand(args ...string) ([]byte, error)
}

// Client is an helper client to work with occtl.
type Client struct {
	cmd Commander
}

// NewClient creates a Client checking availability of occtl tool.
func NewClient(cmd Commander) (*Client, error) {
	ok, err := cmd.Exists()
	if !ok {
		return nil, errors.New("occtl tool not available on this server")
	} else if err != nil {
		return nil, fmt.Errorf("failed to lookup for occtl tool: %v", err)
	}
	return &Client{cmd: cmd}, nil
}

// Exists checks that occtl is installed on the localhost
func (c *Client) Exists() (bool, error) {
	_, err := exec.LookPath("occtl")
	if err != nil {
		return false, err
	}
	return true, nil
}

// RunCommand runs an occtl command and returns its output
func (c *Client) RunCommand(args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "occtl", args...).Output()
	return out, err
}

// ShowStatus decodes output from "occtl --json show status command"
func (c *Client) ShowStatus() (*StatusMessage, error) {
	status := &StatusMessage{}
	out, err := c.cmd.RunCommand("--json", "-n", "show", "status")
	if err != nil {
		return nil, fmt.Errorf("error while running command %v", err)
	}
	if err := json.Unmarshal(out, status); err != nil {
		return nil, fmt.Errorf("error while decoding json %v", err)
	}
	return status, nil
}

// ShowUsers decodes output from "occtl --json show users command"
func (c *Client) ShowUsers() ([]UsersMessage, error) {
	users := []UsersMessage{}
	out, err := c.cmd.RunCommand("--json", "-n", "show", "users")
	if err != nil {
		return nil, fmt.Errorf("error while running command %v", err)
	}
	if err := json.Unmarshal(out, &users); err != nil {
		return nil, fmt.Errorf("error while decoding json %v", err)
	}
	return users, nil
}
