// 排位统计数据
export interface IRankedStats {
    currentSeasonSplitPoints: number;          // 当前赛季分段点数
    earnedRegaliaRewardIds: string[];          // 已获得的荣誉奖励ID列表
    highestCurrentSeasonReachedTierSR: string; // 当前赛季达到的最高段位（单双排）
    highestPreviousSeasonEndDivision: string;  // 上赛季结束时的最高段位分级
    highestPreviousSeasonEndTier: string;      // 上赛季结束时的最高段位
    highestRankedEntry: RankedEntry;           // 最高排位条目
    highestRankedEntrySR: RankedEntry;         // 最高排位条目（单双排）
    previousSeasonSplitPoints: number;         // 上赛季分段点数
    queueMap: { [key: string]: RankedEntry };  // 各队列排位数据映射
    queues: RankedEntry[];                     // 排位队列列表
    rankedRegaliaLevel: number;                // 排位荣誉等级
    seasons: { [key: string]: SeasonInfo };    // 各队列赛季信息
    splitsProgress: { [key: string]: any };    // 分段进度
}

// 排位条目信息
interface RankedEntry {
    currentSeasonWinsForRewards: number;    // 当前赛季奖励所需胜场数
    division: string;                       // 当前段位分级（如I、II、III、IV）
    highestDivision: string;                // 历史最高段位分级
    highestTier: string;                    // 历史最高段位
    isProvisional: boolean;                 // 是否为定级赛状态
    leaguePoints: number;                   // 胜点
    losses: number;                         // 失败场次
    miniSeriesProgress: string;             // 晋级赛进度
    previousSeasonEndDivision: string;      // 上赛季结束时的段位分级
    previousSeasonEndTier: string;          // 上赛季结束时的段位
    previousSeasonHighestDivision: string;  // 上赛季最高段位分级
    previousSeasonHighestTier: string;      // 上赛季最高段位
    previousSeasonWinsForRewards: number;   // 上赛季奖励所需胜场数
    provisionalGameThreshold: number;       // 定级赛阈值
    provisionalGamesRemaining: number;      // 剩余定级赛场次
    queueType: string;                      // 队列类型
    ratedRating: number;                    // 评分
    ratedTier: string;                      // 评分段位
    tier: string;                           // 当前段位
    warnings: any;                          // 警告信息
    wins: number;                           // 胜利场次
}

// 赛季信息
interface SeasonInfo {
    currentSeasonEnd: number;  // 当前赛季结束时间戳
    currentSeasonId: number;   // 当前赛季ID
    nextSeasonStart: number;   // 下赛季开始时间戳
}