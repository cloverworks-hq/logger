package logger

type Options struct {
	Level  string
	Output string
}

type Option func(*Options)

// WithLevel sets the logging level (trace, debug, info, warn, error, fatal).
func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithOutput sets the output target (e.g., "stdout" or a file path).
func WithOutput(output string) Option {
	return func(o *Options) {
		o.Output = output
	}
}

// defaultOptions provides sensible defaults for the logger.
func defaultOptions() *Options {
	return &Options{
		Level:  "info",
		Output: "stdout",
	}
}
