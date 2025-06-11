FROM python

RUN apt-get update && apt-get install -y socat netcat-openbsd && \
    ln -sf /bin/nc.openbsd /usr/bin/nc

WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .

ENTRYPOINT ["python", "playground.py"]
