package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func listener(c *cli.Context) error {
	api_url := "https://redisq.zkillboard.com/listen.php"
	var zkb zKill

	client := http.Client{
		Timeout: time.Duration(20 * time.Second),
	}
	if c.GlobalBool("insecure") {
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client.Transport = tr
	}
	for {
		resp, err := client.Get(api_url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		} else {
			if err := json.Unmarshal(body, &zkb); err != nil {
				return err
			}
			printKill(zkb, c)
		}
	}
	return nil
}

func printKill(z zKill, c *cli.Context) {
	// first check to make sure we got a populated, valid kill
	if z.Payload == nil {
		return // not much else to do, redis occasionally returns nils instead of holding the connection
	}
	kill := z.Payload
	km := kill.Killmail
	zkb := kill.Zkb
	victim := km.Victim

	// items to print
	victim_alliance := ""
	if victim.Corporation != nil {
		victim_alliance = victim.Corporation.Name
	}
	if victim.Alliance != nil {
		victim_alliance = victim.Alliance.Name
	}

	kb_green := false
	for i, _ := range km.Attackers {
		attacker := km.Attackers[i]
		var attacker_corp string

		if attacker.Faction != nil {
			attacker_corp = attacker.Faction.Name
		}
		if attacker.Corporation != nil {
			attacker_corp = attacker.Corporation.Name
		}
		if attacker.Alliance != nil {
			attacker_corp = attacker.Alliance.Name
		}

		if attacker_corp == c.String("alliance") {
			kb_green = true
			break
		}
	}

	print_str := fmt.Sprintf("%v's %v worth %.2f isk was destroyed\n", victim_alliance, victim.Ship.Name, zkb.Value)
	if victim_alliance == c.String("alliance") {
		color.Red(print_str)
	} else if kb_green {
		color.Green(print_str)
	} else if zkb.Value >= c.Float64("isk-threshhold") {
		color.Cyan(print_str)
	} else {
		color.White(print_str)
	}

	if c.GlobalBool("log") {
		log_str := fmt.Sprintf("%v's %v worth %.2f isk was destroyed: https://zkillboard.com/kill/%.f/\n", victim_alliance, victim.Ship.Name, zkb.Value, kill.KillId)
		file, err := os.OpenFile(os.Getenv("HOME")+"/zkill.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err.Error())
		}
		defer file.Close()
		file.WriteString(log_str)
	}
}

type zChar struct {
	Character *struct {
		Id   float64 `json:"id"`
		Name string  `json:"name"`
	} `json:"character"`
	Faction *struct {
		Id   float64 `json:"id"`
		Name string  `json:"name"`
	} `json:"faction"`
	Corporation *struct {
		Id   float64 `json:"id"`
		Name string  `json:"name"`
	} `json:"corporation"`
	Alliance *struct {
		Id   float64 `json:"id"`
		Name string  `json:"name"`
	} `json:"alliance"`
	Ship *struct {
		Id   float64 `json:"id"`
		Name string  `json:"name"`
	} `json:"shipType"`
}

type zKill struct {
	Payload *struct {
		KillId   float64 `json:"killID"`
		Killmail *struct {
			KillId        float64 `json:"killID"`
			KillTime      string  `json:"killTime"`
			AttackerCount float64 `json:"attackerCount"`
			SolarSystem   *struct {
				Id   float64 `json:"id"`
				Name string  `json:"name"`
			} `json:"solarSystem"`
			Attackers []zChar `json:"attackers"`
			Victim    *zChar  `json:"victim"`
		} `json:"killmail"`
		Zkb *struct {
			Value  float64 `json:"totalValue"`
			Points float64 `json:"points"`
			Npc    bool    `json:"npc"`
		} `json:"zkb"`
	} `json:"package"`
}
