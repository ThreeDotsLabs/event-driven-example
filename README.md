# Event-Driven application example

This is an example Event-Driven application written in Go, using [Watermill](https://github.com/ThreeDotsLabs/watermill).

The projects aims to integrate incoming GitHub webhooks with Grafana and Slack, essentially adding annotations and
sending messages when a new commit is pushed. There are also simulated deployment messages sent over RabbitMQ to
demonstrate working with multiple event streams.

![](https://threedots.tech/media/event-driven-applications/diagram.png)

An example result can look like this: 

![](https://threedots.tech/media/event-driven-applications/grafana.png)

![](https://threedots.tech/media/event-driven-applications/slack.png)

## Running

If you'd like to integrate the example with your Slack workspace, copy `.env-example` to `.env` and fill in the
webhook URL in `SLACK_WEBHOOK_URL` variable.

In addition to the application, the docker-compose environment consists of:

* **Kafka** and **ZooKeeper**
* **RabbitMQ**
* **Grafana**
* **Prometheus**

The whole environment can be run with:

```bash
docker-compose up
```

You can now configure your GitHub repository to send webhooks to the application (you need to expose port `8080` to the
external network first).

Alternatively, you can run `./scripts/send-stub-webhook.sh` to send some stub webhooks.

Visit [localhost:3000/d/webhooks](http://localhost:3000/d/webhooks) to see annotations added in Grafana. Use
`admin:secret` as credentials.

## Metrics

You can access the Watermill dashboard at [localhost:3000/d/watermill](http://localhost:3000/d/watermill). See what
changes when you send more webhooks over time.

![](https://threedots.tech/media/event-driven-applications/metrics.png)

## What's next?

See [Watermill's documentation](https://watermill.io/) to learn more.
