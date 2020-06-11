# Kumparan

This project is my submission for Senior Backend Technical Assessment for Kumparan written using.

## Prequisite
Please make sure you have these installed and running
1. Go v1.12 or above
2. RabbitMQ v3.8.3 or above
3. ElasticSearch v7.7.1 or above
4. PostgreSQL

## Environment Variable

This app is using some predefined environment variables that you need to set appropriately before running it.

```bash
export PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=
export DB_PASSWORD=
export DB_NAME=
export MQ_URL=amqp://guest:guest@localhost:5672
export ELASTICSEARCH_URL=http://localhost:9200
```

## Database

Create a Postgres database and import this

```sql
CREATE SEQUENCE id_seq
	START WITH 1
	INCREMENT BY 1
	NO MINVALUE
	NO MAXVALUE
	CACHE 1;

CREATE TABLE news (
	id smallint DEFAULT nextval('id_seq'::regclass) NOT NULL,
	author text DEFAULT '' ::text NOT NULL,
	body text DEFAULT '' ::text NOT NULL,
	created timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE ONLY news
	ADD CONSTRAINT news_pkey PRIMARY KEY (id);
```

## Build and Run

```sh
make start
```

## API Documentation

```http
GET /news?page=1
```

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `page` | `integer` | **Required**|

```http
POST /news
```

| Content-Type | Body |
| :--- | :--- |
| `application/json` | `{ "author": "name","body": "news"}` |
