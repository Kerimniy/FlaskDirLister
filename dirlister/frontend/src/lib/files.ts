import { BACKEND_BASE_URL } from "@/App";
import mime from 'mime';
import {
  File,
  FileText,
  FileImage,
  FileVideo,
  FileAudio,
  FileArchive,
  FileCode,
  FileSpreadsheet,
  FileJson,

} from "lucide-react";
import type { LucideProps } from "lucide-react";

export interface FileItem {
  name: string;
  size: number;
  mimeType: string;
  modifiedAt: string;
  type: string;
  fullName: string;
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 B";
  const units = ["B", "KiB", "MiB", "GiB", "TiB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(i === 0 ? 0 : 1)} ${units[i]}`;
}

export function formatDate(iso: string): string {

  const date = new Date(iso);
  return new Intl.DateTimeFormat("ru-RU", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

export function getMimeTypeIcon(mimeType: string): React.ForwardRefExoticComponent<Omit<LucideProps, "ref"> & React.RefAttributes<SVGSVGElement>> {

  if (mimeType === null || mimeType === undefined) {
    mimeType = ""
  }

  const type = mimeType.toLowerCase();

  if (type.startsWith("image/")) return FileImage;
  if (type.startsWith("video/")) return FileVideo;
  if (type.startsWith("audio/")) return FileAudio;
  if (type.includes("pdf")) return FileText;
  if (type.includes("zip") || type.includes("rar") || type.includes("tar"))
    return FileArchive;
  if (type.includes("json")) return FileJson;
  if (
    type.includes("spreadsheet") ||
    type.includes("excel") ||
    type.includes("csv")
  )
    return FileSpreadsheet;
  if (
    type.includes("javascript") ||
    type.includes("typescript") ||
    type.includes("html") ||
    type.includes("css")
  )
    return FileCode;
  if (type.startsWith("text/")) return FileText;

  return File;
}

export function getMimeColor(mimeType: string): string {

  if (mimeType === null || mimeType === undefined) {
    mimeType = ""
  }

  const type = mimeType.toLowerCase();

  if (type.startsWith("image/")) return "text-pink-500";
  if (type.startsWith("video/")) return "text-purple-500";
  if (type.startsWith("audio/")) return "text-orange-500";
  if (type.includes("pdf")) return "text-red-500";
  if (type.includes("zip") || type.includes("rar")) return "text-yellow-600";
  if (type.includes("json")) return "text-amber-500";
  if (type.includes("spreadsheet") || type.includes("excel"))
    return "text-emerald-500";
  if (type.includes("javascript") || type.includes("typescript"))
    return "text-sky-500";
  return "text-slate-500";
}

export async function searchFiles(query: string, page: number, dim: string, setStatus: React.Dispatch<React.SetStateAction<number>>): Promise<FileItem[]> {
  const url = `${BACKEND_BASE_URL}/search?q=${query.replace(/\/+$/, '')}&p=${page}`;
  let response
  try {
    response = await fetch(url, { credentials: "include" });
    setStatus(response.status)

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


  let files: FileItem[] = []

    for (let el of result) {

      files.push({ name: el.name, type: el.type, size: el.size, mimeType: mime.getType(el.name), modifiedAt: formatDate(el.modTime), fullName: el.fullName })

    }
  
  
  if (dim !== "") {
    files = sortFilesBy(files, dim)
  }
  
  return files


}


export async function getFiles(folder: string, dim: string,page:Number, setStatus: React.Dispatch<React.SetStateAction<number>>): Promise<FileItem[]> {
  const url = `${BACKEND_BASE_URL}/s/${folder.replace(/\/+$/, '')}?p=${page}`;
  let response
  try {
    response = await fetch(url, { credentials: "include" });
    setStatus(response.status)
    if (response.status===204){
      return []
    }

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


  let files: FileItem[] = []
  for (let el of result) {

    files.push({ name: el.name, type: el.type, size: el.size, mimeType: mime.getType(el.name), modifiedAt: formatDate(el.modTime), fullName: el.fullName })

  }

  files = sortFilesBy(files, dim)

  return files


}


export function sortFilesBy(files: FileItem[], dim: string): FileItem[] {

  if (dim === "Alphabet") {
    files.sort((a, b) => a.name.localeCompare(b.name))
  }
  else if (dim === "Size") {
    files.sort((a, b) => a.size - b.size)
  }
  else if (dim === "Age") {
    files.sort((a, b) => b.modifiedAt.localeCompare(a.modifiedAt))
  }
  else if (dim === "Type") {
    files.sort((a, b) => (a.mimeType || "").localeCompare((b.mimeType || "")))
  }

  return files
}