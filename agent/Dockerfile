FROM python

RUN apt-get update && apt-get install -y socat

COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .

CMD ["python", "playground.py"]