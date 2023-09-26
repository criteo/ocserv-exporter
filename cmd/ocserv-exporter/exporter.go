package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/criteo/ocserv-exporter/lib/occtl"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	interval    time.Duration
	listenAddr  string
	occtlCli    *occtl.Client
	promHandler http.Handler

	users []occtl.UsersMessage

	lock sync.Mutex
}

func NewExporter(occtlCli *occtl.Client, listenAddr string, interval time.Duration) *Exporter {
	return &Exporter{
		listenAddr:  listenAddr,
		interval:    interval,
		occtlCli:    occtlCli,
		promHandler: promhttp.Handler(),
	}
}

func (e *Exporter) Run() error {
	// run once to ensure we have data before starting the server
	e.update()

	go func() {
		for range time.Tick(e.interval) {
			e.update()
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", e.metricsHandler)

	log.Infof("Listening on http://%s", e.listenAddr)
	return http.ListenAndServe(e.listenAddr, mux)
}

func (e *Exporter) update() {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.updateStatus()
	e.updateUsers()
}

func (e *Exporter) updateStatus() {
	status, err := e.occtlCli.ShowStatus()
	if err != nil {
		log.Errorf("Failed to get server status: %v", err)
		occtlStatusScrapeError.WithLabelValues().Inc()
		vpnActiveSessions.Reset()
		vpnHandledSessions.Reset()
		vpnIPsBanned.Reset()
		return
	}
	vpnStartTime.WithLabelValues().Set(float64(status.RawUpSince))
	vpnActiveSessions.WithLabelValues().Set(float64(status.ActiveSessions))
	vpnHandledSessions.WithLabelValues().Set(float64(status.HandledSessions))
	vpnIPsBanned.WithLabelValues().Set(float64(status.IPsBanned))
	vpnTotalAuthenticationFailures.WithLabelValues().Set(float64(status.TotalAuthenticationFailures))
	vpnSessionsHandled.WithLabelValues().Set(float64(status.SessionsHandled))
	vpnTimedOutSessions.WithLabelValues().Set(float64(status.TimedOutSessions))
	vpnTimedOutIdleSessions.WithLabelValues().Set(float64(status.TimedOutIdleSessions))
	vpnClosedDueToErrorSessions.WithLabelValues().Set(float64(status.ClosedDueToErrorSessions))
	vpnAuthenticationFailures.WithLabelValues().Set(float64(status.AuthenticationFailures))
	vpnAverageAuthTime.WithLabelValues().Set(float64(status.RawAverageAuthTime))
	vpnMaxAuthTime.WithLabelValues().Set(float64(status.RawMaxAuthTime))
	vpnAverageSessionTime.WithLabelValues().Set(float64(status.RawAverageSessionTime))
	vpnMaxSessionTime.WithLabelValues().Set(float64(status.RawMaxSessionTime))
	vpnTX.WithLabelValues().Set(float64(status.RawTX))
	vpnRX.WithLabelValues().Set(float64(status.RawRX))
}

func (e *Exporter) updateUsers() {
	e.users = nil

	vpnUserTX.Reset()
	vpnUserRX.Reset()
	users, err := e.occtlCli.ShowUsers()
	if err != nil {
		log.Errorf("Failed to get users details: %v", err)
		occtlUsersScrapeError.WithLabelValues().Inc()
		return
	}

	for _, user := range users {
		vpnUserTX.WithLabelValues(user.Username, user.RemoteIP, user.MTU, user.VPNIPv4, user.VPNIPv6, user.Device).Set(float64(user.RawTX))
		vpnUserRX.WithLabelValues(user.Username, user.RemoteIP, user.MTU, user.VPNIPv4, user.VPNIPv6, user.Device).Set(float64(user.RawRX))
	}
}

func (e *Exporter) metricsHandler(rw http.ResponseWriter, r *http.Request) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.promHandler.ServeHTTP(rw, r)
}
