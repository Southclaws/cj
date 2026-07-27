package logs

import (
	"go.uber.org/zap/zapcore"
)

type Core struct {
	zapcore.LevelEnabler
	buffer   *Buffer
	stream   *Stream
	redactor *Redactor
	fields   []zapcore.Field
}

func NewCore(enabler zapcore.LevelEnabler, buffer *Buffer, stream *Stream, redactor *Redactor) *Core {
	return &Core{LevelEnabler: enabler, buffer: buffer, stream: stream, redactor: redactor}
}

func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	merged := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	merged = append(merged, c.fields...)
	merged = append(merged, fields...)
	return &Core{
		LevelEnabler: c.LevelEnabler,
		buffer:       c.buffer,
		stream:       c.stream,
		redactor:     c.redactor,
		fields:       merged,
	}
}

func (c *Core) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	all := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	all = append(all, c.fields...)
	all = append(all, fields...)

	record := Record{
		Timestamp: entry.Time,
		Level:     entry.Level.String(),
		Component: entry.LoggerName,
		Message:   entry.Message,
		Fields:    fieldsToMap(all),
	}
	record = c.redactor.Redact(record)

	c.buffer.Push(record)
	c.stream.Broadcast(record)

	return nil
}

func (c *Core) Sync() error {
	return nil
}
