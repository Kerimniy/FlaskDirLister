from flask import Flask, request, jsonify, render_template, redirect
from flask import url_for
from flask import send_file,send_from_directory, abort
from flask_caching import Cache
from datetime import datetime
import threading  
from io import StringIO
import os
from os import listdir,path
from os.path import isfile, join
from urllib.parse import quote, unquote
import markdown
import pathlib


app = Flask(__name__)
cache = Cache(app, config={'CACHE_TYPE': 'simple'})
PATH = path.dirname(path.realpath(__file__)).replace("\\", "/")+"/static"
PATH_OBJECT = pathlib.Path(PATH)
HOSTNAME = None
SERVICENAME = None

@app.before_request
def store_hostname():
    global HOSTNAME
    global SERVICENAME
    if HOSTNAME == None:
        HOSTNAME = request.scheme +"://" + request.host
    if SERVICENAME ==None:
        SERVICENAME = request.host.split(":")[0]
    
def lowermap(a:str):
    return a.lower()

def GreedUnit(type,name,link,last_mod="——",size="——"):
        if  type=="⮬":
            type = f"""<span class="arrow"></span>"""
        elif type == "folder":
            type = f"""<span class="folder"></span>"""
        else:
            spname =list(map(lowermap, name.split(".")))
            if len(spname)>0 and spname[-1]=="md":
                type = f"""<span class="mdfile"></span>"""
            elif len(spname)>0 and spname[-1]=="iso":
                type = f"""<span class="iso"></span>"""
            elif len(spname)>0 and (spname[-1]=="zip" or spname[-1]=="gz" or spname[-1]=="rar" or spname[-1]=="7z"):
                type = f"""<span class="zip"></span>"""
            else:
                type = f"""<span class="file"></span>"""

        res = f"""
        <a href="{link}">
        <div class="grid">
            <div>{type}</div>
            <div>{name}</div>
            <div style="text-align: right; margin-right: 10%;">{last_mod}</div>
            <div style="text-align: right; margin-right: 10%;">{size}</div>
        </div>
        </a>
        """    
        return res
def GetPrev(folder):
    return '/'.join(folder)

def search_by_name(search_path:str,query:str, host:str) -> list:
    l=[]
    query = query.lower()
    for root, dirs, files in os.walk(search_path):
       
        for i in files:
            if query in i.lower():
                abs_path = os.path.join(root, i)
                size = str(round(os.path.getsize(abs_path)/1024,1))+" KB"
                lmt =str(datetime.fromtimestamp(round(os.path.getmtime(abs_path))))

                l.append((i,host+os.path.relpath(abs_path, search_path).replace("\\","/"),lmt,size))
    return l

@cache.cached(timeout=300, key_prefix="file_list")
@app.route("/")
def index():
        search = request.args.get('search')
        if search!=None:
            search_results = search_by_name(PATH,search,HOSTNAME+"?folder=/")
            theme = request.cookies.get('theme', 'dark') 
            links=[]
            for i in range(len(search_results)):
                link = search_results[i][1]
                links.append(GreedUnit(type="file",name=search_results[i][0],link=link,last_mod=search_results[i][2],size=search_results[i][3]))

            return render_template("index.html", content=' <br> '.join(links),dirs=" ", title=f"{search} : {SERVICENAME}", hostname=HOSTNAME,theme=theme)
        folder = request.args.get('folder')
        if folder!=None:
             folder = folder.strip()
             folder = unquote(folder)
             if '..' in folder:
                abort(403, "Access denied")

        pat = PATH
        f_req = ""
        
        Req_URL = request.url.rstrip('/')
        
        if folder!=None and len(folder)>0 and folder[0]!="/":
            return redirect(f"{HOSTNAME}?folder=/{folder}")

        links = []
        pat+=folder or ""
        
        if pathlib.Path(pat).is_relative_to(PATH_OBJECT)==False:
            abort(403)
        if folder==None:
            f_req = "?folder="
        if path.isfile(pat):
            return send_file(path_or_file=pat,mimetype=None, as_attachment=True)

        items = listdir(pat)

        dirs = [x for x in items if path.isdir(path.join(pat, x))]
        files = [x for x in items if path.isfile(path.join(pat, x))]
        dirs.sort()
        files.sort()

        readme = ""
        
        if f_req == "" and (folder!="" and folder!="/"):
            splitted_url = Req_URL.split("/")
            
            splitted_url.pop()
            previous = GetPrev(splitted_url)
            links.append(GreedUnit(type=f"⮬", name="..", link = previous))
        
        for i in range(len(dirs)):
            link = f"{Req_URL}{f_req}/{quote(dirs[i])}"
            links.append(GreedUnit(type="folder",name=dirs[i],link=link, last_mod="——",size="——"))
        for i in range(len(files)):
            link = f"{Req_URL}{f_req}/{quote(files[i])}"
            filerealpath = pat+"/"+files[i]
            kbsize = round(path.getsize(filerealpath)/1024,1)
            if kbsize>=8192:
                size =str(round(kbsize/1024,1))+ "MB"
            else:
                size = str(kbsize)+" KB"
            lmt =datetime.fromtimestamp(round(path.getmtime(filerealpath)))
            if files[i].lower()=="readme.md":
                with open(filerealpath, "r") as file:
                    readme = markdown.markdown(file.read())
                    readme = f"""
                    <div style="position: relative; left: 50%; background: var(--readme-c); transform: translateX(-50%); width: 95%; border-radius: 4vh; margin-bottom: 4vh;">
                        <div style="height: 5vh; border-top-left-radius: 4vh; border-top-right-radius: 4vh; position: relative; top: 0vh;vertical-align: middle; padding: 3vh; padding-bottom: 0; line-height: 1; background-color: var(--readme-h);">
                            <span style="vertical-align:text-middle;" class="mdreadme">README.md</span>
                        </div>
                        <div style="position: relative; padding: 2vh;padding-top: 0%; width: auto; left: 50%; transform: translateX(-50%);">
                            {readme}
                        </div>
                    </div>
                    """
            links.append(GreedUnit(type="file",name=files[i], link=link,last_mod=lmt,size=size))

        splitted_url = Req_URL.split("/")
        
        dirs = []
        current_dir = f"{splitted_url[0]}//{splitted_url[2]}/{splitted_url[3]}"
        dirs.append(f"""<a  style="color: #fdfeff !important;" href={current_dir}>home</a>""")
        for i in range(4, len(splitted_url)):
            current_dir+=f"/{splitted_url[i]}"
            dirs.append(f"""<a href={current_dir}  style="color: #fdfeff !important;">{unquote(splitted_url[i])}</a>""")
           
        title = ""

        if len(splitted_url)>4:
            title = f"{splitted_url[-1]} · {SERVICENAME}"
        else:
            title = f"Home · {SERVICENAME}"

        theme = request.cookies.get('theme', 'dark')  
        return render_template("index.html", content=' <br> '.join(links), dirs="/".join(dirs), readme=readme, title=unquote(title), hostname=HOSTNAME, theme=theme) 
    



#FOR TEST
if __name__ == '__main__':
   
    app.run(debug=True, host='localhost', port=5000)
    


