package main

// maze object ids (names matching names in IDA)
const (
	MAZEOBJ_TILE_FLOOR = iota
	MAZEOBJ_TILE_STUN
	MAZEOBJ_WALL_REGULAR	// 2
	MAZEOBJ_WALL_MOVABLE
	MAZEOBJ_WALL_SECRET
	MAZEOBJ_WALL_DESTRUCTABLE
	MAZEOBJ_WALL_RANDOM
	MAZEOBJ_WALL_TRAPCYC1
	MAZEOBJ_WALL_TRAPCYC2
	MAZEOBJ_WALL_TRAPCYC3
	MAZEOBJ_TILE_TRAP1
	MAZEOBJ_TILE_TRAP2
	MAZEOBJ_TILE_TRAP3
	MAZEOBJ_DOOR_HORIZ
	MAZEOBJ_DOOR_VERT
	MAZEOBJ_PLAYERSTART
	MAZEOBJ_EXIT
	MAZEOBJ_EXITTO6
	MAZEOBJ_MONST_GHOST
	MAZEOBJ_MONST_GRUNT
	MAZEOBJ_MONST_DEMON
	MAZEOBJ_MONST_LOBBER
	MAZEOBJ_MONST_SORC
	MAZEOBJ_MONST_AUX_GRUNT
	MAZEOBJ_MONST_DEATH
	MAZEOBJ_MONST_ACID
	MAZEOBJ_MONST_SUPERSORC
	MAZEOBJ_MONST_IT
	MAZEOBJ_GEN_GHOST1
	MAZEOBJ_GEN_GHOST2
	MAZEOBJ_GEN_GHOST3
	MAZEOBJ_GEN_GRUNT1
	MAZEOBJ_GEN_GRUNT2
	MAZEOBJ_GEN_GRUNT3
	MAZEOBJ_GEN_DEMON1
	MAZEOBJ_GEN_DEMON2
	MAZEOBJ_GEN_DEMON3
	MAZEOBJ_GEN_LOBBER1
	MAZEOBJ_GEN_LOBBER2
	MAZEOBJ_GEN_LOBBER3
	MAZEOBJ_GEN_SORC1
	MAZEOBJ_GEN_SORC2
	MAZEOBJ_GEN_SORC3
	MAZEOBJ_GEN_AUX_GRUNT1
	MAZEOBJ_GEN_AUX_GRUNT2
	MAZEOBJ_GEN_AUX_GRUNT3
	MAZEOBJ_TREASURE
	MAZEOBJ_TREASURE_LOCKED
	MAZEOBJ_TREASURE_BAG
	MAZEOBJ_FOOD_DESTRUCTABLE
	MAZEOBJ_FOOD_INVULN
	MAZEOBJ_POT_DESTRUCTABLE
	MAZEOBJ_POT_INVULN
	MAZEOBJ_KEY					// 53
	MAZEOBJ_POWER_INVIS
	MAZEOBJ_POWER_REPULSE
	MAZEOBJ_POWER_REFLECT
	MAZEOBJ_POWER_TRANSPORT
	MAZEOBJ_POWER_SUPERSHOT
	MAZEOBJ_POWER_INVULN
	MAZEOBJ_MONST_DRAGON
	MAZEOBJ_HIDDENPOT
	MAZEOBJ_TRANSPORTER
	MAZEOBJ_FORCEFIELDHUB
	MAZEOBJ_MONST_MUGGER		// 64 - these 2 are add-ins to round out monsters
	MAZEOBJ_MONST_THIEF			//		engine is not coded to place them in mazes
	MAZEOBJ_EXTEND
)

// Six: g1 objects
const (
	G1OBJ_TILE_FLOOR = 0
// 1
	G1OBJ_WALL_REGULAR = 2
	G1OBJ_DOOR_HORIZ = 3
	G1OBJ_DOOR_VERT = 4
	G1OBJ_PLAYERSTART = 5
	G1OBJ_EXIT = 6
	G1OBJ_EXIT4 = 7
	G1OBJ_EXIT8 = 8
	G1OBJ_MONST_GHOST1 = 9
	G1OBJ_MONST_GHOST2 = 10
	G1OBJ_MONST_GHOST3 = 11
	G1OBJ_MONST_GRUNT1 = 12
	G1OBJ_MONST_GRUNT2 = 13
	G1OBJ_MONST_GRUNT3 = 14
	G1OBJ_MONST_DEMON1 = 15
	G1OBJ_MONST_DEMON2 = 16
	G1OBJ_MONST_DEMON3 = 17
	G1OBJ_MONST_LOBBER1 = 18
	G1OBJ_MONST_LOBBER2 = 19
	G1OBJ_MONST_LOBBER3 = 20
	G1OBJ_MONST_SORC1 = 21
	G1OBJ_MONST_SORC2 = 22
	G1OBJ_MONST_SORC3 = 23
	G1OBJ_MONST_DEATH = 24
	G1OBJ_GEN_GHOST1 = 25
	G1OBJ_GEN_GHOST2 = 26
	G1OBJ_GEN_GHOST3 = 27
	G1OBJ_GEN_GRUNT1 = 28
	G1OBJ_GEN_GRUNT2 = 29
	G1OBJ_GEN_GRUNT3 = 30
	G1OBJ_GEN_DEMON1 = 31
	G1OBJ_GEN_DEMON2 = 32
	G1OBJ_GEN_DEMON3 = 33
	G1OBJ_GEN_LOBBER1 = 34
	G1OBJ_GEN_LOBBER2 = 35
	G1OBJ_GEN_LOBBER3 = 36
	G1OBJ_GEN_SORC1 = 37
	G1OBJ_GEN_SORC2 = 38
	G1OBJ_GEN_SORC3 = 39
	G1OBJ_TREASURE = 40
// 41
	G1OBJ_FOOD_DESTRUCTABLE = 42
	G1OBJ_FOOD_INVULN = 43
	G1OBJ_POT_DESTRUCTABLE = 44
	G1OBJ_POT_INVULN = 45
	G1OBJ_INVISIBL = 46
	G1OBJ_X_SPEED = 48
	G1OBJ_X_SHOTPW = 50
	G1OBJ_X_SHTSPD = 51
	G1OBJ_X_ARMOR = 47
	G1OBJ_X_FIGHT = 52
	G1OBJ_X_MAGIC = 49
	G1OBJ_KEY = 53
// 54
// 55
	G1OBJ_WALL_DESTRUCTABLE = 56
	G1OBJ_WALL_TRAP1 = 57
	G1OBJ_TILE_TRAP1 = 58
	G1OBJ_TRANSPORTER = 59
// 60
// 61
//	G1OBJ_TILE_STUN = 62				// G¹ had no stun tile
// 63
	G1OBJ_TREASURE_BAG			= 64
	G1OBJ_MONST_THIEF			= 65
	G1OBJ_EXTEND				= 66
	G1OBJ_WIZARD				= 67
// 68
// 69
	GORO_TEST					= 70
// sanctuary engine ops
	SEOBJ_FAKE_BLK				= 90			// fake item, block movement/ weps			will need xb to change appearance
	SEOBJ_FAKE_BLK_SHT			= 91			//  "  , blocks mv, can be shot out
	SEOBJ_FAKE_PAS				= 92			//  "  , pass over (to block weps, xb needed)
	SEOBJ_FAKE_PAS_SHT			= 93			//  "  , pass over, can be shot out (absorbs shot to hit pts)

	SEOBJ_FLOORNUL				= 98			// draw as null square, color = 0
	SEOBJ_FLOORNODRAW			= 99			// dont draw a floor here - colortile / cust floor are 1st, this applies to: std floor, master floor ovrd
	SEOBJ_FLOOR					= 100
// Gauntlet II items implement in SE
	SEOBJ_STUN					= 101
	SEOBJ_PUSHWAL				= 103
	SEOBJ_SECRTWAL				= 104
// 5 is g2 destruct wall
	SEOBJ_RNDWAL				= 106
	SEOBJ_WAL_TRAPCYC1			= 107		// g2 trap & cycle walls
	SEOBJ_WAL_TRAPCYC2			= 108
	SEOBJ_WAL_TRAPCYC3			= 109
	SEOBJ_TILE_TRAP1			= 110
	SEOBJ_TILE_TRAP2			= 111
	SEOBJ_TILE_TRAP3			= 112
	SEOBJ_DOOR_H				= 113
	SEOBJ_DOOR_V				= 114
// g2 start
// g2 exit
	SEOBJ_EXIT6					= 117
	SEOBJ_G2GHOST				= 118
	SEOBJ_G2GRUNT				= 119
	SEOBJ_G2DEMON				= 120
	SEOBJ_G2LOBER				= 121
	SEOBJ_G2SORC				= 122
	SEOBJ_G2AUXGR				= 123
	SEOBJ_G2DEATH				= 124
	SEOBJ_G2ACID				= 125
	SEOBJ_G2SUPSORC				= 126
	SEOBJ_G2IT					= 127
	SEOBJ_G2GN_GST1				= 128
	SEOBJ_G2GN_GST2				= 129
	SEOBJ_G2GN_GST3				= 130
	SEOBJ_G2GN_GR1				= 131
	SEOBJ_G2GN_GR2				= 132
	SEOBJ_G2GN_GR3				= 133
	SEOBJ_G2GN_DM1				= 134
	SEOBJ_G2GN_DM2				= 135
	SEOBJ_G2GN_DM3				= 136
	SEOBJ_G2GN_LB1				= 137
	SEOBJ_G2GN_LB2				= 138
	SEOBJ_G2GN_LB3				= 139
	SEOBJ_G2GN_SORC1			= 140
	SEOBJ_G2GN_SORC2			= 141
	SEOBJ_G2GN_SORC3			= 142
	SEOBJ_G2GN_AUXGR1			= 143
	SEOBJ_G2GN_AUXGR2			= 144
	SEOBJ_G2GN_AUXGR3			= 145
// g2 treasure
	SEOBJ_TREASURE_LOCKED		= 147
// g2 treasure bag
// g2 food destr
// g2 food inv
// g2 pot
// g2 pot inv
// g2 key				153
// g2 power invis
	SEOBJ_SE_ANKH				= 154
	SEOBJ_POWER_REPULSE			= 155
	SEOBJ_POWER_REFLECT			= 156
	SEOBJ_POWER_TRANSPORT		= 157
	SEOBJ_POWER_SUPERSHOT		= 158
	SEOBJ_POWER_INVULN			= 159
	SEOBJ_MONST_DRAGON			= 160
// g2 hidden potion
// g2 transport
	SEOBJ_FORCEFIELDHUB			= 163
	SEOBJ_MONST_MUGGER			= 164
// g2 thief
// extend - sanctuary engine specific
	SEOBJ_FIRE_STICK			= 166		// 26, 33
	SEOBJ_G2_POISPOT			= 167		// 8, 11
	SEOBJ_G2_POISFUD			= 168		// 1, 11		23, 10
	SEOBJ_G2_QFUD				= 169		// 2, 11
	SEOBJ_KEYRING				= 170		// 29, 11

	SEOBJ_MAPPYBDG				= 171		// 33, 23

	SEOBJ_MAPPYARAD				= 172		// 25, 21
	SEOBJ_MAPPYATV				= 173		// 27, 21
	SEOBJ_MAPPYAPC				= 174		// 29, 21
	SEOBJ_MAPPYAART				= 175		// 31, 21
	SEOBJ_MAPPYASAF				= 176		// 33, 21

	SEOBJ_MAPPYRAD				= 177		// 25, 22
	SEOBJ_MAPPYTV				= 178		// 27, 22
	SEOBJ_MAPPYPC				= 179		// 29, 22
	SEOBJ_MAPPYART				= 180		// 31, 22
	SEOBJ_MAPPYSAF				= 181		// 33, 22

	SEOBJ_MAPPYBELL				= 182		// 33, 21
	SEOBJ_MAPPYBAL				= 183		// 33, 22
	SEOBJ_MAPPYGORO				= 184		// 34, 22

	SEOBJ_DETHGEN3				= 185		// 32, 8
	SEOBJ_DETHGEN2				= 186		// 33, 8
	SEOBJ_DETHGEN1				= 187		// 34, 8

	SEOBJ_WATER_POOL			= 188
	SEOBJ_WATER_TOP				= 189
	SEOBJ_WATER_RT				= 190
	SEOBJ_WATER_COR				= 191
	SEOBJ_SLIME_POOL			= 192
	SEOBJ_SLIME_TOP				= 193
	SEOBJ_SLIME_RT				= 194
	SEOBJ_SLIME_COR				= 195
	SEOBJ_LAVA_POOL				= 196
	SEOBJ_LAVA_TOP				= 197
	SEOBJ_LAVA_RT				= 198
	SEOBJ_LAVA_COR				= 199
	SEOBJ_PULS_FLOR				= 200

)

// animated - also se, telepod, ff bunkers, traps,

var animcyc = []int{
	SEOBJ_FIRE_STICK, SEOBJ_MAPPYBELL, SEOBJ_MAPPYBAL, 
	SEOBJ_WATER_POOL, SEOBJ_WATER_TOP, SEOBJ_WATER_RT, SEOBJ_WATER_COR, 
	SEOBJ_SLIME_POOL, SEOBJ_SLIME_TOP, SEOBJ_SLIME_RT, SEOBJ_SLIME_COR,
	SEOBJ_LAVA_POOL,  SEOBJ_LAVA_TOP,  SEOBJ_LAVA_RT,  SEOBJ_LAVA_COR,
	SEOBJ_PULS_FLOR, -1,
}

// contrl var nothing [ no-thing ] that blocks elements display
const (
	NOGEN = 1		// all generators
	NOMON = 2		// all monster, dragon
	NOFUD = 4		// all food
	NOTRS = 8		// treas, locked
	NOPOT = 16		// pots & t.powers
	NODOR = 32		// doors, keys
	NOTRAP = 64		// trap & floor dots, stun, ff tiles
	NOEXP = 128		// exit, push wall, teleporter
	NOTHN = 256		// anything else left			511  - all the items
	NOFLOOR = 512
	NOWALL = 1024	//								1536 - floors & walls
	NOG1W = 2048	// g1 std wall only
	ANIM = 4096		// animated
)
// G1 - list of "wrap levels"
var g1wrp = []int{
// horiz wraps
	7, 15, 26, 32, 34, 36, 38, 39, 40, 54, 74, 80, 97, 98, 116, 118, 121,
	200,
// vert wraps
	32, 33, -1,
}

// text for edit key selection

var g1mapid = map[int]string{
	G1OBJ_TILE_FLOOR:	"TILE_FLOOR",
	1:		"Nothing_1",
	G1OBJ_WALL_REGULAR:	"WALL_REGULAR",
	G1OBJ_DOOR_HORIZ:	"DOOR_HORIZ",
	G1OBJ_DOOR_VERT:	"DOOR_VERT",
	G1OBJ_PLAYERSTART:	"START",
	G1OBJ_EXIT:			"EXIT",
	G1OBJ_EXIT4:		"EXIT4",
	G1OBJ_EXIT8:		"EXIT8",
	G1OBJ_MONST_GHOST1:	"GHOST1",
	G1OBJ_MONST_GHOST2:	"GHOST2",
	G1OBJ_MONST_GHOST3:	"GHOST3",
	G1OBJ_MONST_GRUNT1:	"GRUNT1",
	G1OBJ_MONST_GRUNT2:	"GRUNT2",
	G1OBJ_MONST_GRUNT3:	"GRUNT3",
	G1OBJ_MONST_DEMON1:	"DEMON1",
	G1OBJ_MONST_DEMON2:	"DEMON2",
	G1OBJ_MONST_DEMON3:	"DEMON3",
	G1OBJ_MONST_LOBBER1: "LOBBER1",
	G1OBJ_MONST_LOBBER2: "LOBBER2",
	G1OBJ_MONST_LOBBER3: "LOBBER3",
	G1OBJ_MONST_SORC1:	"SORC1",
	G1OBJ_MONST_SORC2:	"SORC2",
	G1OBJ_MONST_SORC3:	"SORC3",
	G1OBJ_MONST_DEATH:	"DEATH",
	G1OBJ_GEN_GHOST1:	"GEN_GHOST1",
	G1OBJ_GEN_GHOST2:	"GEN_GHOST2",
	G1OBJ_GEN_GHOST3:	"GEN_GHOST3",
	G1OBJ_GEN_GRUNT1:	"GEN_GRUNT1",
	G1OBJ_GEN_GRUNT2:	"GEN_GRUNT2",
	G1OBJ_GEN_GRUNT3:	"GEN_GRUNT3",
	G1OBJ_GEN_DEMON1:	"GEN_DEMON1",
	G1OBJ_GEN_DEMON2:	"GEN_DEMON2",
	G1OBJ_GEN_DEMON3:	"GEN_DEMON3",
	G1OBJ_GEN_LOBBER1:	"GEN_LOBBER1",
	G1OBJ_GEN_LOBBER2:	"GEN_LOBBER2",
	G1OBJ_GEN_LOBBER3:	"GEN_LOBBER3",
	G1OBJ_GEN_SORC1:	"GEN_SORC1",
	G1OBJ_GEN_SORC2:	"GEN_SORC2",
	G1OBJ_GEN_SORC3:	"GEN_SORC3",
	G1OBJ_TREASURE:		"TREASURE",
	41:		"Nothing_41",
	G1OBJ_FOOD_DESTRUCTABLE: "FOOD_DESTRUCTABLE",
	G1OBJ_FOOD_INVULN:	"FOOD_INVULN",
	G1OBJ_POT_DESTRUCTABLE:	"POT_DESTRUCTABLE",
	G1OBJ_POT_INVULN:	"POT_INVULN",
	G1OBJ_INVISIBL:		"INVISIBL",
	G1OBJ_X_ARMOR:		"X_ARMOR",
	G1OBJ_X_SPEED:		"X_SPEED",
	G1OBJ_X_MAGIC:		"X_MAGIC",
	G1OBJ_X_SHOTPW:		"X_SHOTPW",
	G1OBJ_X_SHTSPD:		"X_SHTSPD",
	G1OBJ_X_FIGHT:		"X_FIGHT",
	G1OBJ_KEY:			"KEY",
	54:		"Nothing_54",
	55:		"Nothing_55",
	G1OBJ_WALL_DESTRUCTABLE: "WALL_DESTRUCTABLE",
	G1OBJ_WALL_TRAP1:	"WALL_TRAP1",
	G1OBJ_TILE_TRAP1:	"TILE_TRAP1",
	G1OBJ_TRANSPORTER:	"TRANSPORTER",
	60:		"Nothing_60",
	61:		"Nothing_61",
	62:		"62 - utb stun G¹ no has",
//	G1OBJ_TILE_STUN:	"TILE_STUN",
	63:		"Nothing_63",
	G1OBJ_TREASURE_BAG:	"TREASURE_BAG",
	G1OBJ_MONST_THIEF:	"THIEF",
	G1OBJ_EXTEND: "Extended",
	67:		"Wizard_67",
	68:		"No_thing_68",
	69:		"No_thing_69",
	70:		"Test_70",
	71:		"No_thing_71",
	72:		"No_thing_72",
	73:		"No_thing_73",
	74:		"No_thing_74",
	75:		"No_thing_75",
	76:		"No_thing_76",
	77:		"No_thing_77",
	78:		"No_thing_78",
	79:		"No_thing_79",
	80:		"No_thing_80",
	81:		"No_thing_81",
	82:		"No_thing_82",
	83:		"No_thing_83",
	84:		"No_thing_84",
	85:		"No_thing_85",
	86:		"No_thing_86",
	87:		"No_thing_87",
	88:		"No_thing_88",
	89:		"No_thing_89",
	SEOBJ_FAKE_BLK:		"SE_FAKE_BLK",
	SEOBJ_FAKE_BLK_SHT:		"SE_FAKE_BLK_SHT",
	SEOBJ_FAKE_PAS:		"SE_FAKE_PAS",
	SEOBJ_FAKE_PAS_SHT:		"SE_FAKE_PAS_SHT",
	94:		"No_thing_94",
	95:		"No_thing_95",
	96:		"No_thing_96",
	97:		"No_thing_97",
	SEOBJ_FLOORNUL:		"SE_FLOORNUL",
	SEOBJ_FLOORNODRAW:		"SE_FLOORNODRAW",
	SEOBJ_FLOOR:	"XOV_FLOOR",
	SEOBJ_STUN:	"SE_STUN",
	102:	"No_thing_102",
	SEOBJ_PUSHWAL:	"SE_PUSHWAL",
	SEOBJ_SECRTWAL:	"SE_SECRTWAL",
	105:	"No_thing_105",
	SEOBJ_RNDWAL:	"SE_RNDWAL",
	SEOBJ_WAL_TRAPCYC1:	"SE_WAL_TRAPCYC1",
	SEOBJ_WAL_TRAPCYC2:	"SE_WAL_TRAPCYC2",
	SEOBJ_WAL_TRAPCYC3:	"SE_WAL_TRAPCYC3",
	SEOBJ_TILE_TRAP1:	"SE_TILE_TRAP1",
	SEOBJ_TILE_TRAP2:	"SE_TILE_TRAP2",
	SEOBJ_TILE_TRAP3:	"SE_TILE_TRAP3",
	SEOBJ_DOOR_H:	"SE_DOOR_H",
	SEOBJ_DOOR_V:	"SE_DOOR_V",
	115:	"No_thing_115",
	116:	"No_thing_116",
	SEOBJ_EXIT6:	"SE_EXIT6",
	SEOBJ_G2GHOST:	"SE_G2GHOST",
	SEOBJ_G2GRUNT:	"SE_G2GRUNT",
	SEOBJ_G2DEMON:	"SE_G2DEMON",
	SEOBJ_G2LOBER:	"SE_G2LOBER",
	SEOBJ_G2SORC:	"SE_G2SORC",
	SEOBJ_G2AUXGR:	"SE_G2AUXGR",
	SEOBJ_G2DEATH:	"SE_G2DEATH",
	SEOBJ_G2ACID:	"SE_G2ACID",
	SEOBJ_G2SUPSORC:	"SE_G2SUPSORC",
	SEOBJ_G2IT:	"SE_G2_IT",
	SEOBJ_G2GN_GST1:	"SE_G2GN_GHOST1",
	SEOBJ_G2GN_GST2:	"SE_G2GN_GHOST2",
	SEOBJ_G2GN_GST3:	"SE_G2GN_GHOST3",
	SEOBJ_G2GN_GR1:	"SE_G2GN_GRUNT1",
	SEOBJ_G2GN_GR2:	"SE_G2GN_GRUNT2",
	SEOBJ_G2GN_GR3:	"SE_G2GN_GRUNT3",
	SEOBJ_G2GN_DM1:	"SE_G2GN_DEMON1",
	SEOBJ_G2GN_DM2:	"SE_G2GN_DEMON2",
	SEOBJ_G2GN_DM3:	"SE_G2GN_DEMON3",
	SEOBJ_G2GN_LB1:	"SE_G2GN_LOBBER1",
	SEOBJ_G2GN_LB2:	"SE_G2GN_LOBBER2",
	SEOBJ_G2GN_LB3:	"SE_G2GN_LOBBER3",
	SEOBJ_G2GN_SORC1:	"SE_G2GN_SORCEROR1",
	SEOBJ_G2GN_SORC2:	"SE_G2GN_SORCEROR2",
	SEOBJ_G2GN_SORC3:	"SE_G2GN_SORCEROR3",
	SEOBJ_G2GN_AUXGR1:	"SE_G2GN_AUX_GRUNT1",
	SEOBJ_G2GN_AUXGR2:	"SE_G2GN_AUX_GRUNT2",
	SEOBJ_G2GN_AUXGR3:	"SE_G2GN_AUX_GRUNT3",
	146:	"No_thing_146",
	SEOBJ_TREASURE_LOCKED:	"SE_TREASURE_LOCKED",
	148:	"No_thing_148",
	149:	"No_thing_149",
	150:	"No_thing_150",
	151:	"No_thing_151",
	152:	"No_thing_152",
	153:	"No_thing_153",
	SEOBJ_SE_ANKH:	"SE_ANKH",
	SEOBJ_POWER_REPULSE:	"SE_POWER_REPULSE",
	SEOBJ_POWER_REFLECT:	"SE_POWER_REFLECT",
	SEOBJ_POWER_TRANSPORT:	"SE_POWER_TRANSPORT",
	SEOBJ_POWER_SUPERSHOT:	"SE_POWER_SUPERSHOT",
	SEOBJ_POWER_INVULN:	"SE_POWER_INVULN",
	SEOBJ_MONST_DRAGON:	"SE_MONST_DRAGON",
	161:	"No_thing_161",
	162:	"No_thing_162",
	SEOBJ_FORCEFIELDHUB:	"SE_FORCEFIELDHUB",
	SEOBJ_MONST_MUGGER:	"SE_MONST_MUGGER",
	165:	"No_thing_165",
	SEOBJ_FIRE_STICK:	"SE_FIRE_STICK",
	SEOBJ_G2_POISPOT:	"SE_G2_POISPOT",
	SEOBJ_G2_POISFUD:	"SE_G2_POISFUD",
	SEOBJ_G2_QFUD:		"SE_G2_QFUD",
	SEOBJ_KEYRING:		"SE_KEYRING",		// 29, 11

	SEOBJ_MAPPYBDG:		"SE_MAPPYBDG",		// 33, 23

	SEOBJ_MAPPYARAD:	"SE_MAPPYARAD",		// 25, 21
	SEOBJ_MAPPYATV:		"SE_MAPPYATV",		// 27, 21
	SEOBJ_MAPPYAPC:		"SE_MAPPYAPC",		// 29, 21
	SEOBJ_MAPPYAART:	"SE_MAPPYAART",		// 31, 21
	SEOBJ_MAPPYASAF:	"SE_MAPPYASAF",		// 33, 21

	SEOBJ_MAPPYRAD:		"SE_MAPPYRAD",		// 25, 22
	SEOBJ_MAPPYTV:		"SE_MAPPYTV",		// 27, 22
	SEOBJ_MAPPYPC:		"SE_MAPPYPC",		// 29, 22
	SEOBJ_MAPPYART:		"SE_MAPPYART",		// 31, 22
	SEOBJ_MAPPYSAF:		"SE_MAPPYSAF",		// 33, 22

	SEOBJ_MAPPYBELL:	"SE_MAPPYBELL",		// 35, 21
	SEOBJ_MAPPYBAL:		"SE_MAPPYBAL",		// 35, 22
	SEOBJ_MAPPYGORO:	"SE_MAPPYGORO",		// 35, 22

	SEOBJ_DETHGEN3:		"SE_DETHGEN3",		// 34, 8
	SEOBJ_DETHGEN2:		"SE_DETHGEN2",		// 35, 8
	SEOBJ_DETHGEN1:		"SE_DETHGEN1",		// 36, 8

	SEOBJ_WATER_POOL:	"SE_WATER_POOL",
	SEOBJ_WATER_TOP:	"SE_WATER_TOP",
	SEOBJ_WATER_RT:		"SE_WATER_RT",
	SEOBJ_WATER_COR:	"SE_WATER_COR",
	SEOBJ_SLIME_POOL:	"SE_SLIME_POOL",
	SEOBJ_SLIME_TOP:	"SE_SLIME_TOP",
	SEOBJ_SLIME_RT:		"SE_SLIME_RT",
	SEOBJ_SLIME_COR:	"SE_SLIME_COR",
	SEOBJ_LAVA_POOL:	"SE_LAVA_POOL",
	SEOBJ_LAVA_TOP:		"SE_LAVA_TOP",
	SEOBJ_LAVA_RT:		"SE_LAVA_RT",
	SEOBJ_LAVA_COR:		"SE_LAVA_COR",
	SEOBJ_PULS_FLOR:	"SE_PULS_FLOR",

//	169:	"SE_",

/*
	10:	"No_thing_10",
	11:	"No_thing_11",
	12:	"No_thing_12",
	13:	"No_thing_13",
	14:	"No_thing_14",
	15:	"No_thing_15",
	16:	"No_thing_16",
	17:	"No_thing_17",
	18:	"No_thing_18",
	19:	"No_thing_19",
*/
}

var g2mapid = map[int]string{
	MAZEOBJ_TILE_FLOOR:		"TILE_FLOOR",
	MAZEOBJ_TILE_STUN:		"TILE_STUN",
	MAZEOBJ_WALL_REGULAR:	"WALL_REGULAR",
	MAZEOBJ_WALL_MOVABLE:	"WALL_MOVABLE",
	MAZEOBJ_WALL_SECRET:	"WALL_SECRET",
	MAZEOBJ_WALL_DESTRUCTABLE: "WALL_DESTRUCTABLE",
	MAZEOBJ_WALL_RANDOM:	"WALL_RANDOM",
	MAZEOBJ_WALL_TRAPCYC1:	"WALL_TRAPCYC1",
	MAZEOBJ_WALL_TRAPCYC2:	"WALL_TRAPCYC2",
	MAZEOBJ_WALL_TRAPCYC3:	"WALL_TRAPCYC3",
	MAZEOBJ_TILE_TRAP1:		"TILE_TRAP1",
	MAZEOBJ_TILE_TRAP2:		"TILE_TRAP2",
	MAZEOBJ_TILE_TRAP3:		"TILE_TRAP3",
	MAZEOBJ_DOOR_HORIZ:		"DOOR_HORIZ",
	MAZEOBJ_DOOR_VERT:		"DOOR_VERT",
	MAZEOBJ_PLAYERSTART:	"START",
	MAZEOBJ_EXIT:			"EXIT",
	MAZEOBJ_EXITTO6:		"EXITTO6",
	MAZEOBJ_MONST_GHOST:	"GHOST",
	MAZEOBJ_MONST_GRUNT:	"GRUNT",
	MAZEOBJ_MONST_DEMON:	"DEMON",
	MAZEOBJ_MONST_LOBBER:	"LOBBER",
	MAZEOBJ_MONST_SORC:		"SORC",
	MAZEOBJ_MONST_AUX_GRUNT: "AUX_GRUNT",
	MAZEOBJ_MONST_DEATH:	"DEATH",
	MAZEOBJ_MONST_ACID:		"ACID",
	MAZEOBJ_MONST_SUPERSORC: "SUPERSORC",
	MAZEOBJ_MONST_IT:		"IT",
	MAZEOBJ_GEN_GHOST1:		"GEN_GHOST1",
	MAZEOBJ_GEN_GHOST2:		"GEN_GHOST2",
	MAZEOBJ_GEN_GHOST3:		"GEN_GHOST3",
	MAZEOBJ_GEN_GRUNT1:		"GEN_GRUNT1",
	MAZEOBJ_GEN_GRUNT2:		"GEN_GRUNT2",
	MAZEOBJ_GEN_GRUNT3:		"GEN_GRUNT3",
	MAZEOBJ_GEN_DEMON1:		"GEN_DEMON1",
	MAZEOBJ_GEN_DEMON2:		"GEN_DEMON2",
	MAZEOBJ_GEN_DEMON3:		"GEN_DEMON3",
	MAZEOBJ_GEN_LOBBER1:	"GEN_LOBBER1",
	MAZEOBJ_GEN_LOBBER2:	"GEN_LOBBER2",
	MAZEOBJ_GEN_LOBBER3:	"GEN_LOBBER3",
	MAZEOBJ_GEN_SORC1:		"GEN_SORC1",
	MAZEOBJ_GEN_SORC2:		"GEN_SORC2",
	MAZEOBJ_GEN_SORC3:		"GEN_SORC3",
	MAZEOBJ_GEN_AUX_GRUNT1:	"GEN_AUX_GRUNT1",
	MAZEOBJ_GEN_AUX_GRUNT2:	"GEN_AUX_GRUNT2",
	MAZEOBJ_GEN_AUX_GRUNT3:	"GEN_AUX_GRUNT3",
	MAZEOBJ_TREASURE:		"TREASURE",
	MAZEOBJ_TREASURE_LOCKED: "TREASURE_LOCKED",
	MAZEOBJ_TREASURE_BAG:	"TREASURE_BAG",
	MAZEOBJ_FOOD_DESTRUCTABLE: "FOOD_DESTRUCTABLE",
	MAZEOBJ_FOOD_INVULN:	"FOOD_INVULN",
	MAZEOBJ_POT_DESTRUCTABLE: "POT_DESTRUCTABLE",
	MAZEOBJ_POT_INVULN:		"POT_INVULN",
	MAZEOBJ_KEY:			"KEY",
	MAZEOBJ_POWER_INVIS:	"POWER_INVIS",
	MAZEOBJ_POWER_REPULSE:	"POWER_REPULSE",
	MAZEOBJ_POWER_REFLECT:	"POWER_REFLECT",
	MAZEOBJ_POWER_TRANSPORT: "POWER_TRANSPORT",
	MAZEOBJ_POWER_SUPERSHOT: "POWER_SUPERSHOT",
	MAZEOBJ_POWER_INVULN:	"POWER_INVULN",
	MAZEOBJ_MONST_DRAGON:	"DRAGON",
	MAZEOBJ_HIDDENPOT:		"HIDDENPOT",
	MAZEOBJ_TRANSPORTER:	"TRANSPORTER",
	MAZEOBJ_FORCEFIELDHUB:	"FORCEFIELDHUB",
	MAZEOBJ_MONST_MUGGER:	"MUGGER",
	MAZEOBJ_MONST_THIEF:	"THIEF",
	MAZEOBJ_EXTEND: 		"Extended",
}

// single audio hint for an element

var g1auds = map[int]string{
	G1OBJ_TILE_FLOOR:	"sfx/tile.ogg",
	1:		"",
	G1OBJ_WALL_REGULAR:	"sfx/wall.ogg",
	G1OBJ_DOOR_HORIZ:	"sfx/g1_door.ogg",
	G1OBJ_DOOR_VERT:	"sfx/g1_door.ogg",
	G1OBJ_PLAYERSTART:	"sfx/g1_coindrop.ogg",
	G1OBJ_EXIT:			"sfx/g1_exit.ogg",
	G1OBJ_EXIT4:		"sfx/g1_exit.ogg",
	G1OBJ_EXIT8:		"sfx/g1_exit.ogg",
	G1OBJ_MONST_GHOST1:	"sfx/g1hit_ghost.ogg",
	G1OBJ_MONST_GHOST2:	"sfx/g1hit_ghost.ogg",
	G1OBJ_MONST_GHOST3:	"sfx/g1hit_ghost.ogg",
	G1OBJ_MONST_GRUNT1:	"sfx/g1hit_grunt.ogg",
	G1OBJ_MONST_GRUNT2:	"sfx/g1hit_grunt.ogg",
	G1OBJ_MONST_GRUNT3:	"sfx/g1hit_grunt.ogg",
	G1OBJ_MONST_DEMON1:	"sfx/g1hit_grunt.ogg",
	G1OBJ_MONST_DEMON2:	"sfx/g1hit_grunt.ogg",
	G1OBJ_MONST_DEMON3:	"sfx/g1hit_grunt.ogg",
	G1OBJ_MONST_LOBBER1: "sfx/g2_lobshot.ogg",
	G1OBJ_MONST_LOBBER2: "sfx/g2_lobshot.ogg",
	G1OBJ_MONST_LOBBER3: "sfx/g2_lobshot.ogg",
	G1OBJ_MONST_SORC1:	"sfx/g1hit_grunt.ogg",
	G1OBJ_MONST_SORC2:	"sfx/g1hit_grunt.ogg",
	G1OBJ_MONST_SORC3:	"sfx/g1hit_grunt.ogg",
	G1OBJ_MONST_DEATH:	"sfx/g1_deathtouch.ogg",
	G1OBJ_GEN_GHOST1:	"sfx/g1hit_ghost.ogg",
	G1OBJ_GEN_GHOST2:	"sfx/g1hit_ghost.ogg",
	G1OBJ_GEN_GHOST3:	"sfx/g1hit_ghost.ogg",
	G1OBJ_GEN_GRUNT1:	"sfx/g1hit_grunt.ogg",
	G1OBJ_GEN_GRUNT2:	"sfx/g1hit_grunt.ogg",
	G1OBJ_GEN_GRUNT3:	"sfx/g1hit_grunt.ogg",
	G1OBJ_GEN_DEMON1:	"sfx/g1fire_wiz.ogg",
	G1OBJ_GEN_DEMON2:	"sfx/g1fire_wiz.ogg",
	G1OBJ_GEN_DEMON3:	"sfx/g1fire_wiz.ogg",
	G1OBJ_GEN_LOBBER1:	"sfx/g2_lobshot.ogg",
	G1OBJ_GEN_LOBBER2:	"sfx/g2_lobshot.ogg",
	G1OBJ_GEN_LOBBER3:	"sfx/g2_lobshot.ogg",
	G1OBJ_GEN_SORC1:	"sfx/sorc.ogg",
	G1OBJ_GEN_SORC2:	"sfx/sorc.ogg",
	G1OBJ_GEN_SORC3:	"sfx/sorc.ogg",
	G1OBJ_TREASURE:		"sfx/g1_treaspick.ogg",
	41:		"",
	G1OBJ_FOOD_DESTRUCTABLE: "sfx/g1_foodsnrf.ogg",
	G1OBJ_FOOD_INVULN:	"sfx/g1yum_wiz.ogg",
	G1OBJ_POT_DESTRUCTABLE:	"sfx/g1_potionboom.ogg",
	G1OBJ_POT_INVULN:	"sfx/g1_potionboom.ogg",
	G1OBJ_INVISIBL:		"sfx/g1_potionpick.ogg",
	G1OBJ_X_ARMOR:		"sfx/g1_potionpick.ogg",
	G1OBJ_X_SPEED:		"sfx/g1_potionpick.ogg",
	G1OBJ_X_MAGIC:		"sfx/g1_potionpick.ogg",
	G1OBJ_X_SHOTPW:		"sfx/g1_potionpick.ogg",
	G1OBJ_X_SHTSPD:		"sfx/g1_potionpick.ogg",
	G1OBJ_X_FIGHT:		"sfx/g1_potionpick.ogg",
	G1OBJ_KEY:			"sfx/g1_key.ogg",
	54:		"",
	55:		"",
	G1OBJ_WALL_DESTRUCTABLE: "sfx/crumble.ogg",
	G1OBJ_WALL_TRAP1:	"sfx/g2-wallexit.ogg",
	G1OBJ_TILE_TRAP1:	"sfx/g1_trap.ogg",
	G1OBJ_TRANSPORTER:	"sfx/g1_teleport.ogg",
	60:		"",
	61:		"",
//	G1OBJ_TILE_STUN:	"sfx/g1_stun.ogg",
	62:		"",
	63:		"",
	G1OBJ_TREASURE_BAG:	"sfx/g1_treaspick.ogg",
	G1OBJ_MONST_THIEF:	"sfx/g1thf_lau1.ogg",
	G1OBJ_EXTEND: 		"",
	G1OBJ_WIZARD:		"sfx/g1an_wiz2.ogg",
	68:		"",
	69:		"",
	GORO_TEST:		"sfx/crow.ogg",
	71:		"",
	72:		"",
	73:		"",
	74:		"",
	75:		"",
	76:		"",
	77:		"",
	78:		"",
	79:		"",
	80:		"",
	81:		"",
	82:		"",
	83:		"",
	84:		"",
	85:		"",
	86:		"",
	87:		"",
	88:		"",
	89:		"",
	SEOBJ_FAKE_BLK:		"",
	SEOBJ_FAKE_BLK_SHT:		"",
	SEOBJ_FAKE_PAS:		"",
	SEOBJ_FAKE_PAS_SHT:		"",
	94:		"",
	95:		"",
	96:		"",
	97:		"",
	98:		"",
	SEOBJ_FLOORNODRAW:		"",
	SEOBJ_FLOOR:	"sfx/tile.ogg",
	SEOBJ_STUN:	"sfx/g1_stun.ogg",
	102:	"",
	SEOBJ_PUSHWAL:	"sfx/push_wall.ogg",
	SEOBJ_SECRTWAL:	"sfx/g2_shotwall.ogg",
	105:	"",
	SEOBJ_RNDWAL:	"sfx/wall_rnd.ogg",
	SEOBJ_WAL_TRAPCYC1:	"sfx/g2_wallphase1_f.ogg",
	SEOBJ_WAL_TRAPCYC2:	"sfx/g2_wallphase1_f.ogg",
	SEOBJ_WAL_TRAPCYC3:	"sfx/g2_wallphase1_f.ogg",
	SEOBJ_TILE_TRAP1:	"sfx/g2-trap.ogg",
	SEOBJ_TILE_TRAP2:	"sfx/g2-trap.ogg",
	SEOBJ_TILE_TRAP3:	"sfx/g2-trap.ogg",
	SEOBJ_DOOR_H:	"sfx/g1_door.ogg",
	SEOBJ_DOOR_V:	"sfx/g1_door.ogg",
	115:	"",
	116:	"",
	SEOBJ_EXIT6:	"sfx/g1_exit.ogg",
	SEOBJ_G2GHOST:	"sfx/g1hit_ghost.ogg",
	SEOBJ_G2GRUNT:	"sfx/g1hit_grunt.ogg",
	SEOBJ_G2DEMON:	"sfx/g1hit_grunt.ogg",
	SEOBJ_G2LOBER:	"sfx/g2_lobshot.ogg",
	SEOBJ_G2SORC:	"sfx/g1hit_grunt.ogg",
	SEOBJ_G2AUXGR:	"sfx/g1hit_grunt.ogg",
	SEOBJ_G2DEATH:	"sfx/g1_deathtouch.ogg",
	SEOBJ_G2ACID:	"sfx/g2_pickle.ogg",
	SEOBJ_G2SUPSORC:	"sfx/g1fire_wiz.ogg",
	SEOBJ_G2IT:	"sfx/g2an_nwit.ogg",
	SEOBJ_G2GN_GST1:	"sfx/g1hit_ghost.ogg",
	SEOBJ_G2GN_GST2:	"sfx/g1hit_ghost.ogg",
	SEOBJ_G2GN_GST3:	"sfx/g1hit_ghost.ogg",
	SEOBJ_G2GN_GR1:	"sfx/g1hit_grunt.ogg",
	SEOBJ_G2GN_GR2:	"sfx/g1hit_grunt.ogg",
	SEOBJ_G2GN_GR3:	"sfx/g1hit_grunt.ogg",
	SEOBJ_G2GN_DM1:	"sfx/g1fire_wiz.ogg",
	SEOBJ_G2GN_DM2:	"sfx/g1fire_wiz.ogg",
	SEOBJ_G2GN_DM3:	"sfx/g1fire_wiz.ogg",
	SEOBJ_G2GN_LB1:	"sfx/g2_lobshot.ogg",
	SEOBJ_G2GN_LB2:	"sfx/g2_lobshot.ogg",
	SEOBJ_G2GN_LB3:	"sfx/g2_lobshot.ogg",
	SEOBJ_G2GN_SORC1:	"sfx/sorc.ogg",
	SEOBJ_G2GN_SORC2:	"sfx/sorc.ogg",
	SEOBJ_G2GN_SORC3:	"sfx/sorc.ogg",
	SEOBJ_G2GN_AUXGR1:	"sfx/g1hit_grunt.ogg",
	SEOBJ_G2GN_AUXGR2:	"sfx/g1hit_grunt.ogg",
	SEOBJ_G2GN_AUXGR3:	"sfx/g1hit_grunt.ogg",
	146:	"",
	SEOBJ_TREASURE_LOCKED:	"sfx/g2_unlkchest.ogg",
	148:	"",
	149:	"",
	150:	"",
	151:	"",
	152:	"",
	153:	"",
	SEOBJ_SE_ANKH:	"sfx/g1_potionpick.ogg",
	SEOBJ_POWER_REPULSE:	"sfx/g1_potionpick.ogg",
	SEOBJ_POWER_REFLECT:	"sfx/g2_bouncshot.ogg",
	SEOBJ_POWER_TRANSPORT:	"sfx/g1_teleport.ogg",
	SEOBJ_POWER_SUPERSHOT:	"sfx/g1fire_wiz.ogg",
	SEOBJ_POWER_INVULN:	"sfx/g1_potionpick.ogg",
	SEOBJ_MONST_DRAGON:	"sfx/g2_drag.ogg",
	161:	"",
	162:	"",
	SEOBJ_FORCEFIELDHUB:	"sfx/g2_ffield.ogg",
	SEOBJ_MONST_MUGGER:	"sfx/g2mug_appr.ogg",
	165:	"",
	SEOBJ_FIRE_STICK:	"",
	SEOBJ_G2_POISPOT:	"sfx/g2_slopoisn.ogg",
	SEOBJ_G2_POISFUD:	"sfx/g2_slopoisn.ogg",
	SEOBJ_G2_QFUD:		"sfx/g1_foodsnrf.ogg",
	SEOBJ_KEYRING:		"sfx/g1_key.ogg",		// 29, 11

	SEOBJ_MAPPYBDG:		"",		// 33, 23

	SEOBJ_MAPPYARAD:	"",		// 25, 21
	SEOBJ_MAPPYATV:		"",		// 27, 21
	SEOBJ_MAPPYAPC:		"",		// 29, 21
	SEOBJ_MAPPYAART:	"",		// 31, 21
	SEOBJ_MAPPYASAF:	"",		// 33, 21

	SEOBJ_MAPPYRAD:		"",		// 25, 22
	SEOBJ_MAPPYTV:		"",		// 27, 22
	SEOBJ_MAPPYPC:		"",		// 29, 22
	SEOBJ_MAPPYART:		"",		// 31, 22
	SEOBJ_MAPPYSAF:		"",		// 33, 22

	SEOBJ_MAPPYBELL:	"",		// 35, 21
	SEOBJ_MAPPYBAL:		"",		// 35, 22
	SEOBJ_MAPPYGORO:	"",		// 35, 22

	SEOBJ_DETHGEN3:		"sfx/g1_deathtouch.ogg",		// 34, 8
	SEOBJ_DETHGEN2:		"sfx/g1_deathtouch.ogg",		// 35, 8
	SEOBJ_DETHGEN1:		"sfx/g1_deathtouch.ogg",		// 36, 8

	SEOBJ_WATER_POOL:	"",
	SEOBJ_WATER_TOP:	"",
	SEOBJ_WATER_RT:		"",
	SEOBJ_WATER_COR:	"",
	SEOBJ_SLIME_POOL:	"",
	SEOBJ_SLIME_TOP:	"",
	SEOBJ_SLIME_RT:		"",
	SEOBJ_SLIME_COR:	"",
	SEOBJ_LAVA_POOL:	"",
	SEOBJ_LAVA_TOP:		"",
	SEOBJ_LAVA_RT:		"",
	SEOBJ_LAVA_COR:		"",
	SEOBJ_PULS_FLOR:	"",

/*
	10:	"",
	11:	"",
	12:	"",
	13:	"",
	14:	"",
	15:	"",
	16:	"",
	17:	"",
	18:	"",
	19:	"",
*/
}

var g2auds = map[int]string{
	MAZEOBJ_TILE_FLOOR:		"sfx/tile2.ogg",
	MAZEOBJ_TILE_STUN:		"sfx/g1_stun.ogg",
	MAZEOBJ_WALL_REGULAR:	"sfx/wall.ogg",
	MAZEOBJ_WALL_MOVABLE:	"sfx/push_wall.ogg",
	MAZEOBJ_WALL_SECRET:	"sfx/g2_shotwall.ogg",
	MAZEOBJ_WALL_DESTRUCTABLE: "sfx/crumble.ogg",
	MAZEOBJ_WALL_RANDOM:	"sfx/wall_rnd.ogg",
	MAZEOBJ_WALL_TRAPCYC1:	"sfx/g2_wallphase1_f.ogg",
	MAZEOBJ_WALL_TRAPCYC2:	"sfx/g2_wallphase1_f.ogg",
	MAZEOBJ_WALL_TRAPCYC3:	"sfx/g2_wallphase1_f.ogg",
	MAZEOBJ_TILE_TRAP1:		"sfx/g2-trap.ogg",
	MAZEOBJ_TILE_TRAP2:		"sfx/g2-trap.ogg",
	MAZEOBJ_TILE_TRAP3:		"sfx/g2-trap.ogg",
	MAZEOBJ_DOOR_HORIZ:		"sfx/g1_door.ogg",
	MAZEOBJ_DOOR_VERT:		"sfx/g1_door.ogg",
	MAZEOBJ_PLAYERSTART:	"sfx/g1_coindrop.ogg",
	MAZEOBJ_EXIT:			"sfx/g1_exit.ogg",
	MAZEOBJ_EXITTO6:		"sfx/g1_exit.ogg",
	MAZEOBJ_MONST_GHOST:	"sfx/g1hit_ghost.ogg",
	MAZEOBJ_MONST_GRUNT:	"sfx/g1hit_grunt.ogg",
	MAZEOBJ_MONST_DEMON:	"sfx/g1hit_grunt.ogg",
	MAZEOBJ_MONST_LOBBER:	"sfx/g2_lobshot.ogg",
	MAZEOBJ_MONST_SORC:		"sfx/g1hit_grunt.ogg",
	MAZEOBJ_MONST_AUX_GRUNT: "sfx/g1hit_grunt.ogg",
	MAZEOBJ_MONST_DEATH:	"sfx/g1_deathtouch.ogg",
	MAZEOBJ_MONST_ACID:		"sfx/g2_pickle.ogg",
	MAZEOBJ_MONST_SUPERSORC: "sfx/g1fire_wiz.ogg",
	MAZEOBJ_MONST_IT:		"sfx/g2an_nwit.ogg",
	MAZEOBJ_GEN_GHOST1:		"sfx/g1hit_ghost.ogg",
	MAZEOBJ_GEN_GHOST2:		"sfx/g1hit_ghost.ogg",
	MAZEOBJ_GEN_GHOST3:		"sfx/g1hit_ghost.ogg",
	MAZEOBJ_GEN_GRUNT1:		"sfx/g1hit_grunt.ogg",
	MAZEOBJ_GEN_GRUNT2:		"sfx/g1hit_grunt.ogg",
	MAZEOBJ_GEN_GRUNT3:		"sfx/g1hit_grunt.ogg",
	MAZEOBJ_GEN_DEMON1:		"sfx/g1fire_wiz.ogg",
	MAZEOBJ_GEN_DEMON2:		"sfx/g1fire_wiz.ogg",
	MAZEOBJ_GEN_DEMON3:		"sfx/g1fire_wiz.ogg",
	MAZEOBJ_GEN_LOBBER1:	"sfx/g2_lobshot.ogg",
	MAZEOBJ_GEN_LOBBER2:	"sfx/g2_lobshot.ogg",
	MAZEOBJ_GEN_LOBBER3:	"sfx/g2_lobshot.ogg",
	MAZEOBJ_GEN_SORC1:		"sfx/sorc.ogg",
	MAZEOBJ_GEN_SORC2:		"sfx/sorc.ogg",
	MAZEOBJ_GEN_SORC3:		"sfx/sorc.ogg",
	MAZEOBJ_GEN_AUX_GRUNT1:	"sfx/g1hit_grunt.ogg",
	MAZEOBJ_GEN_AUX_GRUNT2:	"sfx/g1hit_grunt.ogg",
	MAZEOBJ_GEN_AUX_GRUNT3:	"sfx/g1hit_grunt.ogg",
	MAZEOBJ_TREASURE:		"sfx/g1_treaspick.ogg",
	MAZEOBJ_TREASURE_LOCKED: "sfx/g2_unlkchest.ogg",
	MAZEOBJ_TREASURE_BAG:	"sfx/g1_treaspick.ogg",
	MAZEOBJ_FOOD_DESTRUCTABLE: "sfx/g1_foodsnrf.ogg",
	MAZEOBJ_FOOD_INVULN:	"sfx/g1yum_wiz.ogg",
	MAZEOBJ_POT_DESTRUCTABLE: "sfx/g1_potionboom.ogg",
	MAZEOBJ_POT_INVULN:		"sfx/g1_potionboom.ogg",
	MAZEOBJ_KEY:			"sfx/g1_key.ogg",
	MAZEOBJ_POWER_INVIS:	"sfx/g1_potionpick.ogg",
	MAZEOBJ_POWER_REPULSE:	"sfx/g1_potionpick.ogg",
	MAZEOBJ_POWER_REFLECT:	"sfx/g2_bouncshot.ogg",
	MAZEOBJ_POWER_TRANSPORT: "sfx/g1_teleport.ogg",
	MAZEOBJ_POWER_SUPERSHOT: "sfx/g1fire_wiz.ogg",
	MAZEOBJ_POWER_INVULN:	"sfx/g1_potionpick.ogg",
	MAZEOBJ_MONST_DRAGON:	"sfx/g2_drag.ogg",
	MAZEOBJ_HIDDENPOT:		"",
	MAZEOBJ_TRANSPORTER:	"sfx/g1_teleport.ogg",
	MAZEOBJ_FORCEFIELDHUB:	"sfx/g2_ffield.ogg",
	MAZEOBJ_MONST_MUGGER:	"sfx/g2mug_appr.ogg",
	MAZEOBJ_MONST_THIEF:	"sfx/g1thf_lau1.ogg",
	MAZEOBJ_EXTEND: 		"",
}

// translate  g2 map into sanctuary now - have to decide how editor will store, but prob all se from here out

var g2tose = map[int]int{
	MAZEOBJ_TILE_FLOOR:		0,
	MAZEOBJ_TILE_STUN:		SEOBJ_STUN,
	MAZEOBJ_WALL_REGULAR:	G1OBJ_WALL_REGULAR,
	MAZEOBJ_WALL_MOVABLE:	SEOBJ_PUSHWAL,
	MAZEOBJ_WALL_SECRET:	SEOBJ_SECRTWAL,
	MAZEOBJ_WALL_DESTRUCTABLE: G1OBJ_WALL_DESTRUCTABLE,
	MAZEOBJ_WALL_RANDOM:	SEOBJ_RNDWAL,
	MAZEOBJ_WALL_TRAPCYC1:	SEOBJ_WAL_TRAPCYC1,
	MAZEOBJ_WALL_TRAPCYC2:	SEOBJ_WAL_TRAPCYC2,
	MAZEOBJ_WALL_TRAPCYC3:	SEOBJ_WAL_TRAPCYC3,
	MAZEOBJ_TILE_TRAP1:		SEOBJ_TILE_TRAP1,
	MAZEOBJ_TILE_TRAP2:		SEOBJ_TILE_TRAP2,
	MAZEOBJ_TILE_TRAP3:		SEOBJ_TILE_TRAP3,
	MAZEOBJ_DOOR_HORIZ:		G1OBJ_DOOR_HORIZ,
	MAZEOBJ_DOOR_VERT:		G1OBJ_DOOR_VERT,
	MAZEOBJ_PLAYERSTART:	G1OBJ_PLAYERSTART,
	MAZEOBJ_EXIT:			G1OBJ_EXIT,
	MAZEOBJ_EXITTO6:		SEOBJ_EXIT6,
	MAZEOBJ_MONST_GHOST:	SEOBJ_G2GHOST,
	MAZEOBJ_MONST_GRUNT:	SEOBJ_G2GRUNT,
	MAZEOBJ_MONST_DEMON:	SEOBJ_G2DEMON,
	MAZEOBJ_MONST_LOBBER:	SEOBJ_G2LOBER,
	MAZEOBJ_MONST_SORC:		SEOBJ_G2SORC,
	MAZEOBJ_MONST_AUX_GRUNT: SEOBJ_G2AUXGR,
	MAZEOBJ_MONST_DEATH:	SEOBJ_G2DEATH,
	MAZEOBJ_MONST_ACID:		SEOBJ_G2ACID,
	MAZEOBJ_MONST_SUPERSORC: SEOBJ_G2SUPSORC,
	MAZEOBJ_MONST_IT:		SEOBJ_G2IT,
	MAZEOBJ_GEN_GHOST1:		SEOBJ_G2GN_GST1,
	MAZEOBJ_GEN_GHOST2:		SEOBJ_G2GN_GST2,
	MAZEOBJ_GEN_GHOST3:		SEOBJ_G2GN_GST3,
	MAZEOBJ_GEN_GRUNT1:		SEOBJ_G2GN_GR1,
	MAZEOBJ_GEN_GRUNT2:		SEOBJ_G2GN_GR2,
	MAZEOBJ_GEN_GRUNT3:		SEOBJ_G2GN_GR3,
	MAZEOBJ_GEN_DEMON1:		SEOBJ_G2GN_DM1,
	MAZEOBJ_GEN_DEMON2:		SEOBJ_G2GN_DM2,
	MAZEOBJ_GEN_DEMON3:		SEOBJ_G2GN_DM3,
	MAZEOBJ_GEN_LOBBER1:	SEOBJ_G2GN_LB1,
	MAZEOBJ_GEN_LOBBER2:	SEOBJ_G2GN_LB2,
	MAZEOBJ_GEN_LOBBER3:	SEOBJ_G2GN_LB3,
	MAZEOBJ_GEN_SORC1:		SEOBJ_G2GN_SORC1,
	MAZEOBJ_GEN_SORC2:		SEOBJ_G2GN_SORC2,
	MAZEOBJ_GEN_SORC3:		SEOBJ_G2GN_SORC3,
	MAZEOBJ_GEN_AUX_GRUNT1:	SEOBJ_G2GN_AUXGR1,
	MAZEOBJ_GEN_AUX_GRUNT2:	SEOBJ_G2GN_AUXGR2,
	MAZEOBJ_GEN_AUX_GRUNT3:	SEOBJ_G2GN_AUXGR3,
	MAZEOBJ_TREASURE:		G1OBJ_TREASURE,
	MAZEOBJ_TREASURE_LOCKED: SEOBJ_TREASURE_LOCKED,
	MAZEOBJ_TREASURE_BAG:	G1OBJ_TREASURE_BAG,
	MAZEOBJ_FOOD_DESTRUCTABLE: G1OBJ_FOOD_DESTRUCTABLE,
	MAZEOBJ_FOOD_INVULN:	G1OBJ_FOOD_INVULN,
	MAZEOBJ_POT_DESTRUCTABLE: G1OBJ_POT_DESTRUCTABLE,
	MAZEOBJ_POT_INVULN:		G1OBJ_POT_INVULN,
	MAZEOBJ_KEY:			G1OBJ_KEY,
	MAZEOBJ_POWER_INVIS:	G1OBJ_INVISIBL,
	MAZEOBJ_POWER_REPULSE:	SEOBJ_POWER_REPULSE,
	MAZEOBJ_POWER_REFLECT:	SEOBJ_POWER_REFLECT,
	MAZEOBJ_POWER_TRANSPORT: SEOBJ_POWER_TRANSPORT,
	MAZEOBJ_POWER_SUPERSHOT: SEOBJ_POWER_SUPERSHOT,
	MAZEOBJ_POWER_INVULN:	SEOBJ_POWER_INVULN,
	MAZEOBJ_MONST_DRAGON:	SEOBJ_MONST_DRAGON,
	MAZEOBJ_HIDDENPOT:		161,
	MAZEOBJ_TRANSPORTER:	G1OBJ_TRANSPORTER,
	MAZEOBJ_FORCEFIELDHUB:	SEOBJ_FORCEFIELDHUB,
	MAZEOBJ_MONST_MUGGER:	SEOBJ_MONST_MUGGER,
	MAZEOBJ_MONST_THIEF:	G1OBJ_MONST_THIEF,
	MAZEOBJ_EXTEND: 		G1OBJ_EXTEND,
}

// Flags for levels
const (
	LFLAG1_ODDANGLE_GHOSTS     = 0x01000000
	LFLAG1_ODDANGLE_GRUNTS     = 0x02000000
	LFLAG1_ODDANGLE_DEMONS     = 0x04000000
	LFLAG1_ODDANGLE_LOBBERS    = 0x08000000
	LFLAG1_ODDANGLE_SORCERERS  = 0x10000000
	LFLAG1_ODDANGLE_AUX_GRUNTS = 0x20000000
	LFLAG1_ODDANGLE_DEATHS     = 0x40000000
	LFLAG1_INVIS_TRAPWALLS     = 0x80000000

	LFLAG2_FAST_GHOSTS     = 0x010000
	LFLAG2_FAST_GRUNTS     = 0x020000
	LFLAG2_FAST_DEMONS     = 0x040000
	LFLAG2_FAST_LOBBERS    = 0x080000
	LFLAG2_FAST_SORCERERS  = 0x100000
	LFLAG2_FAST_AUX_GRUNTS = 0x200000
	LFLAG2_FAST_DEATHS     = 0x400000
	LFLAG2_INVIS_ALLWALLS  = 0x800000

	LFLAG3_RANDOMFOOD_MASK  = 0x0700
	LFLAG3_WALLS_CYCLIC     = 0x0800
	LFLAG3_WALLS_DELETABLE1 = 0x1000
	LFLAG3_WALLS_DELETABLE2 = 0x2000
	LFLAG3_EXIT_MOVES       = 0x4000
	LFLAG3_EXIT_CHOOSEONE   = 0x8000

	LFLAG4_SHOTS_STUN       = 0x01
	LFLAG4_SHOTS_HURT       = 0x02
	LFLAG4_TRAPS_LOCAL      = 0x04
	LFLAG4_TRAPS_RANDOM     = 0x08
	LFLAG4_WRAP_V           = 0x10
	LFLAG4_WRAP_H           = 0x20
	LFLAG4_EXIT_FAKE        = 0x40
	LFLAG4_PLAYER_OFFSCREEN = 0x80
)

var mazeFlagStrings = map[int]string{
	LFLAG1_ODDANGLE_GHOSTS:     "ODDANGLE_GHOSTS",
	LFLAG1_ODDANGLE_GRUNTS:     "ODDANGLE_GRUNTS",
	LFLAG1_ODDANGLE_DEMONS:     "ODDANGLE_DEMONS",
	LFLAG1_ODDANGLE_LOBBERS:    "ODDANGLE_LOBBERS",
	LFLAG1_ODDANGLE_SORCERERS:  "ODDANGLE_SORCERERS",
	LFLAG1_ODDANGLE_AUX_GRUNTS: "ODDANGLE_GRUNTS",
	LFLAG1_ODDANGLE_DEATHS:     "ODDANGLE_DEATHS",
	LFLAG1_INVIS_TRAPWALLS:     "INVIS_TRAPWALLS",

	LFLAG2_FAST_GHOSTS:     "FAST_GHOSTS",
	LFLAG2_FAST_GRUNTS:     "FAST_GRUNTS",
	LFLAG2_FAST_DEMONS:     "FAST_DEMONS",
	LFLAG2_FAST_LOBBERS:    "FAST_LOBBERS",
	LFLAG2_FAST_SORCERERS:  "FAST_SORCERERS",
	LFLAG2_FAST_AUX_GRUNTS: "FAST_AUX_GRUNTS",
	LFLAG2_FAST_DEATHS:     "FAST_DEATHS",
	LFLAG2_INVIS_ALLWALLS:  "INVIS_ALLWALLS",

	LFLAG3_WALLS_CYCLIC:     "WALLS_CYCLIC",
	LFLAG3_WALLS_DELETABLE1: "WALLS_DELETABLE1",
	LFLAG3_WALLS_DELETABLE2: "WALLS_DELETABLE2",
	LFLAG3_EXIT_MOVES:       "EXIT_MOVES",
	LFLAG3_EXIT_CHOOSEONE:   "EXIT_CHOOSE_ONE",

	LFLAG4_SHOTS_STUN:       "SHOTS_STUN",
	LFLAG4_SHOTS_HURT:       "SHOTS_HURT",
	LFLAG4_TRAPS_LOCAL:      "TRAPS_LOCAL",
	LFLAG4_TRAPS_RANDOM:     "TRAPS_RANDOM",
	LFLAG4_WRAP_V:           "WRAP_V",
	LFLAG4_WRAP_H:           "WRAP_H",
	LFLAG4_EXIT_FAKE:        "EXIT_FAKE",
	LFLAG4_PLAYER_OFFSCREEN: "PLAYER_OFFSCREEN",
}

const (
	TRICK_NONE           = 0x00
	TRICK_TRANSPORT1     = 0x01
	TRICK_TRANSPORT2     = 0x02
	TRICK_TRANSPORT3     = 0x03
	TRICK_TRANSPORT4     = 0x04
	TRICK_WATCHSHOOT1    = 0x05
	TRICK_WATCHSHOOT2    = 0x06
	TRICK_SAVESUPERSHOTS = 0x07
	TRICK_NOUSEINVUL     = 0x08
	TRICK_NOGETHIT       = 0x09
	TRICK_PUSHWALL       = 0x0a
	TRICK_NOFOOLED       = 0x0b
	TRICK_NOGREEDY1      = 0x0c
	TRICK_NOGREEDY2      = 0x0d
	TRICK_DIET           = 0x0e
	TRICK_BEPUSHY        = 0x0f
	TRICK_IT             = 0x10
	TRICK_NOHURTFRIENDS  = 0x11
)

var mazeSecretStrings = map[int]string{
	TRICK_NONE:           "No trick",
	TRICK_TRANSPORT1:     "Try Transportability (onto demon)",
	TRICK_TRANSPORT2:     "Try Transportability (onto death)",
	TRICK_TRANSPORT3:     "Try Transportability (into exit)",
	TRICK_TRANSPORT4:     "Try Transportability (onto secret wall)",
	TRICK_WATCHSHOOT1:    "Watch What You Shoot (shoot foods)",
	TRICK_WATCHSHOOT2:    "Watch What You Shoot (shoot secret walls)",
	TRICK_SAVESUPERSHOTS: "Save Super Shots",
	TRICK_NOUSEINVUL:     "Don't Use Invulnerability",
	TRICK_NOGETHIT:       "Don't Get Hit (while killing a dragon)",
	TRICK_PUSHWALL:       "Try Pushing a Wall",
	TRICK_NOFOOLED:       "Don't Be Fooled",
	TRICK_NOGREEDY1:      "Don't Be Greedy (no keys or potions)",
	TRICK_NOGREEDY2:      "Don't Be Greedy (no treasure)",
	TRICK_DIET:           "Go On a Diet (no food)",
	TRICK_BEPUSHY:        "Be Pushy",
	TRICK_IT:             "IT Could Be Nice",
	TRICK_NOHURTFRIENDS:  "Don't Hurt Friends",
}

// door -> wall overlaps per wall seg and door pos around wall
// 		   wall seq will test for doors in 4 positions
// door over laps - endcap in 4 dir, wall over in 3 dir
//					there is no wall over for dir 2		overlap dirs are from door perspect
// door pieces: overec L, overec R, overec U, overec D, overw L, overw R, overw U 
//				1		  2			3		  4			5		 6		  7
// these point into shtamp shadow set, past the shadows

var dorvwal = [][]int{
	{ 5, 7, 0, 6, },		// 0, a single pillar
	{ 5, 3, 0, 6, },		// 1  largest surf down facing endcap
	{ 0, 7, 0, 2, },		// 2  large surf left facing endcap
	{ 0,0,0,0,},
	{ 5, 0, 4, 6, },		// 4  smol surf up facing endcap
	{ 5, 0, 0, 6, },		// 5  u/d wall doors r/l
	{ 0, 0, 0, 6, },		// 6  ┍ corner doors u/l
	{ 0, 0, 0, 6, },		// 7  ┣ door l
	{ 1, 7, 0, 0, },		// 8  med surf right facing endcap
	{ 5, 7, 0, 0, },		// 9  ┙ corner doors d/r
	{ 0, 7, 0, 0, },		// 10 r/l wall doors u/d
	{ 0, 7, 0, 0, },		// 11 ┻ wall doors d
	{ 5, 0, 0, 0, },		// 12 ┓ door u/r
	{ 5, 0, 0, 0, },		// 13 ┫ door r
	{ 0,0,0,0,},			// 14 ┳ door u - no adj
	{ 0,0,0,0,},
	{ 0,0,0,0,},
	{ 0, 0, 0, 6, },		// 17 ┣ door l
	{ 0,0,0,0,},
	{ 5, 0, 0, 0, },		// 19 ┫ door r
	{ 0,0,0,0,},
	{ 0,0,0,0,},
	{ 0,0,0,0,},
	{ 0,0,0,0,},
	{ 0,0,0,0,},			// 24 ┳ door u - no adj
	{ 0,0,0,0,},
	{ 0,0,0,0,},
	{ 0,0,0,0,},
}
// edit key shortcut list - < 0 means not usable, not reassingable
// most of these will need manually set (note: save to cfg file)

const(
	minkey = 33
	cyckey = 99
	edkdef = 121
	maxkey = 126
)

var g1edit_keymap = []int{
// all non key values
//    0    1    2    3    4    5    6    7    8    9
	 -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1, 	//  units 0 - 9
	 -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1, 	//  unit start 10			valid 33 - 126
	 -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1, 	//  unit start 20
	 -1,  -1,  -1,   0,   0,   0,   0,   0,   0,   0,	//  unit 30, @33 =          !  "  #  $  %  &  '
	  0,   0,   0,   0,   0,   0,   0,   0,   0,   0,	//  unit      40 = (  )  *  +  ,  -  .  /  0  1
	  0,   0,   0,   0,   0,   0,   0,   0,   0,   0,	//  unit      50 = 2  3  4  5  6  7  8  9  :  ;
	  0,   0,   0,  -1,   0,  -1,   0,  -1,   4,   0,	//  unit      60 = <  =  >  ?* @  A* B  C* D  E
	 43,  14,  -1,   0,   0,   0,  -1,  16,   0,  23,	//  unit      70 = F  G  H* I  J  K  L* M  N  O
	 45,   0,   0,  -1,  59,   0,  -1,  56,   0,   0,	//  unit      80 = P  Q  R  S* T  U  V* W  X  Y
	  0,   0,  -1,   0,   0,   0,   0,  -1,  64,   0,	//  unit      90 = Z  [  \* ]  ^  _  `  a* b  c
	  3,   0,  42,  11,  29,  46,   0,  53,  35,  32,	//  unit     100 = d  e  f  g  h  i  j  k  l  m
	 27,  38,  44,  57,  58,   0,  40,   0,   0,   2,	//  unit     110 = n  o  p  q  r  s  t  u  v  w
	  6,   0,  24,   0,   0,   0,   0,   0,   0,   0,	//  unit     120 = x  y  z  {  |  }  ~
														// * - not currently reassignable, they are edit mode ops keys
														// ?, \, C, A #a, ESC, L, S, H, V
}

var g1edit_xbmap = []string{
// store xbuf from key_asgn
//    0    1    2    3    4    5    6    7    8    9
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
	 "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",  "0",
}

var g2edit_keymap = []int{
// all non key values
//    0    1    2    3    4    5    6    7    8    9
	 -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1, 	//  units 0 - 9
	 -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1, 	//  unit start 10			valid 33 - 126
	 -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1,  -1, 	//  unit start 20
	 -1,  -1,  -1,   0,   0,   0,   0,   0,   0,   0,	//  unit 30, @33 =          !  "  #  $  %  &  '
	  0,   0,   0,   0,   0,   0,   0,   0,   0,   0,	//  unit      40 = (  )  *  +  ,  -  .  /  0  1
	  0,   0,   0,   0,   0,   0,   0,   0,   0,   0,	//  unit      50 = 2  3  4  5  6  7  8  9  :  ;
	  0,   0,   0,  -1,   0,  -1,   0,  -1,  14,   0,	//  unit      60 = <  =  >  ?* @  A* B  C* D  E
	 50,   0,  -1,   0,   0,   0,  -1,   0,   0,   0,	//  unit      70 = F  G  H* I  J  K  L* M  N  O
	 52,   0,   6,  -1,  57,   0,  -1,   5,   0,   0,	//  unit      80 = P  Q  R  S* T  U  V* W  X  Y
	  0,   0,  -1,   0,   0,   0,   0,  -1,  63,   0,	//  unit      90 = Z  [  \* ]  ^  _  `  a* b  c
	 13,   0,  49,  18,   0,  54,   0,  53,   0,   3,	//  unit     100 = d  e  f  g  h  i  j  k  l  m
	 25,  26,  51,  10,   7,   1,  46,   0,   0,   2,	//  unit     110 = n  o  p  q  r  s  t  u  v  w
	 16,   0,  24,   0,   0,   0,   0,   0,   0,   0,	//  unit     120 = x  y  z  {  |  }  ~
														// * - not currently reassignable, they are edit mode ops keys
														// ?, \, C, A #a, ESC, L, S, H, V
}

var map_keymap = []string{
// all key indicator for edit mode	-	special chars for SE_LETR 'draw a letter index to map'
//    0    1    2    3    4    5    6    7    8    9		hex		dec
	"א", " ", " ", " ", " ", " ", " ", " ", " ", " ", 	//  00		0	units 0 - 9
	" ", " ", " ", " ", "•", "ƒ", "™", "≈", "¥", "œ", 	//  0A		10	unit start 0A / 10
	"¤", "¦", "¶", "£", "∞", "ø", "«", "»", "§", "˚",	//  14		20
	"∆", "Ω", " ", "!","\"", "#", "$", "%", "&", "'",	//	1E		30			valid keys 33 - 126
	"(", ")", "*", "+", ",", "-", ".", "/", "0", "1",	//	28		40
	"2", "3", "4", "5", "6", "7", "8", "9", ":", ";",	//	32		50
	"<", "=", ">", "?", "@", "A", "B", "C", "D", "E",	//	3C		60
	"F", "G", "H", "I", "J", "K", "L", "M", "N", "O",	//	46		70
	"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y",	//	50		80
	"Z", "[","\\", "]", "^", "_", "`", "a", "b", "c",	//	5A		90
	"d", "e", "f", "g", "h", "i", "j", "k", "l", "m",	//	64		100
	"n", "o", "p", "q", "r", "s", "t", "u", "v", "w",	//	6E		110
	"x", "y", "z", "{", "|", "}", "~", "†", " ", "ח",	//	78		120

}

// some key ops
const (
	NOP		= 0
	COPY	= 1
	CUT		= 2
	PASTE	= 4
	REPLACE	= 8
)
// sanctuary converter values
var sanct_vrt = []int{

	 0xa08060,		// 	G1MP_TILE_FLOOR: 0,
	 0x00001,		// 	G1MP_nospec_1: 1,
	 0x404000,		// 	G1MP_WALL_REGULAR: 2,
	 0xc0c000,		// 	G1MP_DOOR_HORIZ: 3,
	 0xc0c000,		// 	G1MP_DOOR_VERT: 4,
	 0xf000,		// 	G1MP_PLAYERSTART: 5,
	 0x4000,		// 	G1MP_EXIT: 6,
	 0x4010,		// 	G1MP_EXIT4: 7,
	 0x4020,		// 	G1MP_EXIT8: 8,
	 0x4000b0,		// 	G1MP_MONST_GHOST1: 9,
	 0x400060,		// 	G1MP_MONST_GHOST2: 10,
	 0x400000,		// 	G1MP_MONST_GHOST3: 11,
	 0x4000d0,		// 	G1MP_MONST_GRUNT1: 12,
	 0x400080,		// 	G1MP_MONST_GRUNT2: 13,
	 0x400020,		// 	G1MP_MONST_GRUNT3: 14,
	 0x4000c0,		// 	G1MP_MONST_DEMON1: 15,
	 0x400070,		// 	G1MP_MONST_DEMON2: 16,
	 0x400010,		// 	G1MP_MONST_DEMON3: 17,
	 0x4000f0,		// 	G1MP_MONST_LOBBER1: 18,
	 0x4000a0,		// 	G1MP_MONST_LOBBER2: 19,
	 0x400050,		// 	G1MP_MONST_LOBBER3: 20,
	 0x4000e0,		// 	G1MP_MONST_SORC1: 21,
	 0x400090,		// 	G1MP_MONST_SORC2: 22,
	 0x400030,		// 	G1MP_MONST_SORC3: 23,
	 0x400040,		// 	G1MP_MONST_DEATH: 24,
	 0xf000b0,		// 	G1MP_GEN_GHOST1: 25, 
	 0xf00060,		// 	G1MP_GEN_GHOST2: 26,
	 0xf00000,		// 	G1MP_GEN_GHOST3: 27,
	 0xf000d0,		// 	G1MP_GEN_GRUNT1: 28,
	 0xf00080,		// 	G1MP_GEN_GRUNT2: 29,
	 0xf00020,		// 	G1MP_GEN_GRUNT3: 30,
	 0xf000c0,		// 	G1MP_GEN_DEMON1: 31,
	 0xf00070,		// 	G1MP_GEN_DEMON2: 32,
	 0xf00010,		// 	G1MP_GEN_DEMON3: 33,
	 0xf000f0,		// 	G1MP_GEN_LOBBER1: 34,
	 0xf000a0,		// 	G1MP_GEN_LOBBER2: 35,
	 0xf00050,		// 	G1MP_GEN_LOBBER3: 36,
	 0xf000e0,		// 	G1MP_GEN_SORC1: 37,
	 0xf00090,		// 	G1MP_GEN_SORC2: 38,
	 0xf00030,		// 	G1MP_GEN_SORC3: 39,
	 0x8070,		// 	G1MP_TREASURE: 40,
	 0x000041,		// 	G1MP_nospec_41: 41,
	 0x8000,		// 	G1MP_FOOD_DESTRUCTABLE: 42,
	 0x8020,		// 	G1MP_FOOD_INVULN: 43,
	 0x8060,		// 	G1MP_POT_DESTRUCTABLE: 44,
	 0x8061,		// 	G1MP_POT_INVULN: 45,
	 0x80f0,		// 	G1MP_INVISIBL: 46,
	 0x80e3,		// 	G1MP_X_ARMOR: 47,
	 0x80e0,		// 	G1MP_X_SPEED: 48,
	 0x80e5,		// 	G1MP_X_MAGIC: 49,
	 0x80e1,		// 	G1MP_X_SHOTPW: 50,
	 0x80e2,		// 	G1MP_X_SHTSPD: 51,
	 0x80e4,		// 	G1MP_X_FIGHT: 52,
	 0x8050,		// 	G1MP_KEY: 53,
	 0x000054,		// 	G1MP_nospec_54: 54,
	 0x000055,		// 	G1MP_POT_INVULN: 55,
	 0x8210,		// 	G1MP_WALL_DESTRUCTABLE: 56,
	 0x404030,		// 	G1MP_WALL_TRAP1: 57,
	 0x80b0,		// 	G1MP_TILE_TRAP1: 58,
	 0x80a0,		// 	G1MP_TRANSPORTER: 59,
	 0x000060,		// 	G1MP_nospec_60: 60,
	 0x000061,		// 	G1MP_nospec_60: 61,
	 0x80c0,		// 	G1MP_TILE_STUN: 62,
	 0x000063,		// 	G1MP_nospec_60: 63,
	 0x8090,		// 	G1MP_TREASURE_BAG: 64,
	 0x400100,		//  thief is hacked in for score table maze sample area

}

var sanct_vrt2 = []int{

	 0xa08060,		// 	G2MP_TILE_FLOOR: 0,
	 0x80c0,		// 	G2MP_TILE_STUN: 1,
	 0x404000,		// 	G2MP_WALL_REGULAR: 2,
	 0x80d0,		// 	G2MP_WALL_MOVABLE: 3,
	 0x8190,		// 	G2MP_WALL_SECRET: 4,
	 0x8210,		// 	G2MP_WALL_DESTRUCTABLE: 5,
	 0x81d0,		// 	G2MP_WALL_RANDOM: 6,
	 0x404030,		// 	G2MP_WALL_TRAPCYC1: 7,
	 0x404031,		// 	G2MP_WALL_TRAPCYC2: 8,
	 0x404032,		// 	G2MP_WALL_TRAPCYC3: 9,
	 0x80b0,		// 	G2MP_TILE_TRAP1: 10,
	 0x80b1,		// 	G2MP_TILE_TRAP2: 11,
	 0x80b2,		// 	G2MP_TILE_TRAP3: 12,
	 0xc0c000,		// 	G2MP_DOOR_HORIZ: 13,
	 0xc0c000,		// 	G2MP_DOOR_VERT: 14,
	 0xf000,		// 	G2MP_PLAYERSTART: 15,
	 0x4000,		// 	G2MP_EXIT: 16,
	 0x4030,		// 	G2MP_EXIT6: 17,
	 0x400000,		// 	G2MP_MONST_GHOST3: 18,
	 0x400020,		// 	G2MP_MONST_GRUNT3: 19,
	 0x400010,		// 	G2MP_MONST_DEMON3: 20,
	 0x400050,		// 	G2MP_MONST_LOBBER3: 21,
	 0x400030,		// 	G2MP_MONST_SORC3: 22,
		// ! sanctuary does not encode aux grunts yet
	 0x400020,		// 	G2MP_MONST_AUX_GRUNT: 23,
	 0x400040,		// 	G2MP_MONST_DEATH: 24,
	 0x400120,		// 	G2MP_MONST_ACID: 25,
	 0x400130,		// 	G2MP_MONST_SUPERSORC: 26,
	 0x400140,		// 	G2MP_MONST_IT: 27,
	 0xf000b0,		// 	G2MP_GEN_GHOST1: 28,
	 0xf00060,		// 	G2MP_GEN_GHOST2: 29,
	 0xf00000,		// 	G2MP_GEN_GHOST3: 30,
	 0xf000d0,		// 	G2MP_GEN_GRUNT1: 31,
	 0xf00080,		// 	G2MP_GEN_GRUNT2: 32,
	 0xf00020,		// 	G2MP_GEN_GRUNT3: 33,
	 0xf000c0,		// 	G2MP_GEN_DEMON1: 34,
	 0xf00070,		// 	G2MP_GEN_DEMON2: 35,
	 0xf00010,		// 	G2MP_GEN_DEMON3: 36,
	 0xf000f0,		// 	G2MP_GEN_LOBBER1: 37,
	 0xf000a0,		// 	G2MP_GEN_LOBBER2: 38,
	 0xf00050,		// 	G2MP_GEN_LOBBER3: 39,
	 0xf000e0,		// 	G2MP_GEN_SORC1: 40,
	 0xf00090,		// 	G2MP_GEN_SORC2: 41,
	 0xf00030,		// 	G2MP_GEN_SORC3: 42,
	 0xf000d0,		// 	G2MP_GEN_AUX_GRUNT1: 43,
	 0xf00080,		// 	G2MP_GEN_AUX_GRUNT2: 44,
	 0xf00020,		// 	G2MP_GEN_AUX_GRUNT3: 45,
	 0x8070,		// 	G2MP_TREASURE: 46,
	 0x8080,		// 	G2MP_TREASURE_LOCKED: 47,
	 0x8090,		// 	G2MP_TREASURE_BAG: 48,
	 0x8000,		// 	G2MP_FOOD_DESTRUCTABLE: 49,
	 0x8020,		// 	G2MP_FOOD_INVULN: 50,
	 0x8060,		// 	G2MP_POT_DESTRUCTABLE: 51,
	 0x8061,		// 	G2MP_POT_INVULN: 52,
	 0x8050,		// 	G2MP_KEY: 53,
	 0x80f0,		// 	G2MP_POWER_INVIS: 54,
	 0x80f2,		// 	G2MP_POWER_REPULSE: 55,
	 0x80f3,		// 	G2MP_POWER_REFLECT: 56,
	 0x80f5,		// 	G2MP_POWER_TRANSPORT: 57,
	 0x80f4,		// 	G2MP_POWER_SUPERSHOT: 58,
	 0x80f1,		// 	G2MP_POWER_INVULN: 59,
	 0x400150,		// 	G2MP_MONST_DRAGON: 60,
	 0x000061,		// 	G2MP_HIDDENPOT: 61,
	 0x80a0,		// 	G2MP_TRANSPORTER: 62,
	 0x8130,		// 	G2MP_FORCEFIELDHUB: 63,
	 0x000064,		// 	G2MP_nospec_: 64,
	 0x400100,		//  thief is hacked in for score table maze sample area
	 0x000066,		// 	G2MP_nospec_: 66,
	 0x000067,		// 	G2MP_nospec_: 67,
	 0x000068,		// 	G2MP_nospec_: 68,
	 0x000069,		// 	G2MP_nospec_: 69,
	 0x80e0,		// 	G2MP_X_SPEED: 70,
	 0x80e1,		// 	G2MP_X_SHOTPW: 71,
	 0x80e2,		// 	G2MP_X_SHTSPD: 72,
	 0x80e3,		// 	G2MP_X_ARMOR: 73,
	 0x80e4,		// 	G2MP_X_FIGHT: 74,
	 0x80e5,		// 	G2MP_X_MAGIC: 75,
	 0x8090,		// 	G2MP_TREASURE_BAG: 76,

}

// seleector for start level on launch, 0 is default, which is rnd 1 - 7

var lvl_sel = map[string]int{
	"Research":-6,
	"Level 1":1,
	"Level 2":2,
	"Level 3":3,
	"Level 4":4,
	"Level 5":5,
	"Level 6":6,
	"Level 7":7,
	"Level 8...":8,
	"Random 1-7":0,
}

// allow setting of select box from cfg
var lvl_str = map[int]string{
	-6:"Research",
	0:"Random 1-7",
	1:"Level 1",
	2:"Level 2",
	3:"Level 3",
	4:"Level 4",
	5:"Level 5",
	6:"Level 6",
	7:"Level 7",
	8:"Level 8...",
}

// sanctuary custom floor tests
//		if Se_cflr_cnt > 11 { Se_cflr_cnt = 1 }
/*
var Se_maxflr = 11
var Se_cflr = map[int]string{

	1:		"gfx/d3floor_.jpg",
	2:		"gfx/floor009.jpg",
	3:		"gfx/floor011.jpg",
	4:		"gfx/floor014b.jpg",
	5:		"gfx/floor016.jpg",
	6:		"gfx/floor018.jpg",
	7:		"gfx/floor019.jpg",
	8:		"gfx/floor025.jpg",
	9:		"gfx/floor027.jpg",
	10:		"gfx/g1floor7.jpg",
	11:		"gfx/floor013.jpg",
}*/

// font testing
var ld_font = map[int]string{
	0:	"VrBd.ttf",
	1:	"3270-reg.ttf",
	2:	"Atari-reg.ttf",
	3:	"AtariSmol.ttf",
	4:	"C64.ttf",
	5:	"FSEX301-L2.ttf",
	6:	"ganic2.fnt",
	7:	"Gauntlet.ttf",
	8:	"Inconsolata.otf",
	9:	"license.txt",
	10:	"prn3.ttf",
	11:	"ProFont.ttf",
	12:	"ProggyTiny.ttf",
	13:	"PxPlus_bios.ttf",
	14:	"PxPlus_vga8.ttf",
	15:	"Terminus.ttf",
	16:	"Venture.fnt",
	17:	"vlg_reg.ttf",
	18:	"vlpg_reg.ttf",
	19:	"VPPxl.ttf",
	20:	"VrMon.ttf",
	21:	"Vr.ttf",
}
