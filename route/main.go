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

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/knative/pkg/ptr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	v1 "knative.dev/serving/pkg/apis/serving/v1"
	v1alpha1 "knative.dev/serving/pkg/apis/serving/v1alpha1"
	serving_client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1alpha1"
)

// var (
// 	percentage string
// )

// func init() {
// 	flag.StringVar(&percentage, "percentage", "100", "the percentage of the route")
// }

type DataType struct {
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	ApiVersion string `json:"apiVersion"`
	Sink       string `json:"sink"`
}

func create_route(revision string) *v1alpha1.Route {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := serving_client.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	route := constructRoute(revision)
	route, _ = clientset.Routes("default").Create(route)

	return route
}

// Create service struct from provided options
func constructRoute(revision string) *v1alpha1.Route {

	name := fmt.Sprintf("route100-%s", revision)

	atraffic := v1alpha1.TrafficTarget{
		TrafficTarget: v1.TrafficTarget{
			Percent:        ptr.Int64(100),
			LatestRevision: ptr.Bool(false),
			RevisionName:   revision,
		},
	}
	route := v1alpha1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: v1alpha1.RouteSpec{
			Traffic: []v1alpha1.TrafficTarget{atraffic},
		},
	}

	return &route
}

func handler(event cloudevents.Event, resp *cloudevents.EventResponse) error {
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
	route := create_route(data.Name)
	data.Sink = route.Status.URL.String()

	r := cloudevents.Event{
		Context: ctx,
	}

	if err := event.SetData(data); err != nil {
		log.Printf("failed to set data, %v", err)
		os.Exit(1)
	}

	r.SetDataContentType("application/json")

	log.Println("Transform the event to: ")
	log.Printf("[%s] %s %s: %+v", ctx.Time.String(), ctx.GetSource(), ctx.GetType(), data)

	resp.RespondWith(200, &r)
	return nil
}

func main() {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("listening on 8080")
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(context.Background(), handler))
}
