package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
)

type Completion struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Conversation struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func ensureHome() error {
	_, err := os.Stat(os.Getenv("HOME") + "/.jippitty")
	if os.IsNotExist(err) {
		if err := os.MkdirAll(os.Getenv("HOME")+"/.jippitty", os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func converse(file *os.File, conversation Conversation, api_key string) {
	scanner := bufio.NewScanner(os.Stdin)
	client := http.Client{}

	for {
		fmt.Print("You say: ")

		var content string

		for {
			scanner.Scan()
			line := scanner.Text()
			if len(line) == 0 {
				break
			}
			content += line + "\n"
		}

		conversation.Messages = append(conversation.Messages, Message{
			Role:    "user",
			Content: content,
		})

		conversation_data, err := json.Marshal(conversation)

		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", OPENAI_API_URL+"/chat/completions", bytes.NewBuffer(conversation_data))

		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+api_key)
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			log.Fatal(err)
		}

		var completion Completion

		err = json.Unmarshal(body, &completion)

		if err != nil {
			log.Fatal(err)
		}

		if len(completion.Choices) < 1 {
			log.Fatal(errors.New("less than one choice returned"))
		}

		conversation.Messages = append(conversation.Messages, completion.Choices[0].Message)

		fmt.Println("ChatGPT: " + conversation.Messages[len(conversation.Messages)-1].Content)

		data, err := json.Marshal(conversation)

		if err != nil {
			log.Fatal(err)
		}

		err = file.Truncate(0)

		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Seek(0, 0)

		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Write(data)

		if err != nil {
			log.Fatal(err)
		}
	}
}

const OPENAI_API_URL string = "https://api.openai.com/v1"

func main() {
	var file *os.File
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nThank your for using jippitty, goodbye.")
		os.Exit(0)
	}()

	var model string

	err := ensureHome()

	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:        "jippitty",
		Usage:       "A Teletypewriter (TTY) for ChatGPT",
		Description: "JippiTTY is a Teletypewriter (TTY) for ChatGPT hence the name, a combination of 'Jippi' for 'GP' in 'ChatGPT' and 'tty', the abbreviation of Teletypewriter.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "OpenAI API Key",
				EnvVars: []string{"OPENAI_API_KEY"},
			},
			&cli.StringFlag{
				Name:        "model",
				Aliases:     []string{"m"},
				Usage:       "OpenAI model",
				EnvVars:     []string{"OPENAI_MODEL"},
				Value:       "gpt-3.5-turbo",
				Destination: &model,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all conversations",
				Action: func(ctx *cli.Context) error {
					entries, err := os.ReadDir(os.Getenv("HOME") + "/.jippitty/")

					if err != nil {
						log.Fatal(err)
					}

					for _, entry := range entries {
						file_name := entry.Name()
						fmt.Println(file_name[:len(file_name)-len(filepath.Ext(file_name))])
					}

					return nil
				},
			},
			{
				Name:      "new",
				Usage:     "Start a new conversation",
				ArgsUsage: "[name]",
				Action: func(ctx *cli.Context) error {
					OPENAI_API_KEY := ctx.String("key")

					if OPENAI_API_KEY == "" {
						fmt.Println("API key flag not set")
						return nil
					}

					if ctx.Args().Len() < 1 {
						fmt.Println("Please provide a name for the conversation")
						return nil
					}

					file_path := os.Getenv("HOME") + "/.jippitty/" + strings.Join(ctx.Args().Slice(), " ") + ".json"

					var conversation Conversation

					if os.IsNotExist(err) {
						file, err = os.Create(file_path)

						if err != nil {
							log.Fatal(err)
						}

						defer file.Close()

						conversation = Conversation{
							Model: model,
							Messages: []Message{
								{
									Role:    "system",
									Content: "You are a helpful assistant.",
								},
							},
						}

						data, err := json.Marshal(conversation)

						if err != nil {
							log.Fatal(err)
						}

						_, err = file.Write(data)

						if err != nil {
							log.Fatal(err)
						}

						converse(file, conversation, OPENAI_API_KEY)

						return nil
					}

					file, err = os.OpenFile(file_path, os.O_RDWR, fs.ModePerm)

					if err != nil {
						log.Fatal(err)
					}

					defer file.Close()

					content := new(bytes.Buffer)

					_, err = content.ReadFrom(file)

					if err != nil {
						log.Fatal(err)
					}

					err = json.Unmarshal(content.Bytes(), &conversation)

					if err != nil {
						log.Fatal(err)
					}

					converse(file, conversation, OPENAI_API_KEY)

					return nil
				},
			},
			{
				Name:  "join",
				Usage: "Join a conversation",
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() < 1 {
						fmt.Println("Please include the name of the conversation you wish to join")
						return nil
					}

					file_path := os.Getenv("HOME") + "/.jippitty/" + strings.Join(ctx.Args().Slice(), " ") + ".json"

					var conversation Conversation

					OPENAI_API_KEY := ctx.String("key")

					if OPENAI_API_KEY == "" {
						fmt.Println("API key flag not set")
						return nil
					}

					file, err = os.OpenFile(file_path, os.O_RDWR, fs.ModePerm)

					if os.IsNotExist(err) {
						fmt.Println("Conversation not found. Please choose a conversation from the list.")
						return nil
					}

					defer file.Close()

					content := new(bytes.Buffer)

					_, err = content.ReadFrom(file)

					if err != nil {
						log.Fatal(err)
					}

					err = json.Unmarshal(content.Bytes(), &conversation)

					if err != nil {
						log.Fatal(err)
					}

					for _, entry := range conversation.Messages {
						if entry.Role == "system" {
							fmt.Println("System said: " + entry.Content)
						} else if entry.Role == "user" {
							fmt.Println("You said: " + entry.Content)
						} else if entry.Role == "assistant" {
							fmt.Println("ChatGPT said: " + entry.Content)
						}
					}

					converse(file, conversation, OPENAI_API_KEY)

					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a/all conversation(s)",
				Action: func(ctx *cli.Context) error {
					fmt.Println("added task: ", ctx.Args().First())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
