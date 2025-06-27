package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func printHelp() {
	fmt.Println(`
Zengate CLI

Commands:
  set <key> <value>           Set config value (tunnel_url, encrypt_key)
  compose                     Clone zengate repo and run docker compose
  new password [site user pass]  Create a new password entry (interactive if args missing)
  list                        List saved passwords
  remove <id>                 Remove saved password by ID
  help                        Show this help
`)
}

func RunCommand(args []string) error {
	if len(args) == 0 {
		printHelp()
		return nil
	}

	cmd := strings.ToLower(args[0])
	switch cmd {
	case "help":
		printHelp()
		return nil

	case "set":
		if len(args) != 3 {
			return errors.New("usage: zengate set <key> <value>")
		}
		key := strings.ToLower(args[1])
		value := args[2]

		cfg, err := LoadConfig()
		if err != nil {
			return err
		}

		switch key {
		case "tunnel_url", "tunnelurl":
			cfg.TunnelURL = value
		case "encrypt_key", "encryptkey":
			err := ValidateEncryptKey(value)
			if err != nil {
				return fmt.Errorf("invalid encrypt_key: %w", err)
			}
			cfg.EncryptKey = value
		default:
			return fmt.Errorf("unknown config key %s", key)
		}

		err = SaveConfig(cfg)
		if err != nil {
			return err
		}

		fmt.Println("Config saved.")
		return nil

	case "compose":
		fmt.Println("Cloning and running docker-compose...")
		return RunDockerCompose()

	case "new":
		if len(args) < 2 {
			return errors.New(`did you mean "zengate new password"?`)
		}
		subcmd := strings.ToLower(args[1])
		switch subcmd {
		case "password":
			var site, user, pass string
			if len(args) >= 5 {
				site = args[2]
				user = args[3]
				pass = args[4]
			} else {
				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Site? ")
				s, _ := reader.ReadString('\n')
				site = strings.TrimSpace(s)
				if site == "" {
					site = "github.com"
				}

				fmt.Print("Username? ")
				u, _ := reader.ReadString('\n')
				user = strings.TrimSpace(u)
				if user == "" {
					user = "johndoe"
				}

				fmt.Print("Password? ")
				p, _ := reader.ReadString('\n')
				pass = strings.TrimSpace(p)
				if pass == "" {
					pass = "password123"
				}
			}

			err := AddPassword(site, user, pass)
			if err != nil {
				return err
			}
			fmt.Println("Password created! View with \"zengate list\"")
			return nil
		default:
			return fmt.Errorf(`unknown subcommand for new: %s`, subcmd)
		}

	case "list":
		return ListPasswords()

	case "remove":
		if len(args) != 2 {
			return errors.New("usage: zengate remove <id>")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("id must be a number")
		}
		err = RemovePassword(id)
		if err != nil {
			return err
		}
		fmt.Println("Password removed.")
		return nil

	default:
		printHelp()
		return fmt.Errorf("unknown command: %s", cmd)
	}
}
