<script setup lang="ts">
import { computed } from "vue"
import type { types } from "../../wailsjs/go/models"
import { formatPercent } from "@/lib/format"
import RankEmblem from "@/components/RankEmblem.vue"

const props = withDefaults(defineProps<{ rank: types.RankSummary; size?: number; showQueue?: boolean }>(), {
  size: 40,
  showQueue: true,
})

const ranked = computed(() => Boolean(props.rank?.tierKey))
const title = computed(() => (ranked.value ? `${props.rank.tier}${props.rank.division ? ` ${props.rank.division}` : ""}` : "未定级"))
const total = computed(() => (props.rank?.wins ?? 0) + (props.rank?.losses ?? 0))
</script>

<template>
  <div class="flex items-center gap-2">
    <RankEmblem :tier-key="rank?.tierKey" :size="size" class="shrink-0" />
    <div class="min-w-0 leading-tight">
      <p class="truncate text-[13px] font-semibold" :class="ranked ? 'text-gold-bright' : 'text-muted'">
        {{ title }}
        <span v-if="ranked" class="num text-[11px] font-normal text-muted">{{ rank.leaguePoints }} LP</span>
      </p>
      <p class="num whitespace-nowrap text-[11px] text-muted">
        <template v-if="showQueue">{{ rank?.queueName }}</template>
        <template v-if="total > 0">{{ showQueue ? " · " : "" }}{{ rank.wins }}胜{{ rank.losses }}负 {{ formatPercent(rank.winRate) }}</template>
      </p>
    </div>
  </div>
</template>
