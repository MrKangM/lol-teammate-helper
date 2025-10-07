<script setup lang="ts">
import { computed, ref } from "vue"

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

type MockMatch = {
  championName: string
  championIcon?: string
  win: boolean
  kills: number
  deaths: number
  assists: number
  duration: number
  queueLabel: string
}

type Teammate = {
  id: string
  name: string
  tag?: string
  rank: string
  position: string
  championIcon?: string
  matches: MockMatch[]
}

const DEFAULT_ICON =
  "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/Teemo.png"
const EMPTY_ICON =
  "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/profileicon/588.png"

type TeammateSummary = Teammate & {
  averageKda: string
  totalKills: number
  totalDeaths: number
  totalAssists: number
}

const mockTeammates = ref<Teammate[]>([
  {
    id: "1",
    name: "孤城旧梦",
    tag: "70511",
    rank: "钻石 IV",
    position: "Support",
    championIcon: "https://ddragon.leagueoflegends.com/cdn/img/champion/loading/Lux_0.jpg",
    matches: [
      {
        championName: "Lux",
        championIcon: "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/Lux.png",
        win: true,
        kills: 4,
        deaths: 2,
        assists: 18,
        duration: 1724,
        queueLabel: "排位单双"
      },
      {
        championName: "Janna",
        championIcon: "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/Janna.png",
        win: false,
        kills: 1,
        deaths: 4,
        assists: 12,
        duration: 1898,
        queueLabel: "排位单双"
      }
    ]
  },
  {
    id: "2",
    name: "追风少年",
    tag: "82964",
    rank: "铂金 I",
    position: "Mid",
    championIcon: undefined,
    matches: [
      {
        championName: "Ahri",
        championIcon: "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/Ahri.png",
        win: true,
        kills: 11,
        deaths: 3,
        assists: 7,
        duration: 2102,
        queueLabel: "排位单双"
      },
      {
        championName: "LeBlanc",
        championIcon: "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/Leblanc.png",
        win: true,
        kills: 8,
        deaths: 2,
        assists: 5,
        duration: 1954,
        queueLabel: "排位单双"
      },
      {
        championName: "Yone",
        championIcon: "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/Yone.png",
        win: false,
        kills: 6,
        deaths: 7,
        assists: 3,
        duration: 2010,
        queueLabel: "排位单双"
      }
    ]
  },
  {
    id: "3",
    name: "战术鱼",
    tag: "92510",
    rank: "黄金 II",
    position: "Jungle",
    championIcon: "https://ddragon.leagueoflegends.com/cdn/img/champion/loading/LeeSin_0.jpg",
    matches: [
      {
        championName: "LeeSin",
        championIcon: "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/LeeSin.png",
        win: true,
        kills: 9,
        deaths: 3,
        assists: 12,
        duration: 1845,
        queueLabel: "排位单双"
      },
      {
        championName: "Graves",
        championIcon: "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/Graves.png",
        win: false,
        kills: 5,
        deaths: 6,
        assists: 4,
        duration: 1980,
        queueLabel: "排位单双"
      }
    ]
  }
])

const teammateSummaries = computed<TeammateSummary[]>(() => {
  return mockTeammates.value.map((teammate) => {
    const matchCount = teammate.matches.length
    if (matchCount === 0) {
      return {
        ...teammate,
        averageKda: "0.0",
        totalKills: 0,
        totalDeaths: 0,
        totalAssists: 0
      }
    }

    const totals = teammate.matches.reduce(
      (acc, match) => {
        acc.kills += match.kills
        acc.deaths += match.deaths
        acc.assists += match.assists
        return acc
      },
      { kills: 0, deaths: 0, assists: 0 }
    )

    const averageKills = totals.kills / matchCount
    const averageDeaths = totals.deaths / matchCount
    const averageAssists = totals.assists / matchCount
    const rawKdaScore = averageDeaths > 0 ? ((averageKills + averageAssists) / averageDeaths) * 3 : (averageKills + averageAssists) * 3
    const kdaScore = Number.isFinite(rawKdaScore) ? rawKdaScore.toFixed(1) : "0.0"

    return {
      ...teammate,
      averageKda: kdaScore,
      totalKills: totals.kills,
      totalDeaths: totals.deaths,
      totalAssists: totals.assists
    }
  })
})

const formatPlayerName = (teammate: TeammateSummary) => {
  if (teammate.tag) {
    return `${teammate.name} #${teammate.tag}`
  }
  return teammate.name
}

const formatMatchDuration = (seconds: number) => {
  if (!Number.isFinite(seconds) || seconds <= 0) {
    return "--:--"
  }
  const minutes = Math.floor(seconds / 60)
  const remain = seconds % 60
  return `${String(minutes).padStart(2, "0")}:${String(remain).padStart(2, "0")}`
}

const formatKdaLine = (match: MockMatch) => `${match.kills} / ${match.deaths} / ${match.assists}`

const formatWinLabel = (match: MockMatch) => (match.win ? "胜利" : "失利")

const getChampionIcon = (icon?: string) => icon || DEFAULT_ICON
const getSelectedChampionIcon = (icon?: string) => icon || EMPTY_ICON
</script>

<template>
  <div class="flex flex-col gap-8">
    <section class="space-y-4">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 class="text-2xl font-semibold text-foreground">当前队友概览</h2>
          <p class="text-sm text-muted-foreground">展示正在英雄选择阶段的队友信息以及近期表现。</p>
        </div>
        <Button variant="outline" size="sm" disabled>基于示例数据</Button>
      </div>

      <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <Card v-for="teammate in teammateSummaries" :key="teammate.id" class="flex flex-col gap-4">
          <CardHeader class="flex flex-row items-center gap-4">
            <Avatar class="h-16 w-16 border border-border">
              <AvatarImage :src="getSelectedChampionIcon(teammate.championIcon)" :alt="teammate.name" />
              <AvatarFallback>?</AvatarFallback>
            </Avatar>
            <div class="space-y-1">
              <CardTitle class="text-lg font-semibold leading-tight">{{ formatPlayerName(teammate) }}</CardTitle>
              <p class="text-xs text-muted-foreground">位置：{{ teammate.position }} ｜ 段位：{{ teammate.rank }}</p>
            </div>
          </CardHeader>
          <CardContent class="grid grid-cols-3 gap-3 text-center text-sm">
            <div class="rounded-lg bg-muted/40 px-3 py-2">
              <p class="text-xs text-muted-foreground">K</p>
              <p class="text-lg font-semibold text-foreground">{{ teammate.totalKills }}</p>
            </div>
            <div class="rounded-lg bg-muted/40 px-3 py-2">
              <p class="text-xs text-muted-foreground">D</p>
              <p class="text-lg font-semibold text-foreground">{{ teammate.totalDeaths }}</p>
            </div>
            <div class="rounded-lg bg-muted/40 px-3 py-2">
              <p class="text-xs text-muted-foreground">A</p>
              <p class="text-lg font-semibold text-foreground">{{ teammate.totalAssists }}</p>
            </div>
            <div class="col-span-3 rounded-lg bg-muted/20 px-3 py-2">
              <p class="text-xs text-muted-foreground">最近比赛 KDA 评分</p>
              <p class="text-lg font-semibold text-emerald-500">{{ teammate.averageKda }}</p>
            </div>
          </CardContent>
        </Card>
      </div>
    </section>

    <section class="space-y-4">
      <h3 class="text-xl font-semibold text-foreground">队友近期比赛记录</h3>
      <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        <Card v-for="teammate in teammateSummaries" :key="`${teammate.id}-matches`" class="flex flex-col">
          <CardHeader>
            <CardTitle class="text-base font-semibold">{{ formatPlayerName(teammate) }}</CardTitle>
            <p class="text-xs text-muted-foreground">共 {{ teammate.matches.length }} 场示例对局</p>
          </CardHeader>
          <CardContent class="space-y-3">
            <ul class="flex flex-col gap-3">
              <li
                v-for="(match, index) in teammate.matches"
                :key="`${teammate.id}-${index}`"
                class="flex items-center gap-3 rounded-lg border border-border/60 bg-background/60 px-3 py-3 text-sm"
              >
                <Avatar class="h-12 w-12 border border-border">
                  <AvatarImage :src="getChampionIcon(match.championIcon)" :alt="match.championName" />
                  <AvatarFallback>{{ match.championName.slice(0, 1) }}</AvatarFallback>
                </Avatar>
                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-2">
                    <Badge :variant="match.win ? 'default' : 'secondary'" :class="match.win ? 'bg-emerald-500 hover:bg-emerald-500 text-white' : 'bg-muted text-muted-foreground'">
                      {{ formatWinLabel(match) }}
                    </Badge>
                    <span class="truncate font-medium text-foreground">{{ match.championName }}</span>
                  </div>
                  <p class="text-xs text-muted-foreground">{{ match.queueLabel }} ｜ 用时 {{ formatMatchDuration(match.duration) }}</p>
                </div>
                <div class="text-right text-xs">
                  <p class="text-sm font-semibold text-foreground">{{ formatKdaLine(match) }}</p>
                  <p class="text-muted-foreground">K / D / A</p>
                </div>
              </li>
            </ul>
          </CardContent>
        </Card>
      </div>
    </section>
  </div>
</template>