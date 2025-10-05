export namespace controller {
	
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

}

export namespace types {
	
	export class Stats {
	    win: boolean;
	    kills: number;
	    deaths: number;
	    assists: number;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.win = source["win"];
	        this.kills = source["kills"];
	        this.deaths = source["deaths"];
	        this.assists = source["assists"];
	    }
	}
	export class Participant {
	    championId: number;
	    stats: Stats;
	
	    static createFrom(source: any = {}) {
	        return new Participant(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.championId = source["championId"];
	        this.stats = this.convertValues(source["stats"], Stats);
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
	    endOfGameResult: string;
	    gameDuration: number;
	    queueId: number;
	    gameMode: string;
	    participants: Participant[];
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.endOfGameResult = source["endOfGameResult"];
	        this.gameDuration = source["gameDuration"];
	        this.queueId = source["queueId"];
	        this.gameMode = source["gameMode"];
	        this.participants = this.convertValues(source["participants"], Participant);
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

