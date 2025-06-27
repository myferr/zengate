package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type PasswordEntry struct {
	ID       int       `json:"id"`
	Site     string    `json:"site"`
	Username string    `json:"username"`
	Password string    `json:"password"` // encrypted
	Created  time.Time `json:"created"`
}

type PasswordStore struct {
	sync.Mutex
	Passwords []PasswordEntry `json:"passwords"`
}

func passwordsFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".zengate-passwords.json"), nil
}

func LoadPasswords() (*PasswordStore, error) {
	path, err := passwordsFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &PasswordStore{Passwords: []PasswordEntry{}}, nil
		}
		return nil, err
	}
	var store PasswordStore
	err = json.Unmarshal(data, &store)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func SavePasswords(store *PasswordStore) error {
	path, err := passwordsFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func AddPassword(site, username, password string) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	if cfg.EncryptKey == "" {
		return errors.New("no encrypt_key configured. Run `zengate set encrypt_key <base64>` first")
	}

	store, err := LoadPasswords()
	if err != nil {
		return err
	}
	store.Lock()
	defer store.Unlock()

	nextID := 1
	for _, p := range store.Passwords {
		if p.ID >= nextID {
			nextID = p.ID + 1
		}
	}

	encryptedPassword, err := EncryptAESGCM(password, cfg.EncryptKey)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	entry := PasswordEntry{
		ID:       nextID,
		Site:     site,
		Username: username,
		Password: encryptedPassword,
		Created:  time.Now(),
	}
	store.Passwords = append(store.Passwords, entry)
	return SavePasswords(store)
}

func ListPasswords() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	if cfg.EncryptKey == "" {
		return errors.New("no encrypt_key configured. Run `zengate set encrypt_key <base64>` first")
	}

	store, err := LoadPasswords()
	if err != nil {
		return err
	}

	fmt.Printf("ID\tSite\t\tUsername\tPassword\t\tCreated\n")
	fmt.Println("----------------------------------------------------------------------------")
	for _, p := range store.Passwords {
		decryptedPass, err := DecryptAESGCM(p.Password, cfg.EncryptKey)
		if err != nil {
			decryptedPass = "(decryption failed)"
		}
		fmt.Printf("%d\t%s\t%s\t%s\t%s\n", p.ID, p.Site, p.Username, decryptedPass, p.Created.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func RemovePassword(id int) error {
	store, err := LoadPasswords()
	if err != nil {
		return err
	}
	store.Lock()
	defer store.Unlock()

	for i, p := range store.Passwords {
		if p.ID == id {
			store.Passwords = append(store.Passwords[:i], store.Passwords[i+1:]...)
			return SavePasswords(store)
		}
	}
	return errors.New("password id not found")
}
