package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

const (
    SETTINGS_FILE_RELATIVE_PATH string = "./persistent_settings.json"
)

type Settings struct {
	TargetChannel *discordgo.Channel `json:"channel"`
	TargetRole    *discordgo.Role    `json:"role"`
}

func (s *Settings) initFile() (*os.File, error) {
    var file *os.File

    _, err := os.Stat(SETTINGS_FILE_RELATIVE_PATH)
    if err == nil {
        file, err = os.OpenFile(SETTINGS_FILE_RELATIVE_PATH, os.O_RDWR, os.ModePerm)
        if err != nil {
            log.Errorf("Error while opening file: `%s`", err.Error())
            return nil, err
        }
    } else if errors.Is(err, os.ErrNotExist) {
        file, err = os.Create(SETTINGS_FILE_RELATIVE_PATH)
        if err != nil {
            log.Errorf("Error while creating file: `%s`", err.Error())
            return nil, err
        }
    } else {
        log.Errorf("Error while analyzing file: `%s`", err.Error())
        return nil, err
    }

    return file, nil
}

func (s *Settings) SetChannel(channel *discordgo.Channel) error {
    file, err := s.initFile()
    if err != nil {
        return err
    }

    s.TargetChannel = channel

    return s.write(file)
}

func (s *Settings) SetRole(role *discordgo.Role) error {
    file, err := s.initFile()
    if err != nil {
        return err
    }

    s.TargetRole = role

    return s.write(file)
}

func (s *Settings) Channel() (*discordgo.Channel, error) {
    file, err := s.initFile()
    if err != nil {
        return nil, err
    }

    if err := s.read(file); err != nil {
        return nil, err
    }

    return s.TargetChannel, nil
}

func (s *Settings) Role() (*discordgo.Role, error) {
    file, err := s.initFile()
    if err != nil {
        return nil, err
    }

    if err := s.read(file); err != nil {
        return nil, err
    }

    return s.TargetRole, nil
}

func (s *Settings) write(file *os.File) error {
    defer file.Close()

    raw, err := json.Marshal(s)
    if err != nil {
        log.Errorf("Error while encoding json: `%s`", err.Error())
        return err
    }

    if err := os.Remove(SETTINGS_FILE_RELATIVE_PATH); err != nil {
        log.Errorf("Error while removing file: `%s`", err.Error())
        return err
    }

    file, err = os.Create(SETTINGS_FILE_RELATIVE_PATH)
    if err != nil {
        log.Errorf("Error while truncating file: `%s`", err.Error())
        return err
    }

    if _, err := file.Write(raw); err != nil {
        log.Errorf("Error while writing json to file: `%s`", err.Error())
        return err
    }

	return nil
}

func (s *Settings) read(file *os.File) error {
    defer file.Close()

    raw, err := io.ReadAll(file)
    if err != nil {
        log.Errorf("Error while reading contents of file: `%s`", err.Error())
        return err
    }

    if err := json.Unmarshal(raw, &s); err != nil {
        fmt.Println(string(raw))
        log.Errorf("Error while decoding json: `%s`", err.Error())
        return err
    }

    return nil
}
