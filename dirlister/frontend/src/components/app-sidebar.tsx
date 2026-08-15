"use client";

import {
  FileText,
  Folder,
  Home,
  Image as ImageIcon,
  Music,
  Settings,
  Star,
  Trash2,
  Video,
  Signpost,
  UserRound,
  Plus,
  XIcon
} from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Link } from "react-router";

const mainItems = [
  { title: "Home", icon: Home, link: "/" },
  { title: "Route Rules", icon: Signpost, link: "/" },
  { title: "Create", icon: Plus, link: "/.@/create" },

  { title: "Account", icon: UserRound, link: "/.@/account" },
];



export function AppSidebar() {
  return (
    <Sidebar variant="inset">
      <SidebarContent>

        <SidebarGroup>
          <SidebarGroupContent>
            <div className="flex flex-row">
              <div className="flex flex-col justify-center items-center ">
                <XIcon className="w-[2rem] h-[2rem]"></XIcon>
              </div>
              <div className="text-2xl font-semibold">Dirlist</div>
            </div>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarGroup>
          <SidebarGroupLabel>Navigation</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {mainItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton >
                    <Link to={item.link} className="flex flex-row gap-4">
                      <item.icon className="h-4 w-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

      </SidebarContent>
    </Sidebar>
  );
}