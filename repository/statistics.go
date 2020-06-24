package repository

import (
	"fmt"
	"sort"

	"github.com/filhodanuvem/polyglot/language"
)

type Statistics struct {
	counters   []counter
	langs      map[string]int
	reposCount int
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
	s.counters[i], s.counters[j] = s.counters[j], s.counters[i]
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%+v", s.counters)
}

func (s *Statistics) Length() int {
	return s.reposCount
}

func (s *Statistics) FirstLanguages(length int) []string {
	sort.Sort(s)
	langs := make([]string, length)
	j := 0
	for i := range s.counters {
		langs = append(langs, s.counters[i].lang)
		j++
		if j == length {
			break
		}
	}

	return langs
}

func (s *Statistics) Merge(stats *Statistics) {
	stats.reposCount++
	if s.langs == nil {
		s.langs = make(map[string]int)
	}
	for i := range stats.counters {
		lang := stats.counters[i].lang
		if _, ok := s.langs[lang]; !ok {
			s.langs[lang] = len(s.counters)
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
		if lang == "" {
			continue
		}
		if _, ok := stats.langs[lang]; !ok {
			stats.langs[lang] = len(stats.counters)
			c := counter{
				lang:    lang,
				counter: 0,
			}
			stats.counters = append(stats.counters, c)
			continue
		}
		stats.counters[stats.langs[lang]].counter++
	}

	stats.reposCount = 1
	return stats, nil
}
