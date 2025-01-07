# Prometheus Exporter for RabbitMQ

This is a custom Golang prometheus exporter for RabbitMQ

## What is a Prometheus exporter?

A prometheus exporter gathers metrics from target services and converts them to a format that can be utilized by Prometheus.

There are a variety of exporters that can be found on the **[Prometheus documentation](https://prometheus.io/docs/instrumenting/exporters/)** and libraries that you can use to build custom exporters.

I this repository we make use of the **[Go Client library](github.com/prometheus/client_golang/)**.

## RabbitMQ

RabbitMQ is a reliable and mature messaging and streaming broker, which is easy to deploy on cloud environments, on-premises, and on your local machine.

A tutorial on how to use RabbitMQ can be found on their **[docs](https://www.rabbitmq.com/docs)**.

We are making use of RabbitMQ API that helps you query information on exchanges and queues.

To access the API, please ensure you first run **[RabbitMQ docker container](https://hub.docker.com/_/rabbitmq)** the navigate to `/api` and get access to a list of endpoints you can connect to.
