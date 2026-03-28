"""
Минимальный Kafka consumer (Python).
"""

from __future__ import annotations

import json
import logging
import os
import sys

try:
    from kafka import KafkaConsumer
except ImportError:
    print("Install: pip install kafka-python-ng", file=sys.stderr)
    raise

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
)
log = logging.getLogger("consumer")


def getenv(name: str, default: str) -> str:
    value = os.getenv(name)
    return value.strip() if value else default


def main() -> None:
    brokers = getenv("KAFKA_BROKERS", "127.0.0.1:9092")
    topic_user = getenv("KAFKA_TOPIC_USER", "becomeoverman.user.events")
    topic_quest = getenv("KAFKA_TOPIC_QUEST", "becomeoverman.quest.events")

    servers = [b.strip() for b in brokers.split(",") if b.strip()]
    topics = [topic_user, topic_quest]

    log.info("brokers=%s topics=%s mode=no-group", servers, topics)

    consumer = KafkaConsumer(
        *topics,
        bootstrap_servers=servers,
        group_id=None,
        auto_offset_reset="earliest",
        enable_auto_commit=False,
        consumer_timeout_ms=1000,
        value_deserializer=lambda b: b.decode("utf-8", errors="replace"),
    )

    log.info("connected, waiting for messages...")
    try:
        while True:
            records = consumer.poll(timeout_ms=1000)
            for _, messages in records.items():
                for msg in messages:
                    try:
                        data = json.loads(msg.value)
                    except json.JSONDecodeError:
                        log.warning(
                            "non-json value topic=%s partition=%s offset=%s raw=%r",
                            msg.topic,
                            msg.partition,
                            msg.offset,
                            msg.value[:200],
                        )
                        continue

                    log.info(
                        "topic=%s partition=%s key=%s payload=%s",
                        msg.topic,
                        msg.partition,
                        (msg.key or b"").decode("utf-8", errors="replace"),
                        json.dumps(data, ensure_ascii=False),
                    )
    finally:
        consumer.close()


if __name__ == "__main__":
    main()
