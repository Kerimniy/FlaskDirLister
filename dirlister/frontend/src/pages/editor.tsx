import { Editor, type OnMount } from "@monaco-editor/react";
import { useTheme } from "next-themes";
import { useEffect, useState, useRef } from "react";
import type * as Monaco from "monaco-editor";
import { useSearchParams, useLocation, Link, useNavigate } from "react-router";
import { BACKEND_BASE_URL, useAuth } from "@/App";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { Save, Search, Replace, FileCode2, Loader2, HomeIcon } from "lucide-react";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";


const languageMap: Record<string, string> = {
    ts: "typescript",
    tsx: "typescript",
    js: "javascript",
    jsx: "javascript",
    go: "go",
    rs: "rust",
    py: "python",
    cpp: "cpp",
    c: "c",
    h: "cpp",
    html: "html",
    css: "css",
    json: "json",
    yaml: "yaml",
    yml: "yaml",
    md: "markdown",
    sh: "shell",
    sql: "sql",
    xml: "xml",
};

export default function EditPage() {
    const { resolvedTheme } = useTheme();
    const [searchParams, setSearchParams] = useSearchParams();
    const location = useLocation()
    const navigate = useNavigate()

    const [pageParam, setPageParam] = useState(searchParams.get("page") || "");

    const editorRef = useRef<Monaco.editor.IStandaloneCodeEditor | null>(null);
    const monacoRef = useRef<typeof Monaco | null>(null);

    let _ = pageParam.split("/")


    const [content, setContent] = useState("");
    const [filename, setFilename] = useState(_[_.length - 1] || "");
    const [pathname, setPathname] = useState(_.slice(0, _.length - 1).join("/") + "/" || "");



    const [isSaving, setIsSaving] = useState(false);
    const [isLoading, setIsLoading] = useState(false);

    const [fileExtension, setFileExtension] = useState(filename.split(".").pop()?.toLowerCase() || "");
    const [editorLanguage, setEditorLanguage] = useState(languageMap[fileExtension] || "plaintext");
    const { user, updateUser } = useAuth();


    const handleEditorDidMount: OnMount = (editor, monaco) => {
        editorRef.current = editor;
        monacoRef.current = monaco;
    };

    const handleSearch = () => {
        editorRef.current?.trigger("source", "actions.find", null);
    };

    const handleReplace = () => {
        editorRef.current?.trigger("source", "editor.action.startFindReplaceAction", null);
    };

    useEffect(() => {


        if (user===null || user===undefined || user.email===""){
            navigate("/.@/auth/login")
        }

        let pageParam = searchParams.get("page") || ""

        setPageParam(pageParam)
        let _ = pageParam.split("/")
        setFilename(_[_.length - 1] || "")
        setPathname(_.slice(0, _.length - 1).join("/") + "/" || "")

        setFileExtension(filename.split(".").pop()?.toLowerCase() || "");
        setEditorLanguage(languageMap[fileExtension] || "plaintext");


    }, [location.pathname])

    useEffect(() => {
        if (!pageParam) return;

        setIsLoading(true);
        fetch(`${BACKEND_BASE_URL}/s/${pageParam}?open=true`)
            .then((res) => {
                if (res.ok) return res.text();
                throw new Error("Failed to fetch file");
            })
            .then((data) => {
                setContent(data || "");
            })
            .catch((err) => console.error(err))
            .finally(() => setIsLoading(false));
    }, [pageParam]);

    const handleSave = async () => {

        if (location.pathname == "/.@/create") {

            setIsSaving(true);
            try {
                const response = await fetch(`${BACKEND_BASE_URL}/manage/upload?file=${filename}&edit=false`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: content ,
                    credentials: "include"
                });

                if (!response.ok) throw new Error("Save failed");

                console.log("Saved");

                navigate(`/.@/edit?page=${filename}`)

            } catch (error) {
                console.error("Error saving file:", error);
            } finally {
                setIsSaving(false);
            }
        }
        else {

            if (!pageParam) return;

            setIsSaving(true);
            try {
                const response = await fetch(`${BACKEND_BASE_URL}/manage/upload?file=${pageParam}&edit=true`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ newName: (pathname + filename), content: content }),
                    credentials: "include"
                });

                if (!response.ok) throw new Error("Save failed");

                console.log("Saved");

                searchParams.set("page", pathname + filename)
                setSearchParams(searchParams)

            } catch (error) {
                console.error("Error saving file:", error);
            } finally {
                setIsSaving(false);
            }
        }
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center w-full h-[600px] border rounded-md bg-muted/20">
                <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    return (

        <SidebarProvider className="w-full">
            <AppSidebar />

            <SidebarInset className="min-w-0">


                <div className="flex flex-col min-h-[500px] border rounded-md shadow-sm overflow-hidden bg-background">

                    <div className="flex flex-col sm:flex-row items-center justify-between p-2 gap-2 border-b bg-muted/40">

                        <div className="flex items-center gap-2 sm:w-auto flex-1">
                            <SidebarTrigger className="ml-1 text-muted-foreground" />
                            <FileCode2 className="hidden sm:block w-4 h-4 text-muted-foreground ml-1" />
                            <div className="h-8 flex flex-row items-center">{pathname}</div>
                            <Input
                                value={filename}
                                onChange={(e) => setFilename(e.target.value)}
                                className="h-8 bg-background max-w-full sm:max-w-[250px]"
                                placeholder="File name..."
                            />
                        </div>

                        <div className="flex items-center justify-end gap-1.5 w-full sm:w-auto">
                            <Button variant="outline" size="sm" onClick={handleSearch} className="h-8 px-2 sm:px-3" title="Search (Ctrl+F)">
                                <Search className="w-4 h-4 sm:mr-2" />
                                <span className="hidden sm:inline">Search</span>
                            </Button>

                            <Button variant="outline" size="sm" onClick={handleReplace} className="h-8 px-2 sm:px-3" title="Replace (Ctrl+H)">
                                <Replace className="w-4 h-4 sm:mr-2" />
                                <span className="hidden sm:inline">Replace</span>
                            </Button>

                            <Button size="sm" onClick={handleSave} disabled={isSaving} className="h-8 px-2 sm:px-3 min-w-[40px] sm:min-w-[110px]">
                                {isSaving ? (
                                    <Loader2 className="w-4 h-4 animate-spin sm:mr-2" />
                                ) : (
                                    <Save className="w-4 h-4 sm:mr-2" />
                                )}
                                <span className="hidden sm:inline">
                                    {isSaving ? "Saving..." : "Save"}
                                </span>
                            </Button>
                        </div>
                    </div>

                    <div className="flex-grow relative">
                        <Editor
                            onMount={handleEditorDidMount}
                            height="100vh"
                            width="auto"
                            language={editorLanguage}
                            theme={resolvedTheme === "dark" ? "vs-dark" : "light"}
                            value={content}
                            onChange={(value) => setContent(value || "")}
                            options={{
                                minimap: { enabled: false },
                                fontSize: 14,
                                wordWrap: "on",
                                scrollBeyondLastLine: false,
                                automaticLayout: true,
                                padding: { top: 12, bottom: 12 }
                            }}
                        />
                    </div>
                </div>
            </SidebarInset>
        </SidebarProvider>
    );
}