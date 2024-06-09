FROM python:3.11-slim-bookworm

ENV PIP_DEFAULT_TIMEOUT=100 \
    # Allow statements and log messages to immediately appear
    PYTHONUNBUFFERED=1 \
    # disable a pip version check to reduce run-time & log-spam
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    # cache is useless in docker image, so disable to reduce image size
    PIP_NO_CACHE_DIR=1

# Declare build arguments
ARG job_name
ARG partition
ARG N
ARG ntasks_per_node
ARG cpus_per_task
ARG gres

# Set environment variables from build arguments
ENV JOB_NAME=${job_name}
ENV PARTITION=${partition}
ENV N=${N}
ENV NTASKS_PER_NODE=${ntasks_per_node}
ENV CPUS_PER_TASK=${cpus_per_task}
ENV GRES=${gres}

WORKDIR /src

COPY gpu_server.py /src

RUN cd /src
RUN pip install requests
RUN pip install time
RUN pip install paramiko

CMD ["python", "gpu_server.py"]

