package occtl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	validShowStatusJSONOutput = []byte(`{
		"Status":  "online",
		"Server PID":  13385,
		"Sec-mod PID":  13387,
		"Up since":  "2020-05-12 20:50",
		"_Up since":  "55days",
		"raw_up_since":  1589316656,
		"uptime":  4824218,
		"Active sessions":  12,
		"Total sessions":  2251,
		"Total authentication failures":  360,
		"IPs in ban list":  0,
		"Last stats reset":  "2020-07-01 12:37",
		"_Last stats reset":  " 6days",
		"raw_last_stats_reset":  1593607066,
		"Sessions handled":  183,
		"Timed out sessions":  13,
		"Timed out (idle) sessions":  0,
		"Closed due to error sessions":  23,
		"Authentication failures":  73,
		"Average auth time":  "    0s",
		"raw_avg_auth_time":  0,
		"Max auth time":  " 1m:34s",
		"raw_max_auth_time":  94,
		"Average session time":  " 2h:48m",
		"raw_avg_session_time":  10080,
		"Max session time":  "10h:00m",
		"raw_max_session_time":  36000,
		"Min MTU":  1340,
		"Max MTU":  1434,
		"RX":  "2.8 GB",
		"raw_rx":  2786778000,
		"TX":  "69.2 GB",
		"raw_tx":  69237907000
	}`)
	validStatusMessage = &StatusMessage{
		Status:                      "online",
		ServerPID:                   13385,
		SecmodPID:                   13387,
		RawUpSince:                  1589316656,
		Uptime:                      4824218,
		ActiveSessions:              12,
		HandledSessions:             2251,
		TotalAuthenticationFailures: 360,
		IPsBanned:                   0,
		RawLastStatsReset:           1593607066,
		SessionsHandled:             183,
		TimedOutSessions:            13,
		TimedOutIdleSessions:        0,
		ClosedDueToErrorSessions:    23,
		AuthenticationFailures:      73,
		RawAverageAuthTime:          0,
		RawMaxAuthTime:              94,
		RawAverageSessionTime:       10080,
		RawMaxSessionTime:           36000,
		RawRX:                       2786778000,
		RawTX:                       69237907000,
		MinMTU:                      1340,
		MaxMTU:                      1434,
	}

	validShowUsersJSONOutput = []byte(`
	[
		{
			"ID": 4100,
			"Username": "alice",
			"Groupname": "Domain Users",
			"State": "connected",
			"vhost": "default",
			"Device": "vpns9",
			"MTU": "1434",
			"Remote IP": "192.0.2.1",
			"Location": "unknown",
			"Local Device IP": "192.0.2.2",
			"IPv4": "198.51.100.59",
			"P-t-P IPv4": "198.51.100.1",
			"User-Agent": "Open AnyConnect VPN Agent v8.10",
			"RX": "9447052",
			"TX": "328664295",
			"_RX": "9.4 MB",
			"_TX": "328.7 MB",
			"Average RX": "316 bytes/sec",
			"Average TX": "11.0 KB/sec",
			"DPD": "90",
			"KeepAlive": "32400",
			"Hostname": "laptop",
			"Connected at": "2020-07-06 08:08",
			"_Connected at": " 8h:17m",
			"Full session": "930gAhoVKQRcCIOf34rnpRj9Pyg=",
			"Session": "930gAh",
			"TLS ciphersuite": "(TLS1.2)-(ECDHE-ECDSA-SECP384R1)-(AES-256-GCM)",
			"DTLS cipher": "(DTLS1.2)-(PSK)-(AES-256-GCM)",
			"DNS": [
				"10.1.8.53"
			],
			"NBNS": [],
			"Split-DNS-Domains": [
				"example.com"
			],
			"Routes": [
				"10.0.0.0/255.0.0.0"
			],
			"No-routes": [],
			"iRoutes": [],
			"Restricted to routes": "False",
			"Restricted to ports": []
		},
		{
			"ID": 15313,
			"Username": "bob",
			"Groupname": "Domain Users",
			"State": "connected",
			"vhost": "default",
			"Device": "vpns2",
			"MTU": "1340",
			"Remote IP": "192.0.2.2",
			"Location": "unknown",
			"Local Device IP": "192.0.2.2",
			"IPv4": "198.51.100.41",
			"P-t-P IPv4": "198.51.100.1",
			"User-Agent": "OpenConnect-GUI 1.5.3 v7.08",
			"RX": "351042944",
			"TX": "18291460815",
			"_RX": "351.0 MB",
			"_TX": "18.3 GB",
			"Average RX": "10.2 KB/sec",
			"Average TX": "534.1 KB/sec",
			"DPD": "90",
			"KeepAlive": "32400",
			"Hostname": "localhost",
			"Connected at": "2020-07-06 06:55",
			"_Connected at": " 9h:30m",
			"Full session": "8HxOQEoXrjegoTkwioLK41ceims=",
			"Session": "8HxOQE",
			"TLS ciphersuite": "(TLS1.2)-(ECDHE-ECDSA-SECP384R1)-(AES-256-GCM)",
			"DTLS cipher": "(DTLS1.2)-(PSK)-(AES-256-GCM)",
			"DNS": [
				"10.1.8.53"
			],
			"NBNS": [],
			"Split-DNS-Domains": [
				"example.com"
			],
			"Routes": [
				"10.0.0.0/255.0.0.0"
			],
			"No-routes": [],
			"iRoutes": [],
			"Restricted to routes": "False",
			"Restricted to ports": []
		}
	]
	`)
	validShowUsersMessage = []UsersMessage{
		{
			ID:        4100,
			Username:  "alice",
			Vhost:     "default",
			Device:    "vpns9",
			MTU:       "1434",
			RemoteIP:  "192.0.2.1",
			VPNIPv4:   "198.51.100.59",
			VPNIPv6:   "",
			RawRX:     9447052,
			RawTX:     328664295,
			AverageRX: "316 bytes/sec",
			AverageTX: "11.0 KB/sec",
		},
		{
			ID:        15313,
			Username:  "bob",
			Device:    "vpns2",
			Vhost:     "default",
			MTU:       "1340",
			RemoteIP:  "192.0.2.2",
			VPNIPv4:   "198.51.100.41",
			VPNIPv6:   "",
			RawRX:     351042944,
			RawTX:     18291460815,
			AverageRX: "10.2 KB/sec",
			AverageTX: "534.1 KB/sec",
		},
	}
)

// fakeClient allows to mock Exists() and RunCommand() func
type fakeClient struct{}

func (c *fakeClient) Exists() (bool, error) {
	return true, nil
}

func (c *fakeClient) RunCommand(args ...string) ([]byte, error) {
	var output []byte
	command := strings.Join(args, " ")
	switch command {
	case "--json -n show status":
		output = validShowStatusJSONOutput
	case "--json -n show users":
		output = validShowUsersJSONOutput
	default:
		return nil, fmt.Errorf("command not implemented %s", command)
	}
	return output, nil
}

func TestShowStatus(t *testing.T) {
	c, err := NewClient((&fakeClient{}))
	require.NoError(t, err)
	got, err := c.ShowStatus()
	require.NoError(t, err)
	require.Equal(t, validStatusMessage, got)
}

func TestShowUsers(t *testing.T) {
	c, err := NewClient((&fakeClient{}))
	require.NoError(t, err)
	got, err := c.ShowUsers()
	require.NoError(t, err)
	require.Equal(t, validShowUsersMessage, got)
}
