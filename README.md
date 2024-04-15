# Darkan 🐶

Darkan is a Go application that sniffs into libreddit (for now) for a given keyword and returns the HTML content if the keyword is present.

## Tor Installation

Before using Darkan, you need to install Tor on your machine and add it to your $PATH. You can install Tor using Homebrew by running the following command:

```bash
$ brew install tor
$ export PATH="/opt/homebrew/bin/tor:$PATH"
```

**Update:** Now you don't need to open it or establish a connection manually, however, you still need a Tor instance, it can be local (recommended) or remote.

## Running in development

### Setup
The first step you need to do is to create a development database and run the migrations by running:

```
$ go run ./cmd/db create
$ go run ./cmd/db migrate
```

### Running 
To run the app in development you can run the following command:

```
$ go run ./cmd/dev
```

Root path for the endpoints will be available at http://localhost:3000.

### Creating a new Keyword to search

<details>
 <summary><code>POST</code> <code><b>/api/search</b></code> <code>(registers a new keyword to search with a callback URL to be notified)</code></summary>

##### Parameters

> | name      |  type     | data type               | description                                                           |
> |-----------|-----------|-------------------------|-----------------------------------------------------------------------|
> | keyword      |  required | string (JSON)   | Value to search in the Dark Web  |
> | callback_url |  required | string (JSON)   | Endpoint URL to be notified once a keyword in found  |


##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `application/json`                | `Keyword registered successfully`                                   |
> | `500`         | `application/json`                | `{"code":"400","message":"Internal server error saving keyword."`   |
> | `[TODO: Validations]` | N/A| N/A                                                                |

##### Example cURL

> ```javascript
>  curl -X POST -H "Content-Type: application/json" --data @keyword.json http://localhost:3000/api/search
> ```

</details>

## Env Variables

```bash
GO_ENV                          - The environment the app is running in. Defaults "development"
PORT                            - The port the app will run on. Defaults "3000"
TOR_PROXY                       - The Tor proxy you want to use. Defaults "socks5://127.0.0.1:9050"
```

In case you want to use a remote TOR instance, you can set a Tor Relay listed [here](https://metrics.torproject.org/rs.html#toprelays).