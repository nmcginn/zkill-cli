package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"time"
)

func listener(c *cli.Context) error {
	api_url := "https://redisq.zkillboard.com/listen.php"
	var zkb map[string]interface{}

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	if c.GlobalBool("insecure") {
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client.Transport = tr
	}
	for {
		//fmt.Println("request")
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

func printKill(zkb map[string]interface{}, c *cli.Context) {
	// first check to make sure we got a populated, valid kill
	if zkb["package"] == nil {
		return // not much else to do, redis occasionally returns nils instead of holding the connection
	}
	kill := zkb["package"].(map[string]interface{})
	km := kill["killmail"].(map[string]interface{})
	zkb = kill["zkb"].(map[string]interface{})
	victim := km["victim"].(map[string]interface{})

	// items to print
	ship := victim["shipType"].(map[string]interface{})["name"].(string)
	alliance := victim["corporation"].(map[string]interface{})["name"].(string)
	if victim["alliance"] != nil {
		alliance = victim["alliance"].(map[string]interface{})["name"].(string)
	}
	value := zkb["totalValue"].(float64)

	kb_green := false
	attackers := km["attackers"].([]interface{})
	for i, _ := range attackers {
		attacker := attackers[i].(map[string]interface{})
		var attacker_corp string

		if attacker["faction"] != nil {
			attacker_corp = attacker["faction"].(map[string]interface{})["name"].(string)
		}
		if attacker["corporation"] != nil {
			attacker_corp = attacker["corporation"].(map[string]interface{})["name"].(string)
		}
		if attacker["alliance"] != nil {
			attacker_corp = attacker["alliance"].(map[string]interface{})["name"].(string)
		}

		if attacker_corp == alliance {
			kb_green = true
			break
		}
	}

	print_str := fmt.Sprintf("%v's %v worth %.2f isk was destroyed\n", alliance, ship, value)
	if alliance == c.String("alliance") {
		color.Red(print_str)
	} else if kb_green {
		color.Green(print_str)
	} else if value >= c.Float64("isk-threshhold") {
		color.Cyan(print_str)
	} else {
		color.White(print_str)
	}
}

type zKill struct {
	Payload struct {
		KillId   float64 `json:"killID"`
		Killmail struct {
			KillId        float64 `json:"killID"`
			KillTime      string  `json:"killTime"`
			AttackerCount float64 `json:"attackerCount"`
			SolarSystem   struct {
				Id   float64 `json:"id"`
				Name string  `json:"name"`
			} `json:"solarSystem"`
			Attackers []interface{}          `json:"attackers"` // TODO: see if we can strong-type this
			Victim    map[string]interface{} `json:"victim"`
		} `json:"killmail"`
		Zkb struct {
			Value  float64 `json:"totalValue"`
			Points float64 `json:"points"`
			Npc    bool    `json:"npc"`
		} `json:"zkb"`
	} `json:"package"`
}
