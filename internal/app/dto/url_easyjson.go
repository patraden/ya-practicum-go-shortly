// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dto

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

func easyjsonF48b0fb9DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(in *jlexer.Lexer, out *ShortenedURLResponse) {
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
		case "result":
			out.ShortURL = string(in.String())
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
func easyjsonF48b0fb9EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(out *jwriter.Writer, in ShortenedURLResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"result\":"
		out.RawString(prefix[1:])
		out.String(string(in.ShortURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ShortenedURLResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF48b0fb9EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortenedURLResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF48b0fb9EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortenedURLResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF48b0fb9DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortenedURLResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF48b0fb9DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(l, v)
}
func easyjsonF48b0fb9DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(in *jlexer.Lexer, out *ShortenURLRequest) {
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
		case "url":
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
func easyjsonF48b0fb9EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(out *jwriter.Writer, in ShortenURLRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		out.String(string(in.LongURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ShortenURLRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF48b0fb9EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortenURLRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF48b0fb9EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortenURLRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF48b0fb9DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortenURLRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF48b0fb9DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(l, v)
}
