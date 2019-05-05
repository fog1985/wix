# Simple TCP server written in GO:

```
tklym@Tarass-MacBook-Pro-2 wix_test $ tree
.
├── conctcp.go
└── sub
    └── morse.go

1 directory, 2 files

Concurent TCP server sources: https://github.com/mactsouk/opensource.com/blob/master/concTCP.go

tklym@Tarass-MacBook-Pro-2 wix_test $ cat ./conctcp.go
package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
        "./sub"
)

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

        result := morse.EncodeITU(strings.Split(c.RemoteAddr().String(), ":")[0]) + "\n"
        c.Write([]byte(string(result)))
        c.Close()
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
```
*MORSE string encoder sources:* `https://github.com/martinlindhe/morse`
```
tklym@Tarass-MacBook-Pro-2 wix_test $ cat ./sub/morse.go
package morse

import (
	"fmt"
	"strings"
)

// Decode decodes morse code in `s` using `alphabet` mapping
func Decode(s string, alphabet map[string]string, letterSeparator string, wordSeparator string) (string, error) {

	res := ""

	for _, part := range strings.Split(s, letterSeparator) {
		found := false
		for key, val := range alphabet {
			if val == part {
				res += key
				found = true
				break
			}
		}
		if part == wordSeparator {
			res += " "
			found = true
		}
		if found == false {
			return res, fmt.Errorf("unknown character " + part)
		}
	}
	return res, nil
}

// Encode encodes clear text in `s` using `alphabet` mapping
func Encode(s string, alphabet map[string]string, letterSeparator string, wordSeparator string) string {

	res := ""

	for _, part := range s {
		p := string(part)
		if p == " " {
			if wordSeparator != "" {
				res += wordSeparator + letterSeparator
			}
		} else if morseITU[p] != "" {
			res += morseITU[p] + letterSeparator
		}
	}
	return strings.TrimSpace(res)
}

// DecodeITU translates international morse code (ITU) to text
func DecodeITU(s string) (string, error) {
	return Decode(s, morseITU, " ", "/")
}

// EncodeITU translates text to international morse code (ITU)
func EncodeITU(s string) string {
	return Encode(s, morseITU, " ", "/")
}

// LooksLikeMorse returns true if string seems to be a morse encoded string
func LooksLikeMorse(s string) bool {

	if len(s) < 1 {
		return false
	}
	for _, b := range s {
		if b != '-' && b != '.' && b != ' ' {
			return false
		}
	}
	return true
}

var (
	morseITU = map[string]string{
		"a":  ".-",
		"b":  "-...",
		"c":  "-.-.",
		"d":  "-..",
		"e":  ".",
		"f":  "..-.",
		"g":  "--.",
		"h":  "....",
		"i":  "..",
		"j":  ".---",
		"k":  "-.-",
		"l":  ".-..",
		"m":  "--",
		"n":  "-.",
		"o":  "---",
		"p":  ".--.",
		"q":  "--.-",
		"r":  ".-.",
		"s":  "...",
		"t":  "-",
		"u":  "..-",
		"v":  "...-",
		"w":  ".--",
		"x":  "-..-",
		"y":  "-.--",
		"z":  "--..",
		"ä":  ".-.-",
		"ö":  "---.",
		"ü":  "..--",
		"Ch": "----",
		"0":  "-----",
		"1":  ".----",
		"2":  "..---",
		"3":  "...--",
		"4":  "....-",
		"5":  ".....",
		"6":  "-....",
		"7":  "--...",
		"8":  "---..",
		"9":  "----.",
		".":  ".-.-.-",
		",":  "--..--",
		"?":  "..--..",
		"!":  "..--.",
		":":  "---...",
		"\"": ".-..-.",
		"'":  ".----.",
		"=":  "-...-",
	}
)
```
Тетсовий ран:
```
tklym@Tarass-MacBook-Pro-2 Repositories $ nc localhost 8001
.---- ..--- --... .-.-.- ----- .-.-.- ----- .-.-.- .----
tklym@Tarass-MacBook-Pro-2 Repositories $
```
Також перевірив пачкою:
```
tklym@Tarass-MacBook-Pro-2 Repositories $ time for i in {1..1000}; do nc localhost 8001 > /dev/null; done

real	0m12.041s
user	0m4.373s
sys	0m5.074s
```
Тобто 90 реквестів в секунду.
Якщо на ноді(інстанці) заранити 10+ подів то легко буде throughput в 1000+ конекшенів в секунду з одного хоста.

Під час тесту перевірив що малтіпл конекшенів створено а не один ))
```
Every 2.0s: netstat -anp TCP | grep 8001                          Tarass-MacBook-Pro-2.local: Wed May  1 21:37:18 2019

tcp4       0      0  *.8001                 *.*                    LISTEN
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62322        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62323        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62324        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62325        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62326        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62327        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62328        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62329        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62330        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62331        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62332        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62333        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62334        TIME_WAIT
tcp4       0      0  127.0.0.1.8001         127.0.0.1.62335        TIME_WAIT
...
```
Звісно потрібно збілдати сорси в бінарник щоб ранити не юзаючи GO а як екзекютал скрипт.
Обрано метод багатопоточного рана. Щоб хендлити велику кількість реквестів.
На кожен реквест створюється окремий тред.

# Terraform part:
```
// Providers:

# Default Region
provider "aws" {
  region  = "us-east-1"
  version = "~> 1.8"
}

provider "aws" {
  alias  = "us-east-1"
  region = "us-east-1"
}

provider "aws" {
  alias  = "us-east-2"
  region = "us-east-2"
}

provider "aws" {
  alias  = "us-west-1"
  region = "us-west-1"
}

provider "aws" {
  alias  = "us-west-2"
  region = "us-west-2"
}

provider "aws" {
  alias  = "eu-central-1"
  region = "eu-central-1"
}

// VPC:

resource "aws_vpc" "primary" {
  cidr_block           = "10.6.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  providers = { // This will create resource in needed providet for example.
    "aws.region" = "aws.us-east-1"
  }
}

// With 'providers' block within resources we should be able to create resource in any needed reagion. Of corse we will have to double the code or move into modules usage to make code quantity less and better readable and maintable.

resource "aws_elb" "main" {
  name               = "foobar-terraform-elb"
  availability_zones = ["us-east-1a","us-east-1b","us-east-1c"]

  listener {
    instance_port     = 8001
    instance_protocol = "tcp"
    lb_port           = 8001
    lb_protocol       = "tcp"
  }
}

resource "aws_route53_zone" "primary" {
  name = "example.com"
  vpc {
    vpc_id = "${aws_vpc.primary.id}" // 
  }
}

resource "aws_route53_record" "primary-ns" {
  zone_id = "${aws_route53_zone.main.zone_id}"
  name    = "primary.example.com"
  type    = "NS"
  ttl     = "30"

  records = [
    "${aws_route53_zone.primary.name_servers.0}",
    "${aws_route53_zone.primary.name_servers.1}",
    "${aws_route53_zone.primary.name_servers.2}",
    "${aws_route53_zone.primary.name_servers.3}",
  ]
}

resource "aws_route53_record" "www" {
  zone_id = "${aws_route53_zone.primary.zone_id}"
  name    = "example.com"
  type    = "A"

  alias {
    name                   = "${aws_elb.main.dns_name}"
    zone_id                = "${aws_elb.main.zone_id}"
    evaluate_target_health = true
  }
  geolocation_routing_policy = true // this one should do the trick with max speed from any region
}

// now create few records for balancers in different regions:

resource "aws_route53_record" "balanced" {
  zone_id = "${aws_route53_zone.primary.zone_id}"
  name    = "a-example.com"
  type    = "CNAME"

  alias {
    name                   = "${aws_elb.main.dns_name}"
    zone_id                = "${aws_elb.main.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "balanced" {
  zone_id = "${aws_route53_zone.primary.zone_id}"
  name    = "b-example.com"
  type    = "CNAME"

  alias {
    name                   = "${aws_elb.main.dns_name}"
    zone_id                = "${aws_elb.main.zone_id}"
    evaluate_target_health = true
  }
}

// For ASG I suggest using already battle-tested cloudposse code-base:
https://github.com/cloudposse/terraform-aws-ec2-autoscale-group

// For ELB:

resource "aws_elb" "tcpconn" {
  name               = "foobar-terraform-elb"
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

  access_logs {
    bucket        = "foo"
    bucket_prefix = "bar"
    interval      = 60
  }

  listener {
    instance_port     = 8001
    instance_protocol = "tcp"
    lb_port           = 8001
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "TCP:8001/"
    interval            = 30
  }

  instances                   = ["${aws_instance.foo.id}"] // here ASG to be put from Cloudposse module.
  cross_zone_load_balancing   = true
  idle_timeout                = 400
  connection_draining         = true
  connection_draining_timeout = 400

  tags = {
    Name = "foobar-terraform-elb"
  }
}
```
// I will omit security group configuration for both ASG and ELB. IN real env this is vital as well.
NACLs for VPC the same.

// I also suggest using Terragrunt for TF's state locking facility using S3 and DynamoDB. Will help working in a team and to not mess up the stuff.

For Google cloud whole code will be different. Sorry but will not gp into deep. :)
This task might be a lot bigger if we want to take under consioderation all the things and culprits. 

Kindest regards,
Taras.
