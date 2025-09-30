<script setup lang="ts">
import { computed } from "vue"
import type { SidebarProps } from '@/components/ui/sidebar'
import type { IplayerBaseData } from '@/interface/baseData'

import {
  AudioWaveform,
  BookOpen,
  Bot,
  Command,
  Frame,
  GalleryVerticalEnd,
  Map,
  PieChart,
  Settings2,
  SquareTerminal,
} from "lucide-vue-next"
import NavMain from '@/components/NavMain.vue'
import NavProjects from '@/components/NavProjects.vue'
import NavUser from '@/components/NavUser.vue'
import TeamSwitcher from '@/components/TeamSwitcher.vue'

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '@/components/ui/sidebar'

type AppSidebarProps = SidebarProps & {
  avatarSrc?: string
  playerData?: IplayerBaseData
}

const props = withDefaults(defineProps<AppSidebarProps>(), {
  collapsible: "icon",
  avatarSrc: "",
  playerData: undefined,
})

const sanitize = (value?: string | null) => {
  const trimmed = value?.trim()
  return trimmed && trimmed.length > 0 ? trimmed : undefined
}

const resolvePlayerName = (player?: IplayerBaseData) => {
  if (!player) return undefined

  const displayName = sanitize(player.displayName)
  if (displayName) return displayName

  const gameName = sanitize(player.gameName)
  const tagLine = sanitize(player.tagLine)
  if (gameName && tagLine) return `${gameName}#${tagLine}`
  if (gameName) return gameName

  return sanitize(player.internalName)
}

const resolvePlayerIdentifier = (player?: IplayerBaseData) => {
  if (!player) return undefined
  return sanitize(player.puuid) ?? sanitize(player.accountId)
}

const defaultSidebarData = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  teams: [
    {
      name: "Acme Inc",
      logo: GalleryVerticalEnd,
      plan: "Enterprise",
    },
    {
      name: "Acme Corp.",
      logo: AudioWaveform,
      plan: "Startup",
    },
    {
      name: "Evil Corp.",
      logo: Command,
      plan: "Free",
    },
  ],
  navMain: [
    {
      title: "Playground",
      url: "#",
      icon: SquareTerminal,
      isActive: true,
      items: [
        {
          title: "History",
          url: "#",
        },
        {
          title: "Starred",
          url: "#",
        },
        {
          title: "Settings",
          url: "#",
        },
      ],
    },
    {
      title: "Models",
      url: "#",
      icon: Bot,
      items: [
        {
          title: "Genesis",
          url: "#",
        },
        {
          title: "Explorer",
          url: "#",
        },
        {
          title: "Quantum",
          url: "#",
        },
      ],
    },
    {
      title: "Documentation",
      url: "#",
      icon: BookOpen,
      items: [
        {
          title: "Introduction",
          url: "#",
        },
        {
          title: "Get Started",
          url: "#",
        },
        {
          title: "Tutorials",
          url: "#",
        },
        {
          title: "Changelog",
          url: "#",
        },
      ],
    },
    {
      title: "Settings",
      url: "#",
      icon: Settings2,
      items: [
        {
          title: "General",
          url: "#",
        },
        {
          title: "Team",
          url: "#",
        },
        {
          title: "Billing",
          url: "#",
        },
        {
          title: "Limits",
          url: "#",
        },
      ],
    },
  ],
  projects: [
    {
      name: "Design Engineering",
      url: "#",
      icon: Frame,
    },
    {
      name: "Sales & Marketing",
      url: "#",
      icon: PieChart,
    },
    {
      name: "Travel",
      url: "#",
      icon: Map,
    },
  ],
}

const resolvedUser = computed(() => {
  const player = props.playerData
  const fallback = defaultSidebarData.user

  const name = resolvePlayerName(player) ?? fallback.name
  const identifier = resolvePlayerIdentifier(player) ?? fallback.email
  const avatar = sanitize(player?.iconImgSrc) ?? props.avatarSrc ?? fallback.avatar

  return {
    ...fallback,
    name,
    email: identifier,
    avatar,
  }
})

const sidebarData = computed(() => ({
  ...defaultSidebarData,
  user: resolvedUser.value,
}))
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
<!--      <TeamSwitcher :teams="sidebarData.teams" />-->
      <NavUser :user="sidebarData.user" />
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="sidebarData.navMain" />
      <NavProjects :projects="sidebarData.projects" />
    </SidebarContent>
    <SidebarFooter>
<!--      <NavUser :user="sidebarData.user" />-->
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>

