package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	*Config
	HueHost string
	Token   string
	httpCli *http.Client
}

func (c *Client) Put(light int, data string) error {
	url := fmt.Sprintf("http://%s/api/%s/lights/%d/state", c.HueHost, c.Token, light)
	r := strings.NewReader(data)
	req, _ := http.NewRequest("PUT", url, r)
	req.Header.Set("User-Agent", "github.com/j18e/hue-lightswitch")
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpCli.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("got status %d", res.StatusCode)
	}
	return nil
}

type Event struct {
	TimeStr string    `json:"time"`
	Time    time.Time `json:"-"`
	Model   string    `json:"model"`
	BtnData string    `json:"data"`
}

func (c *Client) ProcessEvent(bs []byte) error {
	const timeLayout = `2006-01-02 15:04:05.000000`
	var event Event
	if err := json.Unmarshal(bs, &event); err != nil {
		return fmt.Errorf("unmarshaling json: %w", err)
	}
	ts, err := time.Parse(timeLayout, event.TimeStr)
	if err != nil {
		return fmt.Errorf("parsing time: %w", err)
	}
	event.Time = ts
	mapping := c.MatchEvent(event)
	if mapping == nil {
		return nil
	}
	if ts.Sub(mapping.LastFired) < time.Second {
		mapping.LastFired = ts
		return nil
	}
	mapping.LastFired = ts

	log.Infof("button %s/%s pressed", mapping.Switch, mapping.Button)

	for _, l := range mapping.Lights {
		if err := c.Put(l, mapping.Payload); err != nil {
			log.Error("toggling light %d: %s", l, err)
		}
	}
	return nil
}

func (c *Client) MatchEvent(e Event) *Mapping {
	var sw *Switch
	for _, s := range c.Switches {
		if s.Model == e.Model {
			sw = s
			break
		}
	}
	if sw == nil {
		return nil
	}
	var btn *Button
	for _, b := range sw.Buttons {
		if b.Data == e.BtnData {
			btn = b
			break
		}
	}
	if btn == nil {
		return nil
	}
	for _, m := range c.Mappings {
		if m.Switch == sw.Name && m.Button == btn.Name {
			return m
		}
	}
	return nil
}
