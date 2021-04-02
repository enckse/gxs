package internal

import (
	"fmt"
	"strings"
)

type (
	patternBlock struct {
		lines []string
		mode  string
		err   error
	}
	ParserError struct {
		Error     error
		Backtrace []string
	}
	patternAction struct {
		palette    map[string]string
		stitchMode string
		pattern    []string
	}
)

const (
	parserBlockStart = " => {"
	defaultBlock     = ""
	paletteAssign    = " => "
	noColor          = "NONE"
)

func colors() map[string]string {
	floss := make(map[string]string)
	floss["dustyroseultvydk"] = "rgb(171, 2, 73)"
	floss["dustyrosevrylt"] = "rgb(240, 206, 212)"
	floss["shellpinkmedlight"] = "rgb(226, 160, 153)"
	floss["violetverylight"] = "rgb(230, 204, 217)"
	floss["grapeverydark"] = "rgb(87, 36, 51)"
	floss["bluevioletmeddark"] = "rgb(152, 145, 182)"
	floss["bluevioletmedlt"] = "rgb(163, 174, 209)"
	floss["cornflowerbluevylt"] = "rgb(187, 195, 217)"
	floss["cornflowerblumvd"] = "rgb(76, 82, 110)"
	floss["bluegraylight"] = "rgb(199, 202, 215)"
	floss["bluegraymedium"] = "rgb(153, 159, 183)"
	floss["bluegray"] = "rgb(120, 128, 164)"
	floss["blueultraverylight"] = "rgb(219, 236, 245)"
	floss["celadongreenmd"] = "rgb(77, 131, 97)"
	floss["forestgreenlt"] = "rgb(200, 216, 184)"
	floss["mossgreenvylt"] = "rgb(239, 244, 164)"
	floss["mossgreenmdlt"] = "rgb(192, 200, 64)"
	floss["yellowbeigevdk"] = "rgb(167, 124, 73)"
	floss["pewterverylight"] = "rgb(209, 209, 209)"
	floss["pewterlight"] = "rgb(132, 132, 132)"
	floss["lavenderverydark"] = "rgb(131, 91, 139)"
	floss["lavenderdark"] = "rgb(163, 123, 167)"
	floss["lavendermedium"] = "rgb(195, 159, 195)"
	floss["lavenderlight"] = "rgb(227, 203, 227)"
	floss["shellpinkvydk"] = "rgb(136, 62, 67)"
	floss["shellpinklight"] = "rgb(204, 132, 124)"
	floss["shellpinkverylight"] = "rgb(235, 183, 175)"
	floss["shellpinkultvylt"] = "rgb(255, 223, 213)"
	floss["mahoganyvydk"] = "rgb(111, 47, 0)"
	floss["mahoganymed"] = "rgb(179, 95, 43)"
	floss["khakigreendk"] = "rgb(137, 138, 88)"
	floss["khakigreenmd"] = "rgb(166, 167, 93)"
	floss["khakigreenlt"] = "rgb(185, 185, 130)"
	floss["browngrayvydk"] = "rgb(79, 75, 65)"
	floss["browngraymed"] = "rgb(142, 144, 120)"
	floss["browngraylight"] = "rgb(177, 170, 151)"
	floss["browngrayvylt"] = "rgb(235, 234, 231)"
	floss["mochabrownvydk"] = "rgb(75, 60, 42)"
	floss["mochabrownmed"] = "rgb(179, 159, 139)"
	floss["mochabrownvylt"] = "rgb(227, 216, 204)"
	floss["redmedium"] = "rgb(183, 31, 51)"
	floss["antiquevioletmedium"] = "rgb(149, 111, 124)"
	floss["antiquevioletlight"] = "rgb(183, 157, 167)"
	floss["yellowbeigedk"] = "rgb(188, 150, 106)"
	floss["yellowbeigemd"] = "rgb(216, 188, 154)"
	floss["yellowbeigelt"] = "rgb(231, 214, 193)"
	floss["greengraydk"] = "rgb(95, 102, 72)"
	floss["greengraymd"] = "rgb(136, 146, 104)"
	floss["greengray"] = "rgb(156, 164, 130)"
	floss["desertsand"] = "rgb(196, 142, 112)"
	floss["lemon"] = "rgb(253, 237, 84)"
	floss["beavergrayvylt"] = "rgb(230, 232, 232)"
	floss["goldenyellowvylt"] = "rgb(253, 249, 205)"
	floss["rosedark"] = "rgb(86, 74, 74)"
	floss["black"] = "rgb(0, 0, 0)"
	floss["wedgewoodultvydk"] = "rgb(28, 80, 102)"
	floss["babyblueverydark"] = "rgb(53, 102, 139)"
	floss["antiquemauvemddk"] = "rgb(129, 73, 82)"
	floss["antiquemauvemed"] = "rgb(183, 115, 127)"
	floss["pewtergray"] = "rgb(108, 108, 108)"
	floss["steelgraylt"] = "rgb(171, 171, 171)"
	floss["pistachiogrnvydk"] = "rgb(32, 95, 46)"
	floss["pistachiogreenmed"] = "rgb(105, 136, 90)"
	floss["red"] = "rgb(199, 43, 59)"
	floss["babybluedark"] = "rgb(90, 143, 184)"
	floss["roseverydark"] = "rgb(179, 59, 75)"
	floss["violetdark"] = "rgb(99, 54, 102)"
	floss["babybluelight"] = "rgb(184, 210, 230)"
	floss["roselight"] = "rgb(251, 173, 180)"
	floss["salmondark"] = "rgb(227, 109, 109)"
	floss["bluevioletverydark"] = "rgb(92, 84, 120)"
	floss["babybluemedium"] = "rgb(115, 159, 193)"
	floss["apricotmed"] = "rgb(255, 131, 111)"
	floss["apricot"] = "rgb(252, 171, 152)"
	floss["huntergreendk"] = "rgb(27, 89, 21)"
	floss["huntergreen"] = "rgb(64, 106, 58)"
	floss["yellowgreenmed"] = "rgb(113, 147, 92)"
	floss["yellowgreenlt"] = "rgb(204, 217, 177)"
	floss["rose"] = "rgb(238, 84, 110)"
	floss["dustyroseultradark"] = "rgb(188, 67, 101)"
	floss["dustyroselight"] = "rgb(228, 166, 172)"
	floss["navyblue"] = "rgb(37, 59, 115)"
	floss["pinegreendk"] = "rgb(94, 107, 71)"
	floss["pinegreenmd"] = "rgb(114, 130, 86)"
	floss["pinegreen"] = "rgb(131, 151, 95)"
	floss["blackbrown"] = "rgb(30, 17, 8)"
	floss["bluevioletmedium"] = "rgb(173, 167, 199)"
	floss["bluevioletlight"] = "rgb(183, 191, 221)"
	floss["salmonverydark"] = "rgb(191, 45, 45)"
	floss["coraldark"] = "rgb(210, 16, 53)"
	floss["coralmedium"] = "rgb(224, 72, 72)"
	floss["coral"] = "rgb(233, 106, 103)"
	floss["corallight"] = "rgb(253, 156, 151)"
	floss["peach"] = "rgb(254, 215, 204)"
	floss["terracottadark"] = "rgb(152, 68, 54)"
	floss["terracottamed"] = "rgb(197, 106, 91)"
	floss["plumlight"] = "rgb(197, 73, 137)"
	floss["plumverylight"] = "rgb(234, 156, 196)"
	floss["plumultralight"] = "rgb(244, 174, 213)"
	floss["pistachiogreendk"] = "rgb(97, 122, 82)"
	floss["pistachiogreenlt"] = "rgb(166, 194, 152)"
	floss["mauveverydark"] = "rgb(136, 21, 49)"
	floss["mauve"] = "rgb(201, 107, 112)"
	floss["mauvemedium"] = "rgb(231, 169, 172)"
	floss["mauvelight"] = "rgb(251, 191, 194)"
	floss["pistachiogreenvylt"] = "rgb(215, 237, 204)"
	floss["mustardmedium"] = "rgb(184, 157, 100)"
	floss["melondark"] = "rgb(255, 121, 146)"
	floss["melonmedium"] = "rgb(255, 173, 188)"
	floss["melonlight"] = "rgb(255, 203, 213)"
	floss["mustard"] = "rgb(191, 166, 113)"
	floss["salmonmedium"] = "rgb(241, 135, 135)"
	floss["salmonverylight"] = "rgb(255, 226, 226)"
	floss["dustyrosemedvylt"] = "rgb(255, 189, 189)"
	floss["mustardlt"] = "rgb(204, 183, 132)"
	floss["shellpinkdark"] = "rgb(161, 75, 81)"
	floss["shellpinkmed"] = "rgb(188, 108, 100)"
	floss["antiquemauvedark"] = "rgb(155, 91, 102)"
	floss["antiquemauvelight"] = "rgb(219, 169, 178)"
	floss["dustyroseverydark"] = "rgb(218, 103, 131)"
	floss["dustyrose"] = "rgb(232, 135, 155)"
	floss["antiquevioletdark"] = "rgb(120, 87, 98)"
	floss["antiquevioletvylt"] = "rgb(215, 203, 211)"
	floss["bluevioletdark"] = "rgb(119, 107, 152)"
	floss["bluevioletvylt"] = "rgb(211, 215, 237)"
	floss["antiqueblueverydk"] = "rgb(56, 76, 94)"
	floss["antiqueblueverylt"] = "rgb(199, 209, 219)"
	floss["antiqueblueultvylt"] = "rgb(219, 226, 233)"
	floss["babyblue"] = "rgb(147, 180, 206)"
	floss["babyblueultvylt"] = "rgb(238, 252, 252)"
	floss["wedgewoodmed"] = "rgb(62, 133, 162)"
	floss["skybluelight"] = "rgb(172, 216, 226)"
	floss["peacockbluevydk"] = "rgb(52, 127, 140)"
	floss["peacockbluelight"] = "rgb(153, 207, 217)"
	floss["graygreendark"] = "rgb(101, 127, 127)"
	floss["tawnyvylight"] = "rgb(255, 238, 227)"
	floss["terracottaultvylt"] = "rgb(244, 187, 169)"
	floss["desertsandvydk"] = "rgb(160, 108, 80)"
	floss["desertsanddark"] = "rgb(182, 117, 82)"
	floss["desertsandvylt"] = "rgb(243, 225, 215)"
	floss["mahoganylight"] = "rgb(207, 121, 57)"
	floss["terracottavydk"] = "rgb(134, 48, 34)"
	floss["terracottalight"] = "rgb(217, 137, 120)"
	floss["rosewoodultvylt"] = "rgb(248, 202, 200)"
	floss["mochabrowndk"] = "rgb(107, 87, 67)"
	floss["mochabrownlt"] = "rgb(210, 188, 166)"
	floss["browngraydark"] = "rgb(98, 93, 80)"
	floss["beigegrayultdk"] = "rgb(127, 106, 85)"
	floss["pewtergrayvydk"] = "rgb(66, 66, 66)"
	floss["melonverydark"] = "rgb(231, 73, 103)"
	floss["antiquemauvevydk"] = "rgb(113, 65, 73)"
	floss["mauvedark"] = "rgb(171, 51, 87)"
	floss["cyclamenpinkdark"] = "rgb(224, 40, 118)"
	floss["cyclamenpink"] = "rgb(243, 71, 139)"
	floss["cyclamenpinklight"] = "rgb(255, 140, 174)"
	floss["cornflowerblue"] = "rgb(96, 103, 140)"
	floss["turquoiseultvydk"] = "rgb(54, 105, 112)"
	floss["turquoisevydark"] = "rgb(63, 124, 133)"
	floss["turquoisedark"] = "rgb(72, 142, 154)"
	floss["turquoiseverylight"] = "rgb(188, 227, 230)"
	floss["seagreenvydk"] = "rgb(47, 140, 132)"
	floss["bluegreenlt"] = "rgb(178, 212, 189)"
	floss["aquamarine"] = "rgb(80, 139, 125)"
	floss["celadongreendk"] = "rgb(71, 119, 89)"
	floss["celadongreen"] = "rgb(101, 165, 125)"
	floss["celadongreenlt"] = "rgb(153, 195, 170)"
	floss["emeraldgrnultvdk"] = "rgb(17, 90, 59)"
	floss["mossgreenlt"] = "rgb(224, 232, 104)"
	floss["strawdark"] = "rgb(223, 182, 95)"
	floss["straw"] = "rgb(243, 206, 117)"
	floss["strawlight"] = "rgb(246, 220, 152)"
	floss["yellowultrapale"] = "rgb(255, 253, 227)"
	floss["apricotlight"] = "rgb(254, 205, 194)"
	floss["pumpkinpale"] = "rgb(253, 189, 150)"
	floss["goldenbrown"] = "rgb(173, 114, 57)"
	floss["goldenbrownpale"] = "rgb(247, 187, 119)"
	floss["hazelnutbrown"] = "rgb(183, 139, 97)"
	floss["oldgoldvydark"] = "rgb(169, 130, 4)"
	floss["terracotta"] = "rgb(185, 85, 68)"
	floss["raspberrydark"] = "rgb(179, 47, 72)"
	floss["raspberrymedium"] = "rgb(219, 85, 110)"
	floss["raspberrylight"] = "rgb(234, 134, 153)"
	floss["grapedark"] = "rgb(114, 55, 93)"
	floss["grapemedium"] = "rgb(148, 96, 131)"
	floss["grapelight"] = "rgb(186, 145, 170)"
	floss["lavenderultradark"] = "rgb(108, 58, 110)"
	floss["lavenderbluedark"] = "rgb(92, 114, 148)"
	floss["lavenderbluemed"] = "rgb(123, 142, 171)"
	floss["lavenderbluelight"] = "rgb(176, 192, 218)"
	floss["babybluepale"] = "rgb(205, 223, 237)"
	floss["wedgewoodvrydk"] = "rgb(50, 102, 124)"
	floss["electricblue"] = "rgb(20, 170, 208)"
	floss["turquoisebrightdark"] = "rgb(18, 174, 186)"
	floss["turquoisebrightmed"] = "rgb(4, 196, 202)"
	floss["turquoisebrightlight"] = "rgb(6, 227, 230)"
	floss["tealgreendark"] = "rgb(52, 125, 117)"
	floss["tealgreenmed"] = "rgb(85, 147, 146)"
	floss["tealgreenlight"] = "rgb(82, 179, 164)"
	floss["greenbrightdk"] = "rgb(55, 132, 119)"
	floss["greenbrightlt"] = "rgb(73, 179, 161)"
	floss["strawverydark"] = "rgb(205, 157, 55)"
	floss["autumngolddk"] = "rgb(242, 151, 70)"
	floss["autumngoldmed"] = "rgb(242, 175, 104)"
	floss["autumngoldlt"] = "rgb(250, 211, 150)"
	floss["mahoganyultvylt"] = "rgb(255, 211, 181)"
	floss["rosewooddark"] = "rgb(104, 37, 26)"
	floss["rosewoodmed"] = "rgb(150, 74, 63)"
	floss["rosewoodlight"] = "rgb(186, 139, 124)"
	floss["cocoa"] = "rgb(125, 93, 87)"
	floss["cocoalight"] = "rgb(166, 136, 129)"
	floss["mochabeigedark"] = "rgb(138, 110, 78)"
	floss["mochabeigemed"] = "rgb(164, 131, 92)"
	floss["mochabeigelight"] = "rgb(203, 182, 156)"
	floss["winterwhite"] = "rgb(249, 247, 241)"
	floss["mochabrnultvylt"] = "rgb(250, 246, 240)"
	floss["mahoganydark"] = "rgb(143, 67, 15)"
	floss["mahoganyvylt"] = "rgb(247, 167, 119)"
	floss["desertsandmed"] = "rgb(187, 129, 97)"
	floss["pewtergraydark"] = "rgb(86, 86, 86)"
	floss["steelgraydk"] = "rgb(140, 140, 140)"
	floss["pearlgray"] = "rgb(211, 211, 214)"
	floss["hazelnutbrowndk"] = "rgb(160, 112, 66)"
	floss["hazelnutbrownlt"] = "rgb(198, 159, 123)"
	floss["brownmed"] = "rgb(122, 69, 31)"
	floss["brownlight"] = "rgb(152, 94, 51)"
	floss["brownverylight"] = "rgb(184, 119, 72)"
	floss["tan"] = "rgb(203, 144, 81)"
	floss["tanlight"] = "rgb(228, 187, 142)"
	floss["lemondark"] = "rgb(255, 214, 0)"
	floss["lemonlight"] = "rgb(255, 251, 139)"
	floss["shellgraydark"] = "rgb(145, 123, 115)"
	floss["shellgraymed"] = "rgb(192, 179, 174)"
	floss["shellgraylight"] = "rgb(215, 206, 203)"
	floss["avocadogreen"] = "rgb(114, 132, 60)"
	floss["avocadogrnlt"] = "rgb(148, 171, 79)"
	floss["avocadogrnvlt"] = "rgb(174, 191, 121)"
	floss["avocadogrnult"] = "rgb(216, 228, 152)"
	floss["reddark"] = "rgb(167, 19, 43)"
	floss["bluegreenvydk"] = "rgb(4, 77, 51)"
	floss["bluegreendark"] = "rgb(57, 111, 82)"
	floss["bluegreen"] = "rgb(91, 144, 113)"
	floss["bluegreenmed"] = "rgb(123, 172, 148)"
	floss["bluegreenvylt"] = "rgb(196, 222, 204)"
	floss["jadegreen"] = "rgb(51, 131, 98)"
	floss["wedgewooddark"] = "rgb(59, 118, 143)"
	floss["wedgewoodlight"] = "rgb(79, 147, 167)"
	floss["skyblue"] = "rgb(126, 177, 200)"
	floss["ferngreendark"] = "rgb(102, 109, 79)"
	floss["ferngreen"] = "rgb(150, 158, 126)"
	floss["ferngreenlt"] = "rgb(171, 177, 151)"
	floss["ferngreenvylt"] = "rgb(196, 205, 172)"
	floss["ashgrayvylt"] = "rgb(99, 100, 88)"
	floss["beigebrownultvylt"] = "rgb(242, 227, 206)"
	floss["violetverydark"] = "rgb(92, 24, 78)"
	floss["violetmedium"] = "rgb(128, 58, 107)"
	floss["violet"] = "rgb(163, 99, 139)"
	floss["violetlight"] = "rgb(219, 179, 203)"
	floss["celadongreenvd"] = "rgb(44, 106, 69)"
	floss["jademedium"] = "rgb(83, 151, 106)"
	floss["jadelight"] = "rgb(143, 192, 152)"
	floss["jadeverylight"] = "rgb(167, 205, 175)"
	floss["mossgreendk"] = "rgb(136, 141, 51)"
	floss["mossgreen"] = "rgb(167, 174, 56)"
	floss["turquoise"] = "rgb(91, 163, 179)"
	floss["turquoiselight"] = "rgb(144, 195, 204)"
	floss["cranberryverydark"] = "rgb(205, 47, 99)"
	floss["cranberrydark"] = "rgb(209, 40, 106)"
	floss["cranberrymedium"] = "rgb(226, 72, 116)"
	floss["cranberry"] = "rgb(255, 164, 190)"
	floss["cranberrylight"] = "rgb(255, 176, 190)"
	floss["cranberryverylight"] = "rgb(255, 192, 205)"
	floss["orangeredbright"] = "rgb(250, 50, 3)"
	floss["burntorangebright"] = "rgb(253, 93, 53)"
	floss["drabbrowndk"] = "rgb(121, 96, 71)"
	floss["drabbrown"] = "rgb(150, 118, 86)"
	floss["drabbrownlt"] = "rgb(188, 154, 120)"
	floss["drabbrownvlt"] = "rgb(220, 196, 170)"
	floss["desertsandultvydk"] = "rgb(135, 85, 57)"
	floss["beigegrayvydk"] = "rgb(133, 123, 97)"
	floss["beigegraydark"] = "rgb(164, 152, 120)"
	floss["beigegraymed"] = "rgb(221, 216, 203)"
	floss["beavergrayvydk"] = "rgb(110, 101, 92)"
	floss["beavergraydk"] = "rgb(135, 125, 115)"
	floss["beavergraymed"] = "rgb(176, 166, 156)"
	floss["beavergraylt"] = "rgb(188, 180, 172)"
	floss["brightred"] = "rgb(227, 29, 66)"
	floss["oldgoldlt"] = "rgb(229, 206, 151)"
	floss["oldgoldvylt"] = "rgb(245, 236, 203)"
	floss["oldgolddark"] = "rgb(188, 141, 14)"
	floss["green"] = "rgb(5, 101, 23)"
	floss["greenbright"] = "rgb(7, 115, 27)"
	floss["greenlight"] = "rgb(63, 143, 41)"
	floss["kellygreen"] = "rgb(71, 167, 47)"
	floss["chartreuse"] = "rgb(123, 181, 71)"
	floss["chartreusebright"] = "rgb(158, 207, 52)"
	floss["cream"] = "rgb(255, 251, 239)"
	floss["plum"] = "rgb(156, 36, 98)"
	floss["orangespicedark"] = "rgb(229, 92, 31)"
	floss["orangespicemed"] = "rgb(242, 120, 66)"
	floss["orangespicelight"] = "rgb(247, 151, 111)"
	floss["topazmedlt"] = "rgb(255, 200, 64)"
	floss["topazlight"] = "rgb(253, 215, 85)"
	floss["topazvylt"] = "rgb(255, 241, 175)"
	floss["topaz"] = "rgb(228, 180, 104)"
	floss["oldgoldmedium"] = "rgb(208, 165, 62)"
	floss["olivegreenvdk"] = "rgb(130, 123, 48)"
	floss["olivegreendk"] = "rgb(147, 139, 55)"
	floss["olivegreen"] = "rgb(148, 140, 54)"
	floss["olivegreenmd"] = "rgb(188, 179, 76)"
	floss["olivegreenlt"] = "rgb(199, 192, 119)"
	floss["tanverylight"] = "rgb(236, 204, 158)"
	floss["tanultvylt"] = "rgb(248, 228, 200)"
	floss["tangerine"] = "rgb(255, 139, 0)"
	floss["tangerinemed"] = "rgb(255, 163, 43)"
	floss["tangerinelight"] = "rgb(255, 191, 87)"
	floss["yellowmed"] = "rgb(254, 211, 118)"
	floss["yellowpale"] = "rgb(255, 231, 147)"
	floss["yellowpalelight"] = "rgb(255, 233, 173)"
	floss["offwhite"] = "rgb(252, 252, 238)"
	floss["peacockbluevylt"] = "rgb(229, 252, 253)"
	floss["peachlight"] = "rgb(247, 203, 191)"
	floss["terracottavylt"] = "rgb(238, 170, 155)"
	floss["salmon"] = "rgb(245, 173, 173)"
	floss["salmonlight"] = "rgb(255, 201, 201)"
	floss["pearlgrayvylt"] = "rgb(236, 236, 236)"
	floss["yellowgreenvylt"] = "rgb(228, 236, 212)"
	floss["babyblueverylight"] = "rgb(217, 235, 241)"
	floss["pinkmedium"] = "rgb(252, 176, 185)"
	floss["raspberryverydark"] = "rgb(145, 53, 70)"
	floss["antiquemauvevylt"] = "rgb(223, 179, 187)"
	floss["cocoadark"] = "rgb(98, 75, 69)"
	floss["topazultravydk"] = "rgb(148, 99, 26)"
	floss["topazverydark"] = "rgb(162, 109, 32)"
	floss["topazdark"] = "rgb(174, 119, 32)"
	floss["topazmedium"] = "rgb(206, 145, 36)"
	floss["cornflowerbluevd"] = "rgb(70, 69, 99)"
	floss["cornflowerbluedark"] = "rgb(85, 91, 123)"
	floss["cornflowerbluemed"] = "rgb(112, 125, 162)"
	floss["cornflowerbluelight"] = "rgb(143, 156, 193)"
	floss["royalbluedark"] = "rgb(17, 65, 109)"
	floss["royalblue"] = "rgb(19, 71, 125)"
	floss["delftbluedark"] = "rgb(70, 106, 142)"
	floss["delftbluemedium"] = "rgb(116, 142, 182)"
	floss["delftbluepale"] = "rgb(192, 204, 222)"
	floss["coffeebrowndk"] = "rgb(101, 57, 25)"
	floss["babyblueultvydk"] = "rgb(44, 89, 124)"
	floss["peacockbluedark"] = "rgb(61, 149, 165)"
	floss["peacockblue"] = "rgb(100, 171, 186)"
	floss["delftblue"] = "rgb(148, 168, 198)"
	floss["bluelight"] = "rgb(161, 194, 215)"
	floss["garnetdark"] = "rgb(123, 0, 27)"
	floss["garnetmedium"] = "rgb(135, 7, 31)"
	floss["garnet"] = "rgb(151, 11, 35)"
	floss["coralredverydark"] = "rgb(187, 5, 31)"
	floss["babypink"] = "rgb(255, 223, 217)"
	floss["babypinklight"] = "rgb(255, 238, 235)"
	floss["royalblueverydark"] = "rgb(14, 54, 92)"
	floss["beigegraylight"] = "rgb(231, 226, 211)"
	floss["navybluedark"] = "rgb(33, 48, 99)"
	floss["blueverydark"] = "rgb(57, 105, 135)"
	floss["bluedark"] = "rgb(71, 129, 165)"
	floss["bluemedium"] = "rgb(107, 158, 191)"
	floss["blueverylight"] = "rgb(189, 221, 237)"
	floss["skybluevylt"] = "rgb(197, 232, 237)"
	floss["goldenolivevydk"] = "rgb(126, 107, 66)"
	floss["goldenolivedk"] = "rgb(141, 120, 75)"
	floss["goldenolivemd"] = "rgb(170, 143, 86)"
	floss["goldenolive"] = "rgb(189, 155, 81)"
	floss["goldenolivelt"] = "rgb(200, 171, 108)"
	floss["goldenolivevylt"] = "rgb(219, 190, 127)"
	floss["beigebrownvydk"] = "rgb(89, 73, 55)"
	floss["beigebrowndk"] = "rgb(103, 85, 65)"
	floss["beigebrownmed"] = "rgb(154, 124, 92)"
	floss["beigebrownlt"] = "rgb(182, 155, 126)"
	floss["beigebrownvylt"] = "rgb(209, 186, 161)"
	floss["beavergrayultdk"] = "rgb(72, 72, 72)"
	floss["hazelnutbrownvdk"] = "rgb(131, 94, 57)"
	floss["pistachiogrnultvd"] = "rgb(23, 73, 35)"
	floss["carnationdark"] = "rgb(255, 87, 115)"
	floss["carnationmedium"] = "rgb(255, 121, 140)"
	floss["carnationlight"] = "rgb(252, 144, 162)"
	floss["carnationverylight"] = "rgb(255, 178, 187)"
	floss["huntergreenvydk"] = "rgb(27, 83, 0)"
	floss["coffeebrownvydk"] = "rgb(73, 42, 19)"
	floss["rosemedium"] = "rgb(242, 118, 136)"
	floss["burntorangedark"] = "rgb(209, 88, 7)"
	floss["garnetverydark"] = "rgb(130, 38, 55)"
	floss["parrotgreenvdk"] = "rgb(85, 120, 34)"
	floss["parrotgreendk"] = "rgb(98, 138, 40)"
	floss["parrotgreenmd"] = "rgb(127, 179, 53)"
	floss["parrotgreenlt"] = "rgb(199, 230, 102)"
	floss["emeraldgreenvydk"] = "rgb(21, 111, 73)"
	floss["emeraldgreendark"] = "rgb(24, 126, 86)"
	floss["emeraldgreenmed"] = "rgb(24, 144, 101)"
	floss["emeraldgreenlt"] = "rgb(27, 157, 107)"
	floss["nilegreenmed"] = "rgb(109, 171, 119)"
	floss["plumdark"] = "rgb(130, 0, 67)"
	floss["plummedium"] = "rgb(155, 19, 89)"
	floss["redcopperdark"] = "rgb(130, 52, 10)"
	floss["redcopper"] = "rgb(166, 69, 16)"
	floss["coppermed"] = "rgb(172, 84, 20)"
	floss["copper"] = "rgb(198, 98, 24)"
	floss["copperlight"] = "rgb(226, 115, 35)"
	floss["graygreenvydark"] = "rgb(86, 106, 106)"
	floss["graygreenmed"] = "rgb(152, 174, 174)"
	floss["graygreenlight"] = "rgb(189, 203, 203)"
	floss["graygreenvylt"] = "rgb(221, 227, 227)"
	floss["antiquebluedark"] = "rgb(69, 92, 113)"
	floss["antiquebluemedium"] = "rgb(106, 133, 158)"
	floss["antiquebluelight"] = "rgb(162, 181, 198)"
	floss["avocadogrnblack"] = "rgb(49, 57, 25)"
	floss["avocadogreendk"] = "rgb(66, 77, 33)"
	floss["avocadogrnvdk"] = "rgb(76, 88, 38)"
	floss["avocadogreenmd"] = "rgb(98, 113, 51)"
	floss["coffeebrownultdk"] = "rgb(54, 31, 14)"
	floss["navyblueverydark"] = "rgb(27, 40, 83)"
	floss["greenbrightmd"] = "rgb(61, 147, 132)"
	floss["tawny"] = "rgb(251, 213, 187)"
	floss["burntorangemed"] = "rgb(235, 99, 7)"
	floss["burntorange"] = "rgb(255, 123, 77)"
	floss["peachverylight"] = "rgb(254, 231, 218)"
	floss["desertsandlight"] = "rgb(238, 211, 196)"
	floss["tawnylight"] = "rgb(255, 226, 207)"
	floss["nilegreen"] = "rgb(136, 186, 145)"
	floss["nilegreenlight"] = "rgb(162, 214, 173)"
	floss["geranium"] = "rgb(255, 145, 145)"
	floss["geraniumpale"] = "rgb(253, 181, 181)"
	floss["seagreendark"] = "rgb(62, 182, 161)"
	floss["seagreenmed"] = "rgb(89, 199, 180)"
	floss["dustyrosedark"] = "rgb(207, 115, 115)"
	floss["dustyrosemedium"] = "rgb(230, 138, 138)"
	floss["dustyroseultvylt"] = "rgb(255, 215, 215)"
	floss["seagreenlight"] = "rgb(169, 226, 216)"
	floss["jadeultravylt"] = "rgb(185, 215, 192)"
	floss["apricotverylight"] = "rgb(255, 222, 213)"
	floss["pumpkinlight"] = "rgb(247, 139, 19)"
	floss["pumpkin"] = "rgb(246, 127, 0)"
	floss["canarydeep"] = "rgb(255, 181, 21)"
	floss["canarybright"] = "rgb(255, 227, 0)"
	floss["goldenbrowndk"] = "rgb(145, 79, 18)"
	floss["goldenbrownmed"] = "rgb(194, 129, 66)"
	floss["goldenbrownlight"] = "rgb(220, 156, 86)"
	floss["forestgreenvydk"] = "rgb(64, 82, 48)"
	floss["forestgreendk"] = "rgb(88, 113, 65)"
	floss["forestgreenmed"] = "rgb(115, 139, 91)"
	floss["forestgreen"] = "rgb(141, 166, 117)"
	floss["aquamarinedk"] = "rgb(71, 123, 110)"
	floss["aquamarinelt"] = "rgb(111, 174, 159)"
	floss["aquamarinevylt"] = "rgb(144, 192, 180)"
	floss["electricbluedark"] = "rgb(38, 150, 182)"
	floss["electricbluemedium"] = "rgb(48, 194, 236)"
	floss["snowwhite"] = "rgb(255, 255, 255)"
	floss["ecru"] = "rgb(240, 234, 218)"
	floss["white"] = "rgb(252, 251, 248)"
	return floss
}

func next(stream []string) (patternBlock, int) {
	idx := 0
	inBlock := false
	block := patternBlock{mode: defaultBlock}
	for idx < len(stream) {
		line := strings.TrimSpace(stream[idx])
		if strings.HasPrefix(line, "#") {
			line = ""
		}
		if len(line) > 0 {
			if inBlock {
				if line == "}" {
					if len(block.lines) == 0 {
						return patternBlock{err: fmt.Errorf("empty block found")}, 0
					}
					return block, idx + 1
				}
				block.lines = append(block.lines, line)
			} else {
				if strings.HasSuffix(line, parserBlockStart) {
					inBlock = true
					modeSection := strings.Split(line, parserBlockStart)
					if len(modeSection) != 2 {
						return patternBlock{err: fmt.Errorf("invalid start block")}, 0
					}
					block.mode = modeSection[0]
				} else {
					return patternBlock{err: fmt.Errorf("expected start of block")}, 0
				}
			}
		}
		idx += 1
	}
	if inBlock {
		return patternBlock{err: fmt.Errorf("unclosed block")}, 0
	}
	return block, idx
}

func (b patternBlock) isMatch(is string) bool {
	return len(b.lines) == 1 && b.lines[0] == is
}

func (b patternBlock) toError(message string) *ParserError {
	return &ParserError{Error: fmt.Errorf(message), Backtrace: b.lines}
}

func parseBlocks(blocks []patternBlock) ([]patternAction, *ParserError) {
	var actions []patternAction
	var action patternAction
	colorLookup := colors()
	for _, block := range blocks {
		switch block.mode {
		case "palette":
			action.palette = make(map[string]string)
			for _, line := range block.lines {
				parts := strings.Split(line, paletteAssign)
				if len(parts) != 2 {
					return nil, block.toError("invalid palette assignment")
				}
				char := parts[0]
				color := parts[1]
				if len(char) != 1 {
					return nil, block.toError("only single characters allowed")
				}
				if val, ok := colorLookup[color]; ok {
					color = val
				}
				if _, ok := action.palette[char]; ok {
					return nil, block.toError("character re-used within palette")
				}
				action.palette[char] = color
			}
		case "pattern":
			if len(action.pattern) > 0 {
				return nil, block.toError("pattern not committed")
			}
			action.pattern = block.lines
		case "action":
			if !block.isMatch("commit") {
				return nil, block.toError("unknown action")
			}
			if len(action.pattern) == 0 {
				return nil, block.toError("no pattern")
			}
			switch action.stitchMode {
			case "le", "re", "te", "be", "xs":
				break
			default:
				return nil, block.toError("invalid stitch mode")
			}
			actions = append(actions, action)
			action.pattern = []string{}
			action.stitchMode = ""
		case "mode":
			if len(block.lines) != 1 {
				return nil, block.toError("incorrect stitch mode setting")
			}
			line := block.lines[0]
			if action.stitchMode != "" {
				if action.stitchMode != line {
					return nil, block.toError("stitching not committed")
				}
			}
			action.stitchMode = line
		default:
			return nil, block.toError("unknown mode in block")
		}
	}
	if len(action.pattern) != 0 {
		return nil, &ParserError{Error: fmt.Errorf("uncommitted pattern")}
	}
	return actions, nil
}

func (a patternAction) toPatternError(message string) *ParserError {
	return &ParserError{Error: fmt.Errorf(message), Backtrace: a.pattern}
}

func buildPattern(actions []patternAction) (Pattern, *ParserError) {
	var entries []Entry
	var maxSize = -1
	for _, action := range actions {
		tracking := make(map[string]map[string][]string)
		for height, line := range action.pattern {
			if height > maxSize {
				maxSize = height
			}
			for width, chr := range line {
				if width > maxSize {
					maxSize = width
				}
				symbol := fmt.Sprintf("%c", chr)
				if color, ok := action.palette[symbol]; ok {
					if _, hasColor := tracking[color]; !hasColor {
						tracking[color] = make(map[string][]string)
					}
					curColor := tracking[color]
					if _, hasMode := curColor[action.stitchMode]; !hasMode {
						curColor[action.stitchMode] = []string{}
					}
					modeSet := curColor[action.stitchMode]
					offset := fmt.Sprintf("%dx%d", width+1, height+1)
					modeSet = append(modeSet, offset)
					curColor[action.stitchMode] = modeSet
					tracking[color] = curColor
				} else {
					return Pattern{}, action.toPatternError("symbol unknown")
				}
			}
		}
		for color, modes := range tracking {
			if color == noColor {
				continue
			}
			for mode, cells := range modes {
				entry := Entry{Cells: cells, Mode: mode, Color: color}
				entries = append(entries, entry)
			}
		}
	}
	pattern, err := NewPattern(maxSize + 1)
	if err != nil {
		return pattern, &ParserError{Error: err}
	}
	pattern.Entries = entries
	return pattern, nil
}

func Parse(b []byte) (Pattern, *ParserError) {
	lines := strings.Split(string(b), "\n")
	var blocks []patternBlock
	var pattern Pattern
	for {
		block, read := next(lines)
		if block.err != nil {
			return pattern, &ParserError{Error: block.err, Backtrace: lines}
		}
		if read == 0 {
			break
		}
		if block.mode != defaultBlock {
			blocks = append(blocks, block)
		}
		lines = lines[read:]
	}
	if len(blocks) == 0 {
		return pattern, &ParserError{Error: fmt.Errorf("no blocks found")}
	}
	actions, pErr := parseBlocks(blocks)
	if pErr != nil {
		return pattern, pErr
	}
	if len(actions) == 0 {
		return pattern, &ParserError{Error: fmt.Errorf("no actions, nothing committed?")}
	}
	return buildPattern(actions)
}
