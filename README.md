# datadog-home-project

Website availability &amp; performance monitoring tool
Alert# Website availability &amp; performance monitoring tool

_Overview_

- A console program to monitor the performance and availability of websites
- Websites and check intervals are user-defined

_Statistics_

- Check the different websites with their corresponding check intervals
- Compute a few interesting metrics: availability, max/avg response times, max/avg time to first byte, response codes count

_Alerting_

- When a website availability is below a user-defined threshold for a user-defined interval, an alert message is created: "Website {website} is down. availability={availability}, time={time}" (default config threshold: 80%, interval: 2min)
- When availability resumes, another message is created detailing when the alert recovered

_Dashboard_

- displays stats for a user-defined timeframe, stats are updated following a user-defined interval. Default:
  - Every 10s, display the stats for the past 10 minutes for each website
  - Every minute displays the stats for the past hour for each website
- Show all past alerting messages

### Requirements

- [InfluxDB 2.0](https://www.influxdata.com/) - open source time series database
- [Go 1.14](https://golang.org/) - a systems programming language

### Installation

#### Building from source

Build your Go app:

```sh
$ go build
```

Run an InfluxDB instance:

```sh
$ docker run -p 8086:8086 -v influxdb:/var/lib/influxdb influxdb
```

Run your build file:

```sh
$ ./datadog-home-project
```

### Docker compose

Dillinger is very easy to install and deploy in a Docker container.

By default, the Docker will expose port 8080, so change this within the Dockerfile if necessary. When ready, simply use the Dockerfile to build the image.

```sh
cd dillinger
docker build -t joemccann/dillinger:${package.json.version} .
```

This will create the dillinger image and pull in the necessary dependencies. Be sure to swap out `${package.json.version}` with the actual version of Dillinger.

Once done, run the Docker image and map the port to whatever you wish on your host. In this example, we simply map port 8000 of the host to port 8080 of the Docker (or whatever port was exposed in the Dockerfile):

```sh
docker run -d -p 8000:8080 --restart="always" <youruser>/dillinger:${package.json.version}
```

Verify the deployment by navigating to your server address in your preferred browser.

```sh
127.0.0.1:8000
```

### Possible improvements

- Write MORE Tests
- Add Night Mode

## License

MIT

**Free Software, Hell Yeah!**
