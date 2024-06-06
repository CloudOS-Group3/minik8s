import requests
def main(x, msg):
    request_data = {
        'balance': x
    }
    response = requests.post('http://192.168.3.8:19293/reservehotel', json=request_data)

    response_data = response.json()
    x = response_data['change']
    status = response_data['status']
    uuid = response_data['uuid']
    if status == 'Succeeded':
        msg = uuid
    else:
        msg = f"Failed to buy flight ticket"

    return x, status, msg