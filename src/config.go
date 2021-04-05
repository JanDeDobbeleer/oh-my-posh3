package main

import (
	// "encoding/json"

	"bytes"
	json2 "encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/json"
	"github.com/gookit/config/v2/toml"
	"github.com/gookit/config/v2/yaml"
)

// Config holds all the theme for rendering the prompt
type Config struct {
	FinalSpace           bool              `config:"final_space"`
	OSC99                bool              `config:"osc99"`
	ConsoleTitle         bool              `config:"console_title"`
	ConsoleTitleStyle    ConsoleTitleStyle `config:"console_title_style"`
	ConsoleTitleTemplate string            `config:"console_title_template"`
	TerminalBackground   string            `config:"terminal_background"`
	Blocks               []*Block          `config:"blocks"`
}

// BlockType type of block
type BlockType string

// BlockAlignment aligment of a Block
type BlockAlignment string

const (
	// Prompt writes one or more Segments
	Prompt BlockType = "prompt"
	// LineBreak creates a line break in the prompt
	LineBreak BlockType = "newline"
	// RPrompt a right aligned prompt in ZSH and Powershell
	RPrompt BlockType = "rprompt"
	// Left aligns left
	Left BlockAlignment = "left"
	// Right aligns right
	Right BlockAlignment = "right"
	// EnableHyperlink enable hyperlink
	EnableHyperlink Property = "enable_hyperlink"
)

// Block defines a part of the prompt with optional segments
type Block struct {
	Type             BlockType      `config:"type"`
	Alignment        BlockAlignment `config:"alignment"`
	HorizontalOffset int            `config:"horizontal_offset"`
	VerticalOffset   int            `config:"vertical_offset"`
	Segments         []*Segment     `config:"segments"`
}

// GetConfig returns the default configuration including possible user overrides
func GetConfig(env environmentInfo) *Config {
	cfg, err := loadConfig(env)
	if err != nil {
		return getDefaultConfig(err.Error())
	}
	return cfg
}

func loadConfig(env environmentInfo) (*Config, error) {
	var cfg Config
	configFile := *env.getArgs().Config
	if configFile == "" {
		return nil, errors.New("NO CONFIG")
	}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, errors.New("INVALID CONFIG PATH")
	}

	config.AddDriver(yaml.Driver)
	config.AddDriver(json.Driver)
	config.AddDriver(toml.Driver)
	config.WithOptions(func(opt *config.Options) {
		opt.TagName = "config"
	})

	err := config.LoadFiles(configFile)
	if err != nil {
		return nil, errors.New("UNABLE TO OPEN CONFIG")
	}

	err = config.BindStruct("", &cfg)
	if err != nil {
		return nil, errors.New("INVALID CONFIG")
	}

	return &cfg, nil
}

func exportConfig(configFile, format string) string {
	if len(format) == 0 {
		format = config.JSON
	}

	config.AddDriver(yaml.Driver)
	config.AddDriver(json.Driver)
	config.AddDriver(toml.Driver)

	err := config.LoadFiles(configFile)
	if err != nil {
		return fmt.Sprintf("INVALID CONFIG:\n\n%s", err.Error())
	}

	schemaKey := "$schema"
	if format == config.JSON && !config.Exists(schemaKey) {
		data := config.Data()
		data[schemaKey] = "https://raw.githubusercontent.com/JanDeDobbeleer/oh-my-posh/main/themes/schema.json"
		config.SetData(data)
	}

	buf := new(bytes.Buffer)
	_, err = config.DumpTo(buf, format)
	if err != nil {
		return "UNABLE TO DUMP CONFIG"
	}

	switch format {
	case config.JSON:
		var prettyJSON bytes.Buffer
		err := json2.Indent(&prettyJSON, buf.Bytes(), "", "  ")
		if err == nil {
			unescapeUnicodeCharactersInJSON := func(rawJSON []byte) string {
				str, err := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(rawJSON)), `\\u`, `\u`))
				if err != nil {
					return err.Error()
				}
				return str
			}
			return unescapeUnicodeCharactersInJSON(prettyJSON.Bytes())
		}
	case config.Yaml:
		prefix := "# yaml-language-server: $schema=https://raw.githubusercontent.com/JanDeDobbeleer/oh-my-posh/main/themes/schema.json\n\n"
		content := buf.String()
		return prefix + content

	case config.Toml:
		prefix := "#:schema https://raw.githubusercontent.com/JanDeDobbeleer/oh-my-posh/main/themes/schema.json\n\n"
		content := buf.String()
		return prefix + content
	}

	return buf.String()
}

func getDefaultConfig(info string) *Config {
	cfg := &Config{
		FinalSpace:        true,
		ConsoleTitle:      true,
		ConsoleTitleStyle: FolderName,
		Blocks: []*Block{
			{
				Type:      Prompt,
				Alignment: Left,
				Segments: []*Segment{
					{
						Type:            Session,
						Style:           Diamond,
						Background:      "#c386f1",
						Foreground:      "#ffffff",
						LeadingDiamond:  "\uE0B6",
						TrailingDiamond: "\uE0B0",
					},
					{
						Type:            Path,
						Style:           Powerline,
						PowerlineSymbol: "\uE0B0",
						Background:      "#ff479c",
						Foreground:      "#ffffff",
						Properties: map[Property]interface{}{
							Prefix: " \uE5FF ",
							Style:  "folder",
						},
					},
					{
						Type:            Git,
						Style:           Powerline,
						PowerlineSymbol: "\uE0B0",
						Background:      "#fffb38",
						Foreground:      "#193549",
						Properties: map[Property]interface{}{
							DisplayStashCount:   true,
							DisplayUpstreamIcon: true,
						},
					},
					{
						Type:            Battery,
						Style:           Powerline,
						PowerlineSymbol: "\uE0B0",
						Background:      "#f36943",
						Foreground:      "#193549",
						Properties: map[Property]interface{}{
							ColorBackground:  true,
							ChargedColor:     "#4caf50",
							ChargingColor:    "#40c4ff",
							DischargingColor: "#ff5722",
							Postfix:          "\uF295 ",
						},
					},
					{
						Type:            Node,
						Style:           Powerline,
						PowerlineSymbol: "\uE0B0",
						Background:      "#6CA35E",
						Foreground:      "#ffffff",
						Properties: map[Property]interface{}{
							Prefix:         " \uE718",
							DisplayVersion: false,
						},
					},
					{
						Type:            ShellInfo,
						Style:           Powerline,
						PowerlineSymbol: "\uE0B0",
						Background:      "#0077c2",
						Foreground:      "#ffffff",
						Properties: map[Property]interface{}{
							Prefix: " \uFCB5 ",
						},
					},
					{
						Type:            Root,
						Style:           Powerline,
						PowerlineSymbol: "\uE0B0",
						Background:      "#ffff66",
						Foreground:      "#ffffff",
					},
					{
						Type:            Text,
						Style:           Powerline,
						PowerlineSymbol: "\uE0B0",
						Background:      "#ffffff",
						Foreground:      "#111111",
						Properties: map[Property]interface{}{
							TextProperty: info,
						},
					},
					{
						Type:            Exit,
						Style:           Diamond,
						PowerlineSymbol: "\uE0B0",
						Background:      "#2e9599",
						Foreground:      "#ffffff",
						LeadingDiamond:  "",
						TrailingDiamond: "\uE0B4",
						Properties: map[Property]interface{}{
							DisplayExitCode: false,
							AlwaysEnabled:   true,
							ErrorColor:      "#f1184c",
							ColorBackground: true,
							Prefix:          "<transparent>\uE0B0</> \uE23A",
						},
					},
				},
			},
		},
	}
	return cfg
}
