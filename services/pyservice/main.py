import requests
import asyncio
import os
import signal
import time
# import requests
import json
from gmqtt import Client as MQTTClient
from pprint import pprint
# gmqtt also compatibility with uvloop  
import uvloop
asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())


STOP = asyncio.Event()

def on_connect(client, flags, rc, properties):
    client.subscribe('TEST/#', qos=0)
    headers = {'content-type': 'application/json'}
    url = 'http://localhost:5000/flespi'
    params = {'deviceName': client._client_id, 'deviceToken': token}                                        # url parameters
    data = {"eventType": "AAS_PORTAL_START", "data": {"uid": "hfe3hf45huf33545", "aid": "1", "vid": "1"}}   # body data (json encoded)

    requests.post(url, params=params, data=json.dumps(data), headers=headers)


def on_message(client, topic, payload, qos, properties):
    headers = {'content-type': 'application/json'}
    url = 'http://localhost:5000/flespi'
    params = {'deviceName': client._client_id, 'deviceToken': token}                                        # url parameters
    data = {"eventType": "AAS_PORTAL_START", "data": {"uid": "hfe3hf45huf33545", "aid": "1", "vid": "1"}}   # body data (json encoded)

    requests.post(url, params=params, data=json.dumps(data), headers=headers)
    decoded = payload.decode()
    json_payload = json.loads(decoded)
    print('\t[RECV]')
    pprint(json_payload)

def on_disconnect(client, packet, exc=None):
    print('[Disconnected]')


def on_subscribe(client, mid, qos, properties):
    print('[SUBSCRIBED]')


def ask_exit(*args):
    STOP.set()


async def main(broker_host, token):
    client = MQTTClient("vscode-client")

    client.on_connect = on_connect
    client.on_message = on_message
    client.on_disconnect = on_disconnect
    client.on_subscribe = on_subscribe
    
    client.set_auth_credentials(token, None)
    await client.connect(broker_host)

    # client.publish('TEST/TIME', str(time.time()), qos=0)

    await STOP.wait()
    await client.disconnect()


if __name__ == '__main__':
    loop = asyncio.get_event_loop()

    host = 'mqtt.flespi.io'
    token = os.environ.get('FLESPI_TOKEN')

    loop.add_signal_handler(signal.SIGINT, ask_exit)
    loop.add_signal_handler(signal.SIGTERM, ask_exit)

    loop.run_until_complete(main(host, token))