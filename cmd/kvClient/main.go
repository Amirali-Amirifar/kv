package main

import (
	"bufio"
	"fmt"
	"github.com/Amirali-Amirifar/kv/pkg/kvClient"
	"os"
	"regexp"
	"strings"
)

// parseCommand parses command line input
func parseCommand(input string) (string, []string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil, fmt.Errorf("empty command")
	}

	// Regular expression to match quoted strings and unquoted words
	re := regexp.MustCompile(`"([^"\\]*(\\.[^"\\]*)*)"|(\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	var parts []string
	for _, match := range matches {
		if match[1] != "" {
			// Quoted string - unescape quotes
			unescaped := strings.ReplaceAll(match[1], `\"`, `"`)
			parts = append(parts, unescaped)
		} else if match[3] != "" {
			// Unquoted word
			parts = append(parts, match[3])
		}
	}

	if len(parts) == 0 {
		return "", nil, fmt.Errorf("no command found")
	}

	return strings.ToUpper(parts[0]), parts[1:], nil
}

// executeCommand executes the parsed command
func executeCommand(client *kvClient.Client, cmd string, args []string) error {
	switch cmd {
	case "SET":
		if len(args) < 2 {
			return fmt.Errorf("SET requires key and value: SET \"key\" \"value\"")
		}
		if len(args) > 2 {
			return fmt.Errorf("SET takes exactly 2 arguments: SET \"key\" \"value\"")
		}

		key, value := args[0], args[1]
		if err := client.Set(key, value); err != nil {
			return fmt.Errorf("SET failed: %v", err)
		}
		fmt.Println("OK")
		return nil

	case "GET":
		if len(args) != 1 {
			return fmt.Errorf("GET requires exactly one key: GET \"key\"")
		}

		key := args[0]
		value, err := client.Get(key)
		if err != nil {
			return fmt.Errorf("GET failed: %v", err)
		}
		fmt.Printf("\"%s\"\n", value)
		return nil

	case "DEL":
		if len(args) != 1 {
			return fmt.Errorf("DEL requires exactly one key: DEL \"key\"")
		}

		key := args[0]
		if err := client.Del(key); err != nil {
			return fmt.Errorf("DEL failed: %v", err)
		}
		fmt.Println("OK")
		return nil

	case "QUIT", "EXIT":
		fmt.Println("Goodbye!")
		os.Exit(0)
		return nil

	case "HELP":
		printHelp()
		return nil

	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

// printHelp displays available commands
func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  SET \"key\" \"value\"  - Set a key-value pair")
	fmt.Println("  GET \"key\"           - Get value for a key")
	fmt.Println("  DEL \"key\"           - Delete a key")
	fmt.Println("  HELP                 - Show this help message")
	fmt.Println("  QUIT/EXIT            - Exit the client")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  SET \"mykey\" \"myvalue\"")
	fmt.Println("  GET \"mykey\"")
	fmt.Println("  DEL \"mykey\"")
}

func main() {
	// Get API base URL from command line arguments or environment
	baseURL := "http://localhost:8080"
	if len(os.Args) > 1 {
		baseURL = os.Args[1]
	} else if envURL := os.Getenv("KV_API_URL"); envURL != "" {
		baseURL = envURL
	}

	// Remove trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	client := kvClient.NewClient(baseURL)

	fmt.Printf("KV Database Client\n")
	val, err := client.Connect()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(val)

	fmt.Println("Type HELP for available commands or QUIT to exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("kv> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		cmd, args, err := parseCommand(input)
		if err != nil {
			if err.Error() != "empty command" {
				fmt.Printf("Error: %v\n", err)
			}
			continue
		}

		if err := executeCommand(client, cmd, args); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}
}
