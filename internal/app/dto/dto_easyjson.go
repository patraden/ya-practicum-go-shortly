// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dto

import (
	json "encoding/json"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"

	domain "github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(in *jlexer.Lexer, out *UserSlugBatch) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(UserSlugBatch, 0, 4)
			} else {
				*out = UserSlugBatch{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 domain.Slug
			v1 = domain.Slug(in.String())
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(out *jwriter.Writer, in UserSlugBatch) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			out.String(string(v3))
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v UserSlugBatch) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserSlugBatch) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserSlugBatch) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserSlugBatch) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(in *jlexer.Lexer, out *UserSlug) {
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
		case "Slug":
			out.Slug = domain.Slug(in.String())
		case "UserID":
			if in.IsNull() {
				in.Skip()
			} else {
				copy(out.UserID[:], in.Bytes())
			}
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
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(out *jwriter.Writer, in UserSlug) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Slug\":"
		out.RawString(prefix[1:])
		out.String(string(in.Slug))
	}
	{
		const prefix string = ",\"UserID\":"
		out.RawString(prefix)
		out.Base64Bytes(in.UserID[:])
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserSlug) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserSlug) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserSlug) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserSlug) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto1(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto2(in *jlexer.Lexer, out *URLPairBatch) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(URLPairBatch, 0, 2)
			} else {
				*out = URLPairBatch{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v6 URLPair
			(v6).UnmarshalEasyJSON(in)
			*out = append(*out, v6)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto2(out *jwriter.Writer, in URLPairBatch) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v7, v8 := range in {
			if v7 > 0 {
				out.RawByte(',')
			}
			(v8).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v URLPairBatch) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v URLPairBatch) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *URLPairBatch) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *URLPairBatch) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto2(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto3(in *jlexer.Lexer, out *URLPair) {
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
		case "short_url":
			out.Slug = domain.Slug(in.String())
		case "original_url":
			out.OriginalURL = domain.OriginalURL(in.String())
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
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto3(out *jwriter.Writer, in URLPair) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"short_url\":"
		out.RawString(prefix[1:])
		out.String(string(in.Slug))
	}
	{
		const prefix string = ",\"original_url\":"
		out.RawString(prefix)
		out.String(string(in.OriginalURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v URLPair) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v URLPair) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *URLPair) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *URLPair) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto3(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto4(in *jlexer.Lexer, out *SlugBatch) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(SlugBatch, 0, 2)
			} else {
				*out = SlugBatch{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v9 CorrelatedSlug
			(v9).UnmarshalEasyJSON(in)
			*out = append(*out, v9)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto4(out *jwriter.Writer, in SlugBatch) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v10, v11 := range in {
			if v10 > 0 {
				out.RawByte(',')
			}
			(v11).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v SlugBatch) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SlugBatch) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SlugBatch) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SlugBatch) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto4(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto5(in *jlexer.Lexer, out *ShortenedURLResponse) {
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
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto5(out *jwriter.Writer, in ShortenedURLResponse) {
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
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortenedURLResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortenedURLResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortenedURLResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto5(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto6(in *jlexer.Lexer, out *ShortenURLRequest) {
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
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto6(out *jwriter.Writer, in ShortenURLRequest) {
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
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortenURLRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortenURLRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortenURLRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto6(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto7(in *jlexer.Lexer, out *RepoStats) {
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
		case "urls":
			out.CountSlugs = int64(in.Int64())
		case "users":
			out.CountUsers = int64(in.Int64())
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
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto7(out *jwriter.Writer, in RepoStats) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"urls\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.CountSlugs))
	}
	{
		const prefix string = ",\"users\":"
		out.RawString(prefix)
		out.Int64(int64(in.CountUsers))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RepoStats) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RepoStats) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RepoStats) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RepoStats) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto7(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto8(in *jlexer.Lexer, out *OriginalURLBatch) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(OriginalURLBatch, 0, 2)
			} else {
				*out = OriginalURLBatch{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v12 CorrelatedOriginalURL
			(v12).UnmarshalEasyJSON(in)
			*out = append(*out, v12)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto8(out *jwriter.Writer, in OriginalURLBatch) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v13, v14 := range in {
			if v13 > 0 {
				out.RawByte(',')
			}
			(v14).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v OriginalURLBatch) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OriginalURLBatch) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OriginalURLBatch) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OriginalURLBatch) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto8(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto9(in *jlexer.Lexer, out *CorrelatedSlug) {
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
		case "correlation_id":
			out.CorrelationID = string(in.String())
		case "short_url":
			out.Slug = domain.Slug(in.String())
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
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto9(out *jwriter.Writer, in CorrelatedSlug) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"correlation_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.CorrelationID))
	}
	{
		const prefix string = ",\"short_url\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CorrelatedSlug) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CorrelatedSlug) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CorrelatedSlug) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CorrelatedSlug) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto9(l, v)
}
func easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto10(in *jlexer.Lexer, out *CorrelatedOriginalURL) {
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
		case "correlation_id":
			out.CorrelationID = string(in.String())
		case "original_url":
			out.OriginalURL = domain.OriginalURL(in.String())
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
func easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto10(out *jwriter.Writer, in CorrelatedOriginalURL) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"correlation_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.CorrelationID))
	}
	{
		const prefix string = ",\"original_url\":"
		out.RawString(prefix)
		out.String(string(in.OriginalURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CorrelatedOriginalURL) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CorrelatedOriginalURL) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComPatradenYaPracticumGoShortlyInternalAppDto10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CorrelatedOriginalURL) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CorrelatedOriginalURL) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComPatradenYaPracticumGoShortlyInternalAppDto10(l, v)
}
