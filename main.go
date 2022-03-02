package main

import (
	"bufio"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	print("Connection Checker v1.0 https://dogukanhan.com\n")
	args := os.Args

	if len(args) < 3 {
		directStart()
	}else{
		argumentStart(args)
	}
}

func argumentStart(args  []string){
	style := args[1]

	if style == "-f" {
		fileName := args[2]
		readFile(fileName)
	} else if style == "-u" {
		url := args[2]
		makeRequest(url)
	} else {
		print("Unrecognized Command", style)
	}
}

func directStart(){
	fmt.Println("1 -> Url \n2 -> File \nYour input: ")

	var command string
	fmt.Scanln(&command)


	if command == "1" {
		fmt.Println("Url or domain: ")
		var domain string
		fmt.Scanln(&domain)
		makeRequest(domain)
	}else if command == "2"{
		fmt.Println("File Path: ")
		var fileName string
		fmt.Scanln(&fileName)
		readFile(fileName)
	}else{
		print("Unrecognized Command")
	}

	print("Press enter to exit")
	var waitForInput string
	fmt.Scanln(&waitForInput)
}

func parseDomainName(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
	}

	return u.Hostname()
}

func pingTo(domain string) {

	if strings.Contains(domain, "://") {
		domain = parseDomainName(domain)
	}

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", domain)
	if err != nil {
		fmt.Printf("PING %s -> FAIL %s \n", domain, err)
		return
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("PING %s -> IP Addr: %s receive, RTT: %v\n", domain, addr.String(), rtt)
	}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("PING %s -> FAIL %s \n", domain, err)
	}
}

func makeRequest(url string) {
	pingTo(url)

	if strings.Contains(url,"http"){
		print("HTTP GET ", url, " -> ")

		resp, err := http.Get(url)

		if err == nil {
			print(resp.StatusCode, "\n")
			return
		} else {
			print(" CONNECTION ERROR")
		}

		print("\n")
	}
}

func readFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		makeRequest(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
