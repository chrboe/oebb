package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/aybabtme/rgbterm"
	"github.com/briandowns/spinner"
	oebb "github.com/chrboe/oebb/client"
	"github.com/spf13/cobra"
)

func parseConnTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.999", str)
}

func formatConnTime(str string) (string, error) {
	t, err := parseConnTime(str)
	if err != nil {
		return "", err
	}
	return t.Format("15:04"), nil
}

func formatDuration(dur int) string {
	minutes := dur / 1000 / 60
	durHours := minutes / 60
	durMinutes := minutes % 60
	durStr := fmt.Sprintf("{#ffff00}%02d:%02d{}", durHours, durMinutes)
	return rgbterm.InterpretStr(durStr)
}

func formatDelayTime(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	delayTime, err := formatConnTime(str)
	if err != nil {
		return "", err
	}
	return rgbterm.InterpretStr("{#ff0000}" + delayTime + "{}"), nil
}

func formatDelayLine(depDelay, arrDelay string, dep, arr *string) string {
	delayLine := ""

	if depDelay == "" {
		delayLine += strings.Repeat(" ", len(*dep)+1)
	} else {
		delayLine += depDelay + " "
		*dep = strikethrough(*dep)
	}

	if arrDelay != "" {
		delayLine += arrDelay
		*arr = strikethrough(*arr)
	}

	return delayLine
}

func displaySection(section oebb.Section) error {
	dep, err := formatConnTime(section.From.Departure)
	if err != nil {
		return err
	}

	arr, err := formatConnTime(section.To.Arrival)
	if err != nil {
		return err
	}

	depDelay, err := formatDelayTime(section.From.DepartureDelay)
	if err != nil {
		return err
	}

	arrDelay, err := formatDelayTime(section.To.ArrivalDelay)
	if err != nil {
		return err
	}

	if depDelay != "" || arrDelay != "" {
		fmt.Println("\t" + formatDelayLine(depDelay, arrDelay, &dep, &arr))
	}

	cname := section.Category.DisplayName
	if cname == "" {
		cname = section.Category.ShortName
	}
	category := rgbterm.InterpretStr(fmt.Sprintf("\033[1m{#ffffff,%s}%-3s{}",
		section.Category.BarColor,
		strings.ToUpper(cname)),
	)

	times := rgbterm.InterpretStr(fmt.Sprintf("{#555555}%s{}{#555555}-{}{#555555}%s{}", dep, arr))
	fmt.Printf("\t%s %s %s -> %s\n", times, category, section.From.Name, section.To.Name)
	return nil
}

func displayConnection(conn oebb.Connection) error {
	dep, err := formatConnTime(conn.From.Departure)
	if err != nil {
		return err
	}

	arr, err := formatConnTime(conn.To.Arrival)
	if err != nil {
		return err
	}

	depDelay, err := formatDelayTime(conn.From.DepartureDelay)
	if err != nil {
		return err
	}

	arrDelay, err := formatDelayTime(conn.To.ArrivalDelay)
	if err != nil {
		return err
	}

	if depDelay != "" || arrDelay != "" {
		fmt.Println(formatDelayLine(depDelay, arrDelay, &dep, &arr))
	}

	durStr := formatDuration(conn.Duration)
	fromStr := bold(rgbterm.InterpretStr("{#cc6666}" + conn.From.Name + "{}"))
	toStr := bold(rgbterm.InterpretStr("{#cc6666}" + conn.To.Name + "{}"))

	fmt.Printf("%s-%s (%s) %s -> %s\n", dep, arr, durStr, fromStr, toStr)
	for _, section := range conn.Sections {
		displaySection(section)
	}

	fmt.Println()
	return nil
}

func cacheAuth(auth oebb.AuthInfo, filename string) error {
	bytes, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, bytes, 0600)
	return err
}

func authAndCache(filenameTemplate string) (*oebb.AuthInfo, error) {
	auth, err := oebb.Auth()
	if err != nil {
		return nil, err
	}

	newCache, err := xdg.CacheFile(filenameTemplate)
	if err != nil {
		return nil, err
	}

	err = cacheAuth(auth, newCache)
	return &auth, err
}

// maybeCachedAuth returns possibly cached authentication information
func maybeCachedAuth() (*oebb.AuthInfo, error) {
	cache, err := xdg.SearchCacheFile("oebb-cli/auth.json")
	if err != nil {
		return authAndCache("oebb-cli/auth.json")
	}

	bytes, err := ioutil.ReadFile(cache)

	var newAuth oebb.AuthInfo
	err = json.Unmarshal(bytes, &newAuth)

	// check token expiration date and automatically renew if expired
	stat, err := os.Stat(cache)
	if err != nil {
		return nil, err
	}

	modTime := stat.ModTime()
	if modTime.Add(time.Duration(newAuth.ExpiresIn) * time.Second).Before(time.Now()) {
		return authAndCache("oebb-cli/auth.json")
	}

	return &newAuth, err
}

func handleTimeoutError(e error, auth *oebb.AuthInfo) bool {
	switch e.(type) {
	case *oebb.SessionTimeoutError:
		a, err := authAndCache("oebb-cli/auth.json")
		if err != nil {
			panic(e)
		}

		*auth = *a

		return true
	}

	return false
}

func stopSpinnerOnCtrlC(s *spinner.Spinner) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		s.Stop()
		os.Exit(1)
	}()
}

var searchCmd = &cobra.Command{
	Use:   "search [from] [to]",
	Short: "Search connections",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New([]string{"|", "/", "-", "\\"}, 50*time.Millisecond, spinner.WithHiddenCursor(true))
		s.Prefix = "Searching for connections "
		s.Writer = os.Stderr

		stopSpinnerOnCtrlC(s)

		s.Start()

		numResults, err := cmd.Flags().GetInt("results")
		if err != nil {
			s.Stop()
			panic(err)
		}

		depTimeStr, err := cmd.Flags().GetString("time")
		if err != nil {
			s.Stop()
			panic(err)
		}

		pAuth, err := maybeCachedAuth()
		if err != nil {
			s.Stop()
			panic(err)
		}

		auth := *pAuth

		from := args[0]
		to := args[1]

		fromStation, err := oebb.GetStations(from, auth)
		if err != nil {
			if handleTimeoutError(err, &auth) == true {
				fromStation, err = oebb.GetStations(from, auth)
				if err != nil {
					s.Stop()
					panic(err)
				}
			}
		}

		toStation, err := oebb.GetStations(to, auth)
		if err != nil {
			if handleTimeoutError(err, &auth) {
				toStation, err = oebb.GetStations(to, auth)
				if err != nil {
					s.Stop()
					panic(err)
				}
			}
		}

		var depTime time.Time
		if depTimeStr == "" {
			depTime = time.Now()
		} else {
			depTime, err = time.Parse("15:04", depTimeStr)
			if err != nil {
				s.Stop()
				panic(err)
			}

			now := time.Now()
			depTime = depTime.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)
		}

		connections, err := oebb.GetConnections(fromStation[0], toStation[0], auth, depTime, numResults)

		s.Stop()
		if err != nil {
			panic(err)
		}

		if len(connections) < 1 {
			errFrom := rgbterm.InterpretStr("{#cc6666}" + fromStation[0].Name + "{}")
			errTo := rgbterm.InterpretStr("{#cc6666}" + toStation[0].Name + "{}")
			fmt.Printf("No connections found from %s to %s\n", errFrom, errTo)
		}

		for _, conn := range connections {
			displayConnection(conn)
		}
	},
}
