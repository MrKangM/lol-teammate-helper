<script setup lang="ts">
import { computed, ref, watch } from "vue";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { DonutChart } from "@/components/ui/chart-donut";
import { MatchType, RankTypeEnum } from "@/enum/rankTypeEnum";
import type { Game, IRecentMatchRecord, MatchHistory } from "@/interface/matchHistory";
import type { IPlayerBaseData } from "@/interface/baseData";
import type { IRankedStats } from "@/interface/rankData";
import { GetMatchHistoryHeroesByIds, GetPlayerRankMatches } from "../../wailsjs/go/controller/MatchHistory";

const props = defineProps<{
  playerData?: IPlayerBaseData | null;
  rankData?: IRankedStats | null;
}>();

const matchHistory = ref<MatchHistory | null>(null);

const DEFAULT_CHAMPION_ICON =
  "https://ddragon.leagueoflegends.com/cdn/14.15.1/img/champion/Talon.png";

interface HeroSummary {
  name: string;
  icon: string;
}

const HERO_FALLBACK: HeroSummary = {
  name: "未知英雄",
  icon: DEFAULT_CHAMPION_ICON,
};

const heroInfoCache = new Map<number, HeroSummary>();

const fallbackMatches: IRecentMatchRecord[] = [
  {
    id: 1,
    champion: {
      name: "Talon",
      icon: DEFAULT_CHAMPION_ICON,
    },
    isWin: true,
    kills: 14,
    deaths: 3,
    assists: 9,
    queue: "单双排",
    duration: "28:34",
  },
];

const recentMatches = ref<IRecentMatchRecord[]>([...fallbackMatches]);
const visibleMatches = computed(() => recentMatches.value);

const parseDurationString = (duration: string | undefined): number => {
  if (!duration) {
    return 0;
  }

  const [minute, second] = duration.split(":").map((value) => Number.parseInt(value, 10));
  if (!Number.isFinite(minute) || !Number.isFinite(second)) {
    return 0;
  }

  return minute * 60 + second;
};

const formatSecondsToClock = (totalSeconds: number): string => {
  if (!Number.isFinite(totalSeconds) || totalSeconds <= 0) {
    return "--:--";
  }

  return formatDuration(Math.round(totalSeconds));
};

const recentOverview = computed(() => {
  const matches = visibleMatches.value;
  if (!matches.length) {
    return {
      winRate: "0%",
      wins: 0,
      losses: 0,
      kdaScore: "0.0",
      averageDuration: "--:--",
    };
  }

  const total = matches.length;
  const wins = matches.filter((match) => match.isWin).length;
  const losses = total - wins;

  const aggregate = matches.reduce(
    (accumulator, match) => {
      accumulator.kills += match.kills;
      accumulator.deaths += match.deaths;
      accumulator.assists += match.assists;
      accumulator.duration += parseDurationString(match.duration);
      return accumulator;
    },
    { kills: 0, deaths: 0, assists: 0, duration: 0 }
  );

  const average = {
    kills: aggregate.kills / total,
    deaths: aggregate.deaths / total,
    assists: aggregate.assists / total,
    duration: aggregate.duration / total,
  };

  const rawKdaScore = average.deaths > 0 ? ((average.kills + average.assists) / average.deaths) * 3 : (average.kills + average.assists) * 3;
  const kdaScore = Number.isFinite(rawKdaScore) ? rawKdaScore.toFixed(1) : "0.0";

  const winRate = total > 0 ? `${Math.round((wins / total) * 100)}%` : "0%";

  return {
    winRate,
    wins,
    losses,
    kdaScore,
    averageDuration: formatSecondsToClock(average.duration),
  };
});

const formatDuration = (value: unknown): string => {
  const totalSeconds = typeof value === "number" ? value : Number(value);
  if (!Number.isFinite(totalSeconds) || totalSeconds <= 0) {
    return "";
  }
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = Math.floor(totalSeconds % 60);
  return `${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
};

const queueLabel = (queueId: number | undefined): string => {
  if (queueId === undefined || queueId === null) {
    return "";
  }
  return MatchType[queueId as keyof typeof MatchType] ?? "";
};

const ensureHeroInfos = async (ids: number[]) => {
  const missing = ids.filter((id) => id > 0 && !heroInfoCache.has(id));
  if (!missing.length) {
    return;
  }

  try {
    const heroes = await GetMatchHistoryHeroesByIds(missing);
    Object.entries(heroes ?? {}).forEach(([key, hero]) => {
      const heroId = Number(key);
      if (!Number.isFinite(heroId)) {
        return;
      }

      heroInfoCache.set(heroId, {
        name: hero?.name ?? HERO_FALLBACK.name,
        icon:
          hero?.iconDataURI && hero.iconDataURI.length > 0
            ? hero.iconDataURI
            : hero?.squarePortraitPath ?? HERO_FALLBACK.icon,
      });
    });
  } catch (error) {
    console.error("批量获取英雄信息失败", missing, error);
    missing.forEach((id) => {
      heroInfoCache.set(id, { ...HERO_FALLBACK });
    });
  }
};

const getHeroSummary = (championId: number | undefined): HeroSummary => {
  if (typeof championId !== "number" || championId <= 0) {
    return HERO_FALLBACK;
  }
  return heroInfoCache.get(championId) ?? HERO_FALLBACK;
};

const buildRecentMatchRecord = (game: Game, index: number): IRecentMatchRecord => {
  const participant = game.participants?.[0];
  const stats = participant?.stats ?? {};
  const hero = getHeroSummary(participant?.championId);

  return {
    id: index + 1,
    champion: {
      name: hero.name,
      icon: hero.icon,
    },
    isWin: Boolean(stats.win),
    kills: stats.kills ?? 0,
    deaths: stats.deaths ?? 0,
    assists: stats.assists ?? 0,
    queue: queueLabel(game.queueId),
    duration: formatDuration(game.gameDuration),
  };
};

const updateRecentMatches = async (history: MatchHistory | null) => {
  const games = history?.games?.games ?? [];
  if (!games.length) {
    recentMatches.value = [...fallbackMatches];
    return;
  }

  const gamesWithParticipants = games.filter((game): game is Game =>
    Boolean(game?.participants?.length)
  );

  if (!gamesWithParticipants.length) {
    recentMatches.value = [...fallbackMatches];
    return;
  }

  const missingIds = new Set<number>();
  for (const game of gamesWithParticipants) {
    const championId = game.participants?.[0]?.championId;
    if (typeof championId === "number" && championId > 0 && !heroInfoCache.has(championId)) {
      missingIds.add(championId);
    }
  }

  if (missingIds.size) {
    await ensureHeroInfos(Array.from(missingIds));
  }

  const records = gamesWithParticipants.map((game, idx) => buildRecentMatchRecord(game, idx));

  recentMatches.value = records.length ? records : [...fallbackMatches];
};

const fetchMatchHistory = async (puuid: string | undefined | null) => {
  if (!puuid) {
    matchHistory.value = null;
    recentMatches.value = [...fallbackMatches];
    return;
  }

  try {
    matchHistory.value = await GetPlayerRankMatches(puuid);
  } catch (error) {
    console.error("拉取对局历史失败", error);
    matchHistory.value = null;
    recentMatches.value = [...fallbackMatches];
  }
};

const summoner = computed(() => {
  const baseData = props.playerData;

  const name =
    baseData?.gameName && baseData?.tagLine
      ? `${baseData.gameName}#${baseData.tagLine}`
      : baseData?.gameName ?? "召唤师";

  return {
    name,
    server: baseData?.region ?? "未知",
    level: baseData?.summonerLevel ?? 0,
    avatarUrl: baseData?.iconImgSrc ?? "",
    currentExp: baseData?.xpSinceLastLevel ?? 0,
    nextLevelExp: baseData?.xpUntilNextLevel ?? 0,
  };
});

const seasonStats = computed(() => {
  const soloRankData = props.rankData?.queueMap[RankTypeEnum.SOLO];
  const wins = soloRankData?.wins ?? 0;
  const losses = soloRankData?.losses ?? 0;

  return [
    { name: "胜场", total: wins },
    { name: "败场", total: losses },
  ];
});

const experienceProgress = computed(() => {
  const progress = (summoner.value.currentExp / summoner.value.nextLevelExp) * 100;
  if (!Number.isFinite(progress)) {
    return 0;
  }
  return Math.max(0, Math.min(100, Math.round(progress)));
});

const remainingExp = computed(() =>
  Math.max(summoner.value.nextLevelExp - summoner.value.currentExp, 0)
);

const valueFormatter = (_value: number | Date) => {
  const stats = seasonStats.value;
  const total = stats.reduce((sum, entry) => sum + (entry.total ?? 0), 0);

  if (total <= 0 || typeof _value !== "number") {
    return "胜率 0%";
  }

  const wins = stats[0]?.total ?? 0;
  const winRate = (wins / total) * 100;
  return `胜率 ${winRate.toFixed(1)}%`;
};

watch(
  matchHistory,
  (history) => {
    void updateRecentMatches(history);
  },
  { immediate: true }
);

watch(
  () => props.playerData?.puuid,
  (puuid) => {
    void fetchMatchHistory(puuid);
  },
  { immediate: true }
);
</script>

<template>
  <div class="mx-auto w-full max-w-[960px] space-y-6 p-5 md:p-6">
    <div class="grid gap-6 lg:grid-cols-[1.05fr_1fr]">
      <Card class="h-full">
        <CardHeader class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div class="space-y-1">
            <CardTitle class="text-2xl font-semibold">
              {{ summoner.name }}
            </CardTitle>
            <CardDescription>{{ summoner.server }}</CardDescription>
          </div>
          <Avatar class="size-16 shrink-0">
            <AvatarImage :src="summoner.avatarUrl" alt="召唤师头像" />
            <AvatarFallback>{{ summoner.level }}</AvatarFallback>
          </Avatar>
        </CardHeader>
        <CardContent class="grid gap-4 text-sm sm:grid-cols-2">
          <div class="space-y-2">
            <p class="text-xs text-muted-foreground">当前等级</p>
            <p class="text-3xl font-bold">Lv. {{ summoner.level }}</p>
          </div>
          <div class="space-y-2">
            <p class="text-xs text-muted-foreground">距离升级还需</p>
            <p class="text-lg font-semibold">{{ remainingExp }} XP</p>
          </div>
          <div class="space-y-3 sm:col-span-2">
            <div class="flex items-center justify-between text-xs text-muted-foreground">
              <span>当前经验值</span>
              <span>{{ summoner.currentExp }} / {{ summoner.nextLevelExp }}</span>
            </div>
            <div class="relative h-2 w-full overflow-hidden rounded-full bg-secondary">
              <div class="h-full rounded-full bg-emerald-500 transition-all" :style="{ width: experienceProgress + '%' }" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card class="flex flex-col">
        <CardHeader>
          <CardTitle>赛季概览</CardTitle>
          <CardDescription>单排战绩速览</CardDescription>
        </CardHeader>
        <CardContent class="flex flex-col gap-6">
          <div class="flex flex-col items-center gap-6 sm:flex-row sm:items-start sm:justify-between">
            <DonutChart
              index="name"
              :category="'total'"
              :data="seasonStats"
              :value-formatter="valueFormatter"
              :colors="['green', 'red']"
              class="w-full max-w-[220px]"
            />
            <div class="flex w-full flex-1 flex-col gap-4 text-sm">
              <div class="grid gap-3 sm:grid-cols-2">
                <div class="rounded-lg bg-muted/30 px-4 py-3">
                  <p class="text-xs text-muted-foreground">赛季胜场</p>
                  <p class="text-xl font-semibold text-emerald-500">{{ seasonStats[0].total }}</p>
                </div>
                <div class="rounded-lg bg-muted/30 px-4 py-3">
                  <p class="text-xs text-muted-foreground">赛季败场</p>
                  <p class="text-xl font-semibold text-red-500">{{ seasonStats[1].total }}</p>
                </div>
              </div>
              <Separator />
              <div class="grid gap-3 sm:grid-cols-2">
                <div class="rounded-lg bg-muted/20 px-4 py-3">
                  <p class="text-xs text-muted-foreground">最近胜率</p>
                  <p class="text-xl font-semibold">{{ recentOverview.winRate }}</p>
                </div>
                <div class="rounded-lg bg-muted/20 px-4 py-3">
                  <p class="text-xs text-muted-foreground">KDA</p>
                  <p class="text-lg font-semibold">{{ recentOverview.kdaScore }}</p>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <Card class="flex flex-col">
      <CardHeader class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <CardTitle>近期排位记录</CardTitle>
          <CardDescription>最近对局快速回顾</CardDescription>
        </div>
        <div class="grid gap-3 text-sm text-muted-foreground sm:auto-cols-max sm:grid-flow-col">
          <div class="rounded-md bg-muted/30 px-3 py-2">
            <p class="text-xs">胜负比</p>
            <p class="text-sm font-semibold text-foreground">{{ recentOverview.wins }} : {{ recentOverview.losses }}</p>
          </div>
          <div class="rounded-md bg-muted/30 px-3 py-2">
            <p class="text-xs">平均耗时</p>
            <p class="text-sm font-semibold text-foreground">{{ recentOverview.averageDuration }}</p>
          </div>
        </div>
      </CardHeader>
      <CardContent class="flex flex-col gap-3">
        <div class="max-h-[360px] overflow-y-auto pr-1">
          <ul class="flex flex-col gap-3">
          <li
            v-for="match in visibleMatches"
            :key="match.id"
            class="grid w-full grid-cols-[auto_1fr_auto] items-center gap-3 rounded-lg border border-border/70 bg-background/80 px-3 py-3 text-sm transition hover:bg-muted/60"
          >
            <Avatar class="size-10 shrink-0">
              <AvatarImage :src="match.champion.icon" :alt="match.champion.name" />
              <AvatarFallback>{{ match.champion.name?.charAt(0) ?? "?" }}</AvatarFallback>
            </Avatar>
            <div class="min-w-0 space-y-1">
              <div class="flex items-center gap-2">
                <span class="truncate font-medium leading-none">{{ match.champion.name }}</span>
                <Badge
                  :class="[
                    'border-transparent px-2 py-0 text-xs',
                    match.isWin ? 'bg-emerald-500 text-white' : 'bg-red-500 text-white'
                  ]"
                >
                  {{ match.isWin ? "胜利" : "失败" }}
                </Badge>
              </div>
              <p class="truncate text-xs text-muted-foreground">
                {{ match.queue || "未匹配" }} · {{ match.duration || "--:--" }}
              </p>
            </div>
            <div class="text-right text-xs text-muted-foreground">
              <span class="text-sm font-semibold text-foreground">{{ match.kills }} / {{ match.deaths }} / {{ match.assists }}</span>
              <p class="uppercase tracking-wide">KDA</p>
            </div>
          </li>
          </ul>
        </div>
      </CardContent>
    </Card>
  </div>
</template>






