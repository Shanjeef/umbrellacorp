# umbrellacorp

## Requirements:
* To run this branch, you'll need to setup Golang: Follow instructions at https://golang.org/doc/install
* If you don't have a Golang workspace folder already, create it and set your $GOPATH environment variable to point to this folder: https://github.com/golang/go/wiki/SettingGOPATH
* Checkout this repo to your $GOPATH/src folder

To Start Go Server: Navigate to your $GOPATH/src/umbreallacorp folder and *go run .*

## Details:
This repo provides functionality to start a Go server that allows you to manage customer details containing: name, person of contact, telephone number, location, number of employees. The goal of this server is to provide support for Umbrella Corp's imaginary sales team to notify potential customers of upcoming rain in their location so that we can pitch umbrella sales.


### Outline of pkgs:
* router: Middleware support for API handlers such as standardized Request, Response structs
* models: Shared data models for the application, namely Customer, Address, Weather
* components: pkg to store business logic related to specific concerns, organized in sub-pkgs
* handlers: contains sub pkgs to support specific REST endpoints
* util: common utility methods

## Whats Missing:
* Docker-ize repo to support running application in a container so that clients don't have to setup Golang locally based on the Requirements section above
* Add persistance layer to store customer records in database
* Offload fetching weather details to a cron task running on a background thread. We shouldn't couple management requests to manage Customers with fetching 3rd party data
* Concurrency support for requests to edit the same record






