package main

import (
	"github.com/urfave/cli/v2"
	"net/http"
	"strconv"
)

type Handler struct {
	Ready         bool
	Configuration *SecurityConfiguration
}

func (h *Handler) shouldReady() bool {
	if !h.Ready {
		FailPrintf("Configuration unready")
	}
	return h.Ready
}

func (h *Handler) Setup(c *cli.Context) error {
	newConfiguration := &SecurityConfiguration{
		XAuthEmail: c.String("x-auth-email"),
		XAuthKey:   c.String("x-auth-key"),
	}

	if newConfiguration.XAuthEmail == "" || newConfiguration.XAuthKey == "" {
		return cli.Exit("Invalid configuration", 0)
	}

	return newConfiguration.Save()
}

func (h *Handler) ListDNSRecords(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodGet,
		"/zones/{zone_id}/dns_records",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
		UseQueryParametersWithMap(map[string]string{
			"comment.absent":     c.String("comment-absent"),
			"comment.contains":   c.String("comment-contains"),
			"comment.endswith":   c.String("comment-endswith"),
			"comment.exact":      c.String("comment-exact"),
			"comment.present":    c.String("comment-present"),
			"comment.startswith": c.String("comment-startswith"),
			"content":            c.String("content"),
			"direction":          c.String("direction"),
			"match":              c.String("match"),
			"name":               c.String("name"),
			"order":              c.String("order"),
			"page":               strconv.FormatUint(uint64(c.Uint("page")), 10),
			"per_page":           strconv.FormatUint(uint64(c.Uint("per-page")), 10),
			"proxied":            strconv.FormatBool(c.Bool("proxied")),
			"search":             c.String("search"),
			"tag":                c.String("tag"),
			"tag.absent":         c.String("tag-absent"),
			"tag.contains":       c.String("tag-contains"),
			"tag.endswith":       c.String("tag-endswith"),
			"tag.exact":          c.String("tag-exact"),
			"tag.present":        c.String("tag-present"),
			"tag.startswith":     c.String("tag-startswith"),
			"tag_match":          c.String("tag-match"),
			"type":               c.String("type"),
		}),
	)

	return nil
}

func (h *Handler) CreateDNSRecordA(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodPost,
		"/zones/{zone_id}/dns_records",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
		UseJSONBody(map[string]any{
			"content": c.String("content"),
			"name":    c.String("name"),
			"proxied": c.Bool("proxied"),
			"type":    "A",
			"comment": c.String("comment"),
			"tags":    c.StringSlice("tags"),
			"ttl":     c.Uint64("ttl"),
		}),
	)
	return nil
}

func (h *Handler) CreateDNSRecordAAAA(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodPost,
		"/zones/{zone_id}/dns_records",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
		UseJSONBody(map[string]any{
			"content": c.String("content"),
			"name":    c.String("name"),
			"proxied": c.Bool("proxied"),
			"type":    "AAAA",
			"comment": c.String("comment"),
			"tags":    c.StringSlice("tags"),
			"ttl":     c.Uint64("ttl"),
		}),
	)
	return nil
}

func (h *Handler) ExportDNSRecords(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodGet,
		"/zones/{zone_id}/dns_records/export",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
	)

	return nil
}

func (h *Handler) ScanDNSRecord(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodPost,
		"/zones/{zone_id}/dns_records/scan",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
	)

	return nil
}

func (h *Handler) DeleteDNSRecord(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodDelete,
		"/zones/{zone_id}/dns_records/{dns_record_id}",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
		UsePathParameters("dns_record_id", c.String("record-id")),
	)

	return nil
}

func (h *Handler) DNSRecordDetails(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodGet,
		"/zones/{zone_id}/dns_records/{dns_record_id}",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
		UsePathParameters("dns_record_id", c.String("record-id")),
	)

	return nil
}

func (h *Handler) UpdateDNSRecordA(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodPatch,
		"/zones/{zone_id}/dns_records/{dns_record_id}",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
		UsePathParameters("dns_record_id", c.String("record-id")),
		UseJSONBody(map[string]any{
			"content": c.String("content"),
			"name":    c.String("name"),
			"proxied": c.Bool("proxied"),
			"type":    "A",
			"comment": c.String("comment"),
			"tags":    c.StringSlice("tags"),
			"ttl":     c.Uint64("ttl"),
		}),
	)

	return nil
}

func (h *Handler) UpdateDNSRecordAAAA(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodPatch,
		"/zones/{zone_id}/dns_records/{dns_record_id}",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
		UsePathParameters("dns_record_id", c.String("record-id")),
		UseJSONBody(map[string]any{
			"content": c.String("content"),
			"name":    c.String("name"),
			"proxied": c.Bool("proxied"),
			"type":    "AAAA",
			"comment": c.String("comment"),
			"tags":    c.StringSlice("tags"),
			"ttl":     c.Uint64("ttl"),
		}),
	)

	return nil
}

func (h *Handler) OverwriteDNSRecordA(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodPatch,
		"/zones/{zone_id}/dns_records/{dns_record_id}",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
		UsePathParameters("dns_record_id", c.String("record-id")),
		UseJSONBody(map[string]any{
			"content": c.String("content"),
			"name":    c.String("name"),
			"proxied": c.Bool("proxied"),
			"type":    "A",
			"comment": c.String("comment"),
			"tags":    c.StringSlice("tags"),
			"ttl":     c.Uint64("ttl"),
		}),
	)

	return nil
}

func (h *Handler) OverwriteDNSRecordAAAA(c *cli.Context) error {
	if !h.shouldReady() {
		return nil
	}

	Request(
		http.MethodPatch,
		"/zones/{zone_id}/dns_records/{dns_record_id}",
		UseSecurity(h.Configuration),
		UsePathParameters("zone_id", c.String("zone-id")),
		UsePathParameters("dns_record_id", c.String("record-id")),
		UseJSONBody(map[string]any{
			"content": c.String("content"),
			"name":    c.String("name"),
			"proxied": c.Bool("proxied"),
			"type":    "AAAA",
			"comment": c.String("comment"),
			"tags":    c.StringSlice("tags"),
			"ttl":     c.Uint64("ttl"),
		}),
	)

	return nil
}
