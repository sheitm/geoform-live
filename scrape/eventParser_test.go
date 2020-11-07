package scrape

import (
	"fmt"
	"github.com/sheitm/ofever/contracts"
	"reflect"
	"testing"
	"time"
)

func Test_eventTableParser_parse(t *testing.T) {
	// Arrange
	parser := &eventTableParser{}

	// Act
	results, err := parser.parse(standardHeaders, eventResults)

	// Assert
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if len(results) != 49 {
		t.Errorf("expected 49 result lines, got %d", len(results))
	}
	var result *contracts.Result
	for _, r := range results {
		if r.Athlete == "Heir, Morten" {
			result = r
			break
		}
	}

	if result.Placement != 11 {
		t.Errorf("unexpected placement, got %d", result.Placement)
	}
	if result.Club != "Fossum IF" {
		t.Errorf("unexpected club, got %s", result.Club)
	}
	if result.Points != 145.55 {
		t.Errorf("unexpected points, got %f", result.Points)
	}
	if fmt.Sprintf("%v", result.ElapsedTime) != "54m20s" {
		t.Errorf("unexpected elapsed time, got %v", result.ElapsedTime)
	}
	if result.Disqualified {
		t.Errorf("did not expect qualified")
	}

	// DSQ   Sætran, Bjørn Idar            IF Trauma                 0:41:26  (-1 poster)  95.00
	for _, r := range results {
		if r.Athlete == "Sætran, Bjørn Idar" {
			result = r
			break
		}
	}

	if !result.Disqualified {
		t.Errorf("expected disqualified")
	}
	if result.MissingControls != 1 {
		t.Errorf("expected 1 missing control, got %d", result.MissingControls)
	}
}

func Test_eventTableParser_parse_nonstandard2019(t *testing.T) {
	// Arrange
	parser := &eventTableParser{}

	// Act
	results, err := parser.parse(standardHeaders, eventResults2019_5)

	// Assert
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if len(results) != 25 {
		t.Errorf("expected 25 results, got %d", len(results))
	}

	var result *contracts.Result
	for _, r := range results {
		if r.Placement == 1 {
			result = r
			break
		}
	}

	if result.Athlete != "Grønneberg, Skage" {
		t.Errorf("unexpected athlete name, got %s", result.Athlete)
	}
	if result.Club != "Heming Orientering" {
		t.Errorf("unexpected club, got %s", result.Club)
	}
	if fmt.Sprintf("%v", result.ElapsedTime) != "57m1s" {
		t.Errorf("unexpected elapsed time, got %v", result.ElapsedTime)
	}
	if result.Points != 151.04 {
		t.Errorf("unexpected points, got %f", result.Points)
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

func Test_getResultFromLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *contracts.Result
		wantErr bool
	}{
		{
			name:    "Kostylev",
			args:    args{line:"1    Kostylev, Jegor               Mora                      0:43:46 +  00:00      151.28"},
			want:    &contracts.Result{
				Placement:       1,
				Disqualified:    false,
				Athlete:         "Kostylev, Jegor",
				Club:            "Mora",
				ElapsedTime:     43*time.Minute + 46*time.Second,
				Points:          151.28,
				MissingControls: 0,
			},
			wantErr: false,
		},
		{
			name: "Blom-hagen",
			args: args{line: "4    Blom-hagen, Torbjørn          Fossum IF                 0:47:57 +  04:11      149.01"},
			want:    &contracts.Result{
				Placement:       4,
				Disqualified:    false,
				Athlete:         "Blom-hagen, Torbjørn",
				Club:            "Fossum IF",
				ElapsedTime:     47*time.Minute + 57*time.Second,
				Points:          149.01,
				MissingControls: 0,
			},
			wantErr: false,
		},
		{
			name: "Vister",
			args: args{line: "28    Vister, Hanne Maria                                     1:05:07 +  21:21      139.71"},
			want:    &contracts.Result{
				Placement:       28,
				Disqualified:    false,
				Athlete:         "Vister, Hanne Maria",
				Club:            "",
				ElapsedTime:     1*time.Hour + 5*time.Minute + 7*time.Second,
				Points:          139.71,
				MissingControls: 0,
			},
			wantErr: false,
		},
		{
			name: "Hjelm",
			args: args{line: "DSQ   Hjelm, Morten                 VBIL                      0:59:12  (-1 poster)  95.00"},
			want:    &contracts.Result{
				Placement:       0,
				Disqualified:    true,
				Athlete:         "Hjelm, Morten",
				Club:            "VBIL",
				ElapsedTime:     59*time.Minute + 12*time.Second,
				Points:          95.0,
				MissingControls: 1,
			},
			wantErr: false,
		},
		{
			name: "Feiring",
			args: args{line: "      Feiring, Hege                 IL Tyrving                DELTATT               50.00"},
			want:    &contracts.Result{
				Placement:       0,
				Disqualified:    true,
				Athlete:         "Feiring, Hege",
				Club:            "IL Tyrving",
				ElapsedTime:     0,
				Points:          50.0,
				MissingControls: 0,
			},
			wantErr: false,
		},
		{
			name: "Hobæk",
			args: args{line: "65    Hobæk, Thor                   Gassecure BIL             2:03:41 +1:19:33      108.99"},
			want:    &contracts.Result{
				Placement:       65,
				Disqualified:    false,
				Athlete:         "Hobæk, Thor",
				Club:            "Gassecure BIL",
				ElapsedTime:     2*time.Hour + 3 * time.Minute + 41 * time.Second,
				Points:          108.99,
				MissingControls: 0,
			},
			wantErr: false,
		},
		{
			name: "Blom",
			args: args{line: "26    Blom, Richard                 FBI - Forskningsbedriftenes bedriftsidrettslag0:48:07 +  15:10      128.47"},
			want:    &contracts.Result{
				Placement:       26,
				Disqualified:    false,
				Athlete:         "Blom, Richard",
				Club:            "FBI - Forskningsbedriftenes bedriftsidrettslag",
				ElapsedTime:     48 * time.Minute + 7 * time.Second,
				Points:          128.47,
				MissingControls: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getResultFromLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("getResultFromLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResultFromLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

/////////////////////////////////
/////////////////////////////////

const (
	standardHeaders = `Plass Navn                          Klubb                     Tid                   Poeng`
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
	eventResults2019_5 = `
 1    Grønneberg, Skage             Heming Orientering    0:57:01 +  00:00      151,04
 2    Nipen, Thomas                 Bekkelaget            0:57:21 +  00:20      150,90
 3    Nipen, Jørgen Mathias         Bækkelaget            1:00:36 +  03:35      149,54
 4    Nygård, Svein                 Norges Bank BIL       1:00:52 +  03:51      149,43
 5    Svedberg, Johan               Heming Orientering    1:01:40 +  04:39      149,09
 6    Fremming, Nils Petter         Heming Orientering    1:03:05 +  06:04      148,49
 7    Røberg, Henning               Nittedal OL           1:03:53 +  06:52      148,16
 8    Norman, Niklas                IL GeoForm            1:04:25 +  07:24      147,93
 9    Systad, Rolv Anders           Lyn Ski               1:05:33 +  08:32      147,46
10    Bårtveit, Knut                Bø OL                 1:06:54 +  09:53      146,89
11    Kjølseth, Tore                Lundin                1:07:10 +  10:09      146,78
12    Helland, Knut                 Østmarka OK           1:07:17 +  10:16      146,73
13    Heier, Morten                 Fossum IF             1:09:49 +  12:48      145,67
14    Heir, Marius Borge            NTNUI                 1:11:05 +  14:04      145,13
15    Zeiner-Gundersen, Richard     Aker Brygge Orientering1:14:15 +  17:14      143,80
16    Egge, Guttorm                 ILGeoForm             1:17:43 +  20:42      142,35
17    Heir, Stig                    Asker SK              1:18:01 +  21:00      142,22
18    Roti, Jarle                   Fossum IF             1:20:44 +  23:43      141,08
19    Eriksen, Are                  OSI                   1:23:09 +  26:08      140,06
20    Heitmann, Ståle               Fossum IF             1:26:46 +  29:45      138,54
21    Messel, Espen                 IL Koll               1:29:47 +  32:46      137,28
22    Lahlum, Jon                   IL GeoForm            1:31:28 +  34:27      136,57
23    Syversten, Bjørne             Privat                1:39:00 +  41:59      133,40
24    Grandum, Øyvind               IL GeoForm            1:49:08 +  52:07      129,15
DSQ   Vogelsang, Christian          Nydalens SK           0:29:36  (-9 poster)  71,88`
)

