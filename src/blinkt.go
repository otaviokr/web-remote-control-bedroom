package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/otaviokr/web-remote-control-bedroom/mq"
	log "github.com/sirupsen/logrus"
)

// Pixel represents a LED in Blinkt (which has a total of 8 LEDs)
type Pixel struct {
	R          int
	G          int
	B          int
	Brightness float64
}

// ParseColor will convert the hexadecimal in the string into a valid int.
func ParseColor(c string) int {
	r, err := strconv.ParseInt(c, 16, 32)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ColorComponent": c,
				"error": err.Error(),
			},
		).Error("Could not parse the color component")
		return 0
	}
	return int(r)
}

// SetPixel will parse the incomming details for a LED (a "pixel") to configure it correctly.
func SetPixel(led, rgb string) (int, int, int, int) {
	index, err := strconv.Atoi(led[len(led)-1:])
	if err != nil {
		log.WithFields(
			log.Fields{
				"led": led,
				"rgb": rgb,
				"error": err.Error(),
			},
		).Error("Could not convert LED index")
		return -1, 0, 0, 0
	}

	re := regexp.MustCompile("%23([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})")
	match := re.FindStringSubmatch(rgb)

	red := 0
	green := 0
	blue := 0
	if len(match) == 4 {
		red = ParseColor(match[1])
		green = ParseColor(match[2])
		blue = ParseColor(match[3])
	} else {
		log.WithFields(
			log.Fields{
				"led": led,
				"rgb": rgb,
			},
		).Error("Could not parse RGB. Regex failed.")
	}

	return index - 1, red, green, blue
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.WithFields(
		log.Fields{
			"body": request.Body,
		}).Debug("Request received")

	p := []Pixel{{}, {}, {}, {}, {}, {}, {}, {}}

	for _, entry := range strings.Split(request.Body , "&"){
		entrySplitted := strings.Split(entry, "=")
		name := entrySplitted[0]
		value := entrySplitted[1]
		log.WithFields(
			log.Fields{
				"name": name,
				"value": value,
				"entry": entry,
			}).Debug("Parameter parsed")

		if name[:5] == "input" {
			// Field with name ending with "b" contains brights for that led
			if name[len(name)-1:] == "b" {
				br, err := strconv.Atoi(value)
				if err != nil {
					log.WithFields(
						log.Fields{
							"name": name,
							"value": value,
							"error": err.Error(),
						}).Error("Failed to parse brightness value")
				}
				brightness := float64(br) / 100
				led, err := strconv.Atoi(name[7:len(name)-1])
				if err != nil {
					log.WithFields(
						log.Fields{
							"name": name,
							"led": led,
							"error": err.Error(),
						}).Error("Failed to parse led identifier value")
				}
				led -= 1
				log.WithFields(
					log.Fields{
						"led":        led,
						"brightness": brightness,
					},
				).Info("Set new value to brightness to LED")
				p[led].Brightness = brightness
			} else {
				i, r, g, b := SetPixel(name, value)
				log.WithFields(
					log.Fields{
						"LedIndex": i,
						"Red":      r,
						"Green":    g,
						"Blue":     b,
					},
				).Info("Set new color to LED")
				p[i].R = r
				p[i].G = g
				p[i].B = b
			}
		}
	}

	str := ""
	for _, pixel := range p {
		str = fmt.Sprintf("%s #%d#%d#%d#%f", str, pixel.R, pixel.G, pixel.B, pixel.Brightness)
	}
	log.WithFields(
		log.Fields{
			"Message": str,
		}).Debug("Message ready to be sent to MQ")

	broker := "broker.hivemq.com"
	port := 1883

	mqPublisher, err := mq.NewPublisher(broker, port)
	if err != nil {
		log.WithFields(
			log.Fields{
				"broker": broker,
				"port": port,
				"error": err.Error(),
			}).Errorf("Failed to instantiate new MQ Publisher")
		return &events.APIGatewayProxyResponse {
			StatusCode: 200,
		}, nil
	}

	topic := "test_okr"
	log.WithFields(
		log.Fields{
			"topic": topic,
			"broker": broker,
			"port": port,
		}).Info("Sending to MQ")
	mqPublisher.Publish(topic, str)

	return &events.APIGatewayProxyResponse {
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}