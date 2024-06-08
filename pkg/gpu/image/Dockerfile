FROM python:3.11-slim-bookworm

ENV PIP_DEFAULT_TIMEOUT=100 \
    # Allow statements and log messages to immediately appear
    PYTHONUNBUFFERED=1 \
    # disable a pip version check to reduce run-time & log-spam
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    # cache is useless in docker image, so disable to reduce image size
    PIP_NO_CACHE_DIR=1

WORKDIR /src
ARG SOURCE_DIR
# print the value of SOURCE_DIR
RUN echo $SOURCE_DIR
COPY ${SOURCE_DIR} /src/

# Find the .py file and rename it to main.py
RUN cd /src && \
    py_file=$(ls *.py) && \
    mv "$py_file" my_function.py

COPY gpu_server.py /src

RUN cd /src
RUN pip install flask
RUN pip install requests

CMD ["python", "gpu_server.py"]

