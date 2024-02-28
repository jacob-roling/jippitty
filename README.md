# JippiTTY

JippiTTY is a Teletypewriter (TTY) for ChatGPT hence the name, a combination of 'Jippi' for 'GP' in 'ChatGPT' and 'tty', the abbreviation of Teletypewriter.

## Installation

```sh
go install github.com/jacob-roling/JippiTTY
```

## Usage

```sh
jippitty [global options] command [command options]
```

## Manual

```sh
NAME:
   jippitty - A Teletypewriter (TTY) for ChatGPT

USAGE:
   jippitty [global options] command [command options]

DESCRIPTION:
   JippiTTY is a Teletypewriter (TTY) for ChatGPT hence the name, a combination of 'Jippi' for 'GP' in 'ChatGPT' and 'tty', the abbreviation of Teletypewriter.

COMMANDS:
   list     List all conversations
   new      Start a new conversation
   join     Join a conversation
   delete   Delete a/all conversation(s)
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --key value, -k value    OpenAI API Key [$OPENAI_API_KEY]
   --model value, -m value  OpenAI model (default: "gpt-3.5-turbo") [$OPENAI_MODEL]
   --help, -h               show help
```
