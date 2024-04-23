FROM python:3.9-slim

WORKDIR /app

COPY main.py .

RUN pip install nats-py

CMD ["python", "-u", "main.py"]
