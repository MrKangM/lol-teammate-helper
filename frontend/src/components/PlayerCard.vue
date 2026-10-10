<script setup lang="ts">
import { computed } from "vue"
import type { types } from "../../wailsjs/go/models"
import { RATING_STYLES, formatPercent, kdaRatio } from "@/lib/format"
import MatchDots from "@/components/MatchDots.vue"
import PositionIcon from "@/components/PositionIcon.vue"
import RankLine from "@/components/RankLine.vue"

const props = defineProps<{ member: types.TeamMemberSummary; side: "team" | "enemy" }>()
const emit = defineEmits<{ (e: "open"): void }>()

const stats = computed(() => props.member.stats)
const hasStats = computed(() => (stats.value?.games ?? 0) > 0)
const hasPuuid = computed(() => Boolean(props.member.puuid))
const label = computed(() => props.member.rating?.label ?? "")
const style = computed(() => RATING_STYLES[label.value])
const ratingTitle = computed(() =>
  props.member.rating?.valid ? `综合评分 ${props.member.rating.score.toFixed(0)} / 100` : "数据不足，暂无评价",
)
const spells = computed(() => (props.member.spells ?? []).filter((s) => s && (s.icon || s.name)))
const rateColor = computed(() => ((stats.value?.winRate ?? 0) >= 52 ? "text-win" : "text-loss"))
</script>

<template>
  <button
    type="button"
    class="card group relative block w-full overflow-hidden px-3.5 py-3 pl-4 text-left transition-colors hover:border-[#38405a] hover:bg-panel-2 disabled:hover:bg-panel"
    :disabled="!hasPuuid"
    @click="emit('open')"
  >
    <span class="absolute inset-y-0 left-0 w-1" :class="style?.strip ?? 'bg-line'" />

    <div class="flex items-start gap-3">
      <div class="relative shrink-0">
        <img v-if="member.championIcon" :src="member.championIcon" :alt="member.championName" class="size-[54px] rounded-xl" />
        <div v-else class="size-[54px] rounded-xl bg-panel-3" />
        <span class="absolute -bottom-1 -left-1 grid size-[22px] place-items-center rounded-full border border-line bg-panel"><PositionIcon :position="member.assignedPosition" :size="14" /></span>
        <div v-if="spells.length" class="absolute -bottom-1 -right-1 flex flex-col gap-px">
          <template v-for="(s, i) in spells" :key="i">
            <img v-if="s.icon" :src="s.icon" :alt="s.name" :title="s.name" class="size-[17px] rounded-[4px] border border-bg" />
            <span v-else class="grid size-[17px] place-items-center rounded-[4px] bg-panel-3 text-[9px] text-muted" :title="s.name">{{ s.name.slice(0, 1) }}</span>
          </template>
        </div>
      </div>

      <div class="min-w-0 flex-1">
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <p class="truncate text-sm font-semibold text-ink">
              {{ member.gameName || "未知召唤师" }}<span v-if="member.tagLine" class="font-normal text-muted"> #{{ member.tagLine }}</span>
            </p>
            <div class="flex items-center gap-1.5 overflow-hidden text-xs text-muted">
              <span class="truncate">
                {{ member.championName || "未选英雄" }}
                <template v-if="member.masteryLevel"> · 熟练 {{ member.masteryLevel }} 级</template>
                <template v-if="member.summonerLevel"> · Lv{{ member.summonerLevel }}</template>
              </span>
            </div>
            <div v-if="member.tags?.length" class="mt-1 flex gap-1 overflow-hidden">
              <span
                v-for="t in member.tags"
                :key="t"
                class="shrink-0 rounded bg-panel-3 px-1.5 text-[10px] leading-[16px]"
                :class="/连败|高死亡|低等级|非常用|胜率低/.test(t) ? 'text-loss' : 'text-good'"
                >{{ t }}</span
              >
            </div>
          </div>
          <div class="shrink-0 text-center" :title="ratingTitle">
            <span v-if="label" class="inline-block rounded-md px-2.5 py-0.5 text-xs font-bold" :class="style?.badge">{{ label }}</span>
            <span v-else class="inline-block rounded-md bg-panel-3 px-2.5 py-0.5 text-xs text-muted">暂无评价</span>
            <p v-if="member.rating?.valid" class="num mt-0.5 text-[10px] text-muted">{{ member.rating.score.toFixed(0) }} 分</p>
          </div>
        </div>

      </div>
    </div>

    <div class="mt-2.5 grid grid-cols-[1.05fr_1fr] items-center gap-x-4 border-t border-line pt-2.5">
      <div class="space-y-1.5">
        <RankLine :rank="member.solo" :size="34" />
        <p v-if="member.flex?.tierKey" class="num pl-[44px] text-[11px] text-muted">
          灵活组排 <span class="text-ink">{{ member.flex.tier }} {{ member.flex.division }}</span>
        </p>
      </div>

      <div v-if="hasStats" class="space-y-1.5">
        <div class="flex items-baseline gap-2">
          <span class="num text-base font-bold" :class="rateColor">{{ formatPercent(stats.winRate) }}</span>
          <span class="num text-[11px] text-muted">近{{ stats.games }}场 {{ stats.wins }}胜{{ stats.games - stats.wins }}负</span>
        </div>
        <MatchDots :matches="member.recentMatches ?? []" />
        <p class="num whitespace-nowrap text-[11px] text-muted">
          KDA <b class="text-ink">{{ kdaRatio(stats.avgKills, stats.avgDeaths, stats.avgAssists) }}</b>
          <template v-if="stats.champGames > 0"> · 本英雄 {{ stats.champGames }}场 {{ formatPercent(stats.champWinRate) }}</template>
          <span v-if="stats.streak >= 3" class="text-win"> · {{ stats.streak }}连胜</span>
          <span v-else-if="stats.streak <= -3" class="text-loss"> · {{ -stats.streak }}连败</span>
        </p>
      </div>
      <p v-else class="self-center text-xs text-muted">{{ hasPuuid ? "暂无排位记录" : "信息未公开" }}</p>
    </div>
  </button>
</template>
