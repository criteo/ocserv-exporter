# Prometheus ocserv exporter

This exporter allows to expose statistics from occtl in Prometheus format.
It simply parses the output of `occtl` tool to gather global and user statistics.

# Usage

## Installation

Binaries can be download from the [Github releases](https://github.com/criteo/ocserv-exporter/releases) page.


## Running

Start `ocserv-exporter` as a daemon or from CLI on the same host running `ocserv` (`occtl` must be installed):

```sh
$ ./ocserv-exporter
```

```sh
$ curl -s http://127.0.0.1:8000/metrics
occtl_status_scrape_error_total 4007
occtl_users_scrape_error_total 4007
vpn_active_sessions 29
vpn_authentication_failures 32
vpn_average_auth_time_seconds 0
vpn_average_session_time_seconds 9960
vpn_closed_error_sessions 0
vpn_handled_sessions 589
vpn_ips_banned 0
vpn_max_auth_time_seconds 46
vpn_max_session_time_seconds 43200
vpn_rx_bytes 2.4823e+07
vpn_sessions_handled 139
vpn_timed_out_idle_sessions 0
vpn_timed_out_sessions 0
vpn_total_authentication_failures 171
vpn_tx_bytes 2.84913e+08
vpn_user_rx_bytes{device="vpns0",mtu="1434",remote_ip="203.0.113.0",username="j.doe",vpn_ipv4="10.0.0.1",vpn_ipv6=""} 962053
vpn_user_rx_bytes{device="vpns1",mtu="1434",remote_ip="203.0.113.1",username="a.bob",vpn_ipv4="10.0.0.2",vpn_ipv6=""} 532733
vpn_user_tx_bytes{device="vpns0",mtu="1434",remote_ip="203.0.113.0",username="j.doe",vpn_ipv4="10.0.0.1",vpn_ipv6=""} 3.474418e+06
vpn_user_tx_bytes{device="vpns1",mtu="1434",remote_ip="203.0.113.1",username="a.bob",vpn_ipv4="10.0.0.2",vpn_ipv6=""} 200146
```
## Prometheus Configuration

Example config:
```yaml
scrape_configs:
  - job_name: 'ocserv-exporter'
    scrape_interval: 30s
    static_configs:
      - targets:
        - 127.0.0.1:8000  # The ocserv exporter's real ip:port
    metrics_path: /metrics
```
