package genkey

// mian file for genkey

import (
	. "github.com/roman-koshchei/genkey-api/structs"
)

var Data TextData
var isConfigLoaded bool = false

// Analysis
func Analyze(keys string, fingers string) Analysis {
	if !isConfigLoaded {
		loadConfig()
	}

	layout := loadLayout(keys, fingers)
	var analysis Analysis

	ftri := FastTrigrams(layout, 0)
	ftotal := float64(ftri.Total)

	analysis.LeftInwardRolls = 100 * float64(ftri.LeftInwardRolls) / ftotal
	analysis.LeftOutwardRolls = 100 * float64(ftri.LeftOutwardRolls) / ftotal

	analysis.RightInwardRolls = 100 * float64(ftri.RightInwardRolls) / ftotal
	analysis.RightOutwardRolls = 100 * float64(ftri.RightOutwardRolls) / ftotal

	analysis.Alternates = 100 * float64(ftri.Alternates) / ftotal
	analysis.Onehands = 100 * float64(ftri.Onehands) / ftotal
	analysis.Redirects = 100 * float64(ftri.Redirects) / ftotal

	var weighted []float64
	var unweighted []float64
	weighted = FingerSpeed(&layout, true)
	unweighted = FingerSpeed(&layout, false)

	var highestUnweightedFinger string
	var highestUnweighted float64
	var utotal float64

	var highestWeightedFinger string
	var highestWeighted float64
	var wtotal float64
	for i := 0; i < 8; i++ {
		utotal += unweighted[i]
		if unweighted[i] > highestUnweighted {
			highestUnweighted = unweighted[i]
			highestUnweightedFinger = FingerNames[i]
		}

		wtotal += weighted[i]
		if weighted[i] > highestWeighted {
			highestWeighted = weighted[i]
			highestWeightedFinger = FingerNames[i]
		}
	}

	analysis.WeightedFingerSpeed = weighted
	analysis.UnweightedFingerSpeed = unweighted

	analysis.WeightedHighestSpeed.Value = highestWeighted
	analysis.WeightedHighestSpeed.Finger = highestWeightedFinger
	analysis.UnweightedHighestSpeed.Value = highestUnweighted
	analysis.UnweightedHighestSpeed.Finger = highestUnweightedFinger

	left, right := IndexUsage(layout)
	analysis.LeftIndexUsage = left
	analysis.RightIndexUsage = right

	sfb := SFBs(layout, false)
	sfbs := ListSFBs(layout, false)
	analysis.Sfbs = 100 * sfb / layout.Total
	analysis.Dsfbs = 100 * SFBs(layout, true) / layout.Total
	lsb := float64(LSBs(layout))
	analysis.Lsbs = 100 * lsb / layout.Total

	SortFreqList(sfbs)
	sfbs = sfbs[0:10]
	// into percents
	for i := 0; i < len(sfbs); i++ {
		sfbs[i].Count = 100 * float64(sfbs[i].Count) / float64(Data.TotalBigrams)
	}
	analysis.TopSfbs = sfbs

	worstBigrams := ListWorstBigrams(layout)
	SortFreqList(worstBigrams)
	worstBigrams = worstBigrams[0:10]
	for i := 0; i < len(sfbs); i++ {
		worstBigrams[i].Count = 100 * float64(worstBigrams[i].Count) / float64(Data.TotalBigrams)
	}
	analysis.WorstBigrams = worstBigrams

	analysis.Score = Score(layout)

	//fmt.Println("Top SFBs:")
	//PrintFreqList(sfbs, 8, true)

	return analysis // tmp
}

func loadConfig() {

	Data = loadData() // from text.go
	loadWeights()     // from config.go

	isConfigLoaded = true
	// old loaded layouts
	// we should return layout given to us
	// Layouts = make(map[string]Layout)
	// LoadLayoutDir()

	// was
	// checkLayoutProvided(args)
	// PrintAnalysis(getLayout(args[1])) // args[1] = name of layout
}
