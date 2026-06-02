terraform {
  required_providers {
    cohere = {
      source  = "stark256-spec/cohere"
      version = "~> 1.0"
    }
  }
}

provider "cohere" {
  api_key = var.cohere_api_key   # or COHERE_API_KEY env var
}

resource "cohere_team" "eng" {
  name = "engineering"
}

resource "cohere_api_key" "prod" {
  name = "production"
}

output "prod_key" {
  value     = cohere_api_key.prod.secret_key
  sensitive = true
}