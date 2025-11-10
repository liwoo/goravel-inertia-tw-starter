package requests

// Nationality type
type Nationality string

// Nationality constants - comprehensive list of nationalities
const (
	NationalityAfghan                      Nationality = "Afghan"
	NationalityAlbanian                    Nationality = "Albanian"
	NationalityAlgerian                    Nationality = "Algerian"
	NationalityAmerican                    Nationality = "American"
	NationalityAndorran                    Nationality = "Andorran"
	NationalityAngolan                     Nationality = "Angolan"
	NationalityAnguillan                   Nationality = "Anguillan"
	NationalityAntiguaAndBarbuda           Nationality = "Citizen of Antigua and Barbuda"
	NationalityArgentine                   Nationality = "Argentine"
	NationalityArmenian                    Nationality = "Armenian"
	NationalityAustralian                  Nationality = "Australian"
	NationalityAustrian                    Nationality = "Austrian"
	NationalityAzerbaijani                 Nationality = "Azerbaijani"
	NationalityBahamian                    Nationality = "Bahamian"
	NationalityBahraini                    Nationality = "Bahraini"
	NationalityBangladeshi                 Nationality = "Bangladeshi"
	NationalityBarbadian                   Nationality = "Barbadian"
	NationalityBelarusian                  Nationality = "Belarusian"
	NationalityBelgian                     Nationality = "Belgian"
	NationalityBelizean                    Nationality = "Belizean"
	NationalityBeninese                    Nationality = "Beninese"
	NationalityBermudian                   Nationality = "Bermudian"
	NationalityBhutanese                   Nationality = "Bhutanese"
	NationalityBolivian                    Nationality = "Bolivian"
	NationalityBosniaAndHerzegovina        Nationality = "Citizen of Bosnia and Herzegovina"
	NationalityBotswanan                   Nationality = "Botswanan"
	NationalityBrazilian                   Nationality = "Brazilian"
	NationalityBritish                     Nationality = "British"
	NationalityBritishVirginIslander       Nationality = "British Virgin Islander"
	NationalityBruneian                    Nationality = "Bruneian"
	NationalityBulgarian                   Nationality = "Bulgarian"
	NationalityBurkinan                    Nationality = "Burkinan"
	NationalityBurmese                     Nationality = "Burmese"
	NationalityBurundian                   Nationality = "Burundian"
	NationalityCambodian                   Nationality = "Cambodian"
	NationalityCameroonian                 Nationality = "Cameroonian"
	NationalityCanadian                    Nationality = "Canadian"
	NationalityCapeVerdean                 Nationality = "Cape Verdean"
	NationalityCaymanIslander              Nationality = "Cayman Islander"
	NationalityCentralAfrican              Nationality = "Central African"
	NationalityChadian                     Nationality = "Chadian"
	NationalityChilean                     Nationality = "Chilean"
	NationalityChinese                     Nationality = "Chinese"
	NationalityColombian                   Nationality = "Colombian"
	NationalityComoran                     Nationality = "Comoran"
	NationalityCongoleseCongo              Nationality = "Congolese (Congo)"
	NationalityCongoleseDRC                Nationality = "Congolese (DRC)"
	NationalityCookIslander                Nationality = "Cook Islander"
	NationalityCostaRican                  Nationality = "Costa Rican"
	NationalityCroatian                    Nationality = "Croatian"
	NationalityCuban                       Nationality = "Cuban"
	NationalityCymraes                     Nationality = "Cymraes"
	NationalityCymro                       Nationality = "Cymro"
	NationalityCypriot                     Nationality = "Cypriot"
	NationalityCzech                       Nationality = "Czech"
	NationalityDanish                      Nationality = "Danish"
	NationalityDjiboutian                  Nationality = "Djiboutian"
	NationalityDominican                   Nationality = "Dominican"
	NationalityDominicanRepublic           Nationality = "Citizen of the Dominican Republic"
	NationalityDutch                       Nationality = "Dutch"
	NationalityEastTimorese                Nationality = "East Timorese"
	NationalityEcuadorean                  Nationality = "Ecuadorean"
	NationalityEgyptian                    Nationality = "Egyptian"
	NationalityEmirati                     Nationality = "Emirati"
	NationalityEnglish                     Nationality = "English"
	NationalityEquatorialGuinean           Nationality = "Equatorial Guinean"
	NationalityEritrean                    Nationality = "Eritrean"
	NationalityEstonian                    Nationality = "Estonian"
	NationalityEthiopian                   Nationality = "Ethiopian"
	NationalityFaroese                     Nationality = "Faroese"
	NationalityFijian                      Nationality = "Fijian"
	NationalityFilipino                    Nationality = "Filipino"
	NationalityFinnish                     Nationality = "Finnish"
	NationalityFrench                      Nationality = "French"
	NationalityGabonese                    Nationality = "Gabonese"
	NationalityGambian                     Nationality = "Gambian"
	NationalityGeorgian                    Nationality = "Georgian"
	NationalityGerman                      Nationality = "German"
	NationalityGhanaian                    Nationality = "Ghanaian"
	NationalityGibraltarian                Nationality = "Gibraltarian"
	NationalityGreek                       Nationality = "Greek"
	NationalityGreenlandic                 Nationality = "Greenlandic"
	NationalityGrenadian                   Nationality = "Grenadian"
	NationalityGuamanian                   Nationality = "Guamanian"
	NationalityGuatemalan                  Nationality = "Guatemalan"
	NationalityGuineaBissau                Nationality = "Citizen of Guinea-Bissau"
	NationalityGuinean                     Nationality = "Guinean"
	NationalityGuyanese                    Nationality = "Guyanese"
	NationalityHaitian                     Nationality = "Haitian"
	NationalityHonduran                    Nationality = "Honduran"
	NationalityHongKonger                  Nationality = "Hong Konger"
	NationalityHungarian                   Nationality = "Hungarian"
	NationalityIcelandic                   Nationality = "Icelandic"
	NationalityIndian                      Nationality = "Indian"
	NationalityIndonesian                  Nationality = "Indonesian"
	NationalityIranian                     Nationality = "Iranian"
	NationalityIraqi                       Nationality = "Iraqi"
	NationalityIrish                       Nationality = "Irish"
	NationalityIsraeli                     Nationality = "Israeli"
	NationalityItalian                     Nationality = "Italian"
	NationalityIvorian                     Nationality = "Ivorian"
	NationalityJamaican                    Nationality = "Jamaican"
	NationalityJapanese                    Nationality = "Japanese"
	NationalityJordanian                   Nationality = "Jordanian"
	NationalityKazakh                      Nationality = "Kazakh"
	NationalityKenyan                      Nationality = "Kenyan"
	NationalityKittitian                   Nationality = "Kittitian"
	NationalityKiribati                    Nationality = "Citizen of Kiribati"
	NationalityKosovan                     Nationality = "Kosovan"
	NationalityKuwaiti                     Nationality = "Kuwaiti"
	NationalityKyrgyz                      Nationality = "Kyrgyz"
	NationalityLao                         Nationality = "Lao"
	NationalityLatvian                     Nationality = "Latvian"
	NationalityLebanese                    Nationality = "Lebanese"
	NationalityLiberian                    Nationality = "Liberian"
	NationalityLibyan                      Nationality = "Libyan"
	NationalityLiechtenstein               Nationality = "Liechtenstein citizen"
	NationalityLithuanian                  Nationality = "Lithuanian"
	NationalityLuxembourger                Nationality = "Luxembourger"
	NationalityMacanese                    Nationality = "Macanese"
	NationalityMacedonian                  Nationality = "Macedonian"
	NationalityMalagasy                    Nationality = "Malagasy"
	NationalityMalawian                    Nationality = "Malawian"
	NationalityMalaysian                   Nationality = "Malaysian"
	NationalityMaldivian                   Nationality = "Maldivian"
	NationalityMalian                      Nationality = "Malian"
	NationalityMaltese                     Nationality = "Maltese"
	NationalityMarshallese                 Nationality = "Marshallese"
	NationalityMartiniquais                Nationality = "Martiniquais"
	NationalityMauritanian                 Nationality = "Mauritanian"
	NationalityMauritian                   Nationality = "Mauritian"
	NationalityMexican                     Nationality = "Mexican"
	NationalityMicronesian                 Nationality = "Micronesian"
	NationalityMoldovan                    Nationality = "Moldovan"
	NationalityMonegasque                  Nationality = "Monegasque"
	NationalityMongolian                   Nationality = "Mongolian"
	NationalityMontenegrin                 Nationality = "Montenegrin"
	NationalityMontserratian               Nationality = "Montserratian"
	NationalityMoroccan                    Nationality = "Moroccan"
	NationalityMosotho                     Nationality = "Mosotho"
	NationalityMozambican                  Nationality = "Mozambican"
	NationalityNamibian                    Nationality = "Namibian"
	NationalityNauruan                     Nationality = "Nauruan"
	NationalityNepalese                    Nationality = "Nepalese"
	NationalityNewZealander                Nationality = "New Zealander"
	NationalityNicaraguan                  Nationality = "Nicaraguan"
	NationalityNigerian                    Nationality = "Nigerian"
	NationalityNigerien                    Nationality = "Nigerien"
	NationalityNiuean                      Nationality = "Niuean"
	NationalityNorthKorean                 Nationality = "North Korean"
	NationalityNorthernIrish               Nationality = "Northern Irish"
	NationalityNorwegian                   Nationality = "Norwegian"
	NationalityOmani                       Nationality = "Omani"
	NationalityPakistani                   Nationality = "Pakistani"
	NationalityPalauan                     Nationality = "Palauan"
	NationalityPalestinian                 Nationality = "Palestinian"
	NationalityPanamanian                  Nationality = "Panamanian"
	NationalityPapuaNewGuinean             Nationality = "Papua New Guinean"
	NationalityParaguayan                  Nationality = "Paraguayan"
	NationalityPeruvian                    Nationality = "Peruvian"
	NationalityPitcairnIslander            Nationality = "Pitcairn Islander"
	NationalityPolish                      Nationality = "Polish"
	NationalityPortuguese                  Nationality = "Portuguese"
	NationalityPrydeinig                   Nationality = "Prydeinig"
	NationalityPuertoRican                 Nationality = "Puerto Rican"
	NationalityQatari                      Nationality = "Qatari"
	NationalityRomanian                    Nationality = "Romanian"
	NationalityRussian                     Nationality = "Russian"
	NationalityRwandan                     Nationality = "Rwandan"
	NationalitySalvadorean                 Nationality = "Salvadorean"
	NationalitySammarinese                 Nationality = "Sammarinese"
	NationalitySamoan                      Nationality = "Samoan"
	NationalitySaoTomean                   Nationality = "Sao Tomean"
	NationalitySaudiArabian                Nationality = "Saudi Arabian"
	NationalityScottish                    Nationality = "Scottish"
	NationalitySenegalese                  Nationality = "Senegalese"
	NationalitySerbian                     Nationality = "Serbian"
	NationalitySeychelles                  Nationality = "Citizen of Seychelles"
	NationalitySierraLeonean               Nationality = "Sierra Leonean"
	NationalitySingaporean                 Nationality = "Singaporean"
	NationalitySlovak                      Nationality = "Slovak"
	NationalitySlovenian                   Nationality = "Slovenian"
	NationalitySolomonIslander             Nationality = "Solomon Islander"
	NationalitySomali                      Nationality = "Somali"
	NationalitySouthAfrican                Nationality = "South African"
	NationalitySouthKorean                 Nationality = "South Korean"
	NationalitySouthSudanese               Nationality = "South Sudanese"
	NationalitySpanish                     Nationality = "Spanish"
	NationalitySriLankan                   Nationality = "Sri Lankan"
	NationalityStHelenian                  Nationality = "St Helenian"
	NationalityStLucian                    Nationality = "St Lucian"
	NationalityStateless                   Nationality = "Stateless"
	NationalitySudanese                    Nationality = "Sudanese"
	NationalitySurinamese                  Nationality = "Surinamese"
	NationalitySwazi                       Nationality = "Swazi"
	NationalitySwedish                     Nationality = "Swedish"
	NationalitySwiss                       Nationality = "Swiss"
	NationalitySyrian                      Nationality = "Syrian"
	NationalityTaiwanese                   Nationality = "Taiwanese"
	NationalityTajik                       Nationality = "Tajik"
	NationalityTanzanian                   Nationality = "Tanzanian"
	NationalityThai                        Nationality = "Thai"
	NationalityTogolese                    Nationality = "Togolese"
	NationalityTongan                      Nationality = "Tongan"
	NationalityTrinidadian                 Nationality = "Trinidadian"
	NationalityTristanian                  Nationality = "Tristanian"
	NationalityTunisian                    Nationality = "Tunisian"
	NationalityTurkish                     Nationality = "Turkish"
	NationalityTurkmen                     Nationality = "Turkmen"
	NationalityTurksAndCaicosIslander      Nationality = "Turks and Caicos Islander"
	NationalityTuvaluan                    Nationality = "Tuvaluan"
	NationalityUgandan                     Nationality = "Ugandan"
	NationalityUkrainian                   Nationality = "Ukrainian"
	NationalityUruguayan                   Nationality = "Uruguayan"
	NationalityUzbek                       Nationality = "Uzbek"
	NationalityVaticanCitizen              Nationality = "Vatican citizen"
	NationalityVanuatu                     Nationality = "Citizen of Vanuatu"
	NationalityVenezuelan                  Nationality = "Venezuelan"
	NationalityVietnamese                  Nationality = "Vietnamese"
	NationalityVincentian                  Nationality = "Vincentian"
	NationalityWallisian                   Nationality = "Wallisian"
	NationalityWelsh                       Nationality = "Welsh"
	NationalityYemeni                      Nationality = "Yemeni"
	NationalityZambian                     Nationality = "Zambian"
	NationalityZimbabwean                  Nationality = "Zimbabwean"
)

// AllNationalities returns all nationality values
func AllNationalities() []Nationality {
	return []Nationality{
		NationalityAfghan,
		NationalityAlbanian,
		NationalityAlgerian,
		NationalityAmerican,
		NationalityAndorran,
		NationalityAngolan,
		NationalityAnguillan,
		NationalityAntiguaAndBarbuda,
		NationalityArgentine,
		NationalityArmenian,
		NationalityAustralian,
		NationalityAustrian,
		NationalityAzerbaijani,
		NationalityBahamian,
		NationalityBahraini,
		NationalityBangladeshi,
		NationalityBarbadian,
		NationalityBelarusian,
		NationalityBelgian,
		NationalityBelizean,
		NationalityBeninese,
		NationalityBermudian,
		NationalityBhutanese,
		NationalityBolivian,
		NationalityBosniaAndHerzegovina,
		NationalityBotswanan,
		NationalityBrazilian,
		NationalityBritish,
		NationalityBritishVirginIslander,
		NationalityBruneian,
		NationalityBulgarian,
		NationalityBurkinan,
		NationalityBurmese,
		NationalityBurundian,
		NationalityCambodian,
		NationalityCameroonian,
		NationalityCanadian,
		NationalityCapeVerdean,
		NationalityCaymanIslander,
		NationalityCentralAfrican,
		NationalityChadian,
		NationalityChilean,
		NationalityChinese,
		NationalityColombian,
		NationalityComoran,
		NationalityCongoleseCongo,
		NationalityCongoleseDRC,
		NationalityCookIslander,
		NationalityCostaRican,
		NationalityCroatian,
		NationalityCuban,
		NationalityCymraes,
		NationalityCymro,
		NationalityCypriot,
		NationalityCzech,
		NationalityDanish,
		NationalityDjiboutian,
		NationalityDominican,
		NationalityDominicanRepublic,
		NationalityDutch,
		NationalityEastTimorese,
		NationalityEcuadorean,
		NationalityEgyptian,
		NationalityEmirati,
		NationalityEnglish,
		NationalityEquatorialGuinean,
		NationalityEritrean,
		NationalityEstonian,
		NationalityEthiopian,
		NationalityFaroese,
		NationalityFijian,
		NationalityFilipino,
		NationalityFinnish,
		NationalityFrench,
		NationalityGabonese,
		NationalityGambian,
		NationalityGeorgian,
		NationalityGerman,
		NationalityGhanaian,
		NationalityGibraltarian,
		NationalityGreek,
		NationalityGreenlandic,
		NationalityGrenadian,
		NationalityGuamanian,
		NationalityGuatemalan,
		NationalityGuineaBissau,
		NationalityGuinean,
		NationalityGuyanese,
		NationalityHaitian,
		NationalityHonduran,
		NationalityHongKonger,
		NationalityHungarian,
		NationalityIcelandic,
		NationalityIndian,
		NationalityIndonesian,
		NationalityIranian,
		NationalityIraqi,
		NationalityIrish,
		NationalityIsraeli,
		NationalityItalian,
		NationalityIvorian,
		NationalityJamaican,
		NationalityJapanese,
		NationalityJordanian,
		NationalityKazakh,
		NationalityKenyan,
		NationalityKittitian,
		NationalityKiribati,
		NationalityKosovan,
		NationalityKuwaiti,
		NationalityKyrgyz,
		NationalityLao,
		NationalityLatvian,
		NationalityLebanese,
		NationalityLiberian,
		NationalityLibyan,
		NationalityLiechtenstein,
		NationalityLithuanian,
		NationalityLuxembourger,
		NationalityMacanese,
		NationalityMacedonian,
		NationalityMalagasy,
		NationalityMalawian,
		NationalityMalaysian,
		NationalityMaldivian,
		NationalityMalian,
		NationalityMaltese,
		NationalityMarshallese,
		NationalityMartiniquais,
		NationalityMauritanian,
		NationalityMauritian,
		NationalityMexican,
		NationalityMicronesian,
		NationalityMoldovan,
		NationalityMonegasque,
		NationalityMongolian,
		NationalityMontenegrin,
		NationalityMontserratian,
		NationalityMoroccan,
		NationalityMosotho,
		NationalityMozambican,
		NationalityNamibian,
		NationalityNauruan,
		NationalityNepalese,
		NationalityNewZealander,
		NationalityNicaraguan,
		NationalityNigerian,
		NationalityNigerien,
		NationalityNiuean,
		NationalityNorthKorean,
		NationalityNorthernIrish,
		NationalityNorwegian,
		NationalityOmani,
		NationalityPakistani,
		NationalityPalauan,
		NationalityPalestinian,
		NationalityPanamanian,
		NationalityPapuaNewGuinean,
		NationalityParaguayan,
		NationalityPeruvian,
		NationalityPitcairnIslander,
		NationalityPolish,
		NationalityPortuguese,
		NationalityPrydeinig,
		NationalityPuertoRican,
		NationalityQatari,
		NationalityRomanian,
		NationalityRussian,
		NationalityRwandan,
		NationalitySalvadorean,
		NationalitySammarinese,
		NationalitySamoan,
		NationalitySaoTomean,
		NationalitySaudiArabian,
		NationalityScottish,
		NationalitySenegalese,
		NationalitySerbian,
		NationalitySeychelles,
		NationalitySierraLeonean,
		NationalitySingaporean,
		NationalitySlovak,
		NationalitySlovenian,
		NationalitySolomonIslander,
		NationalitySomali,
		NationalitySouthAfrican,
		NationalitySouthKorean,
		NationalitySouthSudanese,
		NationalitySpanish,
		NationalitySriLankan,
		NationalityStHelenian,
		NationalityStLucian,
		NationalityStateless,
		NationalitySudanese,
		NationalitySurinamese,
		NationalitySwazi,
		NationalitySwedish,
		NationalitySwiss,
		NationalitySyrian,
		NationalityTaiwanese,
		NationalityTajik,
		NationalityTanzanian,
		NationalityThai,
		NationalityTogolese,
		NationalityTongan,
		NationalityTrinidadian,
		NationalityTristanian,
		NationalityTunisian,
		NationalityTurkish,
		NationalityTurkmen,
		NationalityTurksAndCaicosIslander,
		NationalityTuvaluan,
		NationalityUgandan,
		NationalityUkrainian,
		NationalityUruguayan,
		NationalityUzbek,
		NationalityVaticanCitizen,
		NationalityVanuatu,
		NationalityVenezuelan,
		NationalityVietnamese,
		NationalityVincentian,
		NationalityWallisian,
		NationalityWelsh,
		NationalityYemeni,
		NationalityZambian,
		NationalityZimbabwean,
	}
}
