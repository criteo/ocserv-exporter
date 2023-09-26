package main

import (
	"flag"
	"time"

	"github.com/criteo/ocserv-exporter/lib/occtl"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	occtlStatusScrapeError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "occtl_status_scrape_error_total",
		Help: "Total number of errors that occurred when calling occtl show status.",
	}, []string{})
	occtlUsersScrapeError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "occtl_users_scrape_error_total",
		Help: "Total number of errors that occurred when calling occtl show users.",
	}, []string{})
	vpnStartTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_start_time_seconds",
		Help: "Start time of ocserv since unix epoch in seconds.",
	}, []string{})
	vpnActiveSessions = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_active_sessions",
		Help: "Current number of users connected.",
	}, []string{})
	vpnHandledSessions = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_handled_sessions",
		Help: "Total number of sessions handled since server is up.",
	}, []string{})
	vpnIPsBanned = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_ips_banned",
		Help: "Total number of IPs banned.",
	}, []string{})
	vpnTotalAuthenticationFailures = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_total_authentication_failures",
		Help: "Total number of authentication failures since server is up.",
	}, []string{})
	vpnSessionsHandled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_sessions_handled",
		Help: "Total number of sessions handled since last stats reset.",
	}, []string{})
	vpnTimedOutSessions = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_timed_out_sessions",
		Help: "Total number of timed out sessions since last stats reset.",
	}, []string{})
	vpnTimedOutIdleSessions = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_timed_out_idle_sessions",
		Help: "Total number of sessions timed out (idle) since last stats reset.",
	}, []string{})
	vpnClosedDueToErrorSessions = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_closed_error_sessions",
		Help: "Total number of sessions closed due to error since last stats reset.",
	}, []string{})
	vpnAuthenticationFailures = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_authentication_failures",
		Help: "Total number of authentication failures since last stats reset.",
	}, []string{})
	vpnAverageAuthTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_average_auth_time_seconds",
		Help: "Average time in seconds spent to authenticate users since last stats reset.",
	}, []string{})
	vpnMaxAuthTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_max_auth_time_seconds",
		Help: "Maximum time in seconds spent to authenticate users since last stats reset.",
	}, []string{})
	vpnAverageSessionTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_average_session_time_seconds",
		Help: "Average session time in seconds since last stats reset.",
	}, []string{})
	vpnMaxSessionTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_max_session_time_seconds",
		Help: "Max session time in seconds since last stats reset.",
	}, []string{})
	vpnTX = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_tx_bytes",
		Help: "Total TX usage in bytes since last stats reset.",
	}, []string{})
	vpnRX = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_rx_bytes",
		Help: "Total RX usage in bytes since last stats reset.",
	}, []string{})
	vpnUserTX = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_user_tx_bytes",
		Help: "Total TX usage in bytes of a user.",
	}, []string{"username", "remote_ip", "mtu", "vpn_ipv4", "vpn_ipv6", "device"})
	vpnUserRX = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_user_rx_bytes",
		Help: "Total RX usage in bytes of a user.",
	}, []string{"username", "remote_ip", "mtu", "vpn_ipv4", "vpn_ipv6", "device"})
	vpnUserStartTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpn_user_start_time_seconds",
		Help: "Start time of user session since unix epoch in seconds.",
	}, []string{"username", "remote_ip", "mtu", "vpn_ipv4", "vpn_ipv6", "device"})
)

func main() {
	var (
		interval = flag.Duration("interval", 30*time.Second, "Delay between occtl scrape.")
		listen   = flag.String("listen", "127.0.0.1:8000", "Prometheus HTTP listen IP and port.")
	)
	flag.Parse()

	prometheus.MustRegister(
		occtlStatusScrapeError,
		occtlUsersScrapeError,
		vpnStartTime,
		vpnActiveSessions,
		vpnHandledSessions,
		vpnIPsBanned,
		vpnTotalAuthenticationFailures,
		vpnSessionsHandled,
		vpnTimedOutSessions,
		vpnTimedOutIdleSessions,
		vpnClosedDueToErrorSessions,
		vpnAuthenticationFailures,
		vpnAverageAuthTime,
		vpnMaxAuthTime,
		vpnAverageSessionTime,
		vpnMaxSessionTime,
		vpnTX,
		vpnRX,
		vpnUserTX,
		vpnUserRX,
		vpnUserStartTime,
	)

	occtlCli, err := occtl.NewClient(&occtl.OcctlCommander{})
	if err != nil {
		log.Fatalf("Failed to initialize occtl client: %v", err)
	}

	exporter := NewExporter(occtlCli, *listen, *interval)
	err = exporter.Run()
	if err != nil {
		log.Fatal(err)
	}
}
