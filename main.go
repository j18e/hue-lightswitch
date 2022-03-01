package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	configFile := flag.String("config-file", "config.yml", "path to the config file")
	tokenFile := flag.String("token-file", ".token", "path to the file containing the Hue user token")
	hueHost := flag.String("hue-host", "", "IP/hostname of the Hue bridge")
	flag.Parse()

	if *hueHost == "" {
		return errors.New("flag -hue-host is required")
	}
	tokenBS, err := os.ReadFile(*tokenFile)
	if err != nil {
		return err
	}
	if len(flag.Args()) > 0 {
		arg := flag.Arg(0)
		if arg == "lights" {
			cli := &Client{HueHost: *hueHost, Token: string(tokenBS)}
			return cli.PrintLights()
		}
		return fmt.Errorf("unknown arg %s", arg)
	}

	cfg, err := NewConfig(*configFile)
	if err != nil {
		return err
	}

	cli := &Client{
		Config:  cfg,
		HueHost: *hueHost,
		Token:   string(tokenBS),
		httpCli: &http.Client{Timeout: time.Second * 5},
	}

	signals := []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM}
	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	go func() {
		err := webServer(ctx, ":8080")
		if err != nil {
			log.Errorf("web-server: got %s, shutting down", err)
		}
		cancel()
	}()

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := runRTL(ctx, cli); err != nil {
			log.Error(err)
		}
	}
}

func runRTL(ctx context.Context, cli *Client) error {
	cmd := exec.CommandContext(ctx, "rtl_433", "-F", "json", "-M", "time:usec")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	log.Info("started rtl_433 listener")

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		err := cli.ProcessEvent(scanner.Bytes())
		if err != nil {
			log.Error(err)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func webServer(ctx context.Context, listenAddr string) error {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	chErr := make(chan error, 1)
	go func() {
		chErr <- http.ListenAndServe(listenAddr, r)
	}()
	select {
	case <-ctx.Done():
		return nil
	case err := <-chErr:
		return err
	}
}
