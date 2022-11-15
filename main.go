package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gmccue/go-ipset"
)

const IPSETNAME = "whitelist"

func main() {
	iptablesInit()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		clientIP := c.ClientIP()

		fmt.Printf("got client ip: %s", clientIP)
		err := ipsetAppend(clientIP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		}

	})
	r.Run("0.0.0.0:8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func iptablesInit() {
	// ipset create whitelist hash:net
	// iptables -N denywall
	// iptables -A denywall -p tcp --dport 8080 -j DROP
	// iptables -A denywall -p udp -m multiport --dports 450,8080,40243 -j DROP

	// iptables -I denywall -m set --match-set whitelist src -p tcp --destination-port 8080 -j ACCEPT
	// iptables -I denywall -m set --match-set whitelist src -p udp --destination-port 8080 -j ACCEPT
	// iptables -I denywall -m set --match-set whitelist src -p udp --destination-port 450 -j ACCEPT
	// iptables -I denywall -m set --match-set whitelist src -p udp --destination-port 40243 -j ACCEPT
	ipsetS, err := ipset.New()
	if err != nil {
		panic(err)
	}
	ipsetS.Create(IPSETNAME, "hash:net")
}

func ipsetAppend(ip string) error {
	ipsetS, err := ipset.New()
	if err != nil {
		panic(err)
	}
	err = ipsetS.AddUnique(IPSETNAME, ip)
	if err != nil {
		return err
	}
	return nil
}
