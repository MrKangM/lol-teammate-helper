<script lang="ts">
export const description = "A sidebar that collapses to icons."
export const iframeHeight = "800px"
export const containerClass = "w-full h-full"
</script>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue"
import AppSidebar from "@/components/AppSidebar.vue"
import HelloWorld from "@/components/HelloWorld.vue"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Separator } from "@/components/ui/separator"
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { GetImgSrc, Greet } from "../wailsjs/go/main/App"
import { IplayerBaseData } from "@/interface/baseData"
import { GetPlayerRankData } from "../wailsjs/go/controller/PlayerController"
import type { IRankedStats } from "@/interface/rankData"

const playerData = ref<IplayerBaseData | null>(null)
const rankData = ref<IRankedStats>()
const activePanel = ref<"dashboard" | "helloWorld">("dashboard")
const isLoading = ref(true)

const hasPlayerData = computed(() => playerData.value !== null)
const showClientPrompt = computed(() => isLoading.value || !hasPlayerData.value)

async function waitForBackendBridge(maxAttempts = 50, delay = 100) {
  let attempts = 0
  while (!(globalThis as any).go?.main?.App) {
    if (attempts >= maxAttempts) {
      return false
    }
    await new Promise((resolve) => setTimeout(resolve, delay))
    attempts += 1
  }
  return true
}
onMounted(() => {
  loadPlayerData()
})

const parsePlayerData = (raw: unknown): Partial<IplayerBaseData> => {
  if (typeof raw === "string") {
    try {
      return JSON.parse(raw) as Partial<IplayerBaseData>
    } catch (error) {
      console.warn("failed to parse player data", error)
      return {}
    }
  }

  if (typeof raw === "object" && raw !== null) {
    return raw as Partial<IplayerBaseData>
  }

  return {}
}

async function loadPlayerData() {
  isLoading.value = true
  try {
    const bridgeReady = await waitForBackendBridge()
    if (!bridgeReady) {
      console.warn('backend bridge unavailable; run "wails dev" or start the built app to enable backend APIs.')
      return
    }

    const rawBaseData = await Greet("test")
    const parsedBaseData = parsePlayerData(rawBaseData)
    const hasBaseData = Object.keys(parsedBaseData).length > 0

    if (!hasBaseData) {
      playerData.value = null
      rankData.value = undefined
      console.warn("player data unavailable; waiting for League client")
      return
    }

    const baseData = parsedBaseData as IplayerBaseData
    let iconImgSrc: string | undefined
    const profileIconId = typeof baseData.profileIconId === "number" ? baseData.profileIconId : undefined

    try {
      const iconIdForRequest = profileIconId ?? 0
      const iconResponse = await GetImgSrc(iconIdForRequest)
      if (typeof iconResponse === "string" && iconResponse.length > 0) {
        iconImgSrc = iconResponse
      }
    } catch (iconError) {
      console.error("failed to fetch profile icon", iconError)
    }

    playerData.value = {
      ...baseData,
      iconImgSrc,
    }

    console.log("player data", playerData.value)

    const rankIdentifier =
      playerData.value.puuid ??
      playerData.value.accountId ??
      (playerData.value.summonerId != null ? String(playerData.value.summonerId) : undefined)

    if (rankIdentifier) {
      try {
        rankData.value = await GetPlayerRankData(rankIdentifier)
        console.log("player rank data", rankData.value)
      } catch (rankError) {
        console.error("failed to fetch player rank data", rankError)
      }
    } else {
      rankData.value = undefined
      console.warn("missing identifier for rank data")
    }
  } catch (error) {
    console.error("failed to initialise player data", error)
    playerData.value = null
    rankData.value = undefined
  } finally {
    isLoading.value = false
  }
}

type SidebarNavigatePayload = {
  parent: {
    title: string
    url: string
  }
  item: {
    title: string
    url: string
  }
}

const handleSidebarNavigate = (payload: SidebarNavigatePayload) => {
  if (payload.item.url === "hello-world") {
    activePanel.value = "helloWorld"
  } else {
    activePanel.value = "dashboard"
  }
}
</script>

<template>
  <div class="h-full w-full">
    <div
      v-if="showClientPrompt"
      class="flex h-full flex-col items-center justify-center gap-6 bg-muted/50 p-8 text-center"
    >
      <div class="flex flex-col items-center gap-4">
        <div
          class="h-12 w-12 animate-spin rounded-full border-4 border-slate-300 border-t-transparent"
          aria-hidden="true"
        />
        <p class="text-lg font-semibold text-foreground">请打开游戏客户端后再打开本软件</p>
        <p class="text-sm text-muted-foreground">正在尝试获取召唤师数据，请确保英雄联盟客户端已登录。</p>
      </div>
    </div>
    <SidebarProvider v-else>
      <AppSidebar
        :avatar-src="playerData?.iconImgSrc"
        :player-data="playerData ?? undefined"
        :rank-data="rankData"
        @navigate="handleSidebarNavigate"
      />
      <SidebarInset>
        <header
          class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12"
        >
          <div class="flex items-center gap-2 px-4">
            <SidebarTrigger class="-ml-1" />
            <Separator orientation="vertical" class="mr-2 h-4" />
            <Breadcrumb>
              <BreadcrumbList>
                <BreadcrumbItem class="hidden md:block">
                  <BreadcrumbLink href="#">
                    Building Your Application
                  </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator class="hidden md:block" />
                <BreadcrumbItem>
                  <BreadcrumbPage>Data Fetching</BreadcrumbPage>
                </BreadcrumbItem>
              </BreadcrumbList>
            </Breadcrumb>
          </div>
        </header>
        <div class="flex flex-1 flex-col gap-4 p-4 pt-0">
          <template v-if="activePanel === 'helloWorld'">
            <HelloWorld />
          </template>
          <template v-else>
            <div class="grid auto-rows-min gap-4 md:grid-cols-3">
              <div class="aspect-video rounded-xl bg-muted/50" />
              <div class="aspect-video rounded-xl bg-muted/50" />
              <div class="aspect-video rounded-xl bg-muted/50" />
            </div>
            <div class="min-h-[100vh] flex-1 rounded-xl bg-muted/50 md:min-h-min" />
          </template>
        </div>
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>