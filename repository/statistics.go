package repository

import "github.com/filhodanuvem/polyglot/language"

type Statistics struct {
	counters map[string]int
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
	for lang := range stats.counters {
		if _, ok := s.counters[lang]; !ok {
			s.counters = make(map[string]int)
			s.counters[lang] = stats.counters[lang]
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
		}
		stats.counters[lang]++
	}

	return stats, nil
}
