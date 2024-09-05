/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bwmarrin/snowflake"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	node                    *snowflake.Node
	msdConfig               MSDConfig
	allLicenses             map[string]bool
	emailConfig             EmailConfig
	phoneConfig             PhoneConfig
	tokenConfig             TokenConfig
	websiteConfig           WebsiteConfig
	accountConfig           AccountConfig
	randomIdLength          int
	passwordInstance        *passwordImpl
	acceptableAvatarDomains []string
	skipAvatarids           sets.Set[string]
	cdnUrlConfig            string
	allowImageExtension     sets.Set[string]
)

// Init initializes the configuration with the given Config struct.
func Init(cfg *Config) (err error) {
	msdConfig = cfg.MSDConfig
	emailConfig = cfg.EmailConfig
	phoneConfig = cfg.PhoneConfig
	tokenConfig = cfg.TokenConfig
	websiteConfig = cfg.WebsiteConfig
	accountConfig = cfg.AccountConfig
	cdnUrlConfig = cfg.CdnUrlConfig

	m := map[string]bool{}
	for _, v := range cfg.Licenses {
		m[strings.ToLower(v)] = true
	}

	allLicenses = m

	// TODO: node id should be same with replica id
	node, err = snowflake.NewNode(1)

	randomIdLength = cfg.RandomIdLength
	passwordInstance = newPasswordImpl(cfg.PasswordConfig)
	if len(cfg.AcceptableAvatarDomains) <= 0 {
		err = errors.New("no acceptable avatar domains configured")
		return
	}

	acceptableAvatarDomains = cfg.AcceptableAvatarDomains

	skipAvatarids = sets.New[string]()
	for _, v := range cfg.SkipAvatarIds {
		skipAvatarids.Insert(v)
	}

	allowImageExtension = sets.New[string]()
	for _, v := range cfg.AllowImageExtension {
		allowImageExtension.Insert(v)
	}

	return
}

// Config represents the main configuration structure.
type Config struct {
	Licenses                []string       `json:"licenses"                  required:"true"`
	MSDConfig               MSDConfig      `json:"msd"`
	EmailConfig             EmailConfig    `json:"email"`
	PhoneConfig             PhoneConfig    `json:"phone"`
	TokenConfig             TokenConfig    `json:"token"`
	WebsiteConfig           WebsiteConfig  `json:"website"`
	AccountConfig           AccountConfig  `json:"account"`
	RandomIdLength          int            `json:"random_id_length"`
	PasswordConfig          PasswordConfig `json:"password_config"`
	AcceptableAvatarDomains []string       `json:"acceptable_avatar_domains" required:"true"`
	SkipAvatarIds           []string       `json:"skip_avatar_ids"`
	CdnUrlConfig            string         `json:"cdn_url_config"            required:"true"`
	AllowImageExtension     []string       `json:"allow_image_extension"     required:"true"`
}

// SetDefault sets default values for Config if they are not provided.
func (cfg *Config) SetDefault() {
	if cfg.RandomIdLength <= 0 {
		cfg.RandomIdLength = 24
	}
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.MSDConfig,
		&cfg.EmailConfig,
		&cfg.PhoneConfig,
		&cfg.TokenConfig,
		&cfg.WebsiteConfig,
		&cfg.AccountConfig,
		&cfg.PasswordConfig,
	}
}

// MSDConfig represents the configuration for MSD.
type MSDConfig struct {
	NameRegexp        string `json:"msd_name_regexp"          required:"true"`
	MaxNameLength     int    `json:"msd_name_max_length"      required:"true"`
	MinNameLength     int    `json:"msd_name_min_length"      required:"true"`
	MaxDescLength     int    `json:"msd_desc_max_length"      required:"true"`
	MaxFullnameLength int    `json:"msd_fullname_max_length"  required:"true"`

	nameRegexp *regexp.Regexp
}

// Validate check values for MSDConfig whether they are valid.
func (cfg *MSDConfig) Validate() (err error) {
	cfg.nameRegexp, err = regexp.Compile(cfg.NameRegexp)

	return
}

// EmailConfig represents the configuration for Email.
type EmailConfig struct {
	Regexp    string `json:"email_regexp"      required:"true"`
	MaxLength int    `json:"email_max_length"  required:"true"`

	regexp *regexp.Regexp
}

// Validate check values for EmailConfig whether they are valid.
func (cfg *EmailConfig) Validate() (err error) {
	cfg.regexp, err = regexp.Compile(cfg.Regexp)

	return
}

// PhoneConfig represents the configuration for Phone.
type PhoneConfig struct {
	Regexp    string `json:"phone_regexp"      required:"true"`
	MaxLength int    `json:"phone_max_length"  required:"true"`

	regexp *regexp.Regexp
}

// Validate check values for PhoneConfig whether they are valid.
func (cfg *PhoneConfig) Validate() (err error) {
	cfg.regexp, err = regexp.Compile(cfg.Regexp)

	return
}

// TokenConfig represents the configuration for Token.
type TokenConfig struct {
	Regexp        string `json:"token_name_regexp"     required:"true"`
	MaxNameLength int    `json:"token_name_max_length" required:"true"`
	MinNameLength int    `json:"token_name_min_length" required:"true"`

	regexp *regexp.Regexp
}

// Validate check values for TokenConfig whether they are valid.
func (cfg *TokenConfig) Validate() (err error) {
	cfg.regexp, err = regexp.Compile(cfg.Regexp)

	return
}

// WebsiteConfig represents the configuration for Website.
type WebsiteConfig struct {
	Regexp    string `json:"website_regexp"     required:"true"`
	MaxLength int    `json:"website_max_length" required:"true"`

	regexp *regexp.Regexp
}

// Validate check values for WebsiteConfig whether they are valid.
func (cfg *WebsiteConfig) Validate() (err error) {
	cfg.regexp, err = regexp.Compile(cfg.Regexp)

	return
}

// AccountConfig represents the configuration for Account.
type AccountConfig struct {
	NameRegexp        string   `json:"account_name_regexp"         required:"true"`
	MaxNameLength     int      `json:"account_name_max_length"     required:"true"`
	MinNameLength     int      `json:"account_name_min_length"     required:"true"`
	MaxDescLength     int      `json:"account_desc_max_length"     required:"true"`
	ReservedAccounts  []string `json:"reserved_accounts"           required:"true"`
	MinFullnameLength int      `json:"org_fullname_min_length"     required:"true"`
	MaxFullnameLength int      `json:"account_fullname_max_length" required:"true"`

	nameRegexp       *regexp.Regexp
	reservedAccounts sets.Set[string]
}

// Validate check values for AccountConfig whether they are valid.
func (cfg *AccountConfig) Validate() (err error) {
	if cfg.nameRegexp, err = regexp.Compile(cfg.NameRegexp); err != nil {
		return err
	}

	if len(cfg.ReservedAccounts) > 0 {
		cfg.reservedAccounts = sets.New[string]()
		cfg.reservedAccounts.Insert(cfg.ReservedAccounts...)
	}

	return nil
}

// PasswordConfig represents the configuration for password.
type PasswordConfig struct {
	MinLength                int `json:"min_length"`
	MaxLength                int `json:"max_length"`
	MinNumOfCharKind         int `json:"min_num_of_char_kind"`
	MinNumOfConsecutiveChars int `json:"min_num_of_consecutive_chars"`
}

// SetDefault sets default values for PasswordConfig if they are not provided.
func (cfg *PasswordConfig) SetDefault() {
	if cfg.MinLength <= 0 {
		cfg.MinLength = 8
	}

	if cfg.MaxLength <= 0 {
		cfg.MaxLength = 20
	}

	if cfg.MinNumOfCharKind <= 0 {
		cfg.MinNumOfCharKind = 3
	}

	if cfg.MinNumOfConsecutiveChars <= 0 {
		cfg.MinNumOfConsecutiveChars = 2
	}
}
