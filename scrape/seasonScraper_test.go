package scrape

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_startSeason_statusCode500(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
	}))
	defer server.Close()

	url := server.URL + "/rankinglop"
	resultChan := make(chan *SeasonFetch)

	// Act
	startSeason(url, 2020, resultChan, server.Client())

	fetch := <- resultChan

	// Assert
	if fetch.Error == nil {
		t.Error("expected err, got none")
	}
	if fetch.Error.Error() != "500 Internal Server Error" {
		t.Errorf("unexpected error message, got %s", fetch.Error.Error())
	}
}

func Test_startSeason_success(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		p := req.URL.Path
		//fmt.Println(p)
		if  p == "/rankinglop" {
			rw.Write([]byte(htmlSeasonResponse))
			return
		}
		if p == "/rankinglop/res2020-06-17.html" {
			rw.Write([]byte(htmlEventResponse))
			return
		}
		rw.WriteHeader(404)
	}))
	defer server.Close()

	url := server.URL + "/rankinglop"
	resultChan := make(chan *SeasonFetch)

	// Act
	startSeason(url, 2020, resultChan, server.Client())

	fetch := <- resultChan

	// Assert
	if fetch.Error != nil {
		t.Errorf("unexpected error, got %v", fetch.Error)
	}
	if len(fetch.Results) != 24 {
		t.Errorf("unexpected number of results, got %d", len(fetch.Results))
	}

	if fetch.Year != 2020 {
		t.Errorf("unexpected year, got %d", fetch.Year)
	}

	okCount := 0
	failedCount := 0
	for _, result := range fetch.Results {
		if result.Error != nil {
			failedCount++
			continue
		}
		okCount++
	}

	if okCount != 1 {
		t.Errorf("expected 1 ok, got %d", okCount)
	}
	if failedCount != 23 {
		t.Errorf("expected 23 failed, got %d", okCount)
	}
}

//func Test_startSeason_4real(t *testing.T) {
//	// Arrange
//	url := "https://ilgeoform.no/rankinglop/"
//	resultChan := make(chan *SeasonFetch)
//
//	// Act
//	startSeason(url, resultChan, &http.Client{})
//
//	fetch := <- resultChan
//
//	// Assert
//	if fetch.Error != nil {
//		t.Errorf("unexpected error, got %v", fetch.Error)
//	}
//	if len(fetch.Results) != 24 {
//		t.Errorf("unexpected number of results, got %d", len(fetch.Results))
//	}
//
//	okCount := 0
//	failedCount := 0
//	for _, result := range fetch.Results {
//		if result.Error != nil {
//			failedCount++
//			continue
//		}
//		okCount++
//	}
//
//	if okCount != 1 {
//		t.Errorf("expected 1 ok, got %d", okCount)
//	}
//	if failedCount != 23 {
//		t.Errorf("expected 23 failed, got %d", okCount)
//	}
//}

const (
	htmlSeasonResponse = `<html lang="nb"><head><meta charset="utf-8">
	<title>Rankingløp terminliste og resultater</title>
</head>
<body data-gramm="true" data-gramm_editor="true" data-gramm_id="5a119b28-943d-a585-8c39-d5039c2e6051">
<h1 style="color: rgb(255, 0, 0);"><span style="color: rgb(255, 0,
        0);">OSI/GeoForm Rankingløp 2020</span></h1>

<p><font color="#ff6600">Oppdatert 25.10.2020</font></p>

<p><a href="rank-2020.xlsx"><g class="gr_ gr_41 gr-alert gr_spell
          gr_inline_cards gr_run_anim ContextualSpelling ins-del
          multiReplace" data-gr-id="41" id="41">Totalranking</g> 2020</a>&nbsp; (eldre termin- og resultatlister ligger nederst på siden)</p>

<p><b>Onsdager start: kl&nbsp;16:30-18:30. (Fra uke 38: Lørdager kl. 12:00 - 14:00)&nbsp;<span style="color: red;">Målgang innen kl 19:30 (lørdagene 15:30)</span></b></p>

<p>Løypelengdene er: Lang, ca 6 km, Mellom ca 4 km og Kort, ca 2,5 km, alle på A-nivå</p>

<p>Elektronisk tidtaking benyttes på alle løp. Enten <strong>Emit </strong>eller <b>EmiTag - touchfree</b>. Touch-free brikker fås kjøpt hos <a href="https://www.idrettsbutikken.no/andre/2380061027/emitag-brikke">Idrettsbutikken.no</a> eller <a href="http://emit.no/product/ny-emitag-brikke-23">emit.no</a>. Arrangøren har vanligvis noen brikker til utleie (kr. 20).</p>

<p>Forhåndspåmelding i <a href="https://eventor.orientering.no">Eventor</a>. Resultatene finnes her og på <a href="https://eventor.orientering.no">Eventor</a>. Løyper og ruter på <a href="https://www.livelox.com/">Livelox</a>.</p>

<p><b>Startkontingent: 50 kr</b><br>
Medlemmer av OSI Orientering (Årets gruppekontingent for OSI Orientering skal være betalt), IL <span class="SpellE">GeoForm</span> og samarbeidende klubber, arrangører, samt barn/ungdom under 16 år:&nbsp; kr 30,-</p>

<p><strong>NB: Sjekk alltid innbydelsen til løpene. Det kan forekomme avvik fra ovenstående.</strong></p>

<p><strong>Poengberegning og premiering</strong></p>

<p>Rankingpoengene tar utgangspunkt gjennomsnittstiden til&nbsp;de fem beste i hver løype. Denne tiden gir hhv 150, 135 og 120 poeng i de tre løypene. Løperne får høyere eller lavere poengsum basert på differansen. Det gis også poeng til arrangører. Ved sesongslutt teller halvparten av arrangerte løp, rundet oppover. Premier til de tre beste menn og kvinner i kort og mellom og de tre beste menn og kvinner totalt, noe som i&nbsp;praksis betyr langløypa. Dessuten premieres høyest oppnådd poengsum blant&nbsp;menn og kvinner og den såkalte utholdenhetsprisen som går til løperen som har vært lengst i skogen sammenlagt i fullførte og godkjente Rankingløp.</p>

<table border="1" cellpadding="2" cellspacing="0" style="text-align: left; background-color: rgb(255, 255,
      204);" width="887">
	<tbody>
		<tr>
			<td style="vertical-align: top;" valign="top"><span style="font-weight: bold;">Nr</span></td>
			<td style="vertical-align: top;" valign="top"><span style="font-weight: bold;">Dato</span></td>
			<td style="vertical-align: top;" valign="top"><span style="font-weight: bold;">Dag</span></td>
			<td style="vertical-align: top; font-weight: bold;" valign="top">Sted/innbydelse</td>
			<td style="vertical-align: top; font-weight: bold;" valign="top">Løpsområde</td>
			<td style="vertical-align: top; font-weight: bold;" valign="top">Arrangør</td>
			<td style="vertical-align: top; font-weight: bold;" valign="top">Ansvarlig</td>
			<td colspan="1" rowspan="1" style="vertical-align: top; width:
            10px; font-weight: bold;" valign="top">Resultat</td>
			<td style="vertical-align: top;" valign="top"><span style="font-weight: bold;">LL</span></td>
			<td valign="top"><b>Eventor</b></td>
		</tr>
		<tr>
			<td style="vertical-align: top;" valign="top">1</td>
			<td style="vertical-align: top;" valign="top">03.06</td>
			<td style="vertical-align: top;" valign="top">ons</td>
			<td style="vertical-align: top;" valign="top"><a href="innbyd2020-06-03.pdf">Skytterkollen</a></td>
			<td style="vertical-align: top;" valign="top">Skytterkollen</td>
			<td style="vertical-align: top;" valign="top">Fossum IF</td>
			<td colspan="1" rowspan="3" style="vertical-align: top;" valign="top"><strong>Sørkedalskarusellen 2020</strong><br>
			Sigbjørn Modalsli og Richard Zeiner Gundersen<br>
			&nbsp;</td>
			<td style="vertical-align: top;" valign="top"><a href="res2020-06-03.html">Resultater</a></td>
			<td style="vertical-align: top;" valign="top"><a href="https://www.livelox.com/Events/Show/49429" target="_blank">LL</a></td>
			<td valign="top"><a href="res2020-06-03.zip">res-zip</a></td>
		</tr>
		<tr>
			<td style="vertical-align: top;" valign="top">2</td>
			<td style="vertical-align: top;" valign="top">10.06</td>
			<td style="vertical-align: top;" valign="top">ons</td>
			<td style="vertical-align: top; font-family: Times New Roman;" valign="top"><a href="innbyd2020-06-03.pdf">Skansebakken</a></td>
			<td style="vertical-align: top;" valign="top">Bergendal-Lyse</td>
			<td style="vertical-align: top;" valign="top">Fossum IF</td>
			<td style="vertical-align: top;" valign="top"><a href="res2020-06-10.html">Resultater</a></td>
			<td style="vertical-align: top;" valign="top"><a href="https://www.livelox.com/Events/Show/49429/Sorkedalskarusellen-1">LL</a></td>
			<td valign="top"><a href="res2020-06-10.zip">res-zip</a></td>
		</tr>
		<tr>
			<td style="vertical-align: top;" valign="top">3</td>
			<td style="vertical-align: top;" valign="top">17.06</td>
			<td style="vertical-align: top;" valign="top">ons</td>
			<td style="vertical-align: top; font-family: Times New Roman;" valign="top"><a href="innbyd2020-06-17.pdf">Årnes P</a></td>
			<td style="vertical-align: top;" valign="top">Tangenåsen</td>
			<td style="vertical-align: top;" valign="top">Fossum IF</td>
			<td style="vertical-align: top;" valign="top"><a href="res2020-06-17.html">Resultater</a></td>
			<td style="vertical-align: top;" valign="top"><a href="https://www.livelox.com/Events/Show/50225/Sorkedalskarusellen-2">LL</a></td>
			<td valign="top"><a href="res2020-06-17.zip">res-zip</a></td>
		</tr>
		<tr>
			<td style="vertical-align: top;" valign="top">4</td>
			<td style="vertical-align: top;" valign="top">24.06</td>
			<td style="vertical-align: top;" valign="top">ons</td>
			<td style="vertical-align: top;" valign="top"><a href="innbyd2020-06-24.pdf">Stensrudtjern</a></td>
			<td style="vertical-align: top;" valign="top">Slettåsen</td>
			<td style="vertical-align: top;" valign="top">OSI</td>
			<td style="vertical-align: top;" valign="top">Thomas og Mathias Nipen, Mikael Olausson</td>
			<td style="vertical-align: top;" valign="top"><a href="res2020-06-24.html">Resultater</a></td>
			<td style="vertical-align: top;" valign="top"><a href="https://www.livelox.com/Events/Show/50650/OSI-GeoForm-Rankingl-p-4">LL</a></td>
			<td valign="top"><a href="res2020-06-24.zip">res-zip</a></td>
		</tr>
		<tr align="center">
			<td colspan="10" rowspan="1" style="vertical-align: top;
            background-color: rgb(153, 255, 255);" valign="top">Sommerferie. For sommerløp 1, 8, 15. og 22. juli samt 5. august se <a href="http://veritasorientering.com/">http://veritasorientering.com/</a></td>
		</tr>
		<tr>
			<td valign="top">5</td>
			<td valign="top">29.07</td>
			<td valign="top">ons</td>
			<td valign="top"><a href="innbyd2020-07-29.pdf">Sørbråten</a></td>
			<td valign="top">Maridalsalpene</td>
			<td valign="top">GeoForm</td>
			<td valign="top">Guttorm Egge</td>
			<td valign="top"><a href="res2020-07-29.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/51434">LL</a></td>
			<td valign="top"><a href="res2020-06-24.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">6</td>
			<td valign="top">12.08</td>
			<td valign="top">ons</td>
			<td valign="top"><a href="innbyd2020-08-12.pdf">Tryvannstårnet</a></td>
			<td valign="top">Ringeriksflaka</td>
			<td valign="top">OSI</td>
			<td valign="top">Øyvind Due Trier</td>
			<td valign="top"><a href="res2020-08-12.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/52401/OSI-Geoform-rankingl-p-6?culture=nb-NO">LL</a></td>
			<td valign="top"><a href="res2019-08-12.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">7</td>
			<td valign="top">19.08</td>
			<td valign="top">ons</td>
			<td valign="top"><a href="innbyd2020-08-19.pdf">Ellingsrud</a></td>
			<td valign="top">Puttåsen</td>
			<td valign="top">GeoForm</td>
			<td valign="top">Øyvind Grandum</td>
			<td valign="top"><a href="res2020-08-19.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/52566">LL</a></td>
			<td valign="top"><a href="res2019-08-21.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">8</td>
			<td valign="top">26.08</td>
			<td valign="top">ons</td>
			<td valign="top"><a href="innbyd2020-08-26.pdf">Lut/Lutvann</a></td>
			<td valign="top">Puttåsen</td>
			<td valign="top">OSI</td>
			<td valign="top">Are Eriksen og Emilie Sommerstad Vee</td>
			<td valign="top"><a href="res2020-08-26.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/53093">LL</a></td>
			<td valign="top"><a href="res2020-08-26.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">9</td>
			<td valign="top">02.09</td>
			<td valign="top">ons</td>
			<td valign="top"><a href="innbyd2020-09-02.pdf">Sandbakken</a></td>
			<td valign="top">Sandbakken S</td>
			<td valign="top">GeoForm</td>
			<td valign="top">Jon Lahlum</td>
			<td valign="top"><a href="res2020-09-02.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/52821">LL</a></td>
			<td valign="top"><a href="res2020-09-02.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">10</td>
			<td valign="top">09.09</td>
			<td valign="top">ons</td>
			<td valign="top"><a href="innbyd2020-09-09.pdf">Øvreseter</a></td>
			<td valign="top"></td>
			<td valign="top">OSI</td>
			<td valign="top">Magne Vollen og Marie Helene Hansen</td>
			<td valign="top"><a href="res2020-09-09.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/53809">LL</a></td>
			<td valign="top"><a href="res2020-09-09.zip">res-zip</a></td>
		</tr>
		<tr align="center">
			<td bgcolor="#99ffff" colspan="10" rowspan="1" valign="top">Lørdagsløp. For onsdagsløp/nattløp, se <a href="http://veritasorientering.com/">Veritas BIL</a></td>
		</tr>
		<tr>
			<td valign="top">11</td>
			<td valign="top">19.09</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-09-19.pdf">Øyungen</a></td>
			<td valign="top">Maridalsalpene</td>
			<td valign="top">GeoForm</td>
			<td valign="top">Harald Østgaard Lund</td>
			<td valign="top"><a href="res2020-09-19.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/54097">LL</a></td>
			<td valign="top"><a href="res2020-09-19.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">12</td>
			<td valign="top">26.09</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-09-26.pdf">Frognerseteren</a></td>
			<td valign="top">Frønsvollen øst</td>
			<td valign="top">OSI</td>
			<td valign="top">Dagfinn Føllesdal, Christin Vangen, Lenny Enstrøm</td>
			<td valign="top"><a href="res2020-09-26.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/54598/OSI-Geoform-Rankingl-p-12?culture=nb-NO">LL</a></td>
			<td valign="top"><a href="res2020-09-26.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">13</td>
			<td valign="top">03.10</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-10-03.pdf">Sognsvann P</a></td>
			<td valign="top">Sognsvann øst</td>
			<td valign="top">Geoform</td>
			<td valign="top">Audrun Utskarpen og Guttorm Egge</td>
			<td valign="top"><a href="res2020-10-03.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/54871/OSI-Geoform-rankinglop-13-Sognsvann">LL</a></td>
			<td valign="top"><a href="res2020-10-03.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">14</td>
			<td valign="top">10.10</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="https://eventor.orientering.no/Events/Show/12436" target="_blank">Lillomarka arena</a></td>
			<td valign="top">Lillomarka vest</td>
			<td valign="top">Lillomarka OL</td>
			<td valign="top">Lillomarka nord-sør og Gromløpet</td>
			<td valign="top"><a href="res2020-10-10.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/55130/Gromlopet-Lillomarka-Nord-Sor" target="_blank">LL</a></td>
			<td valign="top"><a href="res2020-10-10.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">15</td>
			<td valign="top">17.10</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-10-18.pdf">Elveli</a></td>
			<td valign="top">Heikampen</td>
			<td valign="top">GeoForm</td>
			<td valign="top">Niklas Norman, Skage Grønneberg, Johan Svedberg, Espen B. Prydz</td>
			<td valign="top"><a href="res2020-10-17.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/55233/OSI-GeoForm-Rankingl-p-15" target="_blank">LL</a></td>
			<td valign="top"><a href="res2020-10-18.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">16</td>
			<td valign="top">24.10</td>
			<td valign="top"><font color="#ff0000"><font color="#000000">lør</font></font></td>
			<td valign="top"><a href="innbyd2020-10-24.pdf">Katisa ved Nøklevann</a></td>
			<td valign="top">Slettfjell</td>
			<td valign="top">GeoForm</td>
			<td valign="top">Stig Hultgreen Karlsen</td>
			<td valign="top"><a href="res2020-10-24.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/55556/GeoForm-OSI-Rankingl-p-nr-16">LL</a></td>
			<td valign="top"><a href="res2020-10-24.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">17</td>
			<td valign="top">31.10</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-10-31.pdf">Haraløkka</a></td>
			<td valign="top">Ulsrudvann</td>
			<td valign="top">Nordre Follo orientering/GeoForm</td>
			<td valign="top">Jon og Tor Lahlum</td>
			<td valign="top"><a href="res2020-10-31.html">Resultater</a></td>
			<td valign="top"><a href="https://www.livelox.com/Events/Show/55670/" target="_blank">LL</a></td>
			<td valign="top"><a href="res2020-11-07.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">18</td>
			<td valign="top">07.11</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-11-07.pdf">Tonsenhagen</a></td>
			<td valign="top">Alunsjøen</td>
			<td valign="top">GeoForm</td>
			<td valign="top">Holger Schlaupitz</td>
			<td valign="top"><a href="res2020-11-07.html">Resultater</a></td>
			<td valign="top">LL</td>
			<td valign="top"><a href="res2020-11-07.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">19</td>
			<td valign="top">14.11</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-11-14.pdf">Solemskogen</a></td>
			<td valign="top"></td>
			<td valign="top">GeoForm</td>
			<td valign="top">Bjørn Grinde</td>
			<td valign="top"><a href="res2020-11-14.html">Resultater</a></td>
			<td valign="top">LL</td>
			<td valign="top"><a href="res2020-11-14.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">20</td>
			<td valign="top">21.11</td>
			<td valign="top">lør</td>
			<td valign="top">Låkeberget</td>
			<td valign="top"></td>
			<td valign="top">OSI</td>
			<td valign="top">Øivind Due Trier</td>
			<td valign="top"><a href="res2020-11-21.html">Resultater</a></td>
			<td valign="top">LL</td>
			<td valign="top"><a href="res2020-11-21.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">21</td>
			<td valign="top">28.11</td>
			<td valign="top">lør</td>
			<td valign="top">Badedammen Grorud</td>
			<td valign="top"></td>
			<td valign="top">GeoForm</td>
			<td valign="top">Harald Iwe</td>
			<td valign="top"><a href="res2020-11-28.html">Resultater</a></td>
			<td valign="top">LL</td>
			<td valign="top"><a href="res2020-11-28.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">22</td>
			<td valign="top">05.12</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-12-05.pdf">Urbant</a></td>
			<td valign="top"></td>
			<td valign="top">GeoForm</td>
			<td valign="top">Stig Hultgreen Karlsen</td>
			<td valign="top"><a href="res2020-12-05.html">Resultater</a></td>
			<td valign="top">LL</td>
			<td valign="top"><a href="res2020-12-05.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">23</td>
			<td valign="top">12.12</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-12-12.pdf">Oslo bynært</a></td>
			<td valign="top"></td>
			<td valign="top">Sentrum OK</td>
			<td valign="top"><b>Ribbe og Loff</b>, Helge Gisholt</td>
			<td valign="top"><a href="res2020-12-12.html">Resultater</a></td>
			<td valign="top">LL</td>
			<td valign="top"><a href="res2020-12-12.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top">24</td>
			<td valign="top">19.12</td>
			<td valign="top">lør</td>
			<td valign="top"><a href="innbyd2020-12-19.pdf">Sognsvann P</a></td>
			<td valign="top">Sognsvann</td>
			<td valign="top">GeoForm</td>
			<td valign="top"><b>Jon Tolgensbakks minneløp - Gløggløpet</b><br>
			Oddny Feragen, Ida Tolgensbakk, Eystein Grimstad</td>
			<td valign="top"><a href="res2020-12-19.html">Resultater</a></td>
			<td valign="top">LL</td>
			<td valign="top"><a href="res2020-12-19.zip">res-zip</a></td>
		</tr>
		<tr>
			<td valign="top"></td>
			<td valign="top">28.01</td>
			<td valign="top">ons</td>
			<td valign="top"><a href="geoformersbankett-2020.pdf">Borggata 2B</a></td>
			<td valign="top">Fafo</td>
			<td valign="top">Festkomite</td>
			<td valign="top">Geoformersbanketten 2021</td>
			<td valign="top"></td>
			<td valign="top"></td>
			<td valign="top"></td>
		</tr>
	</tbody>
</table>

<p><a href="rank-2018.xlsx"><g class="gr_ gr_42 gr-alert gr_spell
          gr_inline_cards gr_run_anim ContextualSpelling ins-del
          multiReplace" data-gr-id="42" id="42"></g></a><a href="rank-2019.xlsx"><g class="gr_ gr_42 gr-alert gr_spell gr_inline_cards
              gr_run_anim ContextualSpelling ins-del multiReplace" data-gr-id="42" id="42">Totalranking</g> 2019</a>&nbsp; Totalranking 2018&nbsp;&nbsp;<a href="rank-2017.xlsx"><g class="gr_ gr_43 gr-alert gr_spell gr_inline_cards gr_run_anim
          ContextualSpelling ins-del multiReplace" data-gr-id="43" id="43">Totalranking</g> 2017</a>&nbsp;&nbsp; <a href="rank-2016.xlsx"><g class="gr_ gr_44 gr-alert gr_spell
          gr_inline_cards gr_run_anim ContextualSpelling ins-del
          multiReplace" data-gr-id="44" id="44">Totalranking</g> 2016</a>&nbsp;&nbsp; <a href="rank-2015.xlsx"><g class="gr_ gr_45 gr-alert gr_spell
          gr_inline_cards gr_run_anim ContextualSpelling ins-del
          multiReplace" data-gr-id="45" id="45">Totalranking</g> 2015</a> <a href="rank-2014.xls"><g class="gr_ gr_46 gr-alert gr_spell
          gr_inline_cards gr_run_anim ContextualSpelling ins-del
          multiReplace" data-gr-id="46" id="46">Totalranking</g> 2014</a>&nbsp;&nbsp; <a href="rank-2013.xls"><g class="gr_ gr_47 gr-alert gr_spell
          gr_inline_cards gr_run_anim ContextualSpelling ins-del
          multiReplace" data-gr-id="47" id="47">Totalranking</g> 2013</a>&nbsp;&nbsp; <a href="rank-2012.xls"><g class="gr_ gr_48 gr-alert gr_spell
          gr_inline_cards gr_run_anim ContextualSpelling ins-del
          multiReplace" data-gr-id="48" id="48">Totalranking</g> 2012</a>&nbsp;&nbsp; <a href="rank-2011.xls"><g class="gr_ gr_49 gr-alert gr_spell
          gr_inline_cards gr_run_anim ContextualSpelling ins-del
          multiReplace" data-gr-id="49" id="49">Totalranking</g> 2011</a><br>
<a href="index-2019.html">Resultater 2019</a> &nbsp;&nbsp;&nbsp; <a href="index-2018.html">Resultater 2018</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <a href="index-2017.html">Resultater 2017</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <a href="index-2016.html">Resultater 2016</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <a href="index-2015.html">Resultater 2015</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <a href="index-2014.html">Resultater 2014</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <a href="index-2013.html">Resultater 2013</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <a href="index-2012.html">Resultater 2012</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <a href="index-2011.html">Resultater 2011</a></p>

<p></p>

<p></p>

<p></p>

<p></p>

<p></p>

<p></p>

<p></p>

<p></p>

<p></p>

<p></p>

</body></html>`
)