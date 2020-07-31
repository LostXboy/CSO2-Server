package room

import (
	"log"
	"sync"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/model/packet"
	. "github.com/KouKouChan/CSO2-Server/model/user"
)

//房间信息
type (
	RoomInfo struct {
		Id        uint16
		Lastflags uint64
		Flags     uint64
		//roomName          []byte
		RoomNumber        uint8
		PasswordProtected uint8
		//unk03   = roomid          uint16
		// gameModeID        uint8
		// mapID             uint8
		//maxPlayers uint8
		Unk08        uint8
		HostUserID   uint32
		HostUserName []byte
		Unk11        uint8
		Unk12        uint8
		Unk13        uint32
		Unk14        uint16
		Unk15        uint16
		Unk16        uint32
		Unk17        uint16
		Unk18        uint16
		Unk19        uint8
		Unk20        uint8
		Unk21        uint8
		// roomStatus   uint8
		// enableBots   uint8
		Unk24 uint8
		// startMoney   uint16
		Unk26 uint8
		Unk27 []uint8
		Unk28 uint8
		Unk29 uint8
		Unk30 uint64
		// winLimit          uint8
		// killLimit         uint16
		// forceCamera    uint8
		// botEnabled     uint8
		// botDifficulty  uint8
		// numCtBots      uint8
		// numTrBots      uint8
		Unk31 uint8
		Unk35 uint8
		// nextMapEnabled uint8
		// changeTeams    uint8
		AreFlashesDisabled uint8
		CanSpec            uint8
		IsVipRoom          uint8
		VipRoomLevel       uint8
		// difficulty     uint8

		//设置
		Setting       RoomSettings
		CountingDown  bool
		Countdown     uint8
		NumPlayers    uint8
		UserIDs       []uint32
		ParentChannel uint8
		CtScore       uint8
		TrScore       uint8
		CtKillNum     uint32
		TrKillNum     uint32
		WinnerTeam    uint8

		RoomMutex *sync.Mutex
	}

	//房间请求
	InRoomPaket struct {
		InRoomType uint8
	}

	//房间所属频道，用于请求频道
	RoomsRequestPacket struct {
		ChannelServerIndex uint8
		ChannelIndex       uint8
	}

	//未知，用于请求频道
	LobbyJoinRoom struct {
		Unk00 uint8
		Unk01 uint8
		Unk02 uint8
	}
)

const (
	//房间操作，加入、暂停等

	GameStart         = 0 // when a host starts a new game
	HostJoin          = 1 // when someone joins some host's game
	HostStop          = 3
	LeaveResultWindow = 4

	//频道以及房间
	SendFullRoomList = 0
	JoinRoom         = 1
	UpdateUserInfo   = 2

	//房间操作
	NewRoomRequest            = 0
	JoinRoomRequest           = 1
	LeaveRoomRequest          = 3
	ToggleReadyRequest        = 4
	GameStartRequest          = 5
	UpdateSettings            = 6
	OnCloseResultWindow       = 7
	SetUserTeamRequest        = 9
	GameStartCountdownRequest = 19

	//游戏模式
	original            = 1
	teamdeath           = 2
	zombie              = 3
	stealth             = 4
	gunteamdeath        = 5
	tutorial            = 6
	hide                = 7
	pig                 = 8
	animationtest_vcd   = 9
	gz_survivor         = 10
	devtest             = 11
	originalmr          = 12
	originalmrdraw      = 13
	casualbomb          = 14
	deathmatch          = 15
	scenario_test       = 16
	gz                  = 17
	gz_intro            = 18
	gz_tour             = 19
	gz_pve              = 20
	eventmod01          = 21
	duel                = 22
	gz_ZB               = 23
	heroes              = 24
	eventmod02          = 25
	zombiecraft         = 26
	campaign1           = 27
	campaign2           = 28
	campaign3           = 29
	campaign4           = 30
	campaign5           = 31
	campaign6           = 32
	campaign7           = 33
	campaign8           = 34
	campaign9           = 35
	z_scenario          = 36
	zombie_prop         = 37
	ghost               = 38
	tag                 = 39
	hide_match          = 40
	hide_ice            = 41
	diy                 = 42
	hide_Item           = 43
	zd_boss1            = 44
	zd_boss2            = 45
	zd_boss3            = 46
	practice            = 47
	zombie_commander    = 48
	casualoriginal      = 49
	hide2               = 50
	gunball             = 51
	zombie_zeta         = 53
	tdm_small           = 54
	de_small            = 55
	gunteamdeath_re     = 56
	endless_wave        = 57
	rankmatch_original  = 58
	rankmatch_teamdeath = 59
	play_ground         = 60
	madcity             = 61
	hide_origin         = 62
	teamdeath_mutation  = 63
	giant               = 64
	z_scenario_side     = 65
	hide_multi          = 66
	madcity_team        = 67
	rankmatch_stealth   = 68

	//阵营
	Unknown          = 0
	Terrorist        = 1
	CounterTerrorist = 2

	//房间status
	StatusWaiting = 1
	StatusIngame  = 2

	//队伍平衡
	Disabled   = 0
	Enabled    = 1
	WithBots   = 2
	ByKadRatio = 4

	//房间包表示
	OUTCreateAndJoin  = 0
	OUTPlayerJoin     = 1
	OUTPlayerLeave    = 2
	OUTSetPlayerReady = 3
	OUTUpdateSettings = 4
	OUTSetHost        = 5
	OUTSetGameResult  = 6
	OUTsetUserTeam    = 7
	OUTCountdown      = 14

	//最大房间数
	MAXROOMNUMS         = 0xFF
	DefaultCountdownNum = 7
)

func (dest *InRoomPaket) PraseRoomPacket(p Packet) bool {
	if p.Datalen-HeaderLen < 2 {
		return false
	}
	dest.InRoomType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (rm RoomInfo) isGlobalCountdownInProgress() bool {
	return rm.CountingDown
}

// func (rm roomInfo) toggleUserReadyStatu() {

// }

func (rm RoomInfo) roomGetUser(id uint32) *User {
	if id <= 0 ||
		rm.Id <= 0 ||
		rm.NumPlayers <= 0 {
		return nil
	}
	for k, v := range rm.UserIDs {
		if v == id {
			return &rm.UserIDs[k]
		}
	}
	return nil
}

func (rm *RoomInfo) stopCountdown() {
	if rm == nil {
		return
	}
	(*rm).countdown = DefaultCountdownNum
	(*rm).countingDown = false
}

func (rm *roomInfo) setStatus(status uint8) {
	if rm == nil {
		return
	}
	if status == 1 ||
		status == 2 {
		(*rm).setting.status = status
		(*rm).setting.isIngame = status - 1
	}
}

func (rm roomInfo) canStartGame() bool {
	switch rm.setting.gameModeID {
	case deathmatch, original, originalmr, casualbomb, casualoriginal, eventmod01, eventmod02, diy, campaign1, campaign2, campaign3, campaign4, campaign5, tdm_small, de_small, madcity, madcity_team, gunteamdeath, gunteamdeath_re, stealth, teamdeath, teamdeath_mutation, pig:
		if rm.getNumOfReadyPlayers() < 2 {
			return false
		}
	case giant, hide, hide2, hide_match, hide_origin, hide_Item, hide_multi, ghost, tag, zombie, zombiecraft, zombie_commander, zombie_prop, zombie_zeta:
		if rm.getNumOfRealReadyPlayers() < 2 {
			return false
		}
	}
	return true
}

func (rm *roomInfo) progressCountdown(num uint8) {
	if rm.countdown > DefaultCountdownNum ||
		rm.countdown < 0 {
		(*rm).countdown = 0
	}
	if rm.countingDown == false {
		(*rm).countingDown = true
		(*rm).countdown = DefaultCountdownNum
	}
	(*rm).countdown--
	if rm.countdown != num {
		log.Println("Error : Host is counting", num, "but room is", rm.countdown)
	}
}

func (rm *roomInfo) getCountdown() uint8 {
	if rm.countingDown == false {
		log.Println("Error : tried to get countdown without counting down")
		return 0
	}
	if rm.countdown > DefaultCountdownNum ||
		rm.countdown < 0 {
		(*rm).countdown = DefaultCountdownNum
	}
	return rm.countdown
}

func (rm roomInfo) getAllCtNum() int {
	num := 0
	for _, v := range rm.users {
		if v.getUserTeam() == CounterTerrorist {
			num++
		}
	}
	return num
}

func (rm roomInfo) getAllTrNum() int {
	num := 0
	for _, v := range rm.users {
		if v.getUserTeam() == Terrorist {
			num++
		}
	}
	return num
}

func (rm roomInfo) getFreeSlots() int {
	// u := rm.roomGetUser(rm.hostUserID)
	// if u == nil ||
	// 	u.userid <= 0 {
	// 	return 0
	// }
	// if rm.setting.areBotsEnabled != 0 {
	// 	botsInHostTeam := 0
	// 	humansInHostTeam := 0
	// 	if u.getUserTeam() == CounterTerrorist {
	// 		botsInHostTeam = int(rm.setting.numCtBots)
	// 		humansInHostTeam = rm.getAllCtNum()
	// 	} else if u.getUserTeam() == Terrorist {
	// 		botsInHostTeam = int(rm.setting.numTrBots)
	// 		humansInHostTeam = rm.getAllTrNum()
	// 	}
	// 	return botsInHostTeam - humansInHostTeam
	// }
	return int(rm.setting.maxPlayers - rm.numPlayers)
}

func (rm *roomInfo) joinUser(u *user) bool {
	destTeam := rm.findDesirableTeam()
	if destTeam <= 0 {
		log.Println("Error : Cant add User", string(u.username), "to room", string(rm.setting.roomName))
		return false
	}
	(*rm).numPlayers++
	(*u).currentTeam = uint8(destTeam)
	(*u).setUserStatus(UserNotReady)
	u.setUserRoom(rm.id)
	u.setUserIngame(false)
	(*rm).users = append((*rm).users, *u)
	return true
}

func (rm roomInfo) findDesirableTeam() int {
	trNum := 0
	ctNum := 0
	for _, v := range rm.users {
		if v.getUserTeam() == Terrorist {
			trNum++
		} else if v.getUserTeam() == CounterTerrorist {
			ctNum++
		} else {
			log.Println("Error : User", string(v.username), "is in Unknown team in room", string(rm.setting.roomName))
			return 0
		}
	}
	if rm.setting.areBotsEnabled != 0 {
		u := rm.roomGetUser(rm.hostUserID)
		if u == nil ||
			u.userid <= 0 {
			return 0
		}
		botsInHostTeam := 0
		if u.getUserTeam() == CounterTerrorist {
			botsInHostTeam = int(rm.setting.numCtBots)
			if botsInHostTeam > 0 {
				return CounterTerrorist
			}
		} else if u.getUserTeam() == Terrorist {
			botsInHostTeam = int(rm.setting.numTrBots)
			if botsInHostTeam > 0 {
				return Terrorist
			}
		} else {
			log.Println("Error : Host", string(u.username), "is in Unknown team in room", string(rm.setting.roomName))
			return 0
		}
	}
	if trNum < ctNum {
		return Terrorist
	} else {
		return CounterTerrorist
	}
}

func (rm *roomInfo) CheckIngameStatus() {
	if rm == nil {
		return
	}
	if rm.numPlayers <= 0 {
		rm.setStatus(StatusWaiting)
		return
	}
	for _, v := range rm.users {
		if v.currentIsIngame {
			rm.setStatus(StatusIngame)
			return
		}
	}
	rm.setStatus(StatusWaiting)
}

func (rm roomInfo) getNumOfRealReadyPlayers() int {
	num := 0
	for _, v := range rm.users {
		if v.isUserReady() ||
			v.userid == rm.hostUserID {
			num++
		}
	}
	return num
}
func (rm roomInfo) getNumOfReadyPlayers() int {
	botPlayers := int(rm.setting.numCtBots + rm.setting.numTrBots)
	if rm.setting.teamBalanceType == WithBots {
		numCts := rm.getAllCtNum()
		numTrs := rm.getAllTrNum()
		requiredBalanceBots := IntAbs(numCts - numTrs)
		botPlayers = Ternary(botPlayers > requiredBalanceBots, botPlayers, requiredBalanceBots).(int)
	}
	return botPlayers + rm.getNumOfRealReadyPlayers()
}

func (rm *roomInfo) setRoomScore(ctScore uint8, trScore uint8) {
	if rm == nil {
		return
	}
	(*rm).CtScore = ctScore
	(*rm).TrScore = trScore
}

func (rm *roomInfo) resetRoomScore() {
	if rm == nil {
		return
	}
	(*rm).CtScore = 0
	(*rm).TrScore = 0
}
func (rm *roomInfo) setRoomWinner(Winner uint8) {
	if rm == nil {
		return
	}
	(*rm).WinnerTeam = Winner
}

func (rm *roomInfo) resetRoomWinner() {
	if rm == nil {
		return
	}
	(*rm).WinnerTeam = 0
}
func (rm *roomInfo) CountRoomCtKill() {
	if rm == nil {
		return
	}
	(*rm).CtKillNum++
}

func (rm *roomInfo) CountRoomTrKill() {
	if rm == nil {
		return
	}
	(*rm).TrKillNum++
}

func (rm *roomInfo) resetRoomKillNum() {
	if rm == nil {
		return
	}
	(*rm).CtKillNum = 0
	(*rm).TrKillNum = 0
}