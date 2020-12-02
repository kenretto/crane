// Package areacode International area code and time zone
package areacode

import (
	"sync"
	"time"
)

// Code area code information
type Code struct {
	Country        string // Country or region
	DomainSuffix   string // International domain name suffix
	Code           int16  // International number area code
	TimeDifference int    // Time difference from Beijing time
}

// Filter filter
type Filter struct {
	// International number area code
	Code int16
	// country code
	CountryCode string
}

// Codes code list
type Codes []Code

var (
	// indexes Build an index
	indexes        map[int16][]Code
	ensureIndexing sync.Once
)

func init() {
	ensureIndexing.Do(func() {
		indexes = make(map[int16][]Code)
		for _, code := range AllCode {
			indexes[code.Code] = append(indexes[code.Code], code)
		}
	})
}

// AllCode all code
var AllCode = Codes{
	{
		Country:        "Angola",
		DomainSuffix:   "AO",
		Code:           244,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Afghanistan",
		DomainSuffix:   "AF",
		Code:           93,
		TimeDifference: 0,
	},
	{
		Country:        "Albania",
		DomainSuffix:   "AL",
		Code:           355,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Algeria",
		DomainSuffix:   "DZ",
		Code:           213,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Andorra",
		DomainSuffix:   "AD",
		Code:           376,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Anguilla",
		DomainSuffix:   "AI",
		Code:           1264,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Antigua and Barbuda ",
		DomainSuffix:   "AG",
		Code:           1268,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Argentina",
		DomainSuffix:   "AR",
		Code:           54,
		TimeDifference: -11 * 3600,
	},
	{
		Country:        "Armenia",
		DomainSuffix:   "AM",
		Code:           374,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Ascension",
		DomainSuffix:   "",
		Code:           247,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Australia",
		DomainSuffix:   "AU",
		Code:           61,
		TimeDifference: +2 * 3600,
	},
	{
		Country:        "Austria",
		DomainSuffix:   "AT",
		Code:           43,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Azerbaijan",
		DomainSuffix:   "AZ",
		Code:           994,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Bahamas",
		DomainSuffix:   "BS",
		Code:           1242,
		TimeDifference: -13 * 3600,
	},
	{
		Country:        "Bahrain",
		DomainSuffix:   "BH",
		Code:           973,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Bangladesh",
		DomainSuffix:   "BD",
		Code:           880,
		TimeDifference: -2 * 3600,
	},
	{
		Country:        "Barbados",
		DomainSuffix:   "BB",
		Code:           1246,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Belarus",
		DomainSuffix:   "BY",
		Code:           375,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Belgium",
		DomainSuffix:   "BE",
		Code:           32,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Belize",
		DomainSuffix:   "BZ",
		Code:           501,
		TimeDifference: -14 * 3600,
	},
	{
		Country:        "Benin",
		DomainSuffix:   "BJ",
		Code:           229,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Bermuda Is. ",
		DomainSuffix:   "BM",
		Code:           1441,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Bolivia",
		DomainSuffix:   "BO",
		Code:           591,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Botswana",
		DomainSuffix:   "BW",
		Code:           267,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Brazil",
		DomainSuffix:   "BR",
		Code:           55,
		TimeDifference: -11 * 3600,
	},
	{
		Country:        "Brunei",
		DomainSuffix:   "BN",
		Code:           673,
		TimeDifference: 0 * 3600,
	},
	{
		Country:        "Bulgaria",
		DomainSuffix:   "BG",
		Code:           359,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Burkina-faso ",
		DomainSuffix:   "BF",
		Code:           226,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Burma",
		DomainSuffix:   "MM",
		Code:           95,
		TimeDifference: -1.3 * 3600,
	},
	{
		Country:        "Burundi",
		DomainSuffix:   "BI",
		Code:           257,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Cameroon",
		DomainSuffix:   "CM",
		Code:           237,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Canada",
		DomainSuffix:   "CA",
		Code:           1,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Cayman Is. ",
		DomainSuffix:   "",
		Code:           1345,
		TimeDifference: -13 * 3600,
	},
	{
		Country:        "Central African Republic ",
		DomainSuffix:   "CF",
		Code:           236,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Chad",
		DomainSuffix:   "TD",
		Code:           235,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Chile",
		DomainSuffix:   "CL",
		Code:           56,
		TimeDifference: -13 * 3600,
	},
	{
		Country:        "China",
		DomainSuffix:   "CN",
		Code:           86,
		TimeDifference: 0 * 3600,
	},
	{
		Country:        "Colombia",
		DomainSuffix:   "CO",
		Code:           57,
		TimeDifference: 0 * 3600,
	},
	{
		Country:        "Congo",
		DomainSuffix:   "CG",
		Code:           242,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Cook Is. ",
		DomainSuffix:   "CK",
		Code:           682,
		TimeDifference: -18.3 * 3600,
	},
	{
		Country:        "Costa Rica ",
		DomainSuffix:   "CR",
		Code:           506,
		TimeDifference: -14 * 3600,
	},
	{
		Country:        "Cuba",
		DomainSuffix:   "CU",
		Code:           53,
		TimeDifference: -13 * 3600,
	},
	{
		Country:        "Cyprus",
		DomainSuffix:   "CY",
		Code:           357,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Czech Republic ",
		DomainSuffix:   "CZ",
		Code:           420,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Denmark",
		DomainSuffix:   "DK",
		Code:           45,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Djibouti",
		DomainSuffix:   "DJ",
		Code:           253,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Dominica Rep. ",
		DomainSuffix:   "DO",
		Code:           1890,
		TimeDifference: -13 * 3600,
	},
	{
		Country:        "Ecuador",
		DomainSuffix:   "EC",
		Code:           593,
		TimeDifference: -13 * 3600,
	},
	{
		Country:        "Egypt",
		DomainSuffix:   "EG",
		Code:           20,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "EI Salvador ",
		DomainSuffix:   "SV",
		Code:           503,
		TimeDifference: -14 * 3600,
	},
	{
		Country:        "Estonia",
		DomainSuffix:   "EE",
		Code:           372,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Ethiopia",
		DomainSuffix:   "ET",
		Code:           251,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Fiji",
		DomainSuffix:   "FJ",
		Code:           679,
		TimeDifference: +4 * 3600,
	},
	{
		Country:        "Finland",
		DomainSuffix:   "FI",
		Code:           358,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "France",
		DomainSuffix:   "FR",
		Code:           33,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "French Guiana ",
		DomainSuffix:   "GF",
		Code:           594,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Gabon",
		DomainSuffix:   "GA",
		Code:           241,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Gambia",
		DomainSuffix:   "GM",
		Code:           220,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Georgia",
		DomainSuffix:   "GE",
		Code:           995,
		TimeDifference: 0 * 3600,
	},
	{
		Country:        "Germany",
		DomainSuffix:   "DE",
		Code:           49,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Ghana",
		DomainSuffix:   "GH",
		Code:           233,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Gibraltar",
		DomainSuffix:   "GI",
		Code:           350,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Greece",
		DomainSuffix:   "GR",
		Code:           30,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Grenada",
		DomainSuffix:   "GD",
		Code:           1809,
		TimeDifference: -14 * 3600,
	},
	{
		Country:        "Guam",
		DomainSuffix:   "GU",
		Code:           1671,
		TimeDifference: +2 * 3600,
	},
	{
		Country:        "Guatemala",
		DomainSuffix:   "GT",
		Code:           502,
		TimeDifference: -14 * 3600,
	},
	{
		Country:        "Guinea",
		DomainSuffix:   "GN",
		Code:           224,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Guyana",
		DomainSuffix:   "GY",
		Code:           592,
		TimeDifference: -11,
	},
	{
		Country:        "Haiti",
		DomainSuffix:   "HT",
		Code:           509,
		TimeDifference: -13 * 3600,
	},
	{
		Country:        "Honduras",
		DomainSuffix:   "HN",
		Code:           504,
		TimeDifference: -14 * 3600,
	},
	{
		Country:        "Hongkong",
		DomainSuffix:   "HK",
		Code:           852,
		TimeDifference: 0,
	},
	{
		Country:        "Hungary",
		DomainSuffix:   "HU",
		Code:           36,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Iceland",
		DomainSuffix:   "IS",
		Code:           354,
		TimeDifference: -9 * 3600,
	},
	{
		Country:        "India",
		DomainSuffix:   "IN",
		Code:           91,
		TimeDifference: -2.3 * 3600,
	},
	{
		Country:        "Indonesia",
		DomainSuffix:   "ID",
		Code:           62,
		TimeDifference: -0.3 * 3600,
	},
	{
		Country:        "Iran",
		DomainSuffix:   "IR",
		Code:           98,
		TimeDifference: -4.3 * 3600,
	},
	{
		Country:        "Iraq",
		DomainSuffix:   "IQ",
		Code:           964,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Ireland",
		DomainSuffix:   "IE",
		Code:           353,
		TimeDifference: -4.3 * 3600,
	},
	{
		Country:        "Israel",
		DomainSuffix:   "IL",
		Code:           972,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Italy",
		DomainSuffix:   "IT",
		Code:           39,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Ivory Coast ",
		DomainSuffix:   "",
		Code:           225,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Jamaica",
		DomainSuffix:   "JM",
		Code:           1876,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Japan",
		DomainSuffix:   "JP",
		Code:           81,
		TimeDifference: +1 * 3600,
	},
	{
		Country:        "Jordan",
		DomainSuffix:   "JO",
		Code:           962,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Kampuchea (Cambodia ) ",
		DomainSuffix:   "KH",
		Code:           855,
		TimeDifference: -1 * 3600,
	},
	{
		Country:        "Kazakstan",
		DomainSuffix:   "KZ",
		Code:           327,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Kenya",
		DomainSuffix:   "KE",
		Code:           254,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Korea",
		DomainSuffix:   "KR",
		Code:           82,
		TimeDifference: +1 * 3600,
	},
	{
		Country:        "Kuwait",
		DomainSuffix:   "KW",
		Code:           965,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Kyrgyzstan",
		DomainSuffix:   "KG",
		Code:           331,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Laos",
		DomainSuffix:   "LA",
		Code:           856,
		TimeDifference: -1 * 3600,
	},
	{
		Country:        "Latvia",
		DomainSuffix:   "LV",
		Code:           371,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Lebanon",
		DomainSuffix:   "LB",
		Code:           961,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Lesotho",
		DomainSuffix:   "LS",
		Code:           266,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Liberia",
		DomainSuffix:   "LR",
		Code:           231,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Libya",
		DomainSuffix:   "LY",
		Code:           218,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Liechtenstein",
		DomainSuffix:   "LI",
		Code:           423,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Lithuania",
		DomainSuffix:   "LT",
		Code:           370,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Luxembourg",
		DomainSuffix:   "LU",
		Code:           352,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Macao",
		DomainSuffix:   "MO",
		Code:           853,
		TimeDifference: 0,
	},
	{
		Country:        "Madagascar",
		DomainSuffix:   "MG",
		Code:           261,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Malawi",
		DomainSuffix:   "MW",
		Code:           265,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Malaysia",
		DomainSuffix:   "MY",
		Code:           60,
		TimeDifference: -0.5 * 3600,
	},
	{
		Country:        "Maldives",
		DomainSuffix:   "MV",
		Code:           960,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Mali",
		DomainSuffix:   "ML",
		Code:           223,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Malta",
		DomainSuffix:   "MT",
		Code:           356,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Mariana Is ",
		DomainSuffix:   "",
		Code:           1670,
		TimeDifference: +1 * 3600,
	},
	{
		Country:        "Martinique",
		DomainSuffix:   "",
		Code:           596,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Mauritius",
		DomainSuffix:   "MU",
		Code:           230,
		TimeDifference: -4 * 3600,
	},
	{
		Country:        "Mexico",
		DomainSuffix:   "MX",
		Code:           52,
		TimeDifference: -15 * 3600,
	},
	{
		Country:        "Moldova Republic of ",
		DomainSuffix:   "MD",
		Code:           373,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Monaco",
		DomainSuffix:   "MC",
		Code:           377,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Mongolia",
		DomainSuffix:   "MN",
		Code:           976,
		TimeDifference: 0,
	},
	{
		Country:        "Montserrat Is ",
		DomainSuffix:   "MS",
		Code:           1664,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Morocco",
		DomainSuffix:   "MA",
		Code:           212,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Mozambique",
		DomainSuffix:   "MZ",
		Code:           258,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Namibia",
		DomainSuffix:   "NA",
		Code:           264,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Nauru",
		DomainSuffix:   "NR",
		Code:           674,
		TimeDifference: +4 * 3600,
	},
	{
		Country:        "Nepal",
		DomainSuffix:   "NP",
		Code:           977,
		TimeDifference: -2.3 * 3600,
	},
	{
		Country:        "Netheriands Antilles ",
		DomainSuffix:   " ",
		Code:           599,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Netherlands",
		DomainSuffix:   "NL",
		Code:           31,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "New Zealand ",
		DomainSuffix:   "NZ",
		Code:           64,
		TimeDifference: +4 * 3600,
	},
	{
		Country:        "Nicaragua",
		DomainSuffix:   "NI",
		Code:           505,
		TimeDifference: -14 * 3600,
	},
	{
		Country:        "Niger",
		DomainSuffix:   "NE",
		Code:           227,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Nigeria",
		DomainSuffix:   "NG",
		Code:           234,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "North Korea ",
		DomainSuffix:   "KP",
		Code:           850,
		TimeDifference: +1 * 3600,
	},
	{
		Country:        "Norway",
		DomainSuffix:   "NO",
		Code:           47,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Oman",
		DomainSuffix:   "OM",
		Code:           968,
		TimeDifference: -4 * 3600,
	},
	{
		Country:        "Pakistan",
		DomainSuffix:   "PK",
		Code:           92,
		TimeDifference: -2.3 * 3600,
	},
	{
		Country:        "Panama",
		DomainSuffix:   "PA",
		Code:           507,
		TimeDifference: -13 * 3600,
	},
	{
		Country:        "Papua New Cuinea ",
		DomainSuffix:   "PG",
		Code:           675,
		TimeDifference: +2 * 3600,
	},
	{
		Country:        "Paraguay",
		DomainSuffix:   "PY",
		Code:           595,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Peru",
		DomainSuffix:   "PE",
		Code:           51,
		TimeDifference: -13 * 3600,
	},
	{
		Country:        "Philippines",
		DomainSuffix:   "PH",
		Code:           63,
		TimeDifference: 0,
	},
	{
		Country:        "Poland",
		DomainSuffix:   "PL",
		Code:           48,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "French Polynesia ",
		DomainSuffix:   "PF",
		Code:           689,
		TimeDifference: +3 * 3600,
	},
	{
		Country:        "Portugal",
		DomainSuffix:   "PT",
		Code:           351,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Puerto Rico ",
		DomainSuffix:   "PR",
		Code:           1787,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Qatar",
		DomainSuffix:   "QA",
		Code:           974,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Reunion",
		DomainSuffix:   "",
		Code:           262,
		TimeDifference: -4 * 3600,
	},
	{
		Country:        "Romania",
		DomainSuffix:   "RO",
		Code:           40,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Russia",
		DomainSuffix:   "RU",
		Code:           7,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Saint Lueia ",
		DomainSuffix:   "LC",
		Code:           1758,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Saint Vincent ",
		DomainSuffix:   "VC",
		Code:           1784,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Samoa Eastern ",
		DomainSuffix:   "",
		Code:           684,
		TimeDifference: -19 * 3600,
	},
	{
		Country:        "Samoa Western ",
		DomainSuffix:   "",
		Code:           685,
		TimeDifference: -19 * 3600,
	},
	{
		Country:        "San Marino ",
		DomainSuffix:   "SM",
		Code:           378,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Sao Tome and Principe ",
		DomainSuffix:   "ST",
		Code:           239,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Saudi Arabia ",
		DomainSuffix:   "SA",
		Code:           966,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Senegal",
		DomainSuffix:   "SN",
		Code:           221,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Seychelles",
		DomainSuffix:   "SC",
		Code:           248,
		TimeDifference: -4 * 3600,
	},
	{
		Country:        "Sierra Leone ",
		DomainSuffix:   "SL",
		Code:           232,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Singapore",
		DomainSuffix:   "SG",
		Code:           65,
		TimeDifference: +0.3 * 3600,
	},
	{
		Country:        "Slovakia",
		DomainSuffix:   "SK",
		Code:           421,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Slovenia",
		DomainSuffix:   "SI",
		Code:           386,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Solomon Is ",
		DomainSuffix:   "SB",
		Code:           677,
		TimeDifference: +3 * 3600,
	},
	{
		Country:        "Somali",
		DomainSuffix:   "SO",
		Code:           252,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "South Africa ",
		DomainSuffix:   "ZA",
		Code:           27,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Spain",
		DomainSuffix:   "ES",
		Code:           34,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Sri Lanka ",
		DomainSuffix:   "LK",
		Code:           94,
		TimeDifference: 0,
	},
	{
		Country:        "St.Lucia ",
		DomainSuffix:   "LC",
		Code:           1758,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "St.Vincent ",
		DomainSuffix:   "VC",
		Code:           1784,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Sudan",
		DomainSuffix:   "SD",
		Code:           249,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Suriname",
		DomainSuffix:   "SR",
		Code:           597,
		TimeDifference: -11.3 * 3600,
	},
	{
		Country:        "Swaziland",
		DomainSuffix:   "SZ",
		Code:           268,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Sweden",
		DomainSuffix:   "SE",
		Code:           46,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Switzerland",
		DomainSuffix:   "CH",
		Code:           41,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Syria",
		DomainSuffix:   "SY",
		Code:           963,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Taiwan",
		DomainSuffix:   "TW",
		Code:           886,
		TimeDifference: 0 * 3600,
	},
	{
		Country:        "Tajikstan",
		DomainSuffix:   "TJ",
		Code:           992,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Tanzania",
		DomainSuffix:   "TZ",
		Code:           255,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Thailand",
		DomainSuffix:   "TH",
		Code:           66,
		TimeDifference: -1 * 3600,
	},
	{
		Country:        "Togo",
		DomainSuffix:   "TG",
		Code:           228,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "Tonga",
		DomainSuffix:   "TO",
		Code:           676,
		TimeDifference: +4 * 3600,
	},
	{
		Country:        "Trinidad and Tobago ",
		DomainSuffix:   "TT",
		Code:           1809,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Tunisia",
		DomainSuffix:   "TN",
		Code:           216,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Turkey",
		DomainSuffix:   "TR",
		Code:           90,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Turkmenistan",
		DomainSuffix:   "TM",
		Code:           993,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Uganda",
		DomainSuffix:   "UG",
		Code:           256,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Ukraine",
		DomainSuffix:   "UA",
		Code:           380,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "United Arab Emirates ",
		DomainSuffix:   "AE",
		Code:           971,
		TimeDifference: -4 * 3600,
	},
	{
		Country:        "United Kiongdom ",
		DomainSuffix:   "GB",
		Code:           44,
		TimeDifference: -8 * 3600,
	},
	{
		Country:        "United States of America ",
		DomainSuffix:   "US",
		Code:           1,
		TimeDifference: -12 * 3600,
	},
	{
		Country:        "Uruguay",
		DomainSuffix:   "UY",
		Code:           598,
		TimeDifference: -10.3 * 3600,
	},
	{
		Country:        "Uzbekistan",
		DomainSuffix:   "UZ",
		Code:           233,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Venezuela",
		DomainSuffix:   "VE",
		Code:           58,
		TimeDifference: -12.3 * 3600,
	},
	{
		Country:        "Vietnam",
		DomainSuffix:   "VN",
		Code:           84,
		TimeDifference: -1 * 3600,
	},
	{
		Country:        "Yemen",
		DomainSuffix:   "YE",
		Code:           967,
		TimeDifference: -5 * 3600,
	},
	{
		Country:        "Yugoslavia",
		DomainSuffix:   "YU",
		Code:           381,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Zimbabwe",
		DomainSuffix:   "ZW",
		Code:           263,
		TimeDifference: -6 * 3600,
	},
	{
		Country:        "Zaire",
		DomainSuffix:   "ZR",
		Code:           243,
		TimeDifference: -7 * 3600,
	},
	{
		Country:        "Zambia",
		DomainSuffix:   "ZM",
		Code:           260,
		TimeDifference: -6 * 3600,
	},
}

// Index Back to index
func Index() map[int16][]Code {
	return indexes
}

// GetCodeInfo Obtain specific area information based on filter parameters
func GetCodeInfo(filter Filter) (code Code) {
	var codeList = indexes[filter.Code]
	if codeList == nil {
		// Cannot find the area code and directly default to China
		code = indexes[86][0]
		return
	}

	if len(codeList) == 1 {
		code = codeList[0]
	}

	if filter.CountryCode == "" {
		code = codeList[0]
	}

	for _, c := range codeList {
		if c.DomainSuffix == filter.CountryCode {
			code = c
		}
	}

	if code.Code == 0 {
		code = indexes[86][0]
	}
	return
}

// GetLocalTimeByCode According to the code and datetime, calculate the time of the specified datetime in the code area on this server, based on Beijing time
//  for example, the current time in Japan is 2020-08-18 14:00:00, The area code for Japan is 81, Call this method GetLocalTimeByCode(81, 2020-08-18 14:00:00)
//  will get 2020-08-18 13:00:00 Timestamp
//  datetime only accept yyyy-MM-dd hh:ii:ss
func GetLocalTimeByCode(filter Filter, datetime string) (int64, error) {
	var codeInfo = GetCodeInfo(filter)
	var l, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, err
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", datetime, l)
	if err != nil {
		return 0, nil
	}

	return t.Unix() - int64(codeInfo.TimeDifference), nil
}
