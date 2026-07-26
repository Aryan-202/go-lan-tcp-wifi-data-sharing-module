package security

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// WordList contains a curated list of simple, memorable words for 4-word pairing passphrases
var wordList = []string{
	"amber", "anchor", "apple", "beacon", "breeze", "bridge", "canyon", "castle",
	"cedar", "cipher", "cobalt", "comet", "coral", "cosmos", "crest", "crystal",
	"delta", "desert", "dragon", "eagle", "echo", "emerald", "falcon", "forest",
	"galaxy", "glacier", "harbor", "haven", "horizon", "island", "jaguar", "jungle",
	"knight", "lagoon", "lantern", "legend", "lotus", "lunar", "magnet", "matrix",
	"meadow", "meteor", "mirage", "mountain", "nebula", "nexus", "oasis", "ocean",
	"orbit", "orchid", "osprey", "palace", "panther", "phoenix", "planet", "prism",
	"pulse", "pyramid", "quartz", "quazar", "radar", "radiant", "raven", "river",
	"rocket", "shadow", "shield", "silver", "solar", "spark", "spectrum", "sphere",
	"summit", "thunder", "timber", "titan", "topaz", "torrent", "valley", "vector",
	"velocity", "vortex", "whisper", "zenith", "zephyr", "zodiac", "zone", "zinc",
}

// GeneratePassphrase picks 4 cryptographically secure random words joined by hyphens
func GeneratePassphrase() (string, error) {
	words := make([]string, 4)
	numWords := big.NewInt(int64(len(wordList)))

	for i := 0; i < 4; i++ {
		idx, err := rand.Int(rand.Reader, numWords)
		if err != nil {
			return "", fmt.Errorf("failed to generate random index for passphrase: %w", err)
		}
		words[i] = wordList[idx.Int64()]
	}

	return strings.Join(words, "-"), nil
}

// VerifyPassphrase checks if user input matches expected passphrase (case-insensitive)
func VerifyPassphrase(input string, expected string) bool {
	cleanInput := strings.ToLower(strings.TrimSpace(input))
	cleanExpected := strings.ToLower(strings.TrimSpace(expected))
	return cleanInput == cleanExpected
}
