"use client";

import React, { useState, useRef } from "react";
import type { DragEvent, ChangeEvent, KeyboardEvent } from "react"
import { Upload, X, FileText } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { AppHeader } from "@/components/app-header";
import { useNavigate, useSearchParams } from "react-router";

import { formatFileSize } from "@/lib/files";
import { BACKEND_BASE_URL } from "@/App";
export default function FileUpload() {

    const fileLimit = 10
    let maxSize = 0

    fetch(`${BACKEND_BASE_URL}/info/upload-limit`).then(r=>{
       return r.text()
    }).then(t=>{maxSize=Number(t)})

    const [searchParams, setSearchParams] = useSearchParams()
    const [statusMessage, setStatusMessage] = useState<{ text: string; isError: boolean } | null>(null);

    const [size, setSize] = useState(0);
    const [isLimitExeeded, setIsLimitExeeded] = useState(false);
    const [files, setFiles] = useState<File[]>([]);
    const [isDragging, setIsDragging] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleClick = () => {
        fileInputRef.current?.click();
    };

    const navigate = useNavigate()

    const handleKeyDown = (e: KeyboardEvent<HTMLDivElement>) => {
        if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            handleClick();
        }
    };

    const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (e.target.files) {
            addFiles(Array.from(e.target.files));
        }
    };

    const handleDragEnter = (e: DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(true);
    };

    const handleDragLeave = (e: DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.currentTarget && !e.currentTarget.contains(e.relatedTarget as Node)) {
            setIsDragging(false);
        }
    };

    const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();
    };

    const handleDrop = (e: DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(false);

        if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
            addFiles(Array.from(e.dataTransfer.files));
        }
    };

    const addFiles = (newFiles: File[]) => {

        setFiles((prev) => {
            const existingKeys = new Set(prev.map((f) => `${f.name}-${f.size}`));
            const uniqueNewFiles = newFiles.filter(
                (f) => !existingKeys.has(`${f.name}-${f.size}`)
            );
            let c = fileLimit - files.length - uniqueNewFiles.length

            newFiles = uniqueNewFiles.slice(0, c)

            let _size = 0
            for (let el of newFiles) {
                _size += el.size
            }

            if (size + _size > maxSize) {
                setIsLimitExeeded(true)
                return [...prev]
            } else {
                setIsLimitExeeded(false)
            }

            setSize(size + _size)


            return [...prev, ...newFiles];
        });
    };

    const removeFile = (index: number) => {
        setFiles((prev) => prev.filter((_, i) => i !== index));
    };

    const handleUpload = async () => {
        if (files.length === 0) return;

        const formData = new FormData();
        files.forEach((file) => {
            formData.append("files", file);
        });

        try {
            const response = await fetch(`${BACKEND_BASE_URL}/manage/upload-multiple?dir=${searchParams.get("dir")}`, {
                method: "POST",
                body: formData,
                credentials: "include"
            });

            if (!response.ok) { setStatusMessage({ text: `Error status code: ${response.status}`, isError: true }); throw new Error("Save failed")};

            setStatusMessage({text:"Successfuly saved", isError: false})
            console.log("Saved");

            setFiles([]);
        } catch (error) {
            console.error("error:", error);
        }
    };

    return (

        <SidebarProvider className="w-full">
            <AppSidebar />

            <SidebarInset className="min-w-0">
                <AppHeader
                    onProfileClick={() => navigate("/.@/account")}
                />
                {statusMessage && (
                    <div className={`p-3 rounded-md text-sm font-medium border max-w-sm w-full text-center ${statusMessage.isError
                        ? 'bg-destructive/10 text-destructive border-destructive/20'
                        : 'bg-green-500/10 text-green-600 border-green-500/20'
                        }`}>
                        {statusMessage.text}
                    </div>
                )}

                <div className="p-6">
                    <Card className="mt-6">
                        <CardHeader>
                            <CardTitle>Upload</CardTitle>
                            <CardDescription>
                                Select files
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <input
                                type="file"
                                ref={fileInputRef}
                                onChange={handleFileChange}
                                multiple
                                className="hidden"
                            />

                            <div
                                role="button"
                                tabIndex={0}
                                onClick={handleClick}
                                onKeyDown={handleKeyDown}
                                onDragEnter={handleDragEnter}
                                onDragLeave={handleDragLeave}
                                onDragOver={handleDragOver}
                                onDrop={handleDrop}
                                className={cn(
                                    "flex flex-col items-center justify-center w-full h-48 p-6 border-2 border-dashed rounded-lg cursor-pointer transition-all duration-200",
                                    isDragging
                                        ? "border-primary bg-primary/5 scale-[1.02]"
                                        : "border-muted-foreground/25 hover:border-primary/50 hover:bg-muted/50"
                                )}
                            >
                                <div
                                    className={cn(
                                        "flex items-center justify-center w-12 h-12 rounded-full mb-3 transition-colors",
                                        isDragging ? "bg-primary/10 text-primary" : "bg-muted text-muted-foreground"
                                    )}
                                >
                                    <Upload className="w-6 h-6" />
                                </div>
                                <p className="text-sm font-medium text-foreground">
                                    {isDragging ? "Отпустите файлы для загрузки" : "Перетащите файлы или нажмите для выбора"}
                                </p>
                                <p className="text-xs text-muted-foreground mt-1">
                                    Max 10 files, 1GB limit
                                </p>
                            </div>

                            {files.length > 0 && (
                                <div className="space-y-3">
                                    <div className="flex items-center justify-between">
                                        <h3 className="text-sm font-medium text-foreground">
                                            files selected: {files.length}
                                        </h3>
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => setFiles([])}
                                            className="text-muted-foreground hover:text-destructive"
                                        >
                                            Clear all
                                        </Button>
                                    </div>

                                    <ScrollArea className="h-48 w-full rounded-md border p-2">
                                        <div className="space-y-2">
                                            {files.map((file, index) => (
                                                <div
                                                    key={`${file.name}-${file.size}-${index}`}
                                                    className="flex items-center justify-between p-2 rounded-md bg-muted/50 hover:bg-muted transition-colors group"
                                                >
                                                    <div className="flex items-center gap-3 overflow-hidden">
                                                        <FileText className="w-5 h-5 text-muted-foreground shrink-0" />
                                                        <div className="flex flex-col items-start overflow-hidden">
                                                            <span className="text-sm font-medium truncate">
                                                                {file.name}
                                                            </span>
                                                            <span className="text-xs text-muted-foreground">
                                                                {formatFileSize(file.size)}
                                                            </span>
                                                        </div>
                                                    </div>
                                                    <Button
                                                        variant="ghost"
                                                        size="icon"
                                                        className="h-8 w-8 text-muted-foreground hover:text-destructive opacity-70 group-hover:opacity-100"
                                                        onClick={() => removeFile(index)}
                                                    >
                                                        <X className="h-4 w-4" />
                                                    </Button>
                                                </div>
                                            ))}
                                        </div>
                                    </ScrollArea>

                                    <Button onClick={handleUpload} className="w-full">
                                        Upload {files.length > 0 && `(${files.length})`}
                                    </Button>
                                </div>
                            )}

                            {isLimitExeeded && (<><div>Files size limit exeeded</div></>)}
                        </CardContent>
                    </Card>
                </div>

            </SidebarInset >
        </SidebarProvider >
    );
}