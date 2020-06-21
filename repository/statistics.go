package repository

import (
	"fmt"
	"sort"

	"github.com/filhodanuvem/polyglot/language"
)

type Statistics struct {
	counters []counter
	langs    map[string]int
}

type counter struct {
	lang    string
	counter int
}

// Implementing Sort interface
func (s *Statistics) Len() int {
	return len(s.counters)
}

func (s *Statistics) Less(i, j int) bool {
	return s.counters[i].counter > s.counters[j].counter
}

func (s *Statistics) Swap(i, j int) {
	s.counters[i], s.counters[j] = s.counters[i], s.counters[i]
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%+v", s.counters)
}

func (s *Statistics) FirstLanguage() string {
	sort.Sort(s)

	return s.counters[0].lang
}

func (s *Statistics) Merge(stats *Statistics) {
	for i := range stats.counters {
		lang := stats.counters[i].lang
		if _, ok := s.langs[lang]; !ok {
			s.counters = append(s.counters, stats.counters[i])
			continue
		}

		s.counters[stats.langs[lang]].counter += stats.counters[i].counter
	}
}

func GetStatistics(files []string) (Statistics, error) {
	var stats Statistics
	stats.langs = make(map[string]int)
	for i := range files {
		lang, err := language.DetectByFile(files[i])
		if err != nil {
			return stats, err
		}
		if _, ok := stats.langs[lang]; !ok {
			stats.langs[lang] = len(stats.counters)
			c := counter{
				lang:    lang,
				counter: 0,
			}
			stats.counters = append(stats.counters, c)
		}
		stats.counters[stats.langs[lang]].counter++
		// for j := range stats.counters {
		// 	if stats.counters[j].lang == lang {
		// 		stats.counters[j].counter++
		// 		break
		// 	}
		// }
	}

	return stats, nil
}
