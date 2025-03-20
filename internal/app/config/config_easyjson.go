// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package config

import (
	json "encoding/json"
	time "time"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6615c02eDecodeGithubComPatradenYaPracticumGoShortlyInternalAppConfig(in *jlexer.Lexer, out *Config) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "server_address":
			out.ServerAddr = string(in.String())
		case "base_url":
			out.BaseURL = string(in.String())
		case "file_storage_path":
			out.FileStoragePath = string(in.String())
		case "database_dsn":
			out.DatabaseDSN = string(in.String())
		case "enable_https":
			out.EnableHTTPS = bool(in.Bool())
		case "jwt_secret":
			out.JWTSecret = string(in.String())
		case "tlc_key_path":
			out.TLSKeyPath = string(in.String())
		case "tlc_cert_path":
			out.TLSCertPath = string(in.String())
		case "ConfigJSON":
			out.ConfigJSON = string(in.String())
		case "URLGenTimeout":
			out.URLGenTimeout = time.Duration(in.Int64())
		case "URLGenRetryInterval":
			out.URLGenRetryInterval = time.Duration(in.Int64())
		case "URLsize":
			out.URLsize = int(in.Int())
		case "ServerShutTimeout":
			out.ServerShutTimeout = time.Duration(in.Int64())
		case "ServerReadHeaderTimeout":
			out.ServerReadHeaderTimeout = time.Duration(in.Int64())
		case "ServerWriteTimeout":
			out.ServerWriteTimeout = time.Duration(in.Int64())
		case "ServerIdleTimeout":
			out.ServerIdleTimeout = time.Duration(in.Int64())
		case "ForceEmptyRepo":
			out.ForceEmptyRepo = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeGithubComPatradenYaPracticumGoShortlyInternalAppConfig(out *jwriter.Writer, in Config) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"server_address\":"
		out.RawString(prefix[1:])
		out.String(string(in.ServerAddr))
	}
	{
		const prefix string = ",\"base_url\":"
		out.RawString(prefix)
		out.String(string(in.BaseURL))
	}
	{
		const prefix string = ",\"file_storage_path\":"
		out.RawString(prefix)
		out.String(string(in.FileStoragePath))
	}
	{
		const prefix string = ",\"database_dsn\":"
		out.RawString(prefix)
		out.String(string(in.DatabaseDSN))
	}
	{
		const prefix string = ",\"enable_https\":"
		out.RawString(prefix)
		out.Bool(bool(in.EnableHTTPS))
	}
	{
		const prefix string = ",\"jwt_secret\":"
		out.RawString(prefix)
		out.String(string(in.JWTSecret))
	}
	{
		const prefix string = ",\"tlc_key_path\":"
		out.RawString(prefix)
		out.String(string(in.TLSKeyPath))
	}
	{
		const prefix string = ",\"tlc_cert_path\":"
		out.RawString(prefix)
		out.String(string(in.TLSCertPath))
	}
	{
		const prefix string = ",\"ConfigJSON\":"
		out.RawString(prefix)
		out.String(string(in.ConfigJSON))
	}
	{
		const prefix string = ",\"URLGenTimeout\":"
		out.RawString(prefix)
		out.Int64(int64(in.URLGenTimeout))
	}
	{
		const prefix string = ",\"URLGenRetryInterval\":"
		out.RawString(prefix)
		out.Int64(int64(in.URLGenRetryInterval))
	}
	{
		const prefix string = ",\"URLsize\":"
		out.RawString(prefix)
		out.Int(int(in.URLsize))
	}
	{
		const prefix string = ",\"ServerShutTimeout\":"
		out.RawString(prefix)
		out.Int64(int64(in.ServerShutTimeout))
	}
	{
		const prefix string = ",\"ServerReadHeaderTimeout\":"
		out.RawString(prefix)
		out.Int64(int64(in.ServerReadHeaderTimeout))
	}
	{
		const prefix string = ",\"ServerWriteTimeout\":"
		out.RawString(prefix)
		out.Int64(int64(in.ServerWriteTimeout))
	}
	{
		const prefix string = ",\"ServerIdleTimeout\":"
		out.RawString(prefix)
		out.Int64(int64(in.ServerIdleTimeout))
	}
	{
		const prefix string = ",\"ForceEmptyRepo\":"
		out.RawString(prefix)
		out.Bool(bool(in.ForceEmptyRepo))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Config) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComPatradenYaPracticumGoShortlyInternalAppConfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Config) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComPatradenYaPracticumGoShortlyInternalAppConfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Config) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComPatradenYaPracticumGoShortlyInternalAppConfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Config) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComPatradenYaPracticumGoShortlyInternalAppConfig(l, v)
}
