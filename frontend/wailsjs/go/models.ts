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
	    spells: string[];
	    masteryLevel: number;
	    masteryPoints: number;
	    solo: RankSummary;
	    flex: RankSummary;
	    stats: PlayerStats;
	    rating: Rating;
	    tags: string[];
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
	        this.spells = source["spells"];
	        this.masteryLevel = source["masteryLevel"];
	        this.masteryPoints = source["masteryPoints"];
	        this.solo = this.convertValues(source["solo"], RankSummary);
	        this.flex = this.convertValues(source["flex"], RankSummary);
	        this.stats = this.convertValues(source["stats"], PlayerStats);
	        this.rating = this.convertValues(source["rating"], Rating);
	        this.tags = source["tags"];
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
	export class GameTeam {
	    teamId: number;
	    win: string;
	
	    static createFrom(source: any = {}) {
	        return new GameTeam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.teamId = source["teamId"];
	        this.win = source["win"];
	    }
	}
	export class IdentityPlayer {
	    puuid: string;
	    gameName: string;
	    tagLine: string;
	    summonerName: string;
	    profileIcon: number;
	
	    static createFrom(source: any = {}) {
	        return new IdentityPlayer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.puuid = source["puuid"];
	        this.gameName = source["gameName"];
	        this.tagLine = source["tagLine"];
	        this.summonerName = source["summonerName"];
	        this.profileIcon = source["profileIcon"];
	    }
	}
	export class ParticipantIdentity {
	    participantId: number;
	    player: IdentityPlayer;
	
	    static createFrom(source: any = {}) {
	        return new ParticipantIdentity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.participantId = source["participantId"];
	        this.player = this.convertValues(source["player"], IdentityPlayer);
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
	export class Timeline {
	    lane: string;
	    role: string;
	
	    static createFrom(source: any = {}) {
	        return new Timeline(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lane = source["lane"];
	        this.role = source["role"];
	    }
	}
	export class Stats {
	    win: boolean;
	    kills: number;
	    deaths: number;
	    assists: number;
	    champLevel: number;
	    goldEarned: number;
	    totalDamageDealtToChampions: number;
	    totalDamageTaken: number;
	    totalMinionsKilled: number;
	    neutralMinionsKilled: number;
	    visionScore: number;
	    item0: number;
	    item1: number;
	    item2: number;
	    item3: number;
	    item4: number;
	    item5: number;
	    item6: number;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.win = source["win"];
	        this.kills = source["kills"];
	        this.deaths = source["deaths"];
	        this.assists = source["assists"];
	        this.champLevel = source["champLevel"];
	        this.goldEarned = source["goldEarned"];
	        this.totalDamageDealtToChampions = source["totalDamageDealtToChampions"];
	        this.totalDamageTaken = source["totalDamageTaken"];
	        this.totalMinionsKilled = source["totalMinionsKilled"];
	        this.neutralMinionsKilled = source["neutralMinionsKilled"];
	        this.visionScore = source["visionScore"];
	        this.item0 = source["item0"];
	        this.item1 = source["item1"];
	        this.item2 = source["item2"];
	        this.item3 = source["item3"];
	        this.item4 = source["item4"];
	        this.item5 = source["item5"];
	        this.item6 = source["item6"];
	    }
	}
	export class Participant {
	    participantId: number;
	    teamId: number;
	    championId: number;
	    spell1Id: number;
	    spell2Id: number;
	    stats: Stats;
	    timeline: Timeline;
	
	    static createFrom(source: any = {}) {
	        return new Participant(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.participantId = source["participantId"];
	        this.teamId = source["teamId"];
	        this.championId = source["championId"];
	        this.spell1Id = source["spell1Id"];
	        this.spell2Id = source["spell2Id"];
	        this.stats = this.convertValues(source["stats"], Stats);
	        this.timeline = this.convertValues(source["timeline"], Timeline);
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
	export class Game {
	    gameId: number;
	    gameCreation: number;
	    endOfGameResult: string;
	    gameDuration: number;
	    queueId: number;
	    gameMode: string;
	    participants: Participant[];
	    participantIdentities: ParticipantIdentity[];
	    teams: GameTeam[];
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gameId = source["gameId"];
	        this.gameCreation = source["gameCreation"];
	        this.endOfGameResult = source["endOfGameResult"];
	        this.gameDuration = source["gameDuration"];
	        this.queueId = source["queueId"];
	        this.gameMode = source["gameMode"];
	        this.participants = this.convertValues(source["participants"], Participant);
	        this.participantIdentities = this.convertValues(source["participantIdentities"], ParticipantIdentity);
	        this.teams = this.convertValues(source["teams"], GameTeam);
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
	    spells: string[];
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
	        this.spells = source["spells"];
	        this.items = source["items"];
	        this.isTarget = source["isTarget"];
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
	
	
	
	export class HeroInfo {
	    name: string;
	    squarePortraitPath: string;
	    iconDataURI?: string;
	
	    static createFrom(source: any = {}) {
	        return new HeroInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.squarePortraitPath = source["squarePortraitPath"];
	        this.iconDataURI = source["iconDataURI"];
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
	
	
	export class MatchHistoryGames {
	    gameBeginDate: string;
	    gameCount: number;
	    gameEndDate: string;
	    gameIndexBegin: number;
	    gameIndexEnd: number;
	    games: Game[];
	
	    static createFrom(source: any = {}) {
	        return new MatchHistoryGames(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gameBeginDate = source["gameBeginDate"];
	        this.gameCount = source["gameCount"];
	        this.gameEndDate = source["gameEndDate"];
	        this.gameIndexBegin = source["gameIndexBegin"];
	        this.gameIndexEnd = source["gameIndexEnd"];
	        this.games = this.convertValues(source["games"], Game);
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
	export class MatchHistory {
	    accountId: number;
	    games: MatchHistoryGames;
	    platformId: string;
	
	    static createFrom(source: any = {}) {
	        return new MatchHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.accountId = source["accountId"];
	        this.games = this.convertValues(source["games"], MatchHistoryGames);
	        this.platformId = source["platformId"];
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
	
	
	
	
	
	export class RankedEntry {
	    currentSeasonWinsForRewards: number;
	    division: string;
	    highestDivision: string;
	    highestTier: string;
	    isProvisional: boolean;
	    leaguePoints: number;
	    losses: number;
	    miniSeriesProgress: string;
	    previousSeasonEndDivision: string;
	    previousSeasonEndTier: string;
	    previousSeasonHighestDivision: string;
	    previousSeasonHighestTier: string;
	    previousSeasonWinsForRewards: number;
	    provisionalGameThreshold: number;
	    provisionalGamesRemaining: number;
	    queueType: string;
	    ratedRating: number;
	    ratedTier: string;
	    tier: string;
	    warnings: any;
	    wins: number;
	
	    static createFrom(source: any = {}) {
	        return new RankedEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentSeasonWinsForRewards = source["currentSeasonWinsForRewards"];
	        this.division = source["division"];
	        this.highestDivision = source["highestDivision"];
	        this.highestTier = source["highestTier"];
	        this.isProvisional = source["isProvisional"];
	        this.leaguePoints = source["leaguePoints"];
	        this.losses = source["losses"];
	        this.miniSeriesProgress = source["miniSeriesProgress"];
	        this.previousSeasonEndDivision = source["previousSeasonEndDivision"];
	        this.previousSeasonEndTier = source["previousSeasonEndTier"];
	        this.previousSeasonHighestDivision = source["previousSeasonHighestDivision"];
	        this.previousSeasonHighestTier = source["previousSeasonHighestTier"];
	        this.previousSeasonWinsForRewards = source["previousSeasonWinsForRewards"];
	        this.provisionalGameThreshold = source["provisionalGameThreshold"];
	        this.provisionalGamesRemaining = source["provisionalGamesRemaining"];
	        this.queueType = source["queueType"];
	        this.ratedRating = source["ratedRating"];
	        this.ratedTier = source["ratedTier"];
	        this.tier = source["tier"];
	        this.warnings = source["warnings"];
	        this.wins = source["wins"];
	    }
	}
	export class SeasonInfo {
	    currentSeasonEnd: number;
	    currentSeasonId: number;
	    nextSeasonStart: number;
	
	    static createFrom(source: any = {}) {
	        return new SeasonInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentSeasonEnd = source["currentSeasonEnd"];
	        this.currentSeasonId = source["currentSeasonId"];
	        this.nextSeasonStart = source["nextSeasonStart"];
	    }
	}
	export class RankedStats {
	    currentSeasonSplitPoints: number;
	    earnedRegaliaRewardIds: string[];
	    highestCurrentSeasonReachedTierSR: string;
	    highestPreviousSeasonEndDivision: string;
	    highestPreviousSeasonEndTier: string;
	    highestRankedEntry: RankedEntry;
	    highestRankedEntrySR: RankedEntry;
	    previousSeasonSplitPoints: number;
	    queueMap: Record<string, RankedEntry>;
	    queues: RankedEntry[];
	    rankedRegaliaLevel: number;
	    seasons: Record<string, SeasonInfo>;
	    splitsProgress: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new RankedStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentSeasonSplitPoints = source["currentSeasonSplitPoints"];
	        this.earnedRegaliaRewardIds = source["earnedRegaliaRewardIds"];
	        this.highestCurrentSeasonReachedTierSR = source["highestCurrentSeasonReachedTierSR"];
	        this.highestPreviousSeasonEndDivision = source["highestPreviousSeasonEndDivision"];
	        this.highestPreviousSeasonEndTier = source["highestPreviousSeasonEndTier"];
	        this.highestRankedEntry = this.convertValues(source["highestRankedEntry"], RankedEntry);
	        this.highestRankedEntrySR = this.convertValues(source["highestRankedEntrySR"], RankedEntry);
	        this.previousSeasonSplitPoints = source["previousSeasonSplitPoints"];
	        this.queueMap = this.convertValues(source["queueMap"], RankedEntry, true);
	        this.queues = this.convertValues(source["queues"], RankedEntry);
	        this.rankedRegaliaLevel = source["rankedRegaliaLevel"];
	        this.seasons = this.convertValues(source["seasons"], SeasonInfo, true);
	        this.splitsProgress = source["splitsProgress"];
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

