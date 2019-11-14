/*
Copyright 2019 daisy-ycguo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/kelseyhightower/envconfig"
	"knative.dev/eventing/pkg/kncloudevents"
)

type DataType struct {
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	ApiVersion string `json:"apiVersion"`
	Sink       string `json:"sink"`
}

var (
	emailaddr string
)

type envConfig struct {
	// Sink URL where to send heartbeat cloudevents
	Email string `envconfig:"EMAIL"`
}

func mail(revision string) {
	log.Println("send email revision ", revision, " to ", emailaddr)
	auth := smtp.PlainAuth(
		"",
		"apikey",
		"SG.qt8WVfKKTmSPDGGmFzZ_BA.62XXXXX",
		"smtp.sendgrid.net",
	)
	body := fmt.Sprintf("Subject: Notification\r\n\r\nYour revision %s has been created.", revision)
	err := smtp.SendMail(
		"smtp.sendgrid.net:587",
		auth,
		"knative@sample",
		[]string{emailaddr},
		[]byte(body),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func gotEvent(event cloudevents.Event) error {
	var data DataType
	ctx := event.Context.AsV03()

	dataBytes, err := event.DataBytes()
	if err != nil {
		log.Printf("Got Data Error: %s\n", err.Error())
		return err
	}
	log.Println("Received an event: ")
	log.Printf("[%s] %s %s: %s", ctx.Time.String(), ctx.GetSource(), ctx.GetType(), dataBytes)

	json.Unmarshal(dataBytes, &data)
	mail(data.Name)

	return nil
}

func main() {
	// get variable from ENV
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	if env.Email != "" {
		emailaddr = env.Email
	} else {
		emailaddr = "guoyingc@cn.ibm.com"
	}

	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("listening on 8080")
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(context.Background(), gotEvent))
}
