# README #

This README would normally document whatever steps are necessary to get your application up and running.

### What is this repository for? ###

This Go Implementation is able to monitor [Loxone Miniserver](http://www.loxone.com) and provide the data to [influxdb](https://www.influxdata.com/influxdb/).

### How do I get set up? ###

* Clone the repo
* Set the GOPATH to the cloned repo directory
* Run make 
* Create a Configuration ( acutally only available via editor ... )
* Start the app

### Todo ###
* Create UI to configure the metrics
* Connect to Miniserver, load the configuration, select Controls to Monitor with default query
* Provide a Dialog to select values of the control
* Editing of existing Configuration