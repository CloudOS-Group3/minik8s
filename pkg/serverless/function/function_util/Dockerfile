FROM python:3.11-slim-bookworm

ENV PIP_DEFAULT_TIMEOUT=100 \
    # Allow statements and log messages to immediately appear
    PYTHONUNBUFFERED=1 \
    # disable a pip version check to reduce run-time & log-spam
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    # cache is useless in docker image, so disable to reduce image size
    PIP_NO_CACHE_DIR=1

#RUN set -ex \
#    # Create a non-root user
#    && addgroup --system --gid 1001 appgroup \
#    && adduser --system --uid 1001 --gid 1001 --no-create-home appuser \
#
#RUN chown -R appuser:appuser /your-subdirectory

WORKDIR /app
ARG SOURCE_DIR
# print the value of SOURCE_DIR
RUN echo $SOURCE_DIR
COPY ${SOURCE_DIR} /app/

# Find the .py file and rename it to main.py
RUN cd /app && \
    py_file=$(ls *.py) && \
    mv "$py_file" main.py

COPY server.py /app

RUN cd /app && ls && \
    pip install --no-cache-dir -r requirement.txt
RUN pip install flask
EXPOSE 8080

CMD ["python", "server.py"]

