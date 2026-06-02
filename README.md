# terraform-provider-cohere

Terraform provider for Cohere platform admin API — teams and API keys.

## Usage

```hcl
terraform {
  required_providers {
    cohere = {
      source  = "stark256-spec/cohere"
      version = "~> 1.0"
    }
  }
}

provider "cohere" {
  api_key = var.cohere_api_key
}

resource "cohere_team" "eng" {
  name = "engineering"
}

resource "cohere_api_key" "prod" {
  name = "production"
}
```

## Authentication

Set your API key via the `api_key` argument or the environment variable shown in the provider schema.

## Resources

| Resource | Description |
|----------|-------------|
| `cohere_workspace` / `cohere_project` / `cohere_team` | Isolated environment |
| `cohere_api_key` | API key scoped to a workspace/project |

## License

Apache 2.0
