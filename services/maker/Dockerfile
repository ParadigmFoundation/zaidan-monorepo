FROM python:stretch

WORKDIR /src
COPY . .

ENV ASSET_CONFIG_FILE=services/maker/maker/asset_config.json

RUN pip install -r services/maker/requirements.txt

CMD ["python", "services/maker/maker/server.py"] 