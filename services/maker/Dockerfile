FROM python:stretch

WORKDIR /src
COPY . .

RUN pip install -r services/maker/requirements.txt

CMD ["python", "services/maker/maker/Server.py"] 