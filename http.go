package main

import (
	"fmt"
	"net/http"
	"text/tabwriter"
)

func statHandler(stat *statistic) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		stat.RLock()
		elements := stat.stream.Keys()
		stat.RUnlock()

		w := tabwriter.NewWriter(rw, 10, 4, 5, ' ', tabwriter.Debug)
		defer w.Flush()
		for _, el := range elements {
			fmt.Fprintf(w, "%s\t%d\n", el.Key, el.Count)
		}

		rw.WriteHeader(http.StatusOK)
	}
}
