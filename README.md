# Vault Mapper Sync Server

# Prerequisites

- Deno

# Install

1. Clone the repository
2. Run `deno install` to install the dependencies
3. Run `deno task migrate` to create the database and generate the Prisma client

# Run

## Development

- `deno task dev`

## Production

- `deno task start`

# Development

## Prisma Migrate

- `deno task migrate` to create a new migration and generate the Prisma client

TODO: make a docker image or something
