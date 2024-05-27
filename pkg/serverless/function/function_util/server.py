from flask import Flask, request
import my_function as func
app = Flask(__name__)


# curl -X POST -H "Content-Type: application/json" -d '{"uuid":"1234", "params": {"x": 8, "y": 9}}}' http://localhost:8080/run
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
    except Exception as e:
        return {
            'uuid': uuid,
            'result': '',
            'error': str(e)
        }
    print(result)
    return {
        'uuid': uuid,
        'result': result,
        'error': '',
    }


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=8888)
