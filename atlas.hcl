variable "envfile" {
    type    = string
    default = ".env"
}

locals {
    envfile = {
        for line in split("\n", file(var.envfile)): split("=", line)[0] => regex("=(.*)", line)[0]
        if !startswith(line, "#") && length(split("=", line)) > 1
    }
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./models",
    "--dialect", "postgres",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "postgresql://${local.envfile["POSTGRES_USER"]}:${local.envfile["POSTGRES_PASSWORD"]}@${local.envfile["POSTGRES_HOST"]}:${local.envfile["POSTGRES_PORT"]}/${local.envfile["POSTGRES_DB"]}?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
