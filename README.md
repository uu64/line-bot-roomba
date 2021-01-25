# line-bot-roomba

Simple LINE bot to communicate with Roomba on Google App Engine.

## Requirements

- [Docker Compose](https://docs.docker.com/compose/)

- [Google Cloud SDK](https://cloud.google.com/sdk)

## How to use

1. Create `app.yaml`.

```
$ cp sample.app.yaml app.yaml
```

2. Set the following variables in `app.yaml`.

- GAE_PROJECT_ID: Project ID of the GAE where the LINE bot will be installed.

- LINE_BOT_TOKEN: Access token of the LINE bot.

- LINE_BOT_SECRET: Secret key of the LINE bot.

- LINE_BOT_PRIVATE_ID: ID of the target of the notification. (e.g., LINE room ID)

- IFTTT_KEY: Key of webhooks to request Roomba to clean from IFTTT.

3. Start up dev server.

```
$ docker-compose up
```

4. Deploy to GAE.

```
$ gcloud app deploy
```
