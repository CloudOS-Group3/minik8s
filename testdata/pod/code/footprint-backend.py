from flask import Flask, request

app = Flask(__name__)

dict = {}
@app.route('/store', methods=['POST'])
def store():
    data = request.json
    if data['source'] not in dict:
        dict[data['source']] = [data['time']]
    else:
        dict[data['source']].append(data['time'])

    return "Data stored successfully"


@app.route('/get', methods=['POST'])
def get():
    data = request.json
    result_str = "Pod"+str(data['pod']) + " has footprint: \n" + str(dict)
    return result_str



if __name__ == "__main__":
    app.run(host='0.0.0.0', port=3000)
