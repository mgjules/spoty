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
    url = "https://github.com/JulesMike/spoty.git"
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
