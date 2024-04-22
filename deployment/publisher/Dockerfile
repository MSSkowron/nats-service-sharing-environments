# Use an official Python runtime as a base image
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Copy the Python script into the container
COPY main.py .

# Install the required Python packages
RUN pip install nats-py

# Define the command to run the script
CMD ["python", "main.py"]
