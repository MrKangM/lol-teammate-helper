<script setup lang="ts">
import { computed } from "vue"
import type { types } from "../../wailsjs/go/models"
import { formatPercent, kdaRatio } from "@/lib/format"
import MatchDots from "@/components/MatchDots.vue"
import MatchList from "@/components/MatchList.vue"
import WinRing from "@/components/WinRing.vue"

// Shared body of the player drawer and the career page: form summary,
// averages, champion pool and the clickable ranked history.
const props = defineProps<{ member: types.TeamMemberSummary }>()

const stats = computed(() => props.member.stats)
const hasStats = computed(() => (stats.value?.games ?? 0) > 0)
const pool = computed(() => props.member.pool ?? [])

const tiles = computed(() => {
  const s = stats.value
  if (!s) return []
  return [
    { name: "KDA", value: kdaRatio(s.avgKills, s.avgDeaths, s.avgAssists), hint: `${s.avgKills.toFixed(1)} / ${s.avgDeaths.toFixed(1)} / ${s.avgAssists.toFixed(1)}` },
    { name: "场均补刀", value: s.avgCs.toFixed(0), hint: "含野怪" },
    { name: "场均伤害", value: s.avgDamage >= 1000 ? `${(s.avgDamage / 1000).toFixed(1)}k` : s.avgDamage.toFixed(0), hint: "对英雄" },
    { name: "场均视野", value: s.avgVision.toFixed(1), hint: "视野得分" },
  ]
})
</script>

<template>
  <div class="space-y-4">
    <section v-if="hasStats && stats" class="card flex items-center gap-5 p-4">
      <WinRing :rate="stats.winRate" :size="92" :label="`近${stats.games}场`" />
      <div class="min-w-0 flex-1 space-y-2.5">
        <p class="num text-sm text-ink">
          {{ stats.wins }}胜 {{ stats.games - stats.wins }}负
          <span v-if="stats.streak >= 2" class="ml-2 text-win">{{ stats.streak }}连胜</span>
          <span v-else-if="stats.streak <= -2" class="ml-2 text-loss">{{ -stats.streak }}连败</span>
        </p>
        <MatchDots :matches="member.recentMatches ?? []" />
        <p v-if="stats.champGames > 0" class="num text-xs text-muted">本英雄 {{ stats.champGames }} 场 · 胜率 {{ formatPercent(stats.champWinRate) }}</p>
      </div>
      <div class="grid grid-cols-2 gap-x-6 gap-y-2">
        <div v-for="t in tiles" :key="t.name" class="leading-tight">
          <p class="text-[11px] text-muted">{{ t.name }}</p>
          <p class="num text-lg font-bold text-ink">{{ t.value }}</p>
        </div>
      </div>
    </section>
    <p v-else class="card p-4 text-sm text-muted">暂无单双排或灵活组排的对局记录。</p>

    <section v-if="pool.length" class="card p-4">
      <p class="section-title mb-3">常用英雄</p>
      <div class="grid grid-cols-5 gap-2">
        <div v-for="c in pool" :key="c.championId" class="rounded-lg bg-panel-2 p-2 text-center">
          <img :src="c.championIcon" :alt="c.championName" class="mx-auto size-11 rounded-lg" />
          <p class="mt-1 truncate text-xs font-semibold text-ink">{{ c.championName }}</p>
          <p class="num text-[11px]" :class="c.winRate >= 52 ? 'text-win' : 'text-loss'">{{ formatPercent(c.winRate) }} · {{ c.games }}场</p>
          <p class="num text-[10px] text-muted">KDA {{ c.kda.toFixed(1) }}</p>
        </div>
      </div>
    </section>

    <section class="card p-4">
      <p class="section-title mb-3">近期排位战绩 · 点击查看结局图</p>
      <MatchList :matches="member.recentMatches ?? []" :puuid="member.puuid" />
    </section>
  </div>
</template>
