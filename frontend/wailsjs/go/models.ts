export namespace types {
	
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

