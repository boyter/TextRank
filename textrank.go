package textrank

import (
	"github.com/DavidBelicza/TextRank/rank"
	"github.com/DavidBelicza/TextRank/parse"
	"github.com/DavidBelicza/TextRank/convert"
)

var provider = make(map[int]*rank.Rank)

func AddText(text string, lang string, id int) {
	var ranks *rank.Rank

	if savedTextRank, ok := provider[id]; ok {
		ranks = savedTextRank
	} else {
		ranks = rank.NewRank()
		provider[id] = ranks
	}

	language := convert.NewLanguage()
	language.SetDefaultLanguage(lang)

	parsedText := parse.TokenizeText(text)

	for _, sentence := range parsedText.GetSentences() {
		convert.TextToRank(sentence, language, provider[id])
	}
}

func Calculate(
	id int,
	algorithm func(int, int, int, int, int, int, int, int, int) float32,
) {
	rank.Calculate(provider[id], algorithm)
}

func GetBasicAlgorithm() func(int, int, int, int, int, int, int, int, int) float32 {
	return func(
		word1ID int,
		word2ID int,
		relationQty int,
		relationMin int,
		relationMax int,
		word1Qty int,
		word2Qty int,
		word1Min int,
		word1Max int,
	) float32 {

		if word1ID != 0 && word2ID != 0 {
			return (float32(relationQty) - float32(relationMin)) / (float32(relationMax) - float32(relationMin))
		} else {
			return (float32(word1Qty) - float32(word1Min)) / (float32(word1Max) - float32(word1Min))
		}
	}
}

func GetAlternateAlgorithm() func(int, int, int, int, int, int, int, int, int) float32 {
	return func(
		word1ID int,
		word2ID int,
		relationQty int,
		relationMin int,
		relationMax int,
		word1Qty int,
		word2Qty int,
		word1Min int,
		word1Max int,
	) float32 {

		if word1ID != 0 && word2ID != 0 {
			min := float32(relationMin + word1Min)
			max := float32(relationMax + word1Max)
			qty := float32(relationQty + word1Qty)

			return (qty - min) / (max - min)
		} else {
			return (float32(word1Qty) - float32(word1Min)) / (float32(word1Max) - float32(word1Min))
		}
	}
}

func GetPhrases(id int) []rank.Phrase {
	return rank.GetPhrases(provider[id])
}

func GetSingleWords(id int) []rank.SingleWord {
	return rank.GetSingleWords(provider[id])
}

func GetSentencesByRelationScore(id int, limit int) *[]rank.Sentence {
	return rank.GetSentences(provider[id], rank.ByRelation, limit)
}

func GetSentencesByWordQtyScore(id int, limit int) *[]rank.Sentence {
	return rank.GetSentences(provider[id], rank.ByQty, limit)
}

func GetSentencesByPhrases(id int, phrases []string) *[]rank.Sentence {
	return rank.GetSentencesByPhrases(provider[id], phrases)
}

func GetSentencesFrom(id int, sentenceID int, limit int) *[]rank.Sentence {
	return rank.GetSentencesFrom(provider[id], sentenceID, limit)
}