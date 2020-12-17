import {HasFeature} from '/static/util.js';

export default function() { 
  return {
    "view": function({attrs}) {
      return html`
    <div class="relative bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6">
        <div class="flex justify-between items-center border-b-2 border-gray-100 py-6 md:justify-start md:space-x-10">
          <div class="flex justify-start lg:w-0 lg:flex-1">
            <a href="#">
              <span class="sr-only">Workflow</span>
              <img class="h-8 w-auto sm:h-10" src="https://logo.clearbit.com/${window.data.Demo.Domain}" alt="" />
            </a>
          </div>
          <div class="-mr-2 -my-2 md:hidden">
            <button type="button" class="bg-white rounded-md p-2 inline-flex items-center justify-center text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500">
              <span class="sr-only">Open menu</span>
              <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
          <nav class="hidden md:flex space-x-10">
            ${window.data.Menu.map((item) => html`
              <a href="./${item.Page}" class="text-base font-medium text-gray-500 hover:text-gray-900">
                ${item.Title}
              </a>
            `)} 
      
          </nav>
          <div class="hidden md:flex items-center justify-end md:flex-1 lg:w-0">
            ${(data.Contrib.Auth.IsAuthenticated) ? html`
              ${HasFeature("login") && html`
                <a href="/feature/login/logout?then=${location.pathname}" class="whitespace-nowrap text-base font-medium text-gray-500 hover:text-gray-900 mr-8">
                  Logout
                </a>
              `}
              <button class="bg-white flex text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500" id="user-menu" aria-haspopup="true">
                <img class="h-8 w-8 rounded-full" src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80" alt="" />
              </button>
            `: html`
              ${HasFeature("register") && html`
                <a href="#" class="whitespace-nowrap text-base font-medium text-gray-500 hover:text-gray-900">
                  Register
                </a>              
              `}  
              ${HasFeature("login") && html`
                <a href="./login" class="ml-8 whitespace-nowrap inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-primary-500 hover:bg-primary-500">
                  Login
                </a>
              `}
            `}
                      
          </div>
        </div>
      </div>
        
        
    </div>`
    }
  }
}