import requests
def main(x, msg):
    request_data = {
        'uuid': msg
    }
    response = requests.post('http://192.168.3.8:19293/cancelflight', json=request_data)
    msg = "Fail to order. Train ticket has been canceled."

    return x+800, msg