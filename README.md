# spoty 🎵

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/JulesMike/spoty)
[![Release](https://img.shields.io/github/release/JulesMike/spoty.svg?style=for-the-badge)](https://github.com/JulesMike/spoty/releases/latest)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](LICENSE)

Spoty provides simple REST API endpoints to query the current playing track on [Spotify](https://spotify.com).

## Contents
  - [Getting started](#getting-started)
  - [Configuration](#configuration)
  - [API Documentation](#api-documentation)
  - [About the project](#about-the-project)
  - [Stability](#stability)

## Getting started

1. Download the corresponding binary for your operating system and architecture from the [releases](https://github.com/JulesMike/spoty/releases) page.

2. Create a Spotify application [here](https://developer.spotify.com/dashboard/applications).

    > Take note of the `Client ID` and `Client Secret` for use in the next step.

3. Setup the proper environment variables (Check the [configuration](#configuration) section for references)
    
    Example from [.env.dist](.env.dist):
    ```sh
    PROD=false
    CLIENT_ID=111fedce236f4d34a607711a7ac4a606
    CLIENT_SECRET=fb4d98819b912dfd90d5803bb668ad24
    HOST=localhost
    PORT=13337
    CACHE_MAX_KEYS=64
    CACHE_MAX_COST=1000000
    ```

4. Edit the `Redirect URIs` setting of your application in spotify to match the environment variables:

    Example:
    ```sh
    http://<HOST>:<PORT>/api/callback
    ```

5. Execute the binary and head to the `/api` route and expect this json response with a HTTP status code `200`:

    ```json
    {
        "success": "i'm alright!"
    }
    ```

## Usage

To access most of routes, you'll need to authenticate yourself against spotify by going to the `/api/authenticate` route:

Example:
```sh
http://<HOST>:<PORT>/api/authenticate
```

1. You should be redirected to spotify for authentication. 
2. After which you will be redirected back to the url specified in `REDIRECT_URL`.

## Configuration

| ENV               | Description                               | Required  | Default                           |
| --                | --                                        | --        | --                                |
| CLIENT_ID         | `Client ID` of app created on Spotify     | **Yes**   | <em>empty</em>                    |
| CLIENT_SECRET     | `Client Secret` of app created on Spotify | **Yes**   | <em>empty</em>                    |
| PROD              | Whether running in `PROD` or `DEBUG` mode | No        | false                             |
| HOST              | Host/IP for HTTP server                   | No        | localhost                         |
| PORT              | Port for HTTP server                      | No        | 13337                             |
| CACHE_MAX_KEYS    | Maximum number of keys for cache          | No        | 64                                |
| CACHE_MAX_COST    | Maximum size of cache (in bytes)          | No        | 1000000                           |

## API Documentation

Consult the swagger ui page at the `/swagger/index.html` route:

Example:
```sh
http://<HOST>:<PORT>/swagger/index.html
```

## About the project

This project was inspired by [arwinneil/spotify_chroma](https://github.com/arwinneil/spotify_chroma) and was intially coded in a similar regard as the latter: a fun PoC. However, this project will be maintained until it is deemed feature complete and bug free by the author/maintainer(s) 😊

## Stability

This project follows [SemVer](http://semver.org/) strictly and is not yet `v1`.

Breaking changes might be introduced until `v1` is released.

This project follows the [Go Release Policy](https://golang.org/doc/devel/release.html#policy). Each major version of Go is supported until there are two newer major releases.