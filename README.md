# spoty 🎵

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/mgjules/spoty)
[![Release](https://img.shields.io/github/release/mgjules/spoty.svg?style=for-the-badge)](https://github.com/mgjules/spoty/releases/latest)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](LICENSE)

Spoty provides simple REST API endpoints to query the current playing track on [Spotify](https://spotify.com).

## Contents
  - [Getting started](#getting-started)
  - [API Documentation](#api-documentation)
  - [Configuration](#configuration)
  - [About the project](#about-the-project)
  - [Stability](#stability)

## Getting started

1. Download the corresponding binary for your operating system and architecture from the [releases](https://github.com/mgjules/spoty/releases) page.

2. Create a Spotify application [here](https://developer.spotify.com/dashboard/applications).

    > Take note of the `Client ID` and `Client Secret` for use in the next step.

3. Setup the proper environment variables (Check the [configuration](#configuration) section for references)
    
    Example from [.env.dist](.env.dist):
    ```sh
    SERVICE_NAME=spoty
    PROD=false
    SPOTIFY_CLIENT_ID=111fedce236f4d34a607711a7ac4a606
    SPOTIFY_CLIENT_SECRET=fb4d98819b912dfd90d5803bb668ad24
    HTTP_SERVER_HOST=localhost
    HTTP_SERVER_PORT=13337
    CACHE_MAX_KEYS=64
    CACHE_MAX_COST=1000000
    JAEGER_ENDPOINT=http://localhost:14268/api/traces
    AMQP_URI=amqp://guest:guest@localhost:5672
    ```

4. Edit the `Redirect URIs` setting of your Spotify application to match the environment variables:

    Example:
    ```sh
    http://<HOST>:<PORT>/api/callback
    ```

5. Run the service:

    ```sh
    $ ./spoty serve
    ```

    > You can also run the service in production mode by setting the `PROD` environment variable to `true`.

6. Authenticate against Spotify by heading to `/api/authenticate` route. 
   You should be redirected to the Spotify login page if you not already logged in on Spotify.
   After logging in, you should be redirected back to the service with the following success message:

    ```json
    {
        "message": "welcome, you are now authenticated!"
    }
    ```


7. Head to the `/api` (health-check) route and expect a similar json response with a HTTP status code `200`:

    Example:
    ```json
    {
        "status": "up",
        "details": {
            "spoty": {
                "status": "up",
                "timestamp": "2022-03-31T08:10:31.534317878Z"
            }
        }
    }
    ```

## API Documentation

Consult the swagger ui page at the `/swagger/index.html` route:

Example:
```sh
http://<HOST>:<PORT>/swagger/index.html
```

## Configuration

| ENV                   | Description                               | Required | Default                           |
| --------------------- | ----------------------------------------- | -------- | --------------------------------- |
| SERVICE_NAME          | Name of microservice                      | No       | spoty                             |
| PROD                  | Whether running in `PROD` or `DEBUG` mode | No       | false                             |
| SPOTIFY_CLIENT_ID     | `Client ID` of app created on Spotify     | **Yes**  | <em>empty</em>                    |
| SPOTIFY_CLIENT_SECRET | `Client Secret` of app created on Spotify | **Yes**  | <em>empty</em>                    |
| HTTP_SERVER_HOST      | Host/IP for HTTP server                   | No       | localhost                         |
| HTTP_SERVER_PORT      | Port for HTTP server                      | No       | 13337                             |
| CACHE_MAX_KEYS        | Maximum number of keys for cache          | No       | 64                                |
| CACHE_MAX_COST        | Maximum size of cache (in bytes)          | No       | 1000000                           |
| JAEGER_ENDPOINT       | Jaeger collector endpoint                 | No       | http://localhost:14268/api/traces |
| AMQP_URI              | AMQP 0-9-1 Uniform Resource Identifier    | No       | amqp://guest:guest@localhost:5672 |

## About the project

This project was inspired by [arwinneil/spotify_chroma](https://github.com/arwinneil/spotify_chroma) and was initially coded in a similar regard as the latter: a fun PoC. However, this project will be maintained until it is deemed feature complete and bug free by the author/maintainer(s) 😊

## Stability

This project follows [SemVer](http://semver.org/) strictly and is considered STABLE.

This project follows the [Go Release Policy](https://golang.org/doc/devel/release.html#policy). Each major version of Go is supported until there are two newer major releases.