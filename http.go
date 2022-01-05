package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"

	"github.com/vanhtuan0409/dnsstat/internal/topk"
)

func statHandler(stat *statistic) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		stat.RLock()
		elements := stat.stream.Keys()
		stat.RUnlock()

		var f formatter
		switch r.URL.Query().Get("format") {
		case "csv":
			f = new(csvFormatter)
		default:
			f = new(tabFormatter)
		}

		rw.WriteHeader(http.StatusOK)
		f.Format(rw, elements)
	}
}

type formatter interface {
	Format(w io.Writer, elements []topk.Element) error
}

type tabFormatter struct{}

func (f *tabFormatter) Format(w io.Writer, elements []topk.Element) error {
	tw := tabwriter.NewWriter(w, 10, 4, 5, ' ', tabwriter.Debug)
	defer tw.Flush()
	fmt.Fprintf(tw, "%s\t%s\n", "Domain", "Count")
	for _, el := range elements {
		fmt.Fprintf(tw, "%s\t%d\n", el.Key, el.Count)
	}
	return nil
}

type csvFormatter struct{}

func (f *csvFormatter) Format(w io.Writer, elements []topk.Element) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()
	cw.Write([]string{"Domain", "Count"})
	for _, el := range elements {
		cw.Write([]string{
			el.Key,
			fmt.Sprintf("%d", el.Count),
		})
	}
	return nil
}
