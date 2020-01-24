FROM python:stretch

WORKDIR /src
COPY . .

RUN pip install -r services/hedger/requirements.txt

CMD ["python", "services/hedger/hedger/server.py"] 