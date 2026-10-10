<script setup lang="ts">
import type { types } from "../../wailsjs/go/models"

// One bar per game, oldest on the left, so the recent form reads left to right.
defineProps<{ matches: types.RecentMatchSummary[] }>()
</script>

<template>
  <div class="flex items-end gap-[2px]">
    <span
      v-for="(m, i) in [...matches].reverse()"
      :key="i"
      class="h-[14px] w-[6px] rounded-[2px]"
      :class="m.win ? 'bg-win' : 'bg-loss'"
      :title="`${m.championName} ${m.kills}/${m.deaths}/${m.assists}`"
    />
  </div>
</template>
