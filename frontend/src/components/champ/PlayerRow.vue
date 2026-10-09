<script setup lang="ts">
import { computed } from "vue"
import type { types } from "../../../wailsjs/go/models"
import { RATING_STYLES, formatKda, formatPercent, kdaRatio, positionLabel } from "@/lib/format"
import MatchList from "./MatchList.vue"
import RankBadge from "./RankBadge.vue"

const props = defineProps<{ member: types.TeamMemberSummary; expanded: boolean }>()
const emit = defineEmits<{ (e: "toggle"): void }>()

const stats = computed(() => props.member.stats)
const hasStats = computed(() => (stats.value?.games ?? 0) > 0)
const name = computed(() => props.member.gameName || "未知召唤师")
const rating = computed(() => props.member.rating)
const ratingClass = computed(() => RATING_STYLES[rating.value?.label ?? ""] ?? "bg-muted text-muted-foreground")
const ratingTitle = computed(() =>
  rating.value?.valid ? `综合评分 ${rating.value.score.toFixed(0)} / 100` : "数据不足，暂无评价",
)
const hasPuuid = computed(() => Boolean(props.member.puuid))
</script>

<template>
  <div class="rounded-lg border bg-card">
    <button
      type="button"
      class="grid w-full grid-cols-[52px_minmax(170px,1.5fr)_minmax(170px,1.2fr)_minmax(170px,1.2fr)_minmax(140px,1fr)_84px] items-center gap-3 px-3 py-2 text-left"
      :disabled="!hasPuuid"
      @click="emit('toggle')"
    >
      <span class="text-xs font-medium text-muted-foreground">{{ positionLabel(member.assignedPosition) }}</span>

      <div class="flex min-w-0 items-center gap-2">
        <img v-if="member.championIcon" :src="member.championIcon" :alt="member.championName" class="size-10 rounded-lg" />
        <div v-else class="size-10 rounded-lg bg-muted" />
        <div class="min-w-0">
          <p class="truncate text-sm font-semibold text-foreground">
            {{ name }}<span v-if="member.tagLine" class="font-normal text-muted-foreground"> #{{ member.tagLine }}</span>
          </p>
          <p class="truncate text-xs text-muted-foreground">
            {{ member.championName || "未选英雄" }}
            <template v-if="member.summonerLevel"> · Lv{{ member.summonerLevel }}</template>
            <template v-if="member.spells?.some(Boolean)"> · {{ member.spells.filter(Boolean).join("/") }}</template>
          </p>
          <div v-if="member.tags?.length" class="mt-0.5 flex flex-wrap gap-1">
            <span v-for="t in member.tags" :key="t" class="rounded bg-muted px-1.5 text-[10px] text-muted-foreground">{{ t }}</span>
          </div>
        </div>
      </div>

      <div class="space-y-1">
        <RankBadge :rank="member.solo" />
        <RankBadge v-if="member.flex?.tierKey" :rank="member.flex" compact />
      </div>

      <div v-if="hasStats" class="space-y-1 text-xs">
        <div class="flex items-baseline justify-between gap-2">
          <span class="text-sm font-semibold text-foreground">
            近{{ stats.games }}场 {{ formatPercent(stats.winRate) }}
          </span>
          <span class="text-muted-foreground">{{ stats.wins }}胜{{ stats.games - stats.wins }}负</span>
        </div>
        <div class="h-1.5 overflow-hidden rounded bg-rose-500/40">
          <div class="h-full bg-emerald-500" :style="{ width: `${Math.round(stats.winRate)}%` }" />
        </div>
        <p class="text-muted-foreground">
          KDA <span class="font-medium text-foreground">{{ kdaRatio(stats.avgKills, stats.avgDeaths, stats.avgAssists) }}</span>
          （{{ formatKda(Number(stats.avgKills.toFixed(1)), Number(stats.avgDeaths.toFixed(1)), Number(stats.avgAssists.toFixed(1))) }}）
        </p>
      </div>
      <p v-else class="text-xs text-muted-foreground">{{ hasPuuid ? "暂无排位记录" : "信息未公开" }}</p>

      <div class="space-y-0.5 text-xs">
        <template v-if="member.championId && hasStats">
          <p class="text-foreground">
            本英雄 {{ stats.champGames }}场
            <span v-if="stats.champGames > 0" class="text-muted-foreground">· {{ formatPercent(stats.champWinRate) }}</span>
          </p>
          <p v-if="member.masteryLevel" class="text-muted-foreground">
            熟练度 {{ member.masteryLevel }}级 · {{ Math.round(member.masteryPoints / 1000) }}k
          </p>
        </template>
        <p v-if="stats?.streak >= 3" class="text-emerald-500">{{ stats.streak }}连胜</p>
        <p v-else-if="stats?.streak <= -3" class="text-rose-500">{{ -stats.streak }}连败</p>
      </div>

      <div class="text-right">
        <span class="inline-block rounded px-2 py-0.5 text-xs font-bold" :class="ratingClass" :title="ratingTitle">
          {{ rating?.label || "—" }}
        </span>
        <p v-if="rating?.valid" class="mt-0.5 text-[10px] text-muted-foreground">{{ rating.score.toFixed(0) }} 分</p>
      </div>
    </button>

    <div v-if="expanded" class="border-t px-3 py-2">
      <MatchList :matches="member.recentMatches ?? []" :puuid="member.puuid" />
    </div>
  </div>
</template>
