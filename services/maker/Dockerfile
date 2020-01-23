FROM python:stretch

WORKDIR /src
COPY . .

ENV CONFIG_FILE=services/maker/maker/config.json

RUN pip install -r services/maker/requirements.txt

CMD ["python", "services/maker/maker/server.py"] 