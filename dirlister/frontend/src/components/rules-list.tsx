"use client";

import { Edit, MoreHorizontal, Trash2, Folder, Pencil } from "lucide-react";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import {
    Popover,
    PopoverContent,
    PopoverTrigger,
    PopoverDescription,
    PopoverHeader,
    PopoverTitle
} from "@/components/ui/popover"

import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "./ui/input";
import {
    formatDate,
    getMimeColor,
    getMimeTypeIcon,
} from "@/lib/files";
import { useNavigate } from "react-router-dom";
import type { FileItem } from "@/lib/files";
import { BACKEND_BASE_URL } from "@/App";
import { RESULTS_PER_PAGE } from "@/App";

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
import { useEffect, useRef, useState } from "react";

export interface Rule {
    path: string;
    createdAt: string
}

interface RulesListProps {
    rules: Rule[];
    page: number;
    onDelete: (rule: string[]) => void;
    onChange: (rule: Rule, newPath: string) => void;
    setCheckAll: React.Dispatch<React.SetStateAction<boolean>>;
    checkAll: boolean

    setCheckList: React.Dispatch<React.SetStateAction<Set<string>>>;
    checkList: Set<string>;

    checkCount: number;
    setCheckCount: React.Dispatch<React.SetStateAction<number>>;

}

export function RulesList({ rules, page, onDelete, onChange, checkAll, setCheckAll, setCheckList, checkList, checkCount, setCheckCount }: RulesListProps) {

    const navigate = useNavigate();



    const [newRulePath, setNewRulePath] = useState("")


    if (rules === null) {
        return (
            <div className="flex flex-col items-center justify-center py-20 text-center">
                <div className="rounded-full bg-muted p-4">
                    <MoreHorizontal className="h-6 w-6 text-muted-foreground" />
                </div>
                <h3 className="mt-4 text-lg font-semibold">ERROR</h3>
                <p className="text-sm text-muted-foreground">
                    Backend is down.
                </p>
            </div>)
    }

    if (rules === undefined || rules.length === 0) {
        return (
            <div className="flex flex-col items-center justify-center py-20 text-center">
                <div className="rounded-full bg-muted p-4">
                    <MoreHorizontal className="h-6 w-6 text-muted-foreground" />
                </div>
                <h3 className="mt-4 text-lg font-semibold">No rules</h3>
                <p className="text-sm text-muted-foreground">
                    There is no rules, create first.
                </p>
            </div>
        );
    }


    return (
        <Table >
            <TableHeader>
                <TableRow>
                    <TableHead>
                        <Checkbox checked={checkAll} onClick={() => { let e = !checkAll; setCheckCount(e ? rules.length : 0); checkAllFunc(e, rules, setCheckList) }} onCheckedChange={(e) => { setCheckAll(e); }} />
                    </TableHead>
                    <TableHead className="text-center">Path</TableHead>
                    <TableHead className="hidden md:table-cell text-center">Created</TableHead>
                    <TableHead className="w-24 text-center">Actions</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {rules.slice(page * RESULTS_PER_PAGE, (page + 1) * RESULTS_PER_PAGE).map((rule, i) => {

                    return (
                        <TableRow
                            key={rule.path}
                            className="group"

                        >
                            <TableCell onClick={(e) => e.stopPropagation()}>
                                <Checkbox id={rule.path}
                                    checked={checkList.has(rule.path)}
                                    onCheckedChange={(value) => {

                                        if (checkCount + Number(value) * 2 - 1 == rules.length) {
                                            setCheckAll(true)
                                        }
                                        else {
                                            setCheckAll(false)
                                        }

                                        setCheckCount(checkCount + Number(value) * 2 - 1);

                                        if (value) {
                                            setCheckList(prev => {
                                                const next = new Set(prev);
                                                next.add(rule.path);
                                                return next;
                                            });

                                        } else {
                                            setCheckList(prev => {
                                                const next = new Set(prev);
                                                next.delete(rule.path);
                                                return next;
                                            });
                                        }

                                    }} />
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center gap-3">


                                    <div className="flex flex-col align-left">
                                        <span className="font-medium  max-w-[380px]">
                                            {rule.path}
                                        </span>

                                    </div>
                                </div>
                            </TableCell>
                            <TableCell className="hidden sm:table-cell text-muted-foreground text-sm">
                                {rule.createdAt}
                            </TableCell>
                            <TableCell className="text-right" onClick={(e) => e.stopPropagation()}>
                                <div className="flex items-center justify-end gap-1 focus-within:opacity-100">

                                    <AlertDialog>
                                        <AlertDialogTrigger onClick={() => setNewRulePath(rule.path)} render={<Button variant="ghost" size="icon" title="change" className="h-8 w-8">
                                            <Pencil className="h-4 w-4" />
                                            <span className="sr-only">Change</span>
                                        </Button>}>
                                        </AlertDialogTrigger>
                                        <AlertDialogContent>
                                            <AlertDialogHeader>
                                                <AlertDialogTitle>Change rule</AlertDialogTitle>
                                                <AlertDialogDescription className="w-full">
                                                    <div className="flex flex-col gap-4 w-full">

                                                        <p className="text-center w-full">This action cannot be undone.</p>

                                                        <Input value={newRulePath} onInput={(e) => { setNewRulePath(e.currentTarget.value) }} className="w-full" placeholder="New rulesname" />

                                                    </div>
                                                </AlertDialogDescription>
                                            </AlertDialogHeader>
                                            <AlertDialogFooter>
                                                <AlertDialogCancel>Cancel</AlertDialogCancel>
                                                <AlertDialogCancel variant="default" onClick={() => onChange(rule, newRulePath)}>Continue</AlertDialogCancel>
                                            </AlertDialogFooter>
                                        </AlertDialogContent>
                                    </AlertDialog>


                                    <AlertDialog>
                                        <AlertDialogTrigger render={<Button variant="ghost" size="icon" title="delete" className="h-8 w-8 text-destructive hover:text-destructive">
                                            <Trash2 className="h-4 w-4" />
                                            <span className="sr-only">Delete</span>
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
                                                <AlertDialogCancel variant="default" onClick={() => onDelete([rule.path])}>Continue</AlertDialogCancel>
                                            </AlertDialogFooter>
                                        </AlertDialogContent>
                                    </AlertDialog>



                                   
                                </div>
                            </TableCell>
                        </TableRow>
                    );
                })}
            </TableBody>
        </Table>
    );
}


function checkAllFunc(e, rules: Rule[], setChecked: React.Dispatch<React.SetStateAction<Set<string>>>) {

    let newSet = new Set<string>()

    if (e) {
        for (let el of rules) {
            newSet.add(el.path)
        }
    }

    setChecked(newSet);

}