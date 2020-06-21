package repository

import (
	"fmt"

	"github.com/filhodanuvem/polyglot/language"
)

type Statistics struct {
	counters map[string]int
	langs    []string
}

// Implementing Sort interface
func (s Statistics) Len() int {
	return len(s.counters)
}

func (s Statistics) Less(i, j int) bool {
	return s.counters[s.langs[i]] > s.counters[s.langs[j]]
}

func (s Statistics) Swap(i, j int) {
	s.counters[s.langs[i]], s.counters[s.langs[j]] = s.counters[s.langs[j]], s.counters[s.langs[i]]
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%+v", s.counters)
}

func (s *Statistics) FirstLanguage() string {
	lang := ""
	for i := range s.counters {
		if s.counters[i] >= s.counters[lang] || lang == "" {
			lang = i
		}
	}

	if lang == "" {
		return ""
	}

	return lang
}

func (s *Statistics) Merge(stats *Statistics) {
	if s.counters == nil {
		s.counters = make(map[string]int)
	}
	for lang := range stats.counters {
		if _, ok := s.counters[lang]; !ok {
			s.counters[lang] = stats.counters[lang]
			stats.langs = append(stats.langs, lang)
			continue
		}

		s.counters[lang] += stats.counters[lang]
	}
}

func GetStatistics(files []string) (Statistics, error) {
	var stats Statistics
	stats.counters = make(map[string]int)
	for i := range files {
		lang, err := language.DetectByFile(files[i])
		if err != nil {
			return stats, err
		}
		if _, ok := stats.counters[lang]; !ok {
			stats.counters[lang] = 0
			stats.langs = append(stats.langs, lang)
		}
		stats.counters[lang]++
	}

	return stats, nil
}
