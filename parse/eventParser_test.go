package parse

import (
	"reflect"
	"testing"
)

func Test_eventTableParser_parse(t *testing.T) {
	// Arrange
	parser := &eventTableParser{}

	// Act
	results, err := parser.parse(eventResults)

	// Assert
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	_ = results
}

func Test_getWordIndexes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"1", args{s:"x  y"}, []int{0, 3}},
		{"2", args{s:"19    Heitmann, Ståle               Fossum IF                 0:58:27 +  14:41      143.32"}, []int{0, 6, 37, 63, 74, 85}},
		{"3", args{s:"28    Vister, Hanne Maria                                     1:05:07 +  21:21      139.71"}, []int{0, 6, 63, 74, 85}},
		{"4", args{s:"21    Fløystad, Jostein Bø          Privat                    1:00:52 +  17:06      142.01"}, []int{0, 6, 37, 63, 74, 85}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWordIndexes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWordIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getWords(t *testing.T) {
	s := "19    Heitmann, Ståle               Fossum IF                 0:58:27 +  14:41      143.32"

	words := getWords(s)

	if len(words) != 6 {
		t.Errorf("expected 6 elements, got %d", len(words))
	}

	s = "28    Vister, Hanne Maria                                     1:05:07 +  21:21      139.71"
	words = getWords(s)

	if len(words) != 5 {
		t.Errorf("expected 5 elements, got %d", len(words))
	}
}

/////////////////////////////////
/////////////////////////////////

const (
	eventResults = `
 1    Kostylev, Jegor               Mora                      0:43:46 +  00:00      151.28
 2    Vogelsang, Christian          Nydalens SK               0:44:08 +  00:22      151.08
 3    Prydz, Espen Beer             Heming Orientering        0:44:49 +  01:03      150.71
 4    Blom-hagen, Torbjørn          Fossum IF                 0:47:57 +  04:11      149.01
 5    Nipen, Thomas                 IL Tyrving                0:49:58 +  06:12      147.92
 6    Nipen, Mathias                Bækkelagets SK            0:50:17 +  06:31      147.75
 7    Schlaupitz, Holger            IL GeoForm                0:51:15 +  07:29      147.22
 8    Helland, Knut Edvard          Østmarka OK               0:52:00 +  08:14      146.82
 9    Låg, Steinar                  VBIL                      0:53:12 +  09:26      146.17
10    Nikolaisen, Per-Ivar          Oppsal Orientering        0:53:57 +  10:11      145.76
11    Heir, Morten                  Fossum IF                 0:54:20 +  10:34      145.55
12    Hjermstad, Ragnhild           Fossum IF                 0:54:45 +  10:59      145.33
13    Systad, Rolv Anders           Lyn Ski                   0:54:55 +  11:09      145.24
14    Hæstad, Nils                  Fossum IF                 0:55:41 +  11:55      144.82
15    Grønneberg, Skage             Heming Orientering        0:55:54 +  12:08      144.70
16    Winsvold, Gløer O             Sentrum OK                0:56:22 +  12:36      144.45
17    Frøyd, Jørgen                 Larvik OK                 0:56:38 +  12:52      144.31
18    Östgren, Björn Mo             IL GeoForm                0:56:51 +  13:05      144.19
19    Heitmann, Ståle               Fossum IF                 0:58:27 +  14:41      143.32
20    Kildahl, Øystein              Østmarka OK               0:59:47 +  16:01      142.60
21    Fløystad, Jostein Bø          Privat                    1:00:52 +  17:06      142.01
22    Oram, Louise                  Bækkelagets SK            1:01:07 +  17:21      141.88
23    Iwe, Harald                   IL GeoForm                1:01:40 +  17:54      141.58
24    Kløvrud, Jens Olav            Lillomarka OL             1:01:56 +  18:10      141.43
25    Zeiner-Gundersen, Richard     Lierbygda OL              1:02:32 +  18:46      141.11
26    Olausson, Mikael              Oslostudentenes IK        1:02:38 +  18:52      141.05
27    Teigland, Rune                Østmarka OK               1:03:16 +  19:30      140.71
28    Vister, Hanne Maria                                     1:05:07 +  21:21      139.71
29    Linnebo, Øyvind Johannessen   Asker SK                  1:05:22 +  21:36      139.57
30    Grinde, Bjørn                 IL GeoForm                1:05:26 +  21:40      139.54
31    Pedersen, Atle                Fossum IF                 1:05:45 +  21:59      139.36
32    Reusch, Christian             Heming Orientering        1:06:35 +  22:49      138.91
33    Trier, Øivind Thorvald Due    Oslostudentenes IK        1:06:59 +  23:13      138.70
34    Onsager, Knut                 IL GeoForm                1:08:50 +  25:04      137.69
35    Risvoll, Ketil                Telenor BIL               1:09:24 +  25:38      137.39
36    Refsland, Ivar                IL Tyrving                1:09:59 +  26:13      137.07
37    Tallaksen, Thor Christian     IL GeoForm                1:10:09 +  26:23      136.98
38    Hjermstad, Lars               Fossum IF                 1:10:20 +  26:34      136.88
39    Fixdal, Trude                 IL Koll                   1:11:12 +  27:26      136.41
40    Mygland, Johan                IL GeoForm                1:11:45 +  27:59      136.11
41    Farkas, Lorant                Østmarka OK               1:12:11 +  28:25      135.88
42    Eriksen, Are                  Oslostudentenes IK        1:12:36 +  28:50      135.65
43    Øvergaard, Tormod             Vestre Akers SK           1:19:10 +  35:24      132.09
44    Kippernes, Frank Åge          FFI BIL                   1:19:35 +  35:49      131.87
45    Johansen, Frode               Equinor BIL               1:22:41 +  38:55      130.19
46    Dahlsrud, Per Ole             Nydalens SK               1:24:16 +  40:30      129.33
47    Christensen, Petter           VBIL                      1:31:35 +  47:49      125.36
DSQ   Sætran, Bjørn Idar            IF Trauma                 0:41:26  (-1 poster)  95.00
DSQ   Hjelm, Morten                 VBIL                      0:59:12  (-1 poster)  95.00
`
)

