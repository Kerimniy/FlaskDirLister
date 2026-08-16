"use client";

import { Edit, MoreHorizontal, Trash2, Folder } from "lucide-react";
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
  formatFileSize,
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

interface FileListProps {
  files: FileItem[];
  page: number;
  onEdit: (file: FileItem) => void;
  onDelete: (file: FileItem) => void;
  onRename: (file: FileItem, newName: string) => void;
  setCheckAll: React.Dispatch<React.SetStateAction<boolean>>;
  checkAll: boolean

  setCheckList: React.Dispatch<React.SetStateAction<Set<string>>>;
  checkList: Set<string>;

  checkCount: number;
  setCheckCount: React.Dispatch<React.SetStateAction<number>> ;
}

export function FileList({ files, page, onEdit, onDelete, onRename, checkAll, setCheckAll, setCheckList, checkList, checkCount, setCheckCount }: FileListProps) {

  const navigate = useNavigate();



  const [newName, setNewName] = useState("")


  if (files === null) {
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

  if (files === undefined || files.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-center">
        <div className="rounded-full bg-muted p-4">
          <MoreHorizontal className="h-6 w-6 text-muted-foreground" />
        </div>
        <h3 className="mt-4 text-lg font-semibold">No files</h3>
        <p className="text-sm text-muted-foreground">
          There is no files, upload first.
        </p>
      </div>
    );
  }


  return (
    <Table >
      <TableHeader>
        <TableRow>
          <TableHead>
            <Checkbox checked={checkAll} onClick={() => { let e = !checkAll; setCheckCount(e ? files.length : 0); checkAllFunc(e, files, setCheckList) }} onCheckedChange={(e) => { setCheckAll(e); }} />
          </TableHead>
          <TableHead className="text-center">Name</TableHead>
          <TableHead className="hidden md:table-cell text-center">Size</TableHead>
          <TableHead className="hidden lg:table-cell text-center">Type</TableHead>
          <TableHead className="hidden sm:table-cell text-center">Modified</TableHead>
          <TableHead className="w-24 text-center">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {files.slice(page * RESULTS_PER_PAGE, (page + 1) * RESULTS_PER_PAGE).map((file, i) => {
          const Icon = (file.type === "folder") ? Folder : getMimeTypeIcon(file.mimeType);
          const colorClass = (file.type === "folder") ? "text-slate-500" : getMimeColor(file.mimeType);

          return (
            <TableRow
              key={file.fullName}
              className="group cursor-pointer"
              onClick={() => { if (file.type === "folder") { navigate("/" + file.fullName, { replace: true }) } else { window.location.href = `${BACKEND_BASE_URL}/s/${file.fullName}?open=true` } }}
            >
              <TableCell onClick={(e) => e.stopPropagation()}>
                <Checkbox id={file.fullName}
                  checked={checkList.has(file.fullName)}
                  onCheckedChange={(value) => {

                    if (checkCount + Number(value) * 2 - 1 == files.length) {
                      setCheckAll(true)
                    }
                    else {
                      setCheckAll(false)
                    }

                    setCheckCount(checkCount + Number(value) * 2 - 1);

                    if (value) {
                      setCheckList(prev => {
                        const next = new Set(prev);
                        next.add(file.fullName);
                        return next;
                      });

                    } else {
                      setCheckList(prev => {
                        const next = new Set(prev);
                        next.delete(file.fullName);
                        return next;
                      });
                    }

                  }} />
              </TableCell>
              <TableCell>
                <div className="flex items-center gap-3">
                  <div
                    className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-md bg-muted ${colorClass}`}
                  >
                    <Icon className="h-5 w-5" />
                  </div>
                  <div className="flex flex-col align-left">
                    <span className="font-medium  max-w-[380px]">
                      {file.name}
                    </span>
                    <span className="text-xs text-left text-muted-foreground md:hidden">
                      {(file.type === "folder") ? "" : formatFileSize(file.size)}
                    </span>
                  </div>
                </div>
              </TableCell>
              <TableCell className="hidden md:table-cell text-muted-foreground">
                {(file.type === "folder") ? "" : formatFileSize(file.size)}
              </TableCell>
              <TableCell className="hidden lg:table-cell text-muted-foreground text-sm">
                {(file.type === "folder") ? "folder" : file.mimeType}
              </TableCell>
              <TableCell className="hidden sm:table-cell text-muted-foreground text-sm">
                {file.modifiedAt}
              </TableCell>
              <TableCell className="text-right" onClick={(e) => e.stopPropagation()}>
                <div className="flex items-center justify-end gap-1 focus-within:opacity-100">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8"
                    onClick={() => onEdit(file)}
                    title="edit"
                  >
                    <Edit className="h-4 w-4" />
                    <span className="sr-only">Edit</span>
                  </Button>


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
                        <AlertDialogCancel variant="default" onClick={() => onDelete(file)}>Continue</AlertDialogCancel>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>



                  <Popover>

                    <PopoverTrigger>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                      >
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </PopoverTrigger>

                    <PopoverContent className="w-fit" align="end">

                      <Button variant="ghost" className="p-1">

                        <AlertDialog>
                          <AlertDialogTrigger onClick={() => setNewName(file.fullName)} render={<span title="rename">
                            Rename
                          </span>}>
                          </AlertDialogTrigger>
                          <AlertDialogContent>
                            <AlertDialogHeader>
                              <AlertDialogTitle>Rename</AlertDialogTitle>
                              <AlertDialogDescription className="w-full">
                                <div className="flex flex-col gap-4 w-full">

                                  <p className="text-center w-full">This action cannot be undone.</p>

                                  <Input value={newName} onInput={(e) => { setNewName(e.currentTarget.value) }} className="w-full" placeholder="New filename" />

                                </div>
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogCancel>Cancel</AlertDialogCancel>
                              <AlertDialogCancel variant="default" onClick={() => onRename(file, newName)}>Continue</AlertDialogCancel>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>

                      </Button>

                      <Button variant="ghost" className="p-1" onClick={() => { window.location.href = `${BACKEND_BASE_URL}/s/${file.fullName}?open=false` }}>Download</Button>

                      <Button variant="ghost" className="p-1">

                        <AlertDialog>
                          <AlertDialogTrigger render={<span title="delete" className="text-destructive hover:text-destructive">
                            Delete
                          </span>}>
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
                              <AlertDialogCancel variant="default" onClick={() => onDelete(file)}>Continue</AlertDialogCancel>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>

                      </Button>

                    </PopoverContent>

                  </Popover>

                </div>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}


function checkAllFunc(e, files: FileItem[], setChecked: React.Dispatch<React.SetStateAction<Set<string>>>) {

  let newSet = new Set<string>()

  if (e) {
    for (let el of files) {
      newSet.add(el.fullName)
    }
  }

  setChecked(newSet);

}