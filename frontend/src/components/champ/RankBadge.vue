<script setup lang="ts">
import { computed } from "vue"
import type { types } from "../../../wailsjs/go/models"
import { formatPercent } from "@/lib/format"
import RankEmblem from "./RankEmblem.vue"

const props = withDefaults(defineProps<{ rank: types.RankSummary; compact?: boolean }>(), { compact: false })

const ranked = computed(() => Boolean(props.rank?.tierKey))
const title = computed(() => {
  const r = props.rank
  if (!ranked.value) return "未定级"
  return `${r.tier}${r.division ? ` ${r.division}` : ""}`
})
const total = computed(() => (props.rank?.wins ?? 0) + (props.rank?.losses ?? 0))
</script>

<template>
  <div class="flex items-center gap-2">
    <RankEmblem :tier-key="rank?.tierKey" :size="compact ? 24 : 34" />
    <div class="min-w-0 leading-tight">
      <p class="truncate text-sm font-medium text-foreground">
        {{ title }}
        <span v-if="ranked" class="text-xs font-normal text-muted-foreground">{{ rank.leaguePoints }} LP</span>
      </p>
      <p class="truncate text-xs text-muted-foreground">
        {{ rank?.queueName }}
        <template v-if="total > 0">· {{ rank.wins }}胜{{ rank.losses }}负 {{ formatPercent(rank.winRate) }}</template>
      </p>
    </div>
  </div>
</template>
