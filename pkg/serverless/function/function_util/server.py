import json
import requests
from flask import Flask, request
import my_function as func
app = Flask(__name__)

# curl -X POST -H "Content-Type: application/json" -d '{"uuid":"1234", "params": {"x": 8, "y": 9}}' http://localhost:8080/run
# curl -X POST -H "Content-Type: application/json" -d '{"params": "{\"x\": 8, \"y\": 9}"}' http://localhost:6443/api/v1/namespaces/default/functions/matrix-calculate/run
# {"uuid": "123",
#  "params":
#      {
#          "param1": 8
#      }
#  }
@app.route('/run', methods=['POST'])
def run_script():
    data = request.json
    uuid = data.get('uuid')
    params = data.get('params', {})

    print(f"Received parameters: {params}")

    try:
        result = func.main(**params)
        response_data = {
            'uuid': uuid,
            'result': json.dumps(result),
            'error': '',
        }
    except Exception as e:
        response_data = {
            'uuid': uuid,
            'result': '',
            'error': str(e)
        }
    print(response_data)

    try:
        response = requests.post('http://192.168.3.8:6443/result', json=response_data)
        response.raise_for_status()  # Raise an exception for HTTP errors
    except requests.exceptions.RequestException as e:
        print(f"Error sending data to http://192.168.3.8:6443/result: {e}")
        response_data['error'] = str(e)

    return response_data


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=8080)
