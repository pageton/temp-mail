// Package postfix contains the Postfix configuration for the application.
package postfix

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/pageton/temp-mail/config"
)

type Config struct {
	// General
	Banner             string
	Biff               string
	AppendDotMyDomain  string
	ReadmeDirectory    string
	CompatibilityLevel string

	// Server TLS
	TlSCertFile           string
	TlSKeyFile            string
	TlSSecurityLevel      string
	RelayRestrictions     string
	RecipientRestrictions string

	// Client TLS
	SMTPTLSCAPath        string
	SMTPTLSSecurityLevel string
	SMTPTLSCacheDB       string

	// Mail server
	MyHostname          string
	AliasMaps           string
	AliasDatabase       string
	MyOrigin            string
	MyDestination       string
	RelayHost           string
	Mynetworks          string
	MailboxCommand      string
	MailboxSizeLimit    string
	RecipientDelimiter  string
	VirtualAliasDomains string
	VirtualAliasMaps    string
	InetInterfaces      string
	InetProtocols       string
}

func GenerateMainCF(cfg Config, outputPath string) error {
	tmpl, err := template.ParseFiles("internal/postfix/templates/main.cf.tmpl")
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, cfg)
}

func GenerateVirtualRegexpFile(domains []string, filePath string) error {
	var lines []string

	for _, d := range domains {
		escapedDomain := strings.ReplaceAll(d, ".", `\.`)
		line := fmt.Sprintf(`/.+@%s/ catchall@%s`, escapedDomain, d)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n") + "\n"

	return os.WriteFile(filePath, []byte(content), 0o644)
}

func GenerateAliasesFile(filePath string) error {
	content := `# See man 5 aliases for format
postmaster:    root
catchall: "|/usr/local/bin/forward-to-webhook.sh"
`

	return os.WriteFile(filePath, []byte(content), 0o644)
}

func GenerateForwardScript(filePath string, cfg *config.Config) error {
	content := `#!/bin/bash
sed '/Content-Disposition: attachment/,/^\s*$/d; /Content-Disposition: inline/,/^\s*$/d' | \
curl -X POST -H "Content-Type: text/plain" -H "Secret: %s" --data-binary @- http://localhost:%d/webhook
`
	script := fmt.Sprintf(content, strconv.Itoa(cfg.Server.Secret), cfg.Server.Port)

	if err := os.WriteFile(filePath, []byte(script), 0o755); err != nil {
		return err
	}
	return nil
}

// DefaultConfig returns the default configuration for Postfix.
func DefaultConfig() Config {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		panic(err)
	}
	var domains string
	var mailSubdomains string

	if len(cfg.Domains.Aliases) > 0 {
		domains = strings.Join(cfg.Domains.Aliases, ", ")

		var prefixed []string
		for _, d := range cfg.Domains.Aliases[1:] {
			prefixed = append(prefixed, "mail."+d)
		}
		mailSubdomains = strings.Join(prefixed, ", ")
	}
	err = GenerateVirtualRegexpFile(cfg.Domains.Aliases, "/etc/postfix/virtual_regexp")
	if err != nil {
		panic(err)
	}
	err = GenerateForwardScript("/usr/local/bin/forward-to-webhook.sh", cfg)
	if err != nil {
		panic(err)
	}
	return Config{
		// General
		Banner:             "$myhostname ESMTP $mail_name (Ubuntu)",
		Biff:               "no",
		AppendDotMyDomain:  "no",
		ReadmeDirectory:    "no",
		CompatibilityLevel: "3.6",

		// Server TLS
		TlSCertFile:           "/etc/ssl/certs/ssl-cert-snakeoil.pem",
		TlSKeyFile:            "/etc/ssl/private/ssl-cert-snakeoil.key",
		TlSSecurityLevel:      "may",
		RelayRestrictions:     "permit_mynetworks permit_sasl_authenticated defer_unauth_destination",
		RecipientRestrictions: "permit_mynetworks, reject_unauth_destination",

		// Client TLS
		SMTPTLSCAPath:        "/etc/ssl/certs",
		SMTPTLSSecurityLevel: "may",
		SMTPTLSCacheDB:       "btree:${data_directory}/smtp_scache",

		// Mail server
		MyHostname:    fmt.Sprintf("mail.%s", cfg.Domains.Aliases[0]),
		AliasMaps:     "hash:/etc/aliases",
		AliasDatabase: "hash:/etc/aliases",
		MyOrigin:      "/etc/mailname",
		MyDestination: fmt.Sprintf(
			"$myhostname, %s, %s, localhost.%s, localhost",
			mailSubdomains,
			domains,
			cfg.Domains.Aliases[0],
		),
		RelayHost:           "",
		Mynetworks:          "127.0.0.0/8 [::ffff:127.0.0.0]/104 [::1]/128",
		MailboxCommand:      "procmail -a \"$EXTENSION\"",
		MailboxSizeLimit:    "0",
		RecipientDelimiter:  "+",
		VirtualAliasDomains: domains,
		VirtualAliasMaps:    "regexp:/etc/postfix/virtual_regexp",
		InetInterfaces:      "all",
		InetProtocols:       "all",
	}
}
