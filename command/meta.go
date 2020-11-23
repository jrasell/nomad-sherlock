package command

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/jrasell/nomad-sherlock/api"

	"github.com/mitchellh/cli"
)

type Meta struct {
	Ui      cli.Ui
	address string
}

type FlagSetFlags uint

const (
	FlagSetNone   FlagSetFlags = 0
	FlagSetClient FlagSetFlags = 1 << iota
)

func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	if fs&FlagSetClient != 0 {
		f.StringVar(&m.address, "address", "", "")
	}

	f.SetOutput(&uiErrorWriter{ui: m.Ui})

	return f
}

func (m *Meta) Client() (*api.Client, error) {
	return api.NewClient(m.address)
}

type uiErrorWriter struct {
	ui  cli.Ui
	buf bytes.Buffer
}

func (w *uiErrorWriter) Write(data []byte) (int, error) {
	read := 0
	for len(data) != 0 {
		a, token, err := bufio.ScanLines(data, false)
		if err != nil {
			return read, err
		}

		if a == 0 {
			r, err := w.buf.Write(data)
			return read + r, err
		}

		w.ui.Error(w.buf.String() + string(token))
		data = data[a:]
		w.buf.Reset()
		read += a
	}

	return read, nil
}

func (w *uiErrorWriter) Close() error {
	if w.buf.Len() != 0 {
		w.ui.Error(w.buf.String())
		w.buf.Reset()
	}
	return nil
}
