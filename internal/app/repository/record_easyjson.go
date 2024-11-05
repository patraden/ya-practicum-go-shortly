// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package repository

import (
	json "encoding/json"

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

func easyjson15d5d517DecodeGithubComPatradenYaPracticumGoShortlyInternalAppRepository(in *jlexer.Lexer, out *FileRecord) {
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
		case "uuid":
			out.ID = int(in.Int())
		case "short_url":
			out.ShortURL = string(in.String())
		case "original_url":
			out.LongURL = string(in.String())
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
func easyjson15d5d517EncodeGithubComPatradenYaPracticumGoShortlyInternalAppRepository(out *jwriter.Writer, in FileRecord) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"uuid\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"short_url\":"
		out.RawString(prefix)
		out.String(string(in.ShortURL))
	}
	{
		const prefix string = ",\"original_url\":"
		out.RawString(prefix)
		out.String(string(in.LongURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FileRecord) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson15d5d517EncodeGithubComPatradenYaPracticumGoShortlyInternalAppRepository(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FileRecord) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson15d5d517EncodeGithubComPatradenYaPracticumGoShortlyInternalAppRepository(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FileRecord) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson15d5d517DecodeGithubComPatradenYaPracticumGoShortlyInternalAppRepository(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FileRecord) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson15d5d517DecodeGithubComPatradenYaPracticumGoShortlyInternalAppRepository(l, v)
}
