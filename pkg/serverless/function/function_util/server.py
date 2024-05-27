from flask import Flask, request
import subprocess

app = Flask(__name__)

@app.route('/run', methods=['GET', 'POST'])
def run_script():
    if request.method == 'GET':
        params = request.args
    elif request.method == 'POST':
        params = request.json

    print(f"Received parameters: {params}")

    result = subprocess.run(['python', 'main.py'], capture_output=True, text=True)

    return {
        'output': result.stdout,
        'error': result.stderr,
        'returncode': result.returncode
    }

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=8080)
