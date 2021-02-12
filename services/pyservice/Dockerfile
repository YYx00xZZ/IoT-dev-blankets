FROM python:3.8.0-slim-buster

RUN apt-get update -y
RUN apt-get install -y gcc git build-essential libtool automake
RUN pip3 install Cython
RUN pip3 install git+https://github.com/MagicStack/uvloop
RUN pip3 install gmqtt
# RUN pip3 install paho-mqtt
RUN pip3 install pymongo
RUN pip3 install dnspython
RUN pip3 install asyncio
RUN pip3 install aiohttp

# RUN pip3 install python-telegram-bot --upgrade
# ADD ./requirements.txt /
ADD ./mqtt_db.py /
# CMD [ "python", "-m pip install --upgrade pip"]
# CMD [ "python", "./hi_py_telegram.py" ]
CMD [ "python", "./mqtt_db.py" ]