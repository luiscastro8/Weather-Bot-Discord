# Weather Bot Discord

Language: Go

This is a discord bot for getting U.S. weather from the
[National Weather Service API.](https://www.weather.gov/documentation/services-web-api)

## Instructions

This bot is designed to be run in Docker. It can get the Discord API key either from an environment variable or from
AWS Secret Manager if running on AWS.

### Running in a general Docker container

- DISCORD_TOKEN (Required) - The Discord bot API key

### Running on AWS with Secret Manager

- CLOUD (Required) - Must be set to "aws"
- AWS_REGION (Required) - Region where secret is located
- AWS_SECRET_ID (Required) - Secret ID of API key