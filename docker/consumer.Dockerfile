FROM python:3.12-slim

WORKDIR /app

COPY requirements-kafka-consumer.txt ./
RUN pip install --no-cache-dir -r requirements-kafka-consumer.txt

COPY consumer.py ./

CMD ["python", "consumer.py"]
