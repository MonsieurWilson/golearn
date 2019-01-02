mq-benchmarking
==========================

**Description**
This is a benchmark tool that tests throughput of several types of Message Queues.

**Note**
1. Only Nats-streaming supports async publishing at present.
2. For receive direction, caller should take charge of message sending.

**Usage:**

```bash
Usage: mq-benchmarking <message queue> [options]

Message Queue Options:
        rabbitmq       - RabbitMQ
        nsq            - NSQ
        nats           - NATS
        nats-streaming - NATS Streaming

Benchmarking Options:
        -n,  --num_messages <int> Number of total messages (default: 1000000)
        -s,  --message_size <int> Size of message (default: 1000)
        -a,  --async              Async message publishing (default: false)
        -r,  --random             Send random bytes instead of empty bytes (default: false)
        -u,  --url <string>       Message queue server URL (default: varies from different message queues)
        -t,  --topic <string>     Topic of publish/subscribe (default: test)
        -d,  --direction <string> Direction of test, eg: send|receive (default: sender)
```
