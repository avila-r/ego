package failure

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/avila-r/ego/failure/ctx"
	"github.com/avila-r/ego/failure/property"
	"github.com/avila-r/ego/failure/stacktrace"
	"github.com/avila-r/ego/failure/tags"
	"github.com/avila-r/ego/failure/trail"
	"github.com/avila-r/ego/failure/trait"
)

var (
	_ error          = (*Failure)(nil)
	_ Error          = (*Failure)(nil)
	_ fmt.Formatter  = (*Failure)(nil)
	_ slog.LogValuer = (*Failure)(nil)
)

// Failure represents a rich error with extended functionality.
type Failure struct {
	class *ErrorClass

	message string
	cause   error

	stacktrace    *stacktrace.StackTrace
	properties    *property.List
	transparent   bool
	hasUnderlying bool
	ppc           uint8

	time     time.Time
	duration time.Duration

	domain  string
	tags    tags.Tags
	context ctx.Context

	trace string
	span  string

	hint   string
	public string
	owner  string

	trail *trail.Trail
}

// Error implements the error interface.
// A result is the same as with %s formatter and does not contain a stack trace.
func (e *Failure) Error() string {
	return e.message
}

// Format implements the Formatter interface.
//
// Supported verbs:
//
//	%s		simple message output
//	%v		simple message output
//	%+v		full output complete with a stack trace
//
// In is nearly always preferable to use %+v format.
// If a stack trace is not required, it should be omitted
// at the moment of creation rather in formatting.
func (e *Failure) Format(state fmt.State, verb rune) {
	switch message := e.Summary(); verb {
	case 'v':
		_, _ = io.WriteString(state, message)
		if state.Flag('+') && e.stacktrace != nil {
			e.stacktrace.Format(state, verb)
		}
	case 's':
		_, _ = io.WriteString(state, message)
	}
}

func (e *Failure) LogValue() slog.Value {
	return e.Logs()
}

func (e *Failure) Class() *ErrorClass {
	cause := e
	for cause != nil {
		if !cause.transparent {
			return cause.class
		}

		cause = Cast(cause.cause)
	}

	return foreignClass
}

func (e *Failure) Message() string {
	return e.message
}

func (e *Failure) Cause() error {
	return e.cause
}

func (e *Failure) Time() time.Time {
	return Deep(e, func(e *Failure) time.Time {
		return e.time
	})
}

func (e *Failure) Duration() time.Duration {
	return Deep(e, func(e *Failure) time.Duration {
		return e.duration
	})
}

func (e *Failure) Domain() string {
	return Deep(e, func(e *Failure) string {
		return e.domain
	})
}

func (e *Failure) Tags() tags.Tags {
	result := tags.Tags{}
	Recurse(e, func(e *Failure) {
		tags.Merge(e.tags, &result)
	})
	return result
}

func (e *Failure) Context() ctx.Context {
	context := Gather(e, func(e *Failure) ctx.Context {
		return e.context
	})

	return ctx.Evaluated(context)
}

func (e *Failure) Trace() string {
	trace := Deep(e, func(e *Failure) string {
		return e.trace
	})

	if trace != "" {
		return trace
	}

	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
}

func (e *Failure) Hint() string {
	return Deep(e, func(e *Failure) string {
		return e.hint
	})
}

func (e *Failure) Public() string {
	return Deep(e, func(e *Failure) string {
		return e.public
	})
}

func (e *Failure) Owner() string {
	return Deep(e, func(e *Failure) string {
		return e.owner
	})
}

func (e *Failure) Span() string {
	return e.span
}

func (e *Failure) Trail() string {
	blocks := []string{}
	topFrame := ""

	Recurse(e, func(e *Failure) {
		if e.trail != nil && len(e.trail.Frames) > 0 {
			err := ""
			if e.cause != nil {
				err = e.cause.Error()
			}
			msg := func(values ...string) string {
				var zero string
				for _, v := range values {
					if v != zero {
						return v
					}
				}
				return zero
			}(e.message, err, "Error")
			block := fmt.Sprintf("%s\n%s", msg, e.trail.String(topFrame))
			blocks = append([]string{block}, blocks...)
			topFrame = e.trail.Frames[0].String()
		}
	})

	if len(blocks) == 0 {
		return ""
	}

	return "Failure: " + strings.Join(blocks, "\nThrown: ")
}

func (o *Failure) Sources() string {
	blocks := [][]string{}
	Recurse(o, func(e *Failure) {
		if e.trail != nil && len(e.trail.Frames) > 0 {
			header, body := e.trail.Source()

			if e.message != "" {
				header = fmt.Sprintf("%s\n%s", e.message, header)
			}

			if header != "" && len(body) > 0 {
				blocks = append(
					[][]string{append([]string{header}, body...)},
					blocks...,
				)
			}
		}
	})

	if len(blocks) == 0 {
		return ""
	}

	return "Failure: " + strings.Join(
		func() []string {
			trails := make([]string, len(blocks))
			for i := range blocks {
				trails[i] = strings.Join(blocks[i], "\n")
			}
			return trails
		}(),
		"\n\nThrown: ",
	)
}

func (e *Failure) WithOwner(owner string) Error {
	e.owner = owner
	return e
}

func (e *Failure) WithPublic(public string) Error {
	e.public = public
	return e
}

func (e *Failure) WithHint(hint string) Error {
	e.hint = hint
	return e
}

func (e *Failure) WithSpan(span string) Error {
	e.span = span
	return e
}

func (e *Failure) WithTrace(trace string) Error {
	e.trace = trace
	return e
}

func (e *Failure) WithTags(tags tags.Tags) Error {
	e.tags = tags
	return e
}

func (e *Failure) WithDomain(domain string) Error {
	e.domain = domain
	return e
}

func (e *Failure) WithDuration(duration time.Duration) Error {
	e.duration = duration
	return e
}

func (e *Failure) WithDurationSince(t time.Time) Error {
	e.duration = time.Since(t)
	return e
}

func (e *Failure) WithTime(time time.Time) Error {
	e.time = time
	return e
}

func (e *Failure) WithCause(err error) Error {
	e.cause = err
	return e
}

func (e *Failure) Assert(condition bool, message ...any) Error {
	msg, args := "assertion failed on error's constructor", []any{}
	if len(message) > 0 {
		first, ok := message[0].(string)
		if ok {
			msg = first

			if len(message) > 1 {
				args = append(args, message[1:]...)
			}
		}
	}

	if !condition {
		AssertionFailed.
			New(msg, args...).
			Panic()
	}

	return e
}

func (e *Failure) Recover(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if er, ok := r.(error); ok {
				err = e.Also(er)
			} else {
				err = e.Also(Err("%v", r))
			}
		}
	}()

	f()
	return
}

func (e *Failure) Chain() ErrorChain {
	return ErrorChain{e}
}

func (e *Failure) Decorate() {
	e.stacktrace = BuilderFrom(e).
		StackTrace().
		SetupStackTrace(4)
}

func (e *Failure) Enhance() {
	if e.Cause() != nil {
		e.stacktrace = BuilderFrom(e).
			EnhanceStackTrace().
			SetupStackTrace(3)
	}
}

func (e *Failure) Decorated() Error {
	e.stacktrace = BuilderFrom(e).
		StackTrace().
		SetupStackTrace(4)

	return e
}

func (e *Failure) Enhanced() Error {
	if e.Cause() != nil {
		e.stacktrace = BuilderFrom(e).
			EnhanceStackTrace().
			SetupStackTrace(4)
	}

	return e
}

func (e *Failure) Belongs(err error) bool {
	typed := Cast(err)

	return typed != nil && Extends(e, typed.Class())
}

func (e *Failure) Is(err error) bool {
	if err, ok := err.(*Failure); ok {
		return e.message == err.message
	}

	return e.message == err.Error()
}

func (e *Failure) As(target any) bool {
	t := reflect.Indirect(reflect.ValueOf(target)).Interface()

	if err, ok := t.(*Failure); ok {
		if e.message == err.message {
			reflect.ValueOf(target).Elem().Set(reflect.ValueOf(err))
			return true
		}
	}
	return false
}

func (e *Failure) Has(trait trait.Trait) bool {
	cause := e
	for cause != nil {
		if !cause.transparent {
			return cause.class.Has(trait)
		}
		cause = Cast(cause.cause)
	}

	return false
}

func (e *Failure) Extends(c *ErrorClass) bool {
	cause := e
	for cause != nil {
		if !cause.transparent {
			return cause.Class().Is(c)
		}

		cause = func() *Failure {
			raw := cause.cause
			for raw != nil {
				typed := Cast(raw)
				if typed != nil {
					return typed
				}

				raw = errors.Unwrap(raw)
			}

			return nil
		}()
	}

	return false
}

func (e *Failure) Attribute(key string) property.Result {
	return e.Property(key)
}

func (e *Failure) Field(key string) property.Result {
	return e.Property(key)
}

func (e *Failure) Value(key string) property.Result {
	return e.Property(key)
}

func (e *Failure) Property(key string) property.Result {
	cause := e
	for cause != nil {
		value, ok := cause.properties.Get(key)
		if ok {
			return property.Result{
				Value: value,
				Ok:    true,
			}
		}

		if !cause.transparent {
			break
		}

		cause = Cast(cause.cause)
	}

	return property.Result{
		Value: nil,
		Ok:    false,
	}
}

func (e *Failure) With(key string, value any) Error {
	copy := *e
	copy.properties = copy.properties.Set(key, value)
	if copy.ppc < 255 {
		copy.ppc++
	}
	if key == property.Underlying {
		copy.hasUnderlying = true
	}
	return &copy
}

func (e *Failure) Panic() {
	panic(e)
}

func (e *Failure) Join(errs ...error) Error {
	return e.Also(errs...)
}

func (e *Failure) Also(errs ...error) Error {
	var (
		underlying = e.Underlying()
		new        = underlying
	)

	for _, err := range errs {
		if err == nil {
			continue
		}
		new = append(new, err)
	}

	if len(new) == len(underlying) {
		return e
	}

	l := len(new)
	copy := e.With(property.Underlying, new[:l:l])
	return copy
}

func (e *Failure) Unwrap() error {
	if e != nil && e.cause != nil && e.transparent {
		return e.cause
	} else {
		return nil
	}
}

func (e *Failure) Summary() string {
	var join = func(delimiter string, parts ...string) string {
		switch len(parts) {
		case 0:
			return ""
		case 1:
			return parts[0]
		case 2:
			if len(parts[0]) == 0 {
				return parts[1]
			} else if len(parts[1]) == 0 {
				return parts[0]
			} else {
				return parts[0] + delimiter + parts[1]
			}
		default:
			filtered := make([]string, 0, len(parts))
			for _, part := range parts {
				if len(part) > 0 {
					filtered = append(filtered, part)
				}
			}

			return strings.Join(filtered, delimiter)
		}
	}

	properties := ""
	if e.properties != nil && e.ppc != 0 {
		var (
			uniq = make(map[string]struct{}, e.ppc)
			strs = make([]string, 0, e.ppc)
		)

		for m := e.properties; m != nil; m = m.Next {
			if _, ok := uniq[m.Key]; ok {
				continue
			}
			uniq[m.Key] = struct{}{}
			strs = append(strs, fmt.Sprintf("%s: %v", m.Key, m.Value))
		}

		properties = "{" + strings.Join(strs, ", ") + "}"
	}

	text := join(" ", e.message, properties)
	if cause := e.cause; cause != nil {
		text = join(", cause: ", text, cause.Error())
	}

	underlying := ""
	if e.hasUnderlying {
		details := make([]string, 0, len(e.Underlying()))
		for _, err := range e.Underlying() {
			details = append(details, err.Error())
		}
		underlying = fmt.Sprintf("(hidden: %s)", join(", ", details...))
	}

	if transparent := join(" ", text, underlying); e.transparent {
		return transparent
	} else {
		return join(": ", e.class.Name, transparent)
	}
}

func (e *Failure) Underlying() []error {
	if !e.hasUnderlying {
		return nil
	}
	u, _ := e.properties.Get(property.Underlying)
	return u.([]error)
}

func (e *Failure) Logs() slog.Value {
	attrs := []slog.Attr{slog.String("message", e.message)}

	if err := e.Error(); err != "" {
		attrs = append(attrs, slog.String("err", err))
	}

	if t := e.Time(); t != (time.Time{}) {
		attrs = append(attrs, slog.Time("time", t.In(time.UTC)))
	}

	if duration := e.Duration(); duration != 0 {
		attrs = append(attrs, slog.Duration("duration", duration))
	}

	if domain := e.Domain(); domain != "" {
		attrs = append(attrs, slog.String("domain", domain))
	}

	if tags := e.Tags(); len(tags) > 0 {
		attrs = append(attrs, slog.Any("tags", tags))
	}

	if trace := e.Trace(); trace != "" {
		attrs = append(attrs, slog.String("trace", trace))
	}

	if hint := e.Hint(); hint != "" {
		attrs = append(attrs, slog.String("hint", hint))
	}

	if public := e.Public(); public != "" {
		attrs = append(attrs, slog.String("public", public))
	}

	if owner := e.Owner(); owner != "" {
		attrs = append(attrs, slog.String("owner", owner))
	}

	if context := e.Context(); len(context) > 0 {
		attrs = append(attrs,
			slog.Group(
				"context",
				func() []any {
					collection := func() []slog.Attr {
						result := make([]slog.Attr, 0, len(context))
						for k := range context {
							result = append(result, slog.Any(k, context[k]))
						}
						return result
					}()

					result := make([]any, len(collection))
					for i := range collection {
						result[i] = collection[i]
					}
					return result
				}()...,
			),
		)
	}

	if trail := e.Trail(); trail != "" {
		attrs = append(attrs, slog.String("stacktrace", trail))
	}

	return slog.GroupValue(attrs...)
}
