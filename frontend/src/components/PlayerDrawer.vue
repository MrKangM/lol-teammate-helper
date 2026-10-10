<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted } from "vue"
import { X } from "lucide-vue-next"
import type { types } from "../../wailsjs/go/models"
import { RATING_STYLES } from "@/lib/format"
import PlayerDetail from "@/components/PlayerDetail.vue"
import PositionIcon from "@/components/PositionIcon.vue"
import RankLine from "@/components/RankLine.vue"

const props = defineProps<{ member: types.TeamMemberSummary }>()
const emit = defineEmits<{ (e: "close"): void }>()

const label = computed(() => props.member.rating?.label ?? "")
const style = computed(() => RATING_STYLES[label.value])

const onKey = (e: KeyboardEvent) => {
  if (e.key === "Escape") emit("close")
}
onMounted(() => window.addEventListener("keydown", onKey))
onBeforeUnmount(() => window.removeEventListener("keydown", onKey))
</script>

<template>
  <div class="fixed inset-0 z-40 flex justify-end">
    <div class="absolute inset-0 bg-black/55 backdrop-blur-[2px]" @click="emit('close')" />
    <aside class="relative flex h-full w-[min(800px,100%)] flex-col border-l border-line bg-bg shadow-2xl">
      <header class="flex items-center gap-4 border-b border-line px-5 py-4">
        <div class="relative">
          <img v-if="member.championIcon" :src="member.championIcon" :alt="member.championName" class="size-14 rounded-xl" />
          <div v-else class="size-14 rounded-xl bg-panel-3" />
          <span class="absolute -bottom-1 -left-1 grid size-6 place-items-center rounded-full border border-line bg-panel"><PositionIcon :position="member.assignedPosition" :size="14" /></span>
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-base font-semibold text-ink">
            {{ member.gameName || "未知召唤师" }}<span v-if="member.tagLine" class="font-normal text-muted"> #{{ member.tagLine }}</span>
          </p>
          <p class="text-xs text-muted">
            {{ member.championName }}<template v-if="member.summonerLevel"> · Lv{{ member.summonerLevel }}</template
            ><template v-if="member.masteryLevel"> · 熟练 {{ member.masteryLevel }} 级 {{ Math.round(member.masteryPoints / 1000) }}k</template>
          </p>
          <div v-if="member.tags?.length" class="mt-1.5 flex flex-wrap gap-1">
            <span v-for="t in member.tags" :key="t" class="rounded bg-panel-3 px-1.5 text-[11px] leading-5 text-muted">{{ t }}</span>
          </div>
        </div>
        <div class="space-y-1.5">
          <RankLine :rank="member.solo" :size="32" />
          <RankLine v-if="member.flex?.tierKey" :rank="member.flex" :size="24" />
        </div>
        <div v-if="label" class="text-center">
          <span class="inline-block rounded-md px-2.5 py-1 text-sm font-bold" :class="style?.badge">{{ label }}</span>
          <p v-if="member.rating?.valid" class="num mt-1 text-[11px] text-muted">{{ member.rating.score.toFixed(0) }} 分</p>
        </div>
        <button type="button" class="grid size-8 place-items-center rounded-lg text-muted hover:bg-panel-3 hover:text-ink" aria-label="关闭" @click="emit('close')">
          <X class="size-4" />
        </button>
      </header>
      <div class="min-h-0 flex-1 overflow-y-auto p-5">
        <PlayerDetail :member="member" />
      </div>
    </aside>
  </div>
</template>
