import htm from 'https://unpkg.com/htm?module'
const html = htm.bind(h);

export default function({attrs,hooks}) { 
  var prospectName = attrs.ProspectName;
  var prospectDomain = attrs.ProspectDomain;
  var prospectOID = attrs.ProspectOID;
  var apps = attrs.Apps;
  var pages = attrs.Pages;
  var blocks = attrs.Blocks;
  var pageName = attrs.PageName;
  var pageOID = attrs.PageOID;
  var appName = attrs.AppName;
  var sections = attrs.Sections || [];

  hooks.oncreate = (v) => {
    $( ".pages" ).sortable({
      start: function(event, ui) {
          ui.item.startPos = ui.item.index();
      },
      stop: function(event, ui) {
          if (ui.item.startPos === ui.item.index()) {
            return;
          }
          reorderPage(ui.item[0].dataset["oid"], ui.item.index())
      }
    });

    $( "#sections" ).sortable({
      start: function(event, ui) {
          ui.item.startPos = ui.item.index();
      },
      stop: function(event, ui) {
          if (ui.item.startPos === ui.item.index()) {
            return;
          }
          reorderSection(ui.item[0].dataset["oid"], ui.item.startPos, ui.item.index())
      }
    });

    

    $(document).on('drag dragstart dragend dragover dragenter dragleave drop', function(e) {
      e.preventDefault();
      e.stopPropagation();
    })

    $(".block-drop").on('drop', handleDrop);

  }

  return html`
<div class="h-screen flex overflow-hidden bg-gray-100">      
  <div class="md:flex md:flex-shrink-0">
    <div class="flex flex-col w-64">
      
      <div class="flex flex-col h-0 flex-1 bg-gray-800">
        <div class="flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
          <div class="flex items-center flex-shrink-0 px-4">
              <img class="h-10 rounded-md w-auto mr-3 bg-white p-1" src="http://logo.clearbit.com/${prospectDomain}" /><span class="text-white font-medium text-2xl">${prospectName}</span>
          </div>
          <nav class="mt-5 flex-1 px-2 bg-gray-800 space-y-1">
              <h3 class="px-1 text-xs font-semibold text-gray-500 uppercase tracking-wider" id="projects-headline">
                  Demo Apps
              </h3>

              ${(apps||[]).map((app) => html`
                <div class="space-y-1">
                  <button class="group w-full flex items-center pr-2 py-2 text-sm font-medium rounded-md text-gray-300 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-indigo-500">
                    <svg class="rotate-90 mr-2 h-5 w-5 transform group-hover:text-gray-400 transition-colors ease-in-out duration-150" viewBox="0 0 20 20" aria-hidden="true">
                      <path d="M6 6L14 10L6 14V6Z" fill="currentColor" />
                    </svg>
                    ${app.Name}
                  </button>

                  <div class="space-y-1">
                    <div class="pages">
                    ${(pages[app.Name]||[]).map((page) => html`
                      <a data-oid="${page.OID}" href="?${app.Name}/${page.Name}" class="group w-full flex items-center pl-10 pr-2 py-2 text-sm font-medium text-gray-300 rounded-md hover:text-white hover:bg-gray-700">
                        ${page.Title}
                      </a>
                    `)}
                    </div>
                    
                    <a href="#" onclick=${() => newPage(app.OID, prompt("Name"))} class="group w-full flex items-center pl-10 pr-2 py-2 text-sm font-medium text-gray-500 rounded-md hover:text-white hover:bg-gray-700">
                      New Page
                    </a>

                  </div>
                </div>
              `)}
          
              
          </nav>
        </div>
        <div class="flex-shrink-0 flex bg-gray-700 p-4">
          <a href="#" class="flex-shrink-0 w-full group block">
            <div class="flex items-center">
              <a href="javascript:void" onclick=${() => newApp(prospectOID, prompt("App Name"))}>
              <div class="ml-4">
                  <span class="text-2xl mr-2 text-white font-semibold">${h.trust("&plus;")}</span>
                  <span class="text-lg font-medium text-white">
                      New App
                  </span>
              </div>
              </a>
            </div>
          </a>
        </div>
      </div>
    </div>
  </div>


  <div class="flex flex-col w-0 flex-1 overflow-hidden">
    <div class="md:hidden pl-1 pt-1 sm:pl-3 sm:pt-3">
      <button class="-ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md text-gray-500 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500">
        <span class="sr-only">Open sidebar</span>
        <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
    </div>
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="hidden sm:-my-px sm:ml-6 sm:flex sm:space-x-8">
              <a href="/" class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900">${h.trust("&larr;")} Back to Prospects</a>
            </div>
          </div>
          <div class="hidden sm:ml-6 sm:flex sm:items-center">
              <div class="flex-shrink-0 mr-2">
                <a target="_blank" href="#">
                  <button type="button" class="relative inline-flex items-center px-4 py-1 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-500 hover:bg-indigo-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-indigo-500">
                      <div class="transform rotate-90 text-lg mr-2">${h.trust("&rarrb;")}</div>
                      <span>Download SPA</span>
                  </button>
                </a>
              </div>
              <div class="flex-shrink-0 mr-2">
                <a target="_blank" href="/t/${prospectName}/${appName}/${pageName}">
                  <button type="button" class="relative inline-flex items-center px-4 py-1 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-500 hover:bg-indigo-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-indigo-500">
                      <span class="text-lg mr-2">${h.trust("&starf;")}</span>
                      <span>Visit App</span>
                  </button>
                </a>
              </div>
  
            
            <div class="ml-3 relative">
              <div>
                <button class="bg-white flex text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500" id="user-menu" aria-haspopup="true">
                  <span class="sr-only">Open user menu</span>
                  <img class="h-8 w-8 rounded-full" src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80" alt="" />
                </button>
              </div>
              
              <div class="hidden origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5" role="menu" aria-orientation="vertical" aria-labelledby="user-menu">
                <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" role="menuitem">Sign out</a>
              </div>
            </div>
          </div>
          
        </div>
      </div>
  
      
    </nav>

    <main class="flex-1 relative p-8 z-0 overflow-y-auto focus:outline-none" tabindex="0">

        <div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
          <div class="px-4 sm:px-6 md:px-0">
              <h1 class="text-2xl font-semibold text-gray-900">
                <span class="text-gray-400">${appName}</span> 
                <span class="text-gray-300 mx-1">/</span> 
                ${pageName}
              </h1>
            </div>
         
          <div class="py-4">
            <div class="hidden text-center pt-24 text-gray-300 font-semibold text-4xl border-4 border-dashed border-gray-200 rounded-lg h-64">
                No page sections
            </div>
            <div id="sections">
            ${sections.map((section) => html`
              <div data-oid=${section.OID} style=${{"marginLeft": "-2px"}} class="text-gray-300 mb-2 font-semibold text-4xl border-4 border-indigo-500 rounded-lg">
                <div class="py-1 bg-indigo-500 sm:px-6 h-6">
                  <div class="-mt-4 flex items-center justify-between flex-wrap sm:flex-nowrap">
                    <div class="flex-grow mt-1">
                      <h4 class="text-sm font-medium text-white">
                        ${section.Block.Name}
                      </h4>
                    </div>
                    <div class="-mt-1 flex-shrink-0">
                      <a href="javascript:null" onclick=${() => removeSection(section.OID, section.Key)} class="text-lg text-white">${h.trust("&otimes;")}</a>
                    </div>
                  </div>
                </div>
                <iframe id="${section.Block.Name.split(".")[0]}" onload=${iframeLoad} class="w-full h-0" src="/t/${prospectName}/${appName}/${pageName}?section=${section.Key}&live=0"></iframe>
              </div>
            `)}
            </div>
          </div>
         
          <a href="javascript:null" onclick=${() => $("#picker").show()}>
          <button type="button" class="w-full inline-flex items-center px-4 py-1 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-500 hover:bg-indigo-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-indigo-500">
              <div class="mx-auto">
                  <span class="text-xl mr-2 text-white font-semibold">${h.trust("&plus;")}</span>
                  <span class="text-xl">New Section</span>
              </div>
          </button>
          </a>
        </div>




    </main>
    
  </div>

  <div id="picker" onclick=${() => $("#picker").hide()} class="hidden fixed z-10 inset-0 overflow-y-auto">
    <div  class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 transition-opacity" aria-hidden="true">
        <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
      </div>
  

      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">${h.trust("&#8203;")}</span>

      <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full" role="dialog" aria-modal="true" aria-labelledby="modal-headline">
        <div class="bg-white px-4 py-5 border-b border-gray-200 sm:px-6">
            <h3 class="text-lg leading-6 font-medium text-gray-900">
              Choose Block
            </h3>
        </div>

        
        
        <div class="px-4 py-5 sm:p-6 block-drop" data-pageid=${pageOID}>

          <ul>
            ${(blocks||[]).map((b) => html`
              <li><a href="javascript:void" onclick=${() => addSection(pageOID, b.OID)}>${b.Name}</a></li>
            `)}
          </ul>
            
                        
        </div>
                      
        
          
      </div>
      
    </div>
  </div>
</div>`
}

function iframeLoad(e) {
  var iframe = e.target;
  setTimeout(() => {
    iframe.style.height = iframe.contentWindow.document.body.scrollHeight + 'px';
  }, 500);
}

function removeSection(id, key) {
  m.request({
    method: "POST",
    url: "/c/object.component.remove",
    body: {
        ID: id,
        Component: key
    },
  }).then(() => location.reload());
}

function reorderPage(id, idx) {
  m.request({
    method: "POST",
    url: "/c/object.move",
    body: {
        ID: id,
        Index: idx
    },
  }).then(() => location.reload());
}


function newPage(id, name) {
  m.request({
    method: "POST",
    url: "/c/zt100.new-page",
    body: {
        ID: id,
        Name: name.toLowerCase()
    },
  }).then(() => location.reload());
}

window.addEventListener('resize', () => {
  document.querySelectorAll("iframe").forEach((iframe) => {
    iframe.style.height = "0px";
    iframe.style.height = iframe.contentWindow.document.body.scrollHeight + 'px';
  })
});

function addSection(pageId, blockId) {
  m.request({
    method: "POST",
    url: "/c/zt100.new-section",
    body: {
        PageID: pageId,
        BlockID: blockId
    },
  }).then(() => location.reload());
}

function newApp(id, name) {
  m.request({
    method: "POST",
    url: "/c/zt100.new-app",
    body: {
        ID: id,
        Name: name.toLowerCase()
    },
  }).then(() => location.reload());
}

function reorderSection(id, from, to) {
  m.request({
    method: "POST",
    url: "/c/object.component.move",
    body: {
        ID: id,
        From: from+1,
        To: to+1
    },
  }).then(() => location.reload());
}

function handleDrop(e) {
  let files = e.originalEvent.dataTransfer.files;
  if (files.length === 1) {
    let formData = new FormData()
    let el = e.target.closest(".block-drop");
    formData.append('file', files[0])
    formData.append('PageID', el.dataset["pageid"])
    formData.append('BlockID', "bv7br36g10l83cq306t0")
    fetch("/c/zt100.new-section?upload=1", {
      method: 'POST',
      body: formData
    })
  }
}
