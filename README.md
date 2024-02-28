# JippiTTY

JippiTTY is a Teletypewriter (TTY) for ChatGPT hence the name, a combination of 'Jippi' for 'GP' in 'ChatGPT' and 'tty', the abbreviation of Teletypewriter.

## Installation

```sh
go install github.com/jacob-roling/JippiTTY
```

## Usage

```sh
OPENAI_API_KEY=*Your OpenAI API Key*

jippitty [global options] command [command options]
```

## Examples

```sh
OPENAI_API_KEY=*Your OpenAI API Key*

jippitty new poetry

You say: Tell me your favourite short poem.

ChatGPT: I don't have personal preferences, but here is a short poem by Langston Hughes that many people enjoy:

Hold fast to dreams
For if dreams die
Life is a broken-winged bird
That cannot fly

You say: I love it tell, tell me another.

ChatGPT: I'm glad you enjoyed that one! Here's another short poem for you, by Emily Dickinson:

Hope is the thing with feathers
That perches in the soul,
And sings the tune without the words,
And never stops at all.

You say: ^C

Thank your for using JippiTTY, goodbye.
```

## Conversation Storage

Conversations are stored in a folder named .jippitty in your system's home folder.
For mac and linux users that is $HOME/.jippitty
For windows users that is $USERPROFILE/.jippitty

Each conversation is stored as a JSON file in the format of an OpenAI message list that you might send in a response body to https://api.openai.com/v1/chat/completions

For example:

```JSON
// FILE: $HOME/.jippitty/world series.json
{
    "model": "gpt-3.5-turbo",
    "messages": [
      {
        "role": "system",
        "content": "You are a helpful assistant."
      },
      {
        "role": "user",
        "content": "Who won the world series in 2020?"
      },
      {
        "role": "assistant",
        "content": "The Los Angeles Dodgers won the World Series in 2020."
      },
      {
        "role": "user",
        "content": "Where was it played?"
      }
    ]
  }
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
