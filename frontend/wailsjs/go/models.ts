export namespace diag {
	
	export class Event {
	    // Go type: time
	    time: any;
	    uri: string;
	    type: string;
	    bytes: number;
	    truncated: boolean;
	    data: string;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = this.convertValues(source["time"], null);
	        this.uri = source["uri"];
	        this.type = source["type"];
	        this.bytes = source["bytes"];
	        this.truncated = source["truncated"];
	        this.data = source["data"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Snapshot {
	    connected: boolean;
	    port: number;
	    region: string;
	    lastError: string;
	    // Go type: time
	    connectedAt: any;
	    logPath: string;
	    phase: string;
	    events: Event[];
	
	    static createFrom(source: any = {}) {
	        return new Snapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connected = source["connected"];
	        this.port = source["port"];
	        this.region = source["region"];
	        this.lastError = source["lastError"];
	        this.connectedAt = this.convertValues(source["connectedAt"], null);
	        this.logPath = source["logPath"];
	        this.phase = source["phase"];
	        this.events = this.convertValues(source["events"], Event);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace types {
	
	export class RecentMatchSummary {
	    gameId: number;
	    gameCreation: number;
	    championId: number;
	    championName: string;
	    championIcon: string;
	    position: string;
	    win: boolean;
	    kills: number;
	    deaths: number;
	    assists: number;
	    cs: number;
	    damage: number;
	    gold: number;
	    visionScore: number;
	    queueId: number;
	    gameDuration: number;
	
	    static createFrom(source: any = {}) {
	        return new RecentMatchSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gameId = source["gameId"];
	        this.gameCreation = source["gameCreation"];
	        this.championId = source["championId"];
	        this.championName = source["championName"];
	        this.championIcon = source["championIcon"];
	        this.position = source["position"];
	        this.win = source["win"];
	        this.kills = source["kills"];
	        this.deaths = source["deaths"];
	        this.assists = source["assists"];
	        this.cs = source["cs"];
	        this.damage = source["damage"];
	        this.gold = source["gold"];
	        this.visionScore = source["visionScore"];
	        this.queueId = source["queueId"];
	        this.gameDuration = source["gameDuration"];
	    }
	}
	export class ChampionStat {
	    championId: number;
	    championName: string;
	    championIcon: string;
	    games: number;
	    wins: number;
	    winRate: number;
	    kda: number;
	
	    static createFrom(source: any = {}) {
	        return new ChampionStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.championId = source["championId"];
	        this.championName = source["championName"];
	        this.championIcon = source["championIcon"];
	        this.games = source["games"];
	        this.wins = source["wins"];
	        this.winRate = source["winRate"];
	        this.kda = source["kda"];
	    }
	}
	export class Rating {
	    score: number;
	    valid: boolean;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new Rating(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.score = source["score"];
	        this.valid = source["valid"];
	        this.label = source["label"];
	    }
	}
	export class PlayerStats {
	    games: number;
	    wins: number;
	    winRate: number;
	    avgKills: number;
	    avgDeaths: number;
	    avgAssists: number;
	    kda: number;
	    avgCs: number;
	    avgDamage: number;
	    avgVision: number;
	    streak: number;
	    champGames: number;
	    champWins: number;
	    champWinRate: number;
	    posGames: number;
	
	    static createFrom(source: any = {}) {
	        return new PlayerStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.games = source["games"];
	        this.wins = source["wins"];
	        this.winRate = source["winRate"];
	        this.avgKills = source["avgKills"];
	        this.avgDeaths = source["avgDeaths"];
	        this.avgAssists = source["avgAssists"];
	        this.kda = source["kda"];
	        this.avgCs = source["avgCs"];
	        this.avgDamage = source["avgDamage"];
	        this.avgVision = source["avgVision"];
	        this.streak = source["streak"];
	        this.champGames = source["champGames"];
	        this.champWins = source["champWins"];
	        this.champWinRate = source["champWinRate"];
	        this.posGames = source["posGames"];
	    }
	}
	export class RankSummary {
	    queueType: string;
	    queueName: string;
	    tierKey: string;
	    tier: string;
	    division: string;
	    leaguePoints: number;
	    wins: number;
	    losses: number;
	    winRate: number;
	
	    static createFrom(source: any = {}) {
	        return new RankSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.queueType = source["queueType"];
	        this.queueName = source["queueName"];
	        this.tierKey = source["tierKey"];
	        this.tier = source["tier"];
	        this.division = source["division"];
	        this.leaguePoints = source["leaguePoints"];
	        this.wins = source["wins"];
	        this.losses = source["losses"];
	        this.winRate = source["winRate"];
	    }
	}
	export class SpellInfo {
	    id: number;
	    name: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new SpellInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	    }
	}
	export class TeamMemberSummary {
	    puuid: string;
	    gameName: string;
	    tagLine: string;
	    assignedPosition: string;
	    championId: number;
	    championName: string;
	    championIcon: string;
	    cellId: number;
	    summonerLevel: number;
	    spells: SpellInfo[];
	    masteryLevel: number;
	    masteryPoints: number;
	    solo: RankSummary;
	    flex: RankSummary;
	    stats: PlayerStats;
	    rating: Rating;
	    tags: string[];
	    pool: ChampionStat[];
	    recentMatches: RecentMatchSummary[];
	
	    static createFrom(source: any = {}) {
	        return new TeamMemberSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.puuid = source["puuid"];
	        this.gameName = source["gameName"];
	        this.tagLine = source["tagLine"];
	        this.assignedPosition = source["assignedPosition"];
	        this.championId = source["championId"];
	        this.championName = source["championName"];
	        this.championIcon = source["championIcon"];
	        this.cellId = source["cellId"];
	        this.summonerLevel = source["summonerLevel"];
	        this.spells = this.convertValues(source["spells"], SpellInfo);
	        this.masteryLevel = source["masteryLevel"];
	        this.masteryPoints = source["masteryPoints"];
	        this.solo = this.convertValues(source["solo"], RankSummary);
	        this.flex = this.convertValues(source["flex"], RankSummary);
	        this.stats = this.convertValues(source["stats"], PlayerStats);
	        this.rating = this.convertValues(source["rating"], Rating);
	        this.tags = source["tags"];
	        this.pool = this.convertValues(source["pool"], ChampionStat);
	        this.recentMatches = this.convertValues(source["recentMatches"], RecentMatchSummary);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ChampSelectSnapshot {
	    queueId: number;
	    gameId: number;
	    phase: string;
	    // Go type: time
	    updatedAt: any;
	    team: TeamMemberSummary[];
	    enemy: TeamMemberSummary[];
	
	    static createFrom(source: any = {}) {
	        return new ChampSelectSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.queueId = source["queueId"];
	        this.gameId = source["gameId"];
	        this.phase = source["phase"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.team = this.convertValues(source["team"], TeamMemberSummary);
	        this.enemy = this.convertValues(source["enemy"], TeamMemberSummary);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class GameDetailPlayer {
	    puuid: string;
	    name: string;
	    championId: number;
	    championName: string;
	    championIcon: string;
	    level: number;
	    kills: number;
	    deaths: number;
	    assists: number;
	    cs: number;
	    gold: number;
	    damage: number;
	    damageTaken: number;
	    visionScore: number;
	    spells: SpellInfo[];
	    items: string[];
	    isTarget: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GameDetailPlayer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.puuid = source["puuid"];
	        this.name = source["name"];
	        this.championId = source["championId"];
	        this.championName = source["championName"];
	        this.championIcon = source["championIcon"];
	        this.level = source["level"];
	        this.kills = source["kills"];
	        this.deaths = source["deaths"];
	        this.assists = source["assists"];
	        this.cs = source["cs"];
	        this.gold = source["gold"];
	        this.damage = source["damage"];
	        this.damageTaken = source["damageTaken"];
	        this.visionScore = source["visionScore"];
	        this.spells = this.convertValues(source["spells"], SpellInfo);
	        this.items = source["items"];
	        this.isTarget = source["isTarget"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GameDetailTeam {
	    teamId: number;
	    win: boolean;
	    kills: number;
	    gold: number;
	    players: GameDetailPlayer[];
	
	    static createFrom(source: any = {}) {
	        return new GameDetailTeam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.teamId = source["teamId"];
	        this.win = source["win"];
	        this.kills = source["kills"];
	        this.gold = source["gold"];
	        this.players = this.convertValues(source["players"], GameDetailPlayer);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GameDetail {
	    gameId: number;
	    gameCreation: number;
	    gameDuration: number;
	    queueId: number;
	    teams: GameDetailTeam[];
	
	    static createFrom(source: any = {}) {
	        return new GameDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gameId = source["gameId"];
	        this.gameCreation = source["gameCreation"];
	        this.gameDuration = source["gameDuration"];
	        this.queueId = source["queueId"];
	        this.teams = this.convertValues(source["teams"], GameDetailTeam);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class IReroll {
	    currentPoints: number;
	    maxRolls: number;
	    numberOfRolls: number;
	    pointsCostToRoll: number;
	    pointsToReroll: number;
	
	    static createFrom(source: any = {}) {
	        return new IReroll(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentPoints = source["currentPoints"];
	        this.maxRolls = source["maxRolls"];
	        this.numberOfRolls = source["numberOfRolls"];
	        this.pointsCostToRoll = source["pointsCostToRoll"];
	        this.pointsToReroll = source["pointsToReroll"];
	    }
	}
	export class IPlayerBaseData {
	    accountId?: number;
	    displayName?: string;
	    gameName?: string;
	    internalName?: string;
	    nameChangeFlag?: boolean;
	    percentCompleteForNextLevel?: number;
	    privacy?: string;
	    profileIconId?: number;
	    puuid?: string;
	    rerollPoints?: IReroll;
	    summonerId?: number;
	    summonerLevel?: number;
	    tagLine?: string;
	    unnamed?: boolean;
	    xpSinceLastLevel?: number;
	    xpUntilNextLevel?: number;
	    iconImgSrc?: string;
	    region?: string;
	
	    static createFrom(source: any = {}) {
	        return new IPlayerBaseData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.accountId = source["accountId"];
	        this.displayName = source["displayName"];
	        this.gameName = source["gameName"];
	        this.internalName = source["internalName"];
	        this.nameChangeFlag = source["nameChangeFlag"];
	        this.percentCompleteForNextLevel = source["percentCompleteForNextLevel"];
	        this.privacy = source["privacy"];
	        this.profileIconId = source["profileIconId"];
	        this.puuid = source["puuid"];
	        this.rerollPoints = this.convertValues(source["rerollPoints"], IReroll);
	        this.summonerId = source["summonerId"];
	        this.summonerLevel = source["summonerLevel"];
	        this.tagLine = source["tagLine"];
	        this.unnamed = source["unnamed"];
	        this.xpSinceLastLevel = source["xpSinceLastLevel"];
	        this.xpUntilNextLevel = source["xpUntilNextLevel"];
	        this.iconImgSrc = source["iconImgSrc"];
	        this.region = source["region"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	

}

