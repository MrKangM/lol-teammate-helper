<script setup lang="ts">
import { computed, onMounted, watch } from "vue"
import { RefreshCw } from "lucide-vue-next"
import { career, careerError, careerLoading, loadCareer, status, summoner, summonerIcon } from "@/lib/store"
import { RATING_STYLES, formatPercent, kdaRatio } from "@/lib/format"
import MatchList from "@/components/MatchList.vue"
import MatchDots from "@/components/MatchDots.vue"
import RankLine from "@/components/RankLine.vue"
import WinRing from "@/components/WinRing.vue"

const stats = computed(() => career.value?.stats)
const hasStats = computed(() => (stats.value?.games ?? 0) > 0)
const label = computed(() => career.value?.rating?.label ?? "")
const style = computed(() => RATING_STYLES[label.value])
const pool = computed(() => career.value?.pool ?? [])

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

onMounted(() => {
  if (status.value?.connected && !career.value) loadCareer()
})
watch(
  () => status.value?.connected,
  (connected) => {
    if (connected && !career.value) loadCareer()
  },
)
</script>

<template>
  <div class="space-y-4">
    <div v-if="!status?.connected && !career" class="panel px-6 py-16 text-center">
      <p class="title-gold text-lg">尚未连接到英雄联盟客户端</p>
      <p class="mt-2 text-sm text-muted">连接后这里会显示你自己的段位和近 20 场排位战绩。</p>
    </div>

    <template v-else>
      <section class="panel panel-gold flex flex-wrap items-center gap-5 p-4">
        <div class="relative">
          <img v-if="summonerIcon" :src="summonerIcon" alt="" class="size-[72px] rounded-full border-2 border-gold object-cover" />
          <div v-else class="size-[72px] rounded-full border-2 border-gold-dim bg-panel-3" />
          <span v-if="summoner?.summonerLevel" class="num absolute -bottom-1 left-1/2 -translate-x-1/2 rounded-sm border border-line-gold bg-bg px-1.5 text-[11px] text-gold">{{ summoner.summonerLevel }}</span>
        </div>
        <div class="min-w-0 flex-1">
          <p class="title-gold truncate text-xl">
            {{ summoner?.gameName || summoner?.displayName || "召唤师" }}<span v-if="summoner?.tagLine" class="text-sm font-normal text-muted"> #{{ summoner.tagLine }}</span>
          </p>
          <p class="text-xs text-muted">{{ summoner?.region }}</p>
          <span v-if="label" class="mt-2 inline-block rounded-sm border px-2 py-0.5 text-xs font-bold" :class="style?.badge" title="按近期战绩估算的综合水平">{{ label }}</span>
        </div>
        <button type="button" class="btn" :disabled="careerLoading" @click="loadCareer">
          <RefreshCw class="size-3.5" :class="careerLoading ? 'animate-spin' : ''" />
          刷新
        </button>
      </section>

      <p v-if="careerError" class="panel border-loss/40 px-4 py-3 text-sm text-loss">{{ careerError }}</p>
      <p v-else-if="careerLoading && !career" class="py-10 text-center text-sm text-muted">正在加载生涯数据…</p>

      <template v-if="career">
        <div class="grid gap-4 lg:grid-cols-2">
          <section class="panel p-4">
            <p class="eyebrow mb-3">排位段位</p>
            <div class="space-y-4">
              <RankLine :rank="career.solo" :size="56" />
              <RankLine :rank="career.flex" :size="56" />
            </div>
          </section>

          <section class="panel p-4">
            <p class="eyebrow mb-3">近 {{ stats?.games ?? 0 }} 场排位</p>
            <div v-if="hasStats && stats" class="flex items-center gap-5">
              <WinRing :rate="stats.winRate" :size="92" />
              <div class="min-w-0 flex-1 space-y-2">
                <p class="num text-sm text-ink">{{ stats.wins }}胜 {{ stats.games - stats.wins }}负</p>
                <MatchDots :matches="career.recentMatches ?? []" />
                <p v-if="stats.streak >= 2" class="text-xs text-win">当前 {{ stats.streak }} 连胜</p>
                <p v-else-if="stats.streak <= -2" class="text-xs text-loss">当前 {{ -stats.streak }} 连败</p>
              </div>
            </div>
            <p v-else class="text-sm text-muted">暂无单双排或灵活组排的对局记录。</p>
          </section>
        </div>

        <div v-if="hasStats" class="grid grid-cols-2 gap-3 sm:grid-cols-4">
          <div v-for="t in tiles" :key="t.name" class="panel px-4 py-3">
            <p class="text-xs text-muted">{{ t.name }}</p>
            <p class="num mt-1 text-2xl font-bold text-gold-bright">{{ t.value }}</p>
            <p class="num text-[11px] text-muted">{{ t.hint }}</p>
          </div>
        </div>

        <section v-if="pool.length" class="panel p-4">
          <p class="eyebrow mb-3">常用英雄</p>
          <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
            <div v-for="c in pool" :key="c.championId" class="flex items-center gap-3 rounded-md border border-line bg-panel-2/60 p-2">
              <img :src="c.championIcon" :alt="c.championName" class="size-11 rounded-md border border-line-gold" />
              <div class="min-w-0 leading-tight">
                <p class="truncate text-sm font-semibold text-gold-bright">{{ c.championName }}</p>
                <p class="num text-xs" :class="c.winRate >= 55 ? 'text-win' : c.winRate < 45 ? 'text-loss' : 'text-muted'">{{ formatPercent(c.winRate) }} · {{ c.games }}场</p>
                <p class="num text-[11px] text-muted">KDA {{ c.kda.toFixed(1) }}</p>
              </div>
            </div>
          </div>
        </section>

        <section class="panel p-4">
          <p class="eyebrow mb-3">近期排位战绩 · 点击查看结局图</p>
          <MatchList :matches="career.recentMatches ?? []" :puuid="career.puuid" />
        </section>
      </template>
    </template>
  </div>
</template>
