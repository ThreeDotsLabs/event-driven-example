# Event-Driven application example

This is an example project of building Event-Driven services in Go, using [Watermill](https://github.com/ThreeDotsLabs/watermill).

In addition to the application, the environment consists of:

* **Kafka** and **ZooKeeper**
* **RabbitMQ**
* **Grafana**
* **Prometheus**

## Running

If you'd like to integrate the example with your Slack workspace, copy `.env-example` to `.env` and fill in the
webhook URL in `SLACK_WEBHOOK_URL` variable.

The whole environment can be run with:

```bash
docker-compose up
```

You can now configure your GitHub repository to send webhooks to the application (you need to expose the port to the
external network first).

Alternatively, you can run `./scripts/send-stub-webhook.sh` to send some stub webhooks.

Visit [localhost:3000/d/webhooks](http://localhost:3000/d/webhooks) to see annotations added in Grafana. Use
`admin:secret` as credentials.

## Metrics

You can access the Watermill dashboard at [localhost:3000/d/watermill](http://localhost:3000/d/watermill). See what
changes when you send more webhooks over time.

## What's next?

See [Watermill's documentation](https://watermill.io/) to learn more.
