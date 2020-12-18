package http

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/progrium/watcher"
	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
)

type FileServer struct {
	Filesystem  afero.Fs
	LiveUpdates bool
	Debug       bool

	liveClients sync.Map
	watcher     watchable
	unwatch     func()
}

type watchable interface {
	Watch(name string, watch func(watcher.Event)) (func(), error)
}

func (c *FileServer) ComponentEnable() {
	if c.Filesystem == nil {
		return
	}
}

func (c *FileServer) Open(name string) (f http.File, err error) {
	f, err = c.Filesystem.Open(name)
	if err != nil {
		return
	}
	c.watcher, _ = c.Filesystem.(watchable)
	if c.LiveUpdates && c.watcher != nil {
		// TODO: keep returned unwatches and call on ComponentDisable
		_, err = c.watcher.Watch(name, c.fsWatch)
		if err != nil {
			log.Println(err)
			return
		}
		// TODO: detect directory/index and if no body tag, just append
		if filepath.Ext(name) == ".html" {
			b, err := afero.ReadAll(f)
			if err != nil {
				return f, err
			}
			re := regexp.MustCompile(`(?i)</body>`)
			content := re.ReplaceAllString(string(b), "<script src=\"/.live-updates.js\"></script>\n</body>")
			ff := mem.NewFileHandle(mem.CreateFile(name))
			_, err = ff.Write([]byte(content))
			if err != nil {
				return f, err
			}
			ff.Seek(0, 0)
			return ff, nil
		}
	}
	return
}

func (c *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.Method, r.URL)
	switch r.URL.Path {
	case "/.live-updates.js":
		c.serveClient(w, r)
	case "/.live-updates":
		c.serveUpdatesWebsocket(w, r)
	default:
		// TODO: serve root
		// http.StripPrefix("/files/", http.FileServer(c.FileSystem)).ServeHTTP(w, r)
		http.FileServer(c).ServeHTTP(w, r)
	}
}

func (c *FileServer) serveClient(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("client").Parse(liveUpdateClientSrc))
	w.Header().Set("content-type", "text/javascript")
	if err := tmpl.Execute(w, map[string]interface{}{
		"Debug":    c.Debug,
		"Endpoint": fmt.Sprintf("ws://%s/.live-updates", r.Host),
	}); err != nil {
		log.Println(err)
	}
}

func (c *FileServer) serveUpdatesWebsocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()
	ch := make(chan string)
	c.liveClients.Store(ch, struct{}{})
	if c.Debug {
		log.Println("new liveupdate connection")
	}

	for filepath := range ch {
		log.Println("RECEIVING", filepath)
		err := conn.WriteJSON(map[string]interface{}{
			"path": filepath,
		})
		log.Println("WRITTEN", filepath, err)
		if err != nil {
			c.liveClients.Delete(ch)
			if !strings.Contains(err.Error(), "broken pipe") {
				log.Println(err)
			}
			return
		}
	}
}

func (c *FileServer) fsWatch(event watcher.Event) {
	c.liveClients.Range(func(k, v interface{}) bool {
		log.Println("SENDING", event.Path)
		k.(chan string) <- event.Path
		return true
	})
}

var liveUpdateClientSrc = `
let listeners = {};
let refreshers = [];
let ws = undefined;
let debug = {{if .Debug}}true{{else}}false{{end}};

const scheduleRetry = (fn, r) => setTimeout(() => fn(r), {1:200,2:200,3:300,4:500}[r]||1000);

(function connect(retry=0) {
	ws = new WebSocket("{{.Endpoint}}");
	this.retry = retry;
    if (debug) {
        ws.onopen = () => {
			this.retry = 0;
			console.debug("liveupdates websocket open");
		}
        ws.onclose = () => {
			console.debug("liveupdates websocket closed, retrying["+this.retry+"]...");
			scheduleRetry(connect, this.retry+1);
		}
    }
    //ws.onerror = (err) => console.debug("liveupdates websocket error: ", err);
    ws.onmessage = async (event) => {
        let msg = JSON.parse(event.data);
        if (debug) {
            console.debug("liveupdates trigger:", msg.path);
        }
        let paths = Object.keys(listeners);
        paths.sort((a, b) => b.length - a.length);
        for (const idx in paths) {
            let path = paths[idx];
            if (msg.path.startsWith(path)) {
                for (const i in listeners[path]) {
                    await listeners[path][i]((new Date()).getTime(), msg.path);
                }
            }
        }
        // wtf why aren't refreshers consistently 
        // run after listeners are called.
        // setTimeout workaround seems ok for now
        setTimeout(() => refreshers.forEach((cb) => cb()), 20);
    }; 
})();  

function accept(path, cb) {
    if (listeners[path] === undefined) {
        listeners[path] = [];
    }
    listeners[path].push(cb);
}

function refresh(cb) {
    refreshers.push(cb);
    cb();
}

(function watchHTML() {
    let withIndex = "";
    if (location.pathname[location.pathname.length-1] == "/") {
        withIndex = location.pathname + "index.html";
    } else {
        withIndex = location.pathname + "/index.html";
    }
    accept(location.pathname, (ts, path) => {
        if (path == location.pathname || path == withIndex) {
            location.reload();
        }
    });
})();

// TODO: this loads any css changed on the watched filesystem, including non-linked stylesheets!
(function watchCSS() {
    accept("", (ts, path) => {
        if (path.endsWith(".css")) {
            let link = document.createElement('link');
            link.setAttribute('rel', 'stylesheet');
            link.setAttribute('type', 'text/css');
            link.setAttribute('href', path+'?'+(new Date()).getTime());
            document.getElementsByTagName('head')[0].appendChild(link);
            let styles = document.getElementsByTagName("link");
            for (let i=0; i<styles.length; i++) {
                if (i < styles.length-1 && styles[i].getAttribute("href").startsWith(path)) {
                    styles[i].remove();
                }
            }
        }
    });
})();

`
