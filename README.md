# Darkan 🐶

Darkan is a Go application that sniffs into libreddit (for now) for a given keyword and returns the HTML content if the keyword is present.

## Requirements

- Go installed on your machine
- Tor installed in your machine

## Tor Installation

Before using Darkan, you need to install Tor on your machine and add it to your $PATH. You can install Tor using Homebrew by running the following command:

```bash
$ brew install tor
$ export PATH="/opt/homebrew/bin/tor:$PATH"
```

**Update:** Now you don't need to open it or establish a connection manually.

**TODO:** Running a remote Tor instance.

## Usage

1. Clone the repository:

   ```bash
   $ git clone https://github.com/wawandco/Darkan.git
2. Navigate into the Darkan directory:

   ```bash
   $ cd Darkan
3. Run the application with the following command:
    ```bash
    $ go run main.go --keyword="your_keyword"
Built as a starting point for exploration and searching of relevant terms on the Dark Web.