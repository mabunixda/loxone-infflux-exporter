# README #

This README would normally document whatever steps are necessary to get your application up and running.

### What is this repository for? ###

This Go Implementation is able to monitor [Loxone Miniserver](http://www.loxone.com) and provide the data to [prometheus](https://prometheus.io/).

### How do I get set up? ###

* Clone the repo
* Set the GOPATH to the cloned repo directory
* From bitbucket-pipelines run the commands under the "script" node
* Create a Configuration ( acutally only available via editor - UI is coming ... )
* Start the app
* Configure prometheus to scrape the http://$HOSTNAME:8080/metrics 

### Todo ###
* Create UI to configure the metrics
* Connect to Miniserver, load the configuration, select Controls to Monitor with default query
* Provide a Dialog to select values of the control
* Editing of existing Configuration