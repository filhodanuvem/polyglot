package repository

import "testing"

func Test_MergeOneStatsShouldCountCorrectly(t *testing.T) {
	stats := Statistics{
		counters: []counter{counter{lang: "PHP", counter: 2}},
		langs:    map[string]int{"PHP": 0},
	}

	resultStats := Statistics{}
	resultStats.Merge(&stats)
	langs := resultStats.FirstLanguages(1)

	if langs[0].lang != "PHP" && langs[0].counter != 1 {
		t.Errorf("Expected PHP:1 found %+v", langs)
	}
}

func Test_MergeTwoStatsShouldCountCorrectly(t *testing.T) {
	stats1 := Statistics{
		counters: []counter{counter{lang: "PHP", counter: 2}},
		langs:    map[string]int{"PHP": 0},
	}

	stats2 := Statistics{
		counters: []counter{counter{lang: "Go", counter: 3}, counter{lang: "PHP", counter: 5}},
		langs:    map[string]int{"Go": 0, "PHP": 1},
	}

	resultStats := Statistics{}
	resultStats.Merge(&stats1)
	resultStats.Merge(&stats2)
	langs := resultStats.FirstLanguages(2)

	if langs[0].lang != "PHP" && langs[0].counter != 6 {
		t.Errorf("Expected PHP:6 found %+v", langs)
	}

	if langs[1].lang != "Go" && langs[0].counter != 3 {
		t.Errorf("Expected PHP:6 found %+v", langs)
	}
}
