// Print the local time in a loop near the top of each second.
//
// I use this for setting mechanical wristwatches.  This is of limited
// accuracy, but then, so are my watches.

package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/sean-/sysexits"
	"github.com/spf13/pflag"
)

var (
	help     bool
	paranoid bool
	useUTC   bool
	format   string

	formats = map[string]string{
		"1123":  time.RFC1123,
		"1123z": time.RFC1123Z,
		"3339":  time.RFC3339,
		"822":   time.RFC822,
		"822z":  time.RFC822Z,
		"c":     time.ANSIC,
		"date":  time.UnixDate,
	}
)

func init() {
	pflag.BoolVarP(&paranoid, "paranoid", "p", false,
		"be paranoid and look for time to get squirrely")
	pflag.BoolVarP(&useUTC, "utc", "u", false, "use UTC (instead of local time)")
	pflag.BoolVarP(&help, "help", "h", false, "this help")
	// overengineered?  who, me?
	pflag.StringVar(&format, "format", "1123z",
		fmt.Sprintf("set output format, one of %v",
			keys(formats)))
}

func keys(m map[string]string) []string {
	r := []string{}
	for k, _ := range m {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}

func now() time.Time {
	t := time.Now()
	if useUTC {
		t = t.UTC()
	}
	return t
}

func main() {
	pflag.Parse()
	if help {
		pflag.Usage()
		return
	}
	timeFmt, ok := formats[format]
	if !ok {
		fmt.Printf("unknown --format=%q\n", format)
		os.Exit(sysexits.Usage)
	}
	for {
		n := now()
		trunc := n.Truncate(time.Second)
		z := (1 * time.Second) - n.Sub(trunc)
		ts := trunc.Format(timeFmt)
		fmt.Printf("\r%v (sleep %v)      ", ts, z)
		time.Sleep(z)
		if paranoid {
			later := now()
			lt := later.Truncate(time.Second)
			if lt != trunc.Add(1*time.Second) {
				fmt.Printf("\n slept %v and now it's %v\n", later.Sub(n), later)
			}
		}
	}
}
