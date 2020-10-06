package repository

import "testing"

func est_MergeOneStatsShouldCountCorrectly(t *testing.T) {
	stats := Statistics{
		counters: []Counter{Counter{Lang: "PHP", Counter: 2}},
		langs:    map[string]int{"PHP": 0},
	}

	resultStats := Statistics{}
	resultStats.Merge(&stats)
	langs := resultStats.FirstLanguages(1)

	if langs[0].Lang != "PHP" && langs[0].Counter != 1 {
		t.Errorf("Expected PHP:1 found %+v", langs)
	}
}

func Test_MergeTwoStatsShouldCountCorrectly(t *testing.T) {
	stats1 := Statistics{
		counters: []Counter{Counter{Lang: "PHP", Counter: 2}},
		langs:    map[string]int{"PHP": 0},
	}

	stats2 := Statistics{
		counters: []Counter{Counter{Lang: "Go", Counter: 3}, Counter{Lang: "PHP", Counter: 5}},
		langs:    map[string]int{"Go": 0, "PHP": 1},
	}

	resultStats := Statistics{}
	resultStats.Merge(&stats1)
	resultStats.Merge(&stats2)
	langs := resultStats.FirstLanguages(2)

	if langs[0].Lang != "PHP" && langs[0].Counter != 6 {
		t.Errorf("Expected PHP:6 found %+v", langs)
	}

	if langs[1].Lang != "Go" && langs[0].Counter != 3 {
		t.Errorf("Expected PHP:6 found %+v", langs)
	}
}
