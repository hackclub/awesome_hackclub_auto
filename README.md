<!-- DO NOT REMOVE - contributor_list:data:start:["cjdenio", "Matt-Gleich"]:end -->

# awesome_hackclub_auto

> Automation service for [awesome-hackclub](https://github.com/hackclub/awesome-hackclub)

## Running locally

After installing [Docker](https://docker.com) and [Docker Compose](https://docs.docker.com/compose), run `docker-compose up` in your favorite terminal.

### Environment variables

Create a `.env` file, and stick the following env variables in:

```
SLACK_TOKEN=a slack bot token
SLACK_SIGNING_SECRET=signing secret, NOT verification token
REVIEW_CHANNEL=channel ID to post review messages to
AIRTABLE_API_KEY=airtable api key
AIRTABLE_BASE_ID=id of the airtable base
GH_APP_ID=id of github app
GH_INSTALLATION_ID=the installation id
```

### Commands

| Command                       | Description                                              |
| ----------------------------- | -------------------------------------------------------- |
| `docker-compose up`           | Starts the dev environment locally                       |
| `docker-compose up -d`        | Starts the dev environment, then detaches from the shell |
| `docker-compose logs -f main` | View logs from the main process                          |
| `docker-compose down`         | Shut down the dev environment                            |
| `docker-compose restart`      | Restarts the dev environment                             |
| `docker ps`                   | View running services                                    |

<!-- DO NOT REMOVE - contributor_list:start -->

## 👥 Contributors

- **[@cjdenio](https://github.com/cjdenio)**

- **[@Matt-Gleich](https://github.com/Matt-Gleich)**

<!-- DO NOT REMOVE - contributor_list:end -->
