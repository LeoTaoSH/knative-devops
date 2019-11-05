# eventdetails

A Knative Service to send Revision information to an email address.

## How to build

```
docker build -t daisyycguo/eventdetails .
docker push daisyycguo/eventdetails:latest
```

## Run on Knative

```
kn service create --image docker.io/daisyycguo/eventdetails event-details --env EMAIL=daisy.ycguo@gmail.com
```

## Delete from Knative

```
kn service delete event-details
```