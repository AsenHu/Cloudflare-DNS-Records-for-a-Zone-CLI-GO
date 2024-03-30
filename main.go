package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

const (
	BaseAPI = "https://api.cloudflare.com/client/v4"
)

func main() {
	handler := &Handler{}
	configuration, err := OpenSecurityConfiguration()
	handler.Ready = err == nil && configuration != nil
	handler.Configuration = configuration

	app := &cli.App{
		Name:                 "cf-cli",
		Version:              "0.0.1",
		Usage:                "Cloudflare DNS Records for a Zone shell",
		UsageText:            `cf-cli COMMAND [OPTIONS]`,
		Suggest:              true,
		EnableBashCompletion: true,
		CommandNotFound: func(context *cli.Context, s string) {
			FailPrintf("command '%s' not found", s)
		},
		Commands: []*cli.Command{
			// setup
			{
				Name:    "setup",
				Aliases: []string{"init"},
				Usage:   "Setup security configuration, Eg. 'X-Auth-Email', 'X-Auth-Key'",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "x-auth-email",
						Usage: "Cloudflare auth email",
					},
					&cli.StringFlag{
						Name:  "x-auth-key",
						Usage: "Cloudflare auth key",
					},
					&cli.StringFlag{
						Name:  "api-token",
						Usage: "Cloudflare API token",
					},
				},
				Action: handler.Setup,
			},

			// list
			{
				Name:  "list",
				Usage: "List, search, sort, and filter a zones' DNS records.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "zone-id",
						Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
					},
					&cli.StringFlag{
						Name:  "comment-exact",
						Usage: "Exact value of the DNS record comment. Comment filters are case-insensitive. Eg. Hello, world",
					},
					&cli.StringFlag{
						Name:  "comment-absent",
						Usage: "If this parameter is present, only records without a comment are returned.",
					},
					&cli.StringFlag{
						Name:  "comment-contains",
						Usage: "Substring of the DNS record comment. Comment filters are case-insensitive. Eg. ello, worl",
					},
					&cli.StringFlag{
						Name:  "comment-endswith",
						Usage: "Suffix of the DNS record comment. Comment filters are case-insensitive. Eg. o, world",
					},
					&cli.StringFlag{
						Name:  "comment-present",
						Usage: "If this parameter is present, only records with a comment are returned.",
					},
					&cli.StringFlag{
						Name:  "comment-startswith",
						Usage: "Prefix of the DNS record comment. Comment filters are case-insensitive. Eg. Hello, w",
					},
					&cli.StringFlag{
						Name:  "content",
						Usage: "DNS record content. Eg. 127.0.0.1",
					},
					&cli.StringFlag{
						Name:  "direction",
						Value: "asc",
						Usage: "Direction to order DNS records in. Allowed values: asc, desc",
					},
					&cli.StringFlag{
						Name:  "match",
						Value: "all",
						Usage: "Whether to match all search requirements or at least one (any). If set to all, acts like a logical AND between filters. If set to any, acts like a logical OR instead. Note that the interaction between tag filters is controlled by the tag-match parameter instead. Allowed values: any, all",
					},
					&cli.StringFlag{
						Name:  "name",
						Usage: "DNS record name (or @ for the zone apex) in Punycode. Eg. example.com",
					},
					&cli.StringFlag{
						Name:  "order",
						Value: "type",
						Usage: "Field to order DNS records by. Allowed values: type, name, content, ttl, proxied",
					},
					&cli.UintFlag{
						Name:  "page",
						Value: 1,
						Usage: "Page number of paginated results.",
					},
					&cli.UintFlag{
						Name:  "per-page",
						Value: 100,
						Usage: "Number of DNS records per page. Eg. 5",
					},
					&cli.BoolFlag{
						Name:  "proxied",
						Usage: "Whether the record is receiving the performance and security benefits of Cloudflare.",
					},
					&cli.StringFlag{
						Name:  "search",
						Usage: "Allows searching in multiple properties of a DNS record simultaneously. This parameter is intended for human users, not automation. Its exact behavior is intentionally left unspecified and is subject to change in the future. This parameter works independently of the match setting. For automated searches, please use the other available parameters. Eg. www.cloudflare.com",
					},
					&cli.StringFlag{
						Name:  "tag",
						Usage: "Condition on the DNS record tag. Eg. team:DNS",
					},
					&cli.StringFlag{
						Name:  "tag-absent",
						Usage: "Name of a tag which must not be present on the DNS record. Tag filters are case-insensitive. Eg. important",
					},
					&cli.StringFlag{
						Name:  "tag-endswith",
						Usage: "A tag and value, of the form <tag-name>:<tag-value>. The API will only return DNS records that have a tag named <tag-name> whose value ends with <tag-value>. Tag filters are case-insensitive.",
					},
					&cli.StringFlag{
						Name:  "tag-exact",
						Usage: "A tag and value, of the form <tag-name>:<tag-value>. The API will only return DNS records that have a tag named <tag-name> whose value is <tag-value>. Tag filters are case-insensitive.",
					},
					&cli.StringFlag{
						Name:  "tag-present",
						Usage: "Name of a tag which must be present on the DNS record. Tag filters are case-insensitive.",
					},
					&cli.StringFlag{
						Name:  "tag-startswith",
						Usage: "A tag and value, of the form <tag-name>:<tag-value>. The API will only return DNS records that have a tag named <tag-name> whose value starts with <tag-value>. Tag filters are case-insensitive.",
					},
					&cli.StringFlag{
						Name:  "tag-match",
						Value: "all",
						Usage: "Whether to match all tag search requirements or at least one (any). If set to all, acts like a logical AND between tag filters. If set to any, acts like a logical OR instead. Note that the regular match parameter is still used to combine the resulting condition with other filters that aren't related to tags. Allowed values: any, all",
					},
					&cli.StringFlag{
						Name:  "type",
						Usage: "Record type. Allowed values: A, AAAA, CAA, CERT, CNAME, DNSKEY, DS, HTTPS, LOC, MX, NAPTR, NS, PTR, SMIMEA, SRV, SSHFP, SVCB, TLSA, TXT, URI",
					},
				},
				Action: handler.ListDNSRecords,
			},

			// create
			{
				Name:  "create",
				Usage: "Create a new DNS record for a zone.",
				Subcommands: []*cli.Command{
					{
						Name:    "A",
						Aliases: []string{"a"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "zone-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "content",
								Usage: "A valid IPv4 address.",
							},
							&cli.StringFlag{
								Name:  "name",
								Usage: "DNS record name (or @ for the zone apex) in Punycode.",
							},
							&cli.BoolFlag{
								Name:  "proxied",
								Usage: "Whether the record is receiving the performance and security benefits of Cloudflare.",
							},
							&cli.StringFlag{
								Name:  "comment",
								Usage: "Comments or notes about the DNS record. This field has no effect on DNS responses.",
							},
							&cli.StringSliceFlag{
								Name:  "tags",
								Usage: "Custom tags for the DNS record. This field has no effect on DNS responses.",
							},
							&cli.Uint64Flag{
								Name:  "ttl",
								Usage: "Time To Live (TTL) of the DNS record in seconds. Setting to 1 means 'automatic'. Value must be between 60 and 86400, with the minimum reduced to 30 for Enterprise zones.",
							},
						},
						Action: handler.CreateDNSRecordA,
					},
					{
						Name:    "AAAA",
						Aliases: []string{"aaaa"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "zone-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "content",
								Usage: "A valid IPv4 address.",
							},
							&cli.StringFlag{
								Name:  "name",
								Usage: "DNS record name (or @ for the zone apex) in Punycode.",
							},
							&cli.BoolFlag{
								Name:  "proxied",
								Usage: "Whether the record is receiving the performance and security benefits of Cloudflare.",
							},
							&cli.StringFlag{
								Name:  "comment",
								Usage: "Comments or notes about the DNS record. This field has no effect on DNS responses.",
							},
							&cli.StringSliceFlag{
								Name:  "tags",
								Usage: "Custom tags for the DNS record. This field has no effect on DNS responses.",
							},
							&cli.Uint64Flag{
								Name:  "ttl",
								Usage: "Time To Live (TTL) of the DNS record in seconds. Setting to 1 means 'automatic'. Value must be between 60 and 86400, with the minimum reduced to 30 for Enterprise zones.",
							},
						},
						Action: handler.CreateDNSRecordAAAA,
					},
				},
			},

			// export
			{
				Name:  "export",
				Usage: "You can export your BIND config through this endpoint.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "zone-id",
						Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
					},
				},
				Action: handler.ExportDNSRecords,
			},

			// scan
			{
				Name:  "scan",
				Usage: "Scan for common DNS records on your domain and automatically add them to your zone. Useful if you haven't updated your nameservers yet.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "zone-id",
						Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
					},
				},
				Action: handler.ScanDNSRecord,
			},

			// delete
			{
				Name: "delete",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "zone-id",
						Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
					},
					&cli.StringFlag{
						Name:  "record-id",
						Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
					},
				},
				Action: handler.DeleteDNSRecord,
			},
			{
				Name: "details",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "zone-id",
						Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
					},
					&cli.StringFlag{
						Name:  "record-id",
						Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
					},
				},
				Action: handler.DNSRecordDetails,
			},

			// update
			{
				Name: "update",
				Usage: `Update an existing DNS record. Notes:
A/AAAA records cannot exist on the same name as CNAME records.
NS records cannot exist on the same name as any other record type.
Domain names are always represented in Punycode, even if Unicode characters were used when creating the record.`,
				Subcommands: []*cli.Command{
					{
						Name:    "A",
						Aliases: []string{"a"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "zone-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "record-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "content",
								Usage: "A valid IPv4 address.",
							},
							&cli.StringFlag{
								Name:  "name",
								Usage: "DNS record name (or @ for the zone apex) in Punycode.",
							},
							&cli.BoolFlag{
								Name:  "proxied",
								Usage: "Whether the record is receiving the performance and security benefits of Cloudflare.",
							},
							&cli.StringFlag{
								Name:  "comment",
								Usage: "Comments or notes about the DNS record. This field has no effect on DNS responses.",
							},
							&cli.StringSliceFlag{
								Name:  "tags",
								Usage: "Custom tags for the DNS record. This field has no effect on DNS responses.",
							},
							&cli.Uint64Flag{
								Name:  "ttl",
								Usage: "Time To Live (TTL) of the DNS record in seconds. Setting to 1 means 'automatic'. Value must be between 60 and 86400, with the minimum reduced to 30 for Enterprise zones.",
							},
						},
						Action: handler.UpdateDNSRecordA,
					},
					{
						Name:    "AAAA",
						Aliases: []string{"aaaa"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "zone-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "record-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "content",
								Usage: "A valid IPv4 address.",
							},
							&cli.StringFlag{
								Name:  "name",
								Usage: "DNS record name (or @ for the zone apex) in Punycode.",
							},
							&cli.BoolFlag{
								Name:  "proxied",
								Usage: "Whether the record is receiving the performance and security benefits of Cloudflare.",
							},
							&cli.StringFlag{
								Name:  "comment",
								Usage: "Comments or notes about the DNS record. This field has no effect on DNS responses.",
							},
							&cli.StringSliceFlag{
								Name:  "tags",
								Usage: "Custom tags for the DNS record. This field has no effect on DNS responses.",
							},
							&cli.Uint64Flag{
								Name:  "ttl",
								Usage: "Time To Live (TTL) of the DNS record in seconds. Setting to 1 means 'automatic'. Value must be between 60 and 86400, with the minimum reduced to 30 for Enterprise zones.",
							},
						},
						Action: handler.UpdateDNSRecordAAAA,
					},
				},
			},

			// overwrite
			{
				Name: "overwrite",
				Usage: `Overwrite an existing DNS record. Notes:
A/AAAA records cannot exist on the same name as CNAME records.
NS records cannot exist on the same name as any other record type.
Domain names are always represented in Punycode, even if Unicode characters were used when creating the record.`,
				Subcommands: []*cli.Command{
					{
						Name:    "A",
						Aliases: []string{"a"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "zone-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "record-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "content",
								Usage: "A valid IPv4 address.",
							},
							&cli.StringFlag{
								Name:  "name",
								Usage: "DNS record name (or @ for the zone apex) in Punycode.",
							},
							&cli.BoolFlag{
								Name:  "proxied",
								Usage: "Whether the record is receiving the performance and security benefits of Cloudflare.",
							},
							&cli.StringFlag{
								Name:  "comment",
								Usage: "Comments or notes about the DNS record. This field has no effect on DNS responses.",
							},
							&cli.StringSliceFlag{
								Name:  "tags",
								Usage: "Custom tags for the DNS record. This field has no effect on DNS responses.",
							},
							&cli.Uint64Flag{
								Name:  "ttl",
								Usage: "Time To Live (TTL) of the DNS record in seconds. Setting to 1 means 'automatic'. Value must be between 60 and 86400, with the minimum reduced to 30 for Enterprise zones.",
							},
						},
						Action: handler.OverwriteDNSRecordA,
					},
					{
						Name:    "AAAA",
						Aliases: []string{"aaaa"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "zone-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "record-id",
								Usage: "Identifier, <= 32 characters, Eg. 023e105f4ecef8ad9ca31a8372d0c353",
							},
							&cli.StringFlag{
								Name:  "content",
								Usage: "A valid IPv4 address.",
							},
							&cli.StringFlag{
								Name:  "name",
								Usage: "DNS record name (or @ for the zone apex) in Punycode.",
							},
							&cli.BoolFlag{
								Name:  "proxied",
								Usage: "Whether the record is receiving the performance and security benefits of Cloudflare.",
							},
							&cli.StringFlag{
								Name:  "comment",
								Usage: "Comments or notes about the DNS record. This field has no effect on DNS responses.",
							},
							&cli.StringSliceFlag{
								Name:  "tags",
								Usage: "Custom tags for the DNS record. This field has no effect on DNS responses.",
							},
							&cli.Uint64Flag{
								Name:  "ttl",
								Usage: "Time To Live (TTL) of the DNS record in seconds. Setting to 1 means 'automatic'. Value must be between 60 and 86400, with the minimum reduced to 30 for Enterprise zones.",
							},
						},
						Action: handler.OverwriteDNSRecordAAAA,
					},
				},
			},
		},
	}

	if err = app.Run(os.Args); err != nil {
		FailPrintf(err.Error())
		return
	}
}
