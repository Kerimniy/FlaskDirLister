"use client";

import { useEffect, useRef, useState } from "react";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { AppHeader } from "@/components/app-header";
import { FileList } from "@/components/file-list";
import { Button } from "@/components/ui/button";
import { Plus, ChevronLeft, ChevronRight } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card"
import { formatDate } from "@/lib/files";

import { useLocation } from 'react-router-dom';

import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"

import { Trash2 } from "lucide-react";

import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { Link } from "react-router-dom";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { RESULTS_PER_PAGE, useAuth } from "@/App";

import { BACKEND_BASE_URL } from "@/App";

import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { RulesList, type Rule } from "@/components/rules-list";


export default function RulesPage() {

    const [searchParams, setSearchParams] = useSearchParams();
    const navigate = useNavigate()
    const pageParam = searchParams.get("page");

    const [rules, setRules] = useState<Rule[]>([]);
    const [page, setPage] = useState<number>(Number(pageParam) || 0);


    const [checkAll, setCheckAll] = useState(false)
    const [checksList, setChecksList] = useState(new Set<string>)
    const [checkCount, setCheckCount] = useState(0)
    const [newRulePath, setNewRulePath] = useState("")


    const leftArrowPageButton = useRef(null)
    const rightArrowPageButton = useRef(null)

    const location = useLocation()
    const [sortBy, setSortBy] = useState(localStorage.getItem("sortRulesBy") || "Alphabet")


    useEffect(() => {
        if (page === 0) {
            searchParams.delete('page');

            setSearchParams(searchParams);
        }
        else {
            setSearchParams({ "page": String(page) })
        }
        getRules(page, sortBy).then(rules => setRules(rules))

    }, [page])


    const handleDelete = (rules: string[]) => {
        let arr = Array.from(rules)

        for (let i=0; i<arr.length;i++){
            arr[i] = arr[i].replace(/^\/|\/$/g, '')
        }

        fetch(`${BACKEND_BASE_URL}/rules/delete`, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(arr)

        }).then((res) => {
            console.log("Deletion status: ", res.status)

            if (res.ok) {
                getRules(page, sortBy).then((res) => {

                    setRules(res);

                });

                if (rules.length === 1) {
                    let nchl = new Set(checksList)
                    nchl.delete(rules[0])
                    setChecksList(nchl)
                }
                else {
                    setChecksList(new Set<string>())

                }

            }
        });



        return
    };

    const handleCreate = (rule: string) => {

        fetch(`${BACKEND_BASE_URL}/rules/create?p=${rule.replace(/^\/|\/$/g, '')}`, {
            method: "PUT",
            credentials: "include",

        }).then((res) => {
            console.log("Creation status: ", res.status)

            if (res.ok) {
                getRules(page, sortBy).then((res) => {

                    setRules(res);

                });

            }
        });



        return
    };


    const handleChange = (rule: Rule, newPath: string) => {

        fetch(`${BACKEND_BASE_URL}/rules/change?r=${rule.path.replace(/^\/|\/$/g, '') }n=${newPath}`, {
            method: "PATCH",
            credentials: "include"

        }).then((res) => {
            console.log("Rename status: ", res.status)

            if (res.ok) {
                getRules(page, sortBy).then((res) => {

                    setRules(res);

                });

            }
        });



        return
    };


    const { user } = useAuth()

    return (
        <SidebarProvider className="w-full">
            <AppSidebar />



            <SidebarInset className="flex flex-col min-w-0">
                <AppHeader
                    onProfileClick={() => navigate("/.@/account")}
                />



                <main className="flex-1 overflow-auto p-4 md:p-6">

                    <div className="mb-6 flex items-center justify-between">
                        <div>
                            <h1 className="text-lg font-bold tracking-tight">Rules</h1>

                        </div>


                    </div>

                    <div className="flex flex-row pb-4 gap-4">
                        <Input value={newRulePath} onInput={(e)=>setNewRulePath(e.currentTarget.value)} placeholder="/folder1"></Input>
                        <Button onClick={()=>handleCreate(newRulePath)} className="pl-3 pr-3">Create</Button>
                    </div>

                    <div className="mb-3 flex flex-row justify-between">

                        <Select value={sortBy} onValueChange={(e) => { localStorage.setItem("sortRulesBy", e); setRules(sortRulesBy(rules, e)); setSortBy(e) }}>
                            <SelectTrigger className="w-[240px]">
                                <SelectValue placeholder="Theme" />
                            </SelectTrigger>
                            <SelectContent alignItemWithTrigger={false}>
                                <SelectGroup>
                                    <SelectItem key="Alphabet" value="Alphabet">
                                        Alphabet
                                    </SelectItem>

                                    <SelectItem key="Age" value="Age">
                                        Age
                                    </SelectItem>

                                </SelectGroup>
                            </SelectContent>
                        </Select>


                        <AlertDialog>
                            <AlertDialogTrigger render={<Button variant="destructive" title="delete" className={(checkCount) > 1 ? "opacity-100" : "opacity-0"}>
                                <Trash2></Trash2>
                                Delete All
                            </Button>}>
                            </AlertDialogTrigger>
                            <AlertDialogContent>
                                <AlertDialogHeader>
                                    <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                                    <AlertDialogDescription>
                                        This action cannot be undone.
                                    </AlertDialogDescription>
                                </AlertDialogHeader>
                                <AlertDialogFooter>
                                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                                    <AlertDialogCancel variant="default" onClick={() => handleDelete(Array.from(checksList))}>Continue</AlertDialogCancel>
                                </AlertDialogFooter>
                            </AlertDialogContent>
                        </AlertDialog>


                    </div>


                    <div className="rounded-lg border bg-card">
                        <RulesList
                            rules={rules}
                            page={page}
                            onDelete={handleDelete}
                            onChange={handleChange}
                            checkAll={checkAll}
                            setCheckAll={setCheckAll}

                            setCheckList={setChecksList}
                            checkList={checksList}
                            checkCount={checkCount}
                            setCheckCount={setCheckCount}
                        />
                    </div>
                    {(rules !== null && rules !== undefined) &&
                        <Card className="flex-row justify-center mt-4">
                            <Button disabled={page === 0} ref={leftArrowPageButton} variant="outline" onClick={() => { if (page > 0) { rightArrowPageButton.current.disabled = false; setPage(page - 1); if (page - 1 === 0) { leftArrowPageButton.current.disabled = true } } }}><ChevronLeft /></Button>

                            <Input min={0} max={Math.floor(rules.length / RESULTS_PER_PAGE)} style={{ width: `${String(page).length + 6}ch` }} type="number" value={page} onInput={(e) => { setPage(Number(e.currentTarget.value)) }}></Input>

                            <Button disabled={page === Math.floor(rules.length / RESULTS_PER_PAGE)} ref={rightArrowPageButton} variant="outline" onClick={() => { let maxPage = Math.floor(rules.length / RESULTS_PER_PAGE); if (page < maxPage) { leftArrowPageButton.current.disabled = false; setPage(page + 1); if (page + 1 === maxPage) { rightArrowPageButton.current.disabled = true } } }}><ChevronRight /></Button>

                        </Card>
                    }
                </main>
            </SidebarInset>
        </SidebarProvider>
    );
}

export async function getRules(page: Number, dim: string): Promise<Rule[]> {
    const url = `${BACKEND_BASE_URL}/rules/get?p=${page}`;
    let response
    try {
        response = await fetch(url);
        if (!response.ok) {
            console.log(response.status)
            return null
        }
        
    }
    catch (err) {
        console.log(err)
        return null
    }
    const result = await response.json();


    let rules: Rule[] = []
    for (let el of result) {

        rules.push({ path: el.path, createdAt: formatDate(el.createdAt) })

    }

    rules = sortRulesBy(rules, dim)

    return rules


}

export function sortRulesBy(rules: Rule[], dim: string): Rule[] {

    if (dim === "Alphabet") {
        rules.sort((a, b) => a.path.localeCompare(b.path))
    }

    else if (dim === "Age") {
        rules.sort((a, b) => b.createdAt.localeCompare(a.createdAt))
    }
    return rules
}