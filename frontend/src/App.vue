<script lang="ts">
export const description = "A sidebar that collapses to icons."
export const iframeHeight = "800px"
export const containerClass = "w-full h-full"
</script>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue"
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
import { IPlayerBaseData } from "@/interface/baseData"
import { GetPlayerRankData } from "../wailsjs/go/controller/PlayerController"
import type { IRankedStats } from "@/interface/rankData"

const playerData = ref<IPlayerBaseData | null>(null)
const rankData = ref<IRankedStats>()
const activePanel = ref<"dashboard" | "helloWorld">("dashboard")
const isLoading = ref(true)
const retryTimer = ref<number | null>(null)
const retryCount = ref(0)

const hasPlayerData = computed(() => playerData.value !== null)
const showClientPrompt = computed(() => isLoading.value || !hasPlayerData.value)

const RETRY_DELAY_MS = 2000
const overlayTitle = "请打开游戏客户端后再启动本软件"
const overlaySubtitle = "正在尝试获取召唤师资料，请确认英雄联盟客户端已登录"
const overlayRetryPrefix = "已尝试重新连接"
const overlayRetrySuffix = "次，将继续自动重试"

onMounted(() => {
  loadPlayerData()
})

onBeforeUnmount(() => {
  if (retryTimer.value !== null) {
    clearTimeout(retryTimer.value)
    retryTimer.value = null
  }
})

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

const parsePlayerData = (raw: unknown): Partial<IPlayerBaseData> => {
  if (typeof raw === "string") {
    try {
      return JSON.parse(raw) as Partial<IPlayerBaseData>
    } catch (error) {
      console.warn("failed to parse player data", error)
      return {}
    }
  }

  if (typeof raw === "object" && raw !== null) {
    return raw as Partial<IPlayerBaseData>
  }

  return {}
}

function scheduleReload() {
  if (retryTimer.value !== null) {
    clearTimeout(retryTimer.value)
  }

  retryTimer.value = window.setTimeout(() => {
    retryCount.value += 1
    loadPlayerData()
  }, RETRY_DELAY_MS)
}

async function loadPlayerData() {
  isLoading.value = true

  try {
    const bridgeReady = await waitForBackendBridge()
    if (!bridgeReady) {
      console.warn("backend bridge unavailable; run wails dev or start the packaged app.")
      scheduleReload()
      return
    }

    let baseData: IPlayerBaseData
    try {
      baseData = await Greet("load-player")
    } catch (requestError) {
      console.error("failed to fetch summoner data", requestError)
      scheduleReload()
      return
    }

    // const parsedBaseData = rawBaseData
    // const hasBaseData = Object.keys(parsedBaseData).length > 0

    // if (!hasBaseData) {
    //   console.warn("player data unavailable; waiting for League client")
    //   scheduleReload()
    //   return
    // }

     // baseData = parsedBaseData as IplayerBaseData
    let iconImgSrc: string | undefined
    const profileIconId = typeof baseData.profileIconId === "number" ? baseData.profileIconId : undefined

    try {
      const iconIdForRequest = profileIconId ?? 0
      const iconResponse = await GetImgSrc(iconIdForRequest)
      if (typeof iconResponse === "string" && iconResponse.startsWith("data:image")) {
        iconImgSrc = iconResponse
      }
    } catch (iconError) {
      console.error("failed to fetch profile icon", iconError)
    }

    playerData.value = {
      ...baseData,
      iconImgSrc,
    }

    const puuid = playerData.value?.puuid ?? ""

    if (puuid.length > 0) {
      try {
        rankData.value = await GetPlayerRankData(puuid)
        console.log(`API获取的排位数据:${JSON.stringify(rankData)}`)
      } catch (rankError) {
        console.error("failed to fetch player rank data", rankError)
      }
    } else {
      rankData.value = undefined
      console.warn("missing puuid for rank data")
    }
    if (retryTimer.value !== null) {
      clearTimeout(retryTimer.value)
      retryTimer.value = null
    }

    retryCount.value = 0
    isLoading.value = false
  } catch (error) {
    console.error("failed to initialise player data", error)
    playerData.value = null
    rankData.value = undefined
    scheduleReload()
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
  <div class="relative h-full w-full">
    <div
      v-if="showClientPrompt"
      class="absolute inset-0 z-50 flex flex-col items-center justify-center gap-6 bg-background/95 px-6 text-center backdrop-blur-sm"
    >
      <div class="flex flex-col items-center gap-4">
        <div
          class="h-14 w-14 animate-spin rounded-full border-4 border-muted-foreground/40 border-t-primary"
          aria-hidden="true"
        />
        <div class="space-y-2">
          <p class="text-xl font-semibold text-foreground">
            {{ overlayTitle }}
          </p>
          <p class="text-sm text-muted-foreground">
            {{ overlaySubtitle }}
          </p>
          <p v-if="retryCount > 0" class="text-xs text-muted-foreground/70">
            {{ overlayRetryPrefix }} {{ retryCount }} {{ overlayRetrySuffix }}
          </p>
        </div>
      </div>
    </div>

    <SidebarProvider v-else :default-open="true">
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
            <HelloWorld :player-data="playerData" :rank-data="rankData" />
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






