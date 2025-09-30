export interface IplayerBaseData{
    accountId?:string,
    displayName?:string,
    gameName?:string,
    internalName?:string,
    nameChangeFlag?:boolean,
    percentCompleteForNextLevel?:number, //下一个等级百分比
    privacy?:string,
    profileIconId?:number, //头像ID
    puuid?:string,
    rerollPoints?:Ireroll,
    summonerId?:number,
    summonerLevel?:number, //召唤师等级
    tagLine?:string,
    unnamed?:boolean,
    xpSinceLastLevel?:number, //当前经验
    xpUntilNextLevel?:number, //下一等级经验
    iconImgSrc?:string,       //头像URL
}

interface Ireroll{
    currentPoints:number,
    maxRolls:number,
    numberOfRolls:number,
    pointsCostToRoll:number,
    pointsToReroll:number,
}