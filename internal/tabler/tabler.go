package tabler

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// RenderTable renders a table based on the provided TableOptions.
func PrintAsTable(opts ...TableOption) error {
	cfg := defaultTableConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	if err := cfg.validate(); err != nil {
		return fmt.Errorf("invalid table configuration: %w", err)
	}

	table := tablewriter.NewWriter(cfg.writer)
	table.Options(
		tablewriter.WithRendition(cfg.rendition),
		tablewriter.WithAlignment(cfg.estimateAlignment()),
		tablewriter.WithHeader(cfg.header),
	)

	for _, row := range cfg.rows {
		table.Append(row)
	}

	table.Render()

	return nil
}

// TableOption allows callers to customize tablewriter.Table before rendering.
type TableOption func(*TableConfig)

// TableConfig describes the data and presentation settings used when rendering a table.
type TableConfig struct {
	writer    io.Writer
	header    []string
	rows      [][]string
	rendition tw.Rendition
}

func (cfg *TableConfig) validate() error {
	if cfg.writer == nil {
		return fmt.Errorf("table writer must not be nil")
	}

	if len(cfg.rows) == 0 {
		return fmt.Errorf("table must have at least one row")
	}

	if len(cfg.rows[0]) == 0 {
		return fmt.Errorf("table rows must have at least one column")
	}

	if len(cfg.header) != 0 && len(cfg.header) != len(cfg.rows[0]) {
		return fmt.Errorf("table header length (%d) does not match number of columns in rows (%d)",
			len(cfg.header), len(cfg.rows[0]))
	}

	return nil
}

// estimateAlignment estimates the alignment for each column
// based on the number of columns in the rows.
func (cfg *TableConfig) estimateAlignment() []tw.Align {
	if len(cfg.rows) == 0 {
		return []tw.Align{}
	}

	numCols := len(cfg.rows[0])
	alignment := make([]tw.Align, numCols)
	for i := 0; i < numCols; i++ {
		alignment[i] = tw.AlignLeft
	}
	return alignment
}

func defaultTableConfig() TableConfig {
	return TableConfig{
		writer: os.Stdout,
		header: []string{},
		rows:   [][]string{},
		rendition: tw.Rendition{
			Borders: tw.BorderNone,
			Settings: tw.Settings{
				Separators: tw.Separators{
					BetweenColumns: tw.Off,
				},
				Lines: tw.Lines{
					ShowHeaderLine: tw.Off,
				},
			},
		},
	}
}

// WithWriter sets the writer for the table output.
func WithWriter(w io.Writer) TableOption {
	return func(cfg *TableConfig) {
		cfg.writer = w
	}
}

// WithHeader sets the header row for the table.
func WithHeader(header []string) TableOption {
	return func(cfg *TableConfig) {
		cfg.header = header
	}
}

// WithRows sets the data rows for the table.
func WithRows(rows [][]string) TableOption {
	return func(cfg *TableConfig) {
		cfg.rows = rows
	}
}
