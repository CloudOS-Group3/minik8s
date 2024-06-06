import requests
def main(x, msg):
    request_data = {
        'uuid': msg
    }
    response = requests.post('http://192.168.3.8:19293/canceltrain', json=request_data)
    msg = "Fail to order. Flight ticket has been canceled."

    return x, msg