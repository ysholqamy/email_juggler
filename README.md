# Email Juggler

  

A simple API written in Golang, for sending emails.

The service provides an abstraction over several external email providers. Emails are routed to the external providers in a roundrobin fashion, fails over to next provider in case an external provider was down.


## Getting started

To run the service locally, update your environment to include the required keys:

    MAILGUN_KEY=YOUR_KEY
    MAILGUN_DOMAIN=YOUR_MAILGUN_DOMAIN
    
    SENDGRID_KEY=YOUR_SENDGRID_KEY
    
then run main.go:

    go run main.go

and you are done!

## Architecture and Design

### Fault tolerance
RoundRobin routing is used between different providers, on failure the request is attempted using the next provider, until the request succeeds or all providers are tried for the given request. 

### Queue backed Provider
Email providers usually use background processing for sending emails. A new implementation of the Provider interface can be submitting the messages to a job queue, while workers process the submitted jobs. 

**Advantages:**
1) Control retry policy.
2) Lower response times.
3) Producers and Consumers can be scaled independently.

**Disadvantages:**

1) Potential client confusion, as successful response will just mean that the message was queued rather than sent to sub provider.


### Adding additional Providers
Rolling additional providers is quite simple. A provider just needs to implement the following interface:

    type Provider interface{
    	func Send(m Message) error
    }

The new provider can be passed to a RoundRobin Provider or even directly to the service.

### Scalability
The service can be scaled either by:

1) Running multiple instances of the whole service concurrently behind a load balancer.
2) Using a background job queue, where the producers and consumers can be scaled independently.


### Message Status
In a better version, the delivery status of each message should be persisted, and clients should be able to query the status of a message in the system.

# Security
Currently, the service does not imploy any security measures like authenticating the sender. A bearer authentication scheme would be suitable.


## Production API
The service is running in production on Heroku at 
https://email-juggler.herokuapp.com/messages.

You can pass the components of the messages such as  `To`,  `From`,  `Subject`, `Text` and Juggler will pass the message to one of its sub providers to send it.


To send via API, run the following:

    curl -XPOST -i https://email-juggler.herokuapp.com/messages \
    -d from="example@email.com" \
    -d to="youssefsholqamy@gmail.com" \
    -d subject="Hello from tests" \
    -d text="yep, running"


Output:

    HTTP/1.1 200 OK 
    Server: Cowboy
    Connection: keep-alive
    Date: Tue, 22 May 2018 07:30:51 GMT
    Content-Length: 26
    Content-Type: text/plain; charset=utf-8
    Via: 1.1 vegur

    message sent successfully
    
The service exposes a single endpoint `/messages`, which accepts only `POST` requests and supports `Content-Type`  of `application/json` and `application/x-www-form-urlencoded`. 

## Documentation

For more information, checkout the [docs](https://godoc.org/github.com/ysholqamy/email_juggler/email).

## Tests
Tests are separated into two folds using build tags.
1) **Internal** which uses mocked providers.
2) **external** which uses real email providers in tests.

To run internal tests:

    go tests ./email_test -tags internal

To run external tests:

    go test ./email_test -tags external

To run all tests:

    go test ./email_test -tags all

## Future Work

1) Refactor the code to use a message queue instead of processing emails synchronously.
2) Add rate limiting.
3) Add additional failover strategies
4) Rewrite the tests.














