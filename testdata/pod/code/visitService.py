from flask import Flask, request, jsonify
import requests

app = Flask(__name__)

@app.route('/get', methods=['GET'])
def get_ip():
    ip = request.args.get('ip')
    if not ip:
        return jsonify({"error": "IP address is required"}), 400

    try:
        response = requests.get(f'http://{ip}')
        return jsonify({
            "status_code": response.status_code,
            "content": response.text
        })
    except requests.RequestException as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0',port=8080)
