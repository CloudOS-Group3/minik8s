import datetime

import requests
from flask import Flask, request

app = Flask(__name__)


@app.route('/footprint', methods=['GET'])
def run_script():
    request_data = {
        'pod': 1,
    }
    try:
        response = requests.post('http://127.0.0.1:3000/get', json=request_data)
        response.raise_for_status()  # Raise an exception for HTTP errors
    except requests.exceptions.RequestException as e:
        print(f"Error sending data to http://127.0.0.1:3000/get: {e}")
        return f"Error sending data to http://127.0.0.1:3000/get: {e}"

    return response.text


@app.route('/', methods=['GET'])
def greeting():
    source_ip = request.remote_addr
    # send to backend, localhost:3000
    # send Post request to backend
    request_data = {
        'source': source_ip,
        'time': datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    }
    try:
        response = requests.post('http://127.0.0.1:3000/store', json=request_data)
        response.raise_for_status()  # Raise an exception for HTTP errors
    except requests.exceptions.RequestException as e:
        print(f"Error sending data to http://127.0.0.1:3000/store: {e}")

    return f'Hello From Pod1! Visit /footprint to see your footprint.\n'


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=8888)
