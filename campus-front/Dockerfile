FROM python:3.10-alpine as builder
WORKDIR /usr/src/app
RUN apk add --no-cache --virtual .build-deps gcc musl-dev
COPY requirements.txt .
RUN python -m venv /venv
RUN /venv/bin/pip install --upgrade pip
RUN /venv/bin/pip install --no-cache-dir -r requirements.txt
RUN apk del .build-deps

FROM python:3.10-alpine
COPY --from=builder /venv /venv
ENV PATH="/venv/bin:$PATH"
ENV FLASK_APP=main.py
ENV FLASK_RUN_HOST=0.0.0.0
WORKDIR /usr/src/app
COPY . .
EXPOSE 5000
CMD ["flask", "run"]
