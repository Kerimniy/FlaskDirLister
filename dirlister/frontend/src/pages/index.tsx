"use client";

import { useEffect, useRef, useState } from "react";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { AppHeader } from "@/components/app-header";
import { FileList } from "@/components/file-list";
import { Button } from "@/components/ui/button";
import { Plus, ChevronLeft, ChevronRight } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card"
import { getFiles, sortFilesBy, type FileItem } from "@/lib/files";

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

interface PathEl {
  name: string,
  path: string
}

function getPathBarLinks() {

  if (location.pathname === "/") {
    return []
  }

  let pathElements = location.pathname.replace(/^\/|\/$/g, '').split("/")

  let paths: PathEl[] = []

  let currentPath = ""

  for (let el of pathElements) {
    currentPath += "/" + el
    paths.push({ path: currentPath, name: el })
  }

  return paths
}


export default function IndexPage() {

  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate()
  const pageParam = searchParams.get("page");

  const [files, setFiles] = useState<FileItem[]>([]);
  const [paths, setPaths] = useState<PathEl[]>([]);
  const [page, setPage] = useState<number>(Number(pageParam) || 0);


  const [checkAll, setCheckAll] = useState(false)
  const [checksList, setChecksList] = useState(new Set<string>)
  const [checkCount, setCheckCount] = useState(0)


  const leftArrowPageButton = useRef(null)
  const rightArrowPageButton = useRef(null)

  const location = useLocation()

  useEffect(() => {
    if (page === 0) {
      searchParams.delete('page');

      setSearchParams(searchParams);
    }
    else {
      setSearchParams({ "page": String(page) })
    }
  }, [page])

  useEffect(() => {
    setPage(0)
    let folder = location.pathname

    let isMounted = true;

    getFiles(folder, sortBy).then((res) => {

      if (isMounted) {
        setFiles(res);
      }
    });

    setPaths(getPathBarLinks())

    return () => {
      isMounted = false;
    };

  }, [location.pathname])



  const handleDelete = (file: FileItem) => {

    fetch(`${BACKEND_BASE_URL}/manage/delete?file=${file.fullName}`, {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      credentials: "include"

    }).then((res) => {
      console.log("Deletion status: ", res.status)

      if (res.ok) {
        getFiles(location.pathname, sortBy).then((res) => {

          setFiles(res);

        });

        let nchl = new Set(checksList)
        nchl.delete(file.fullName)
        setChecksList(nchl)

        setPaths(getPathBarLinks())
      }
    });



    return
  };
  const handleDeleteAll = () => {

   

    fetch(`${BACKEND_BASE_URL}/manage/delete-all`, {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(Array.from(checksList)),
      credentials: "include"
    }).then((res) => {
      console.log("Deletion status: ", res.status)

      if (res.ok) {
        getFiles(location.pathname, sortBy).then((res) => {

          setFiles(res);

        });
        setChecksList(new Set<string>())

        setPaths(getPathBarLinks())
      }
    });



    return
  };

  const handleEdit = (file: FileItem) => {
    navigate(`/.@/edit?page=${file.fullName}`)
  };


  const handleRename = (file: FileItem) => {

    fetch(`${BACKEND_BASE_URL}/manage/rename?file=${file.fullName}`, {
      method: "PATCH",
      credentials: "include"

    }).then((res) => {
      console.log("Rename status: ", res.status)

      if (res.ok) {
        getFiles(location.pathname, sortBy).then((res) => {

          setFiles(res);

        });

        setPaths(getPathBarLinks())
      }
    });



    return
  };


  const [sortBy, setSortBy] = useState(localStorage.getItem("sortBy") || "Alphabet")
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
              <h1 className="text-2xl font-bold tracking-tight">Files</h1>

            </div>
            <Button>
              <Link className="flex-row flex justify-center items-center" to="/.@/create">
                <Plus className="mr-2 h-4 w-4" />
                New File</Link>
            </Button>

          </div>

          <div className="mb-3 flex flex-row justify-between">

            <Select value={sortBy} onValueChange={(e) => { localStorage.setItem("sortBy", e); setFiles(sortFilesBy(files, e)); setSortBy(e) }}>
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
                  <SelectItem key="Size" value="Size">
                    Size
                  </SelectItem>

                  <SelectItem key="Type" value="Type">
                    Type
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
                  <AlertDialogCancel variant="default" onClick={() => handleDeleteAll()}>Continue</AlertDialogCancel>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>


          </div>

          <Card className="overflow-auto p-2 md:p-3 mb-3 flex flex-row justify-start align-center">
            <CardContent className="flex flex-row gap-[0.125rem] items-center">


              <><Link to="/">home</Link><span>/</span></>
              {paths.map((el) => {
                return <><Link to={el.path}>{el.name}</Link><span>/</span></>
              })}
            </CardContent>
          </Card>

          <div className="rounded-lg border bg-card">
            <FileList
              files={files}
              page={page}
              onEdit={handleEdit}
              onDelete={handleDelete}
              onRename={handleRename}
              checkAll={checkAll}
              setCheckAll={setCheckAll}

              setCheckList={setChecksList}
              checkList={checksList}
              checkCount={checkCount}
              setCheckCount={setCheckCount}
            />
          </div>
          {(files !== null && files !== undefined) &&
            <Card className="flex-row justify-center mt-4">
              <Button disabled={page === 0} ref={leftArrowPageButton} variant="outline" onClick={() => { if (page > 0) { rightArrowPageButton.current.disabled = false; setPage(page - 1); if (page - 1 === 0) { leftArrowPageButton.current.disabled = true } } }}><ChevronLeft /></Button>

              <Input min={0} max={Math.floor(files.length / RESULTS_PER_PAGE)} style={{ width: `${String(page).length + 6}ch` }} type="number" value={page} onInput={(e) => { setPage(Number(e.currentTarget.value)) }}></Input>

              <Button disabled={page === Math.floor(files.length / RESULTS_PER_PAGE)} ref={rightArrowPageButton} variant="outline" onClick={() => { let maxPage = Math.floor(files.length / RESULTS_PER_PAGE); if (page < maxPage) { leftArrowPageButton.current.disabled = false; setPage(page + 1); if (page + 1 === maxPage) { rightArrowPageButton.current.disabled = true } } }}><ChevronRight /></Button>

            </Card>
          }
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}