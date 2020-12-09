
import htm from 'https://unpkg.com/htm?module'
const html = htm.bind(h);

export default function() { 
  return html`
    <div class="relative bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6">
        <div class="flex justify-between items-center border-b-2 border-gray-100 py-6 md:justify-start md:space-x-10">
          <div class="flex justify-start lg:w-0 lg:flex-1">
            <a href="#">
              <span class="sr-only">Workflow</span>
              <img class="h-8 w-auto sm:h-10" src="https://logo.clearbit.com/${window.config.ProspectDomain}" alt="" />
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
            ${window.config.PageMenu.map((item) => html`
              <a href="./${item.Page}" class="text-base font-medium text-gray-500 hover:text-gray-900">
                ${item.Title}
              </a>
            `)} 
      
          </nav>
          <div class="hidden md:flex items-center justify-end md:flex-1 lg:w-0">
            <a href="#" class="whitespace-nowrap text-base font-medium text-gray-500 hover:text-gray-900">
              Sign in
            </a>
            <a href="#" class="ml-8 whitespace-nowrap inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-primary-500 hover:bg-primary-500">
              Sign up
            </a>
          </div>
        </div>
      </div>
        
        
      <div class="absolute top-0 inset-x-0 p-2 transition transform origin-top-right md:hidden">
        <div class="rounded-lg shadow-lg ring-1 ring-black ring-opacity-5 bg-white divide-y-2 divide-gray-50">
          <div class="pt-5 pb-6 px-5">
            <div class="flex items-center justify-between">
              <div>
                <img class="h-8 w-auto" src="https://tailwindui.com/img/logos/workflow-mark-primary-600.svg" alt="Workflow" />
              </div>
              <div class="-mr-2">
                <button type="button" class="bg-white rounded-md p-2 inline-flex items-center justify-center text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500">
                  <span class="sr-only">Close menu</span>
                  <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
          <div class="py-6 px-5 space-y-6">
            <div class="grid grid-cols-2 gap-y-4 gap-x-8">
              <a href="#" class="text-base font-medium text-gray-900 hover:text-gray-700">
                Pricing
              </a>
      
              <a href="#" class="text-base font-medium text-gray-900 hover:text-gray-700">
                Docs
              </a>
      
              <a href="#" class="text-base font-medium text-gray-900 hover:text-gray-700">
                Enterprise
              </a>
      
              <a href="#" class="text-base font-medium text-gray-900 hover:text-gray-700">
                Blog
              </a>
      
              <a href="#" class="text-base font-medium text-gray-900 hover:text-gray-700">
                Help Center
              </a>
      
              <a href="#" class="text-base font-medium text-gray-900 hover:text-gray-700">
                Guides
              </a>
      
              <a href="#" class="text-base font-medium text-gray-900 hover:text-gray-700">
                Security
              </a>
      
              <a href="#" class="text-base font-medium text-gray-900 hover:text-gray-700">
                Events
              </a>
            </div>
            <div>
              <a href="#" class="w-full flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-primary-600 hover:bg-primary-700">
                Sign up
              </a>
              <p class="mt-6 text-center text-base font-medium text-gray-500">
                Existing customer?
                <a href="#" class="text-primary-600 hover:text-primary-500">
                  Sign in
                </a>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>`
}