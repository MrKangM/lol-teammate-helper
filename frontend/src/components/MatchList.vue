<script setup lang="ts">
import { ref } from "vue"
import { GetGameDetail } from "../../wailsjs/go/controller/MatchHistory"
import type { types } from "../../wailsjs/go/models"
import { compactNumber, formatDuration, formatKda, kdaRatio, positionLabel, queueLabel, timeAgo } from "@/lib/format"
import GameDetailPanel from "@/components/GameDetailPanel.vue"

const props = defineProps<{ matches: types.RecentMatchSummary[]; puuid: string }>()

const openGameId = ref<number | null>(null)
const loadingId = ref<number | null>(null)
const errorText = ref("")
const details = ref<Record<number, types.GameDetail>>({})

const toggle = async (gameId: number) => {
  errorText.value = ""
  if (openGameId.value === gameId) {
    openGameId.value = null
    return
  }
  openGameId.value = gameId
  if (details.value[gameId]) return

  loadingId.value = gameId
  try {
    details.value[gameId] = await GetGameDetail(gameId, props.puuid)
  } catch (error) {
    console.error("failed to load game detail", error)
    errorText.value = "对局详情加载失败，可能已超出客户端保留的历史。"
  } finally {
    loadingId.value = null
  }
}
</script>

<template>
  <div class="space-y-1">
    <p v-if="matches.length === 0" class="py-3 text-center text-xs text-muted">暂无排位对局记录</p>

    <div v-for="m in matches" :key="m.gameId">
      <button
        type="button"
        class="grid w-full grid-cols-[4px_30px_minmax(110px,1fr)_88px_58px_minmax(96px,auto)_minmax(104px,auto)_64px] items-center gap-2.5 overflow-hidden rounded-sm bg-panel-2/70 py-1.5 pr-3 text-left text-xs transition-colors hover:bg-panel-3"
        @click="toggle(m.gameId)"
      >
        <span class="h-full min-h-8 w-1" :class="m.win ? 'bg-win' : 'bg-loss'" />
        <img v-if="m.championIcon" :src="m.championIcon" :alt="m.championName" class="size-[30px] rounded" />
        <div v-else class="size-[30px] rounded bg-panel-3" />

        <span class="truncate text-ink">
          <b :class="m.win ? 'text-win' : 'text-loss'">{{ m.win ? "胜" : "负" }}</b>
          {{ m.championName || `英雄${m.championId}` }}
          <span class="text-muted">· {{ positionLabel(m.position) }}</span>
        </span>
        <span class="num font-semibold text-gold-bright">{{ formatKda(m.kills, m.deaths, m.assists) }}</span>
        <span class="num text-muted">{{ kdaRatio(m.kills, m.deaths, m.assists) }} KDA</span>
        <span class="num whitespace-nowrap text-muted">补刀 {{ m.cs }} · 伤害 {{ compactNumber(m.damage) }}</span>
        <span class="whitespace-nowrap text-muted">{{ queueLabel(m.queueId) }} {{ formatDuration(m.gameDuration) }}</span>
        <span class="text-right text-muted">{{ timeAgo(m.gameCreation) }}</span>
      </button>

      <div v-if="openGameId === m.gameId" class="py-2">
        <p v-if="loadingId === m.gameId" class="text-xs text-muted">正在加载对局详情…</p>
        <p v-else-if="errorText" class="text-xs text-loss">{{ errorText }}</p>
        <GameDetailPanel v-else-if="details[m.gameId]" :detail="details[m.gameId]" />
      </div>
    </div>
  </div>
</template>
