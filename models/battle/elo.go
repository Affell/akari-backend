package battle

import "math"

func ExpectedScoreWithFactors(elo1, elo2 int) float64 {
	return 1 / (1 + math.Pow(10, float64(elo2-elo1)/float64(D)))
}

func RatingDeltaWithFactors(elo1, elo2 int, score float64) int {
	return int(float64(K) * (score - ExpectedScoreWithFactors(elo1, elo2)))
}

func ComputeResult(elo1, elo2 int, score float64) (newElo1, newElo2 int) {
	delta := RatingDeltaWithFactors(elo1, elo2, score)
	newElo1 = elo1 + delta
	newElo2 = elo2 - delta
	return
}
