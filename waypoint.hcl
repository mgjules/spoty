project = "spoty"


variable "docker_username" {
  type = string
}

variable "docker_password" {
  type = string
}

runner {
  enabled = true
  profile = "kubernetes"

  data_source "git" {
    url = "https://github.com/mgjules/spoty.git"
    ref = "main"
  }

  poll {
    enabled = true
  }
}

app "api" {
  labels = {
    "lang" = "go"
  }

  config {
    env = {
      SPOTIFY_CLIENT_ID = dynamic("kubernetes", {
        name   = "spoty-creds"
        key    = "spotifyClientID"
        secret = true
      })

      SPOTIFY_CLIENT_SECRET = dynamic("kubernetes", {
        name   = "spoty-creds"
        key    = "spotifyClientSecret"
        secret = true
      })

      SPOTIFY_REDIRECT_URI = dynamic("kubernetes", {
        name = "spoty-config"
        key  = "spotifyRedirectUri"
      })

      SERVICE_NAME = dynamic("kubernetes", {
        name = "spoty-config"
        key  = "serviceName"
      })

      PROD = dynamic("kubernetes", {
        name = "spoty-config"
        key  = "prod"
      })

      HTTP_SERVER_HOST = dynamic("kubernetes", {
        name = "spoty-config"
        key  = "httpServerHost"
      })

      HTTP_SERVER_PORT = dynamic("kubernetes", {
        name = "spoty-config"
        key  = "httpServerPort"
      })

      CACHE_MAX_KEYS = dynamic("kubernetes", {
        name = "spoty-config"
        key  = "cacheMaxKeys"
      })

      CACHE_MAX_COST = dynamic("kubernetes", {
        name = "spoty-config"
        key  = "cacheMaxCost"
      })

      JAEGER_ENDPOINT = dynamic("kubernetes", {
        name = "spoty-config"
        key  = "jaegerEndpoint"
      })

      AMQP_URI = dynamic("kubernetes", {
        name = "spoty-config"
        key  = "amqpUri"
      })
    }
  }

  build {
    use "docker" {
      buildkit = false
    }

    registry {
      use "docker" {
        image    = "julesmike/spoty"
        tag      = gitrefpretty()
        username = var.docker_username
        password = var.docker_password
      }
    }
  }

  deploy {
    use "kubernetes" {
      service_port = 13337
    }
  }

  release {
    use "kubernetes" {
      port = 8080
    }
  }
}
