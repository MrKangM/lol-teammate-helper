<script setup lang="ts">
import { computed, onMounted, watch } from "vue"
import { RefreshCw } from "lucide-vue-next"
import { career, careerError, careerLoading, loadCareer, status, summoner, summonerIcon } from "@/lib/store"
import { RATING_STYLES } from "@/lib/format"
import PlayerDetail from "@/components/PlayerDetail.vue"
import RankLine from "@/components/RankLine.vue"

const label = computed(() => career.value?.rating?.label ?? "")
const style = computed(() => RATING_STYLES[label.value])

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
  <div class="mx-auto max-w-[920px] space-y-4">
    <div v-if="!status?.connected && !career" class="card px-6 py-20 text-center">
      <p class="text-lg font-semibold text-ink">尚未连接到英雄联盟客户端</p>
      <p class="mt-2 text-sm text-muted">连接后这里会显示你自己的段位和近 20 场排位战绩。</p>
    </div>

    <template v-else>
      <section class="card flex flex-wrap items-center gap-5 p-5">
        <div class="relative">
          <img v-if="summonerIcon" :src="summonerIcon" alt="" class="size-[76px] rounded-2xl object-cover" />
          <div v-else class="size-[76px] rounded-2xl bg-panel-3" />
          <span v-if="summoner?.summonerLevel" class="num absolute -bottom-2 left-1/2 -translate-x-1/2 rounded-full bg-panel-3 px-2 py-0.5 text-[11px] font-semibold text-ink ring-2 ring-panel">{{ summoner.summonerLevel }}</span>
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-xl font-bold text-ink">
            {{ summoner?.gameName || summoner?.displayName || "召唤师" }}<span v-if="summoner?.tagLine" class="text-sm font-normal text-muted"> #{{ summoner.tagLine }}</span>
          </p>
          <p class="text-xs text-muted">{{ summoner?.region }}</p>
          <span v-if="label" class="mt-2 inline-block rounded-md px-2.5 py-0.5 text-xs font-bold" :class="style?.badge" title="按近期战绩估算的综合水平">{{ label }}</span>
        </div>
        <div v-if="career" class="flex gap-6">
          <RankLine :rank="career.solo" :size="52" />
          <RankLine :rank="career.flex" :size="52" />
        </div>
        <button type="button" class="btn" :disabled="careerLoading" @click="loadCareer">
          <RefreshCw class="size-3.5" :class="careerLoading ? 'animate-spin' : ''" />
          刷新
        </button>
      </section>

      <p v-if="careerError" class="card border-loss/40 px-4 py-3 text-sm text-loss">{{ careerError }}</p>
      <p v-else-if="careerLoading && !career" class="py-10 text-center text-sm text-muted">正在加载生涯数据…</p>

      <PlayerDetail v-if="career" :member="career" />
    </template>
  </div>
</template>
