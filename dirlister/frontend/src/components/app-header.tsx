"use client";

import { LogIn, Search, Upload, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Input } from "@/components/ui/input";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Link, useLocation } from "react-router";
import { useEffect, useState } from "react";

interface AppHeaderProps {
  onProfileClick?: () => void;
}

export function AppHeader({ onProfileClick }: AppHeaderProps) {

  const location = useLocation();

  const [dir, setDir] = useState("/")

  useEffect(() => {

    if (location.pathname.startsWith("/.@/")) {
      setDir("/")
    } else {
      setDir(location.pathname)
    }

  }, [location.pathname])


  return (
    <header className="sticky top-0 z-10 flex h-16 items-center rounded-full gap-4 border-b bg-background/80 px-4 backdrop-blur-sm md:px-6">
      <SidebarTrigger className="-ml-1" />

      <div className="relative flex-1 max-w-md">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          placeholder="Поиск файлов..."
          className="pl-9"
        />
      </div>

      <div className="ml-auto flex items-center gap-2">
        <Link to={`/.@/upload?dir=${dir}`}>
          <Button variant="outline" size="sm" className="hidden sm:flex">
            <Upload className="mr-2 h-4 w-4" />
            Upload
          </Button>
        </Link>

        <Button
          variant="ghost"
          size="icon"
          onClick={onProfileClick}
          className="rounded-full"
        >
          <User className="h-8 w-8" />

        </Button>

      </div>
    </header>
  );
}