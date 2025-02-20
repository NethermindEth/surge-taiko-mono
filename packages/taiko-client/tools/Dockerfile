FROM python:3.9-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY tx_spammer.py .

ENV TX_FREQUENCY_SECONDS=900
ENV TX_AMOUNT_WEI=1
ENV PYTHONUNBUFFERED=1

CMD ["python", "-u", "tx_spammer.py"]