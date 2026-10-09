<script setup lang="ts">
import { computed } from "vue"
import { ChevronDown } from "lucide-vue-next"
import type { types } from "../../wailsjs/go/models"
import { RATING_STYLES, formatPercent, kdaRatio } from "@/lib/format"
import MatchDots from "@/components/MatchDots.vue"
import MatchList from "@/components/MatchList.vue"
import PositionIcon from "@/components/PositionIcon.vue"
import RankLine from "@/components/RankLine.vue"

const props = defineProps<{ member: types.TeamMemberSummary; expanded: boolean }>()
const emit = defineEmits<{ (e: "toggle"): void }>()

const stats = computed(() => props.member.stats)
const hasStats = computed(() => (stats.value?.games ?? 0) > 0)
const hasPuuid = computed(() => Boolean(props.member.puuid))
const label = computed(() => props.member.rating?.label ?? "")
const style = computed(() => RATING_STYLES[label.value])
const ratingTitle = computed(() =>
  props.member.rating?.valid ? `综合评分 ${props.member.rating.score.toFixed(0)} / 100` : "数据不足，暂无评价",
)
const spells = computed(() => (props.member.spells ?? []).filter((s) => s && (s.icon || s.name)))
const pool = computed(() => (props.member.pool ?? []).slice(0, 3))
</script>

<template>
  <div class="panel overflow-hidden transition-colors" :class="expanded ? 'panel-gold' : ''">
    <button
      type="button"
      class="grid w-full grid-cols-[4px_34px_62px_minmax(160px,1.1fr)_minmax(180px,1fr)_minmax(230px,1.3fr)_minmax(90px,auto)_72px_16px] items-center gap-3 pr-3 text-left"
      :disabled="!hasPuuid"
      @click="emit('toggle')"
    >
      <span class="h-full min-h-[76px] w-1" :class="style?.strip ?? 'bg-line'" />

      <div class="flex justify-center"><PositionIcon :position="member.assignedPosition" /></div>

      <div class="relative size-[58px] py-0">
        <img v-if="member.championIcon" :src="member.championIcon" :alt="member.championName" class="size-[58px] rounded-md border border-line-gold" />
        <div v-else class="size-[58px] rounded-md border border-line bg-panel-3" />
        <div v-if="spells.length" class="absolute -bottom-1 -right-1.5 flex flex-col gap-px">
          <template v-for="(s, i) in spells" :key="i">
            <img v-if="s.icon" :src="s.icon" :alt="s.name" :title="s.name" class="size-[19px] rounded-sm border border-bg" />
            <span v-else class="grid size-[19px] place-items-center rounded-sm bg-panel-3 text-[9px] text-muted" :title="s.name">{{ s.name.slice(0, 1) }}</span>
          </template>
        </div>
      </div>

      <div class="min-w-0 py-2">
        <p class="truncate text-sm font-semibold text-gold-bright">
          {{ member.gameName || "未知召唤师" }}<span v-if="member.tagLine" class="font-normal text-muted"> #{{ member.tagLine }}</span>
        </p>
        <p class="truncate text-xs text-muted">
          {{ member.championName || "未选英雄" }}
          <template v-if="member.summonerLevel"> · Lv{{ member.summonerLevel }}</template>
          <template v-if="member.masteryLevel"> · 熟练 {{ member.masteryLevel }}级</template>
        </p>
        <div v-if="member.tags?.length" class="mt-1 flex flex-wrap gap-1">
          <span
            v-for="t in member.tags"
            :key="t"
            class="rounded-sm border border-line px-1 text-[10px] leading-4"
            :class="/连败|高死亡|低等级|非常用|胜率低/.test(t) ? 'text-loss/90' : 'text-teal'"
            >{{ t }}</span
          >
        </div>
      </div>

      <div class="space-y-1.5 py-2">
        <RankLine :rank="member.solo" :size="34" />
        <RankLine v-if="member.flex?.tierKey" :rank="member.flex" :size="22" class="opacity-80" />
      </div>

      <div v-if="hasStats" class="space-y-1.5 py-2">
        <div class="flex items-baseline gap-2">
          <span class="num text-sm font-bold" :class="stats.winRate >= 55 ? 'text-win' : stats.winRate < 45 ? 'text-loss' : 'text-gold-bright'">
            {{ formatPercent(stats.winRate) }}
          </span>
          <span class="num text-xs text-muted">近{{ stats.games }}场 {{ stats.wins }}胜{{ stats.games - stats.wins }}负</span>
        </div>
        <MatchDots :matches="member.recentMatches ?? []" />
        <p class="num whitespace-nowrap text-xs text-muted">
          KDA <b class="text-ink">{{ kdaRatio(stats.avgKills, stats.avgDeaths, stats.avgAssists) }}</b>
          <span v-if="stats.streak >= 3" class="text-win"> · {{ stats.streak }}连胜</span>
          <span v-else-if="stats.streak <= -3" class="text-loss"> · {{ -stats.streak }}连败</span>
        </p>
        <p v-if="stats.champGames > 0" class="num whitespace-nowrap text-[11px] text-muted">
          本英雄 {{ stats.champGames }}场 · 胜率 {{ formatPercent(stats.champWinRate) }}
        </p>
      </div>
      <p v-else class="text-xs text-muted">{{ hasPuuid ? "暂无排位记录" : "信息未公开" }}</p>

      <div class="flex gap-1" title="近期常用英雄">
        <img
          v-for="c in pool"
          :key="c.championId"
          :src="c.championIcon"
          :alt="c.championName"
          :title="`${c.championName} ${c.games}场 胜率${Math.round(c.winRate)}%`"
          class="size-7 rounded-sm border border-line"
        />
      </div>

      <div class="text-center">
        <span v-if="label" class="inline-block rounded-sm border px-2 py-0.5 text-xs font-bold" :class="style?.badge" :title="ratingTitle">{{ label }}</span>
        <span v-else class="text-xs text-muted" :title="ratingTitle">—</span>
        <p v-if="member.rating?.valid" class="num mt-0.5 text-[10px] text-muted">{{ member.rating.score.toFixed(0) }} 分</p>
      </div>

      <ChevronDown v-if="hasPuuid" class="size-4 text-gold-dim transition-transform" :class="expanded ? 'rotate-180' : ''" />
    </button>

    <div v-if="expanded" class="border-t border-line bg-bg/40 px-3 py-3">
      <MatchList :matches="member.recentMatches ?? []" :puuid="member.puuid" />
    </div>
  </div>
</template>
