/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/kelseyhightower/envconfig"
	sendgrid "github.com/sendgrid/sendgrid-go"
	mailhelper "github.com/sendgrid/sendgrid-go/helpers/mail"
	v1alpha1 "knative.dev/serving/pkg/apis/serving/v1alpha1"
)

const (
	Subject       = "subject"
	SubjectPrefix = "/apis/serving.knative.dev/v1alpha1/namespaces/default/services/"
)

var (
	emailaddr string
)

type envConfig struct {
	// Sink URL where to send heartbeat cloudevents
	Email string `envconfig:"EMAIL"`
}

/*func mail(content string) {
	log.Println("send email to ", emailaddr)
	auth := smtp.PlainAuth(
		"",
		"apikey",
		"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",
		"smtp.sendgrid.net",
	)
	err := smtp.SendMail(
		"smtp.sendgrid.net:587",
		auth,
		"knative@sample",
		[]string{emailaddr},
		[]byte(content),
	)
	if err != nil {
		log.Fatal(err)
	}
}*/

func _sendmail(content string) {
	log.Println("send email to ", emailaddr)
	from := mailhelper.NewEmail("Knative Demo", "knative@demo.com")
	subject := "Notification"
	to := mailhelper.NewEmail("Knative User", emailaddr)
	plainTextContent := content
	htmlContent := "<strong>" + content + "</strong>"
	message := mailhelper.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient("XXXXXXXXXXXXXXXx")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

type EmailType struct {
	// Address string `json:"address"`
	Content string `json:"content"`
}

func handler(event cloudevents.Event, resp *cloudevents.EventResponse) error {
	var data v1alpha1.Service
	var responseData EmailType
	var subject string
	ctx := event.Context.AsV03()

	log.Printf("Received an event: %s, %s \n", ctx.GetSource(), ctx.GetType())

	err := event.ExtensionAs(Subject, &subject)
	if err != nil {
		log.Printf("Got Subject Error: %s\n", err.Error())
		return err
	}

	log.Printf("Event Subject: %s\n", subject)

	// Test the subject is a Service
	if !strings.HasPrefix(subject, SubjectPrefix) {
		log.Printf("Not Service update, return directly \n")
		return nil
	}

	dataBytes, err := event.DataBytes()
	if err != nil {
		log.Printf("Got Data Error: %s\n", err.Error())
		return err
	}

	json.Unmarshal(dataBytes, &data)

	if data.Status.Conditions == nil {
		log.Printf("Service Status Conditions Nil, return directly \n")
		return nil
	}

	for _, cond := range data.Status.Conditions {
		if cond.Status != "True" {
			log.Printf("Service Status Conditions Not Ready, return directly \n")
			return nil
		}
	}

	url := data.Status.RouteStatusFields.URL.String()
	log.Printf("Service is Ready. url is :%s \n", url)

	responseData.Content = fmt.Sprintf("Your service %s is ready at %s.", data.ObjectMeta.Name, url)
	_sendmail(responseData.Content)

	r := cloudevents.Event{
		Context: ctx,
	}

	if err := event.SetData(responseData); err != nil {
		log.Printf("failed to set data, %v", err)
		os.Exit(1)
	}

	r.SetDataContentType("application/json")

	resp.RespondWith(200, &r)
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

	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("listening on 8080")
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(context.Background(), handler))
}
