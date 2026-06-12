package helper

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetAddress() (ipport string, network string) {
	port := os.Getenv("PORT")
	network = "tcp4"
	if port == "" {
		port = ":8080"
	} else if port[0:1] != ":" {
		ip := os.Getenv("IP")
		if ip == "" {
			ipport = ":" + port
		} else {
			if strings.Contains(ip, ".") {
				ipport = ip + ":" + port
			} else {
				ipport = "[" + ip + "]" + ":" + port
				network = "tcp6"
			}
		}
	}
	return
}

func GetIPaddress() string {
	resp, err := http.Get("https://icanhazip.com/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func SRVLookup(srvuri string) (mongouri string) {
	atsplits := strings.Split(srvuri, "@")
	if len(atsplits) < 2 {
		return srvuri
	}
	userpassSplits := strings.Split(atsplits[0], "//")
	if len(userpassSplits) < 2 {
		return srvuri
	}
	userpass := userpassSplits[1]
	mongouri = "mongodb://" + userpass + "@"
	slashsplits := strings.Split(atsplits[1], "/")
	if len(slashsplits) < 2 {
		return srvuri
	}
	domain := slashsplits[0]
	path := slashsplits[1]
	dbname := path
	params := ""
	if strings.Contains(path, "?") {
		psplits := strings.Split(path, "?")
		dbname = psplits[0]
		params = psplits[1]
	}

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "8.8.8.8:53")
		},
	}
	_, srvs, err := r.LookupSRV(context.Background(), "mongodb", "tcp", domain)
	if err != nil {
		return srvuri
	}
	var srvlist string
	for _, srv := range srvs {
		srvlist += strings.TrimSuffix(srv.Target, ".") + ":" + strconv.FormatUint(uint64(srv.Port), 10) + ","
	}

	txtrecords, _ := r.LookupTXT(context.Background(), domain)
	var txtlist string
	for _, txt := range txtrecords {
		txtlist += txt
	}
	if params != "" {
		params = params + "&"
	}
	mongouri = mongouri + strings.TrimSuffix(srvlist, ",") + "/" + dbname + "?ssl=true&" + params + txtlist
	return
}