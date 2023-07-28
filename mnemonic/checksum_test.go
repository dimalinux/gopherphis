package mnemonic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test dataset was generated using monero-wallet-rpc's
// restore_deterministic_wallet method with 24 words, then grabbing the 25-word
// seed phrase from the response. The seeds are all at index multiples of 65
// from each language and all generate the same key:
// bfdd3f0a82de3f0a45df3f0a08e03f0acbe03f0a8ee13f0a51e23f0a14e33f0a
//
//nolint:misspell,lll
var checkSumTests = []struct {
	lang      string
	mnemonic  []string
	checkSeed string
}{
	{
		lang:      "English",
		mnemonic:  []string{"always", "awesome", "boss", "cohesive", "diplomat", "eight", "fatal", "gained", "gypsy", "icing", "jeopardy", "lazy", "mammal", "narrate", "october", "pastry", "pruned", "ridges", "seeded", "soprano", "technical", "tumbling", "utmost", "website"},
		checkSeed: "fatal",
	},
	{
		lang:      "Russian",
		mnemonic:  []string{"арест", "благо", "вена", "выгодный", "гулять", "дрянь", "занимать", "икра", "клиент", "левый", "манера", "мысль", "нужный", "орбита", "пенсия", "пурга", "рубль", "синий", "стыд", "титул", "улыбка", "форма", "честь", "штора"}, //nolint:lll
		checkSeed: "нужный",
	},
	{
		lang:      "Spanish",
		mnemonic:  []string{"alcalde", "apuesta", "ave", "bondad", "camino", "choza", "corona", "dental", "edición", "eterno", "fijo", "ganso", "hacer", "imperio", "laguna", "llegar", "masa", "monarca", "nido", "olivo", "papel", "pichón", "príncipe", "rebote"}, //nolint:lll
		checkSeed: "dental",
	},
	{
		lang:      "Portuguese",
		mnemonic:  []string{"alpiste", "audivel", "biunivoco", "chuvoso", "damista", "dueto", "eolico", "faquir", "fucsia", "grotesco", "iansa", "isento", "laico", "luvas", "muarra", "nublar", "opio", "pezada", "rabujice", "rouxinol", "sirio", "tegumento", "uiste", "viavel"}, //nolint:lll
		checkSeed: "iansa",
	},
	{
		lang:      "German",
		mnemonic:  []string{"Alpental", "Apfel", "Baum", "Biologe", "Buslinie", "Dialekt", "Eingang", "Ernte", "Faulheit", "Fokus", "Gefieder", "Grünalge", "Hochform", "Junimond", "Kloster", "Leder", "Maßkrug", "Narbe", "Paket", "Quellsee", "Ruhe", "Songtext", "Torlinie", "Walhai"}, //nolint:lll
		checkSeed: "Junimond",
	},
	{
		lang:      "French",
		mnemonic:  []string{"ample", "azur", "braise", "certes", "connu", "difficile", "enfin", "fente", "fuir", "habiter", "ivre", "ligoter", "membre", "muret", "odeur", "partir", "pliage", "puiser", "rejeter", "rouler", "seuil", "subir", "toiser", "veau"}, //nolint:lll
		checkSeed: "ligoter",
	},
	{
		lang:      "Italian",
		mnemonic:  []string{"america", "assalire", "benzina", "cacciare", "cercare", "conforto", "derivare", "educare", "famiglia", "foro", "gioiello", "inferno", "lievito", "metodo", "notare", "panorama", "pista", "radice", "ritmo", "scrupolo", "sondare", "sussurro", "treccia", "venire"}, //nolint:lll
		checkSeed: "gioiello",
	},
	{
		lang:      "Dutch",
		mnemonic:  []string{"arganolie", "betichten", "buma", "dagprijs", "duwwerk", "ertussen", "follikel", "geseald", "hazebroek", "invitatie", "klotefilm", "laxeerpil", "lurven", "napels", "omroep", "pakzadel", "publiceer", "roon", "sjezen", "subregent", "topvrouw", "velgrem", "wadvogel", "wrijf"}, //nolint:lll
		checkSeed: "sjezen",
	},
	{
		lang:      "Japanese",
		mnemonic:  []string{"いけばな", "いらすと", "うらない", "おおよそ", "かいつう", "きくばり", "ぎんいろ", "けいれき", "げんそう", "こつこつ", "さうな", "しあさって", "しゃたい", "ずほう", "せんぞ", "そよかぜ", "たそがれ", "ちあい", "つなみ", "でんあつ", "ともる", "にまめ", "ねんぐ", "はたん"}, //nolint:lll
		checkSeed: "はたん",
	},
	{
		lang:      "Chinese (simplified)",
		mnemonic:  []string{"之", "样", "位", "南", "具", "复", "院", "含", "兵", "未", "盐", "息", "香", "亮", "弱", "纳", "猪", "您", "辩", "浪", "琴", "奶", "壳", "乳"},
		checkSeed: "琴",
	},
}

func TestGetChecksumWord(t *testing.T) {
	for _, tc := range checkSumTests {
		expected := tc.checkSeed
		checkSeed, wl, err := GetChecksumWord(tc.mnemonic)
		require.NoError(t, err)
		require.Equal(t, tc.lang, wl.EnglishName)
		require.Equal(t, expected, checkSeed)
	}
}

func TestWordList_HasWord(t *testing.T) {
	// Test 1st, middle and last
	for _, wl := range WordLists {
		for _, w := range wl.Entries {
			require.True(t, wl.HasWord(w))
		}
		require.False(t, wl.HasWord("x"))
	}
}
