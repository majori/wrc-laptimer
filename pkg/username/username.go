package username

import (
	"hash/fnv"
	"math/rand"
)

var (
	adjective = []string{
		"Ruosteinen", "Bensanhajuinen", "Dieselmäinen", "Kireä", "Taivaallinen",
		"Mutainen", "Pölyinen", "Rämisevä", "Legendaarinen", "Tahmanen",
		"Yliohjautuva", "Aliohjautuva", "Savuava", "Öljyinen", "Likainen",
		"Kiiltävä", "Ruosteeton", "Kolhiintunut", "Naarmuinen", "Äänekäs",
		"Tehokas", "Voimakas", "Nopea", "Kevyt", "Ketterä",
		"Raskas", "Kevyt", "Virtaviivainen", "Sporttinen", "Klassinen",
		"Vintage", "Moderni", "Retrohenkinen", "Ahdettu", "Virityskelpoinen",
	}
	epithet = []string{
		"Fullsend", "Lintta", "Kelirikko", "Eri Nopee", "Vaihdekeppi",
		"V8", "Dacia", "Nahkapuku", "Saabisti", "Terminal damage",
		"Sytytystulppa", "Vale", "Sierra", "Neliveto", "Erikoiskoe",
		"Täyskaasu", "Kaasujalka", "Turboahdettu", "Siirtymätaival",
		"Apukuski", "Maximum Attack",
	}
	name = []string{
		"Loeb",	"McRae", "Tommi", "Timo", "Marcus",
		"Kalle", "Grönholm", "Carlos", "Solberg", "Burns",
		"Kankkunen", "Vatanen", "Väinö", "Olavi", "Kalevi",
		"Paavo", "Jorma", "Stefa", "Arvo", "Reino",
		"Aino", "Helmi", "Martta", "Tyyne", "Hilja",
		"Ahti",	"Lempi", "Lalli", "Kyllikki", "Aili",
		"Saima", "Ester", "Hilma", "Bertta", "Lyyli",
		"Hilda", "Kerttu", "Elsa", "Sylvi", "Hillervo",
		"Eeva", "Kaarina", "Kirsti", "Bensalenkkari", "Bensalenkkar",
		"Teuvo", "Orvokki", "Jallu",
	}
)

func getFNVHash(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

func GenerateFromSeed(seed string) string {
	seededRand := rand.New(rand.NewSource(getFNVHash(seed)))
	randomisedEpithet := adjective[seededRand.Intn(len(adjective))]
	randomisedNickname := epithet[seededRand.Intn(len(epithet))]
	randomisedName := name[seededRand.Intn(len(name))]
	return randomisedEpithet + " " + randomisedNickname + " " + randomisedName
}
