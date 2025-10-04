export enum RankTypeEnum {
  SOLO = "RANKED_SOLO_5x5", // 单双排
  FLEX = "RANKED_FLEX_SR", // 灵活排位
}

export const MatchType: Record<number, string> = {
  420: "单双排",
  430: "匹配模式",
  440: "灵活排位",
  450: "大乱斗",
}
